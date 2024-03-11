package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

const imageTag = "(image"

type layout byte

const (
	_ layout = iota
	embedded
	surrounded
	floating
)

type image struct {
	rectangle
	layout layout
}

func paseLayout(s string) (layout, error) {
	switch s {
	case "embedded":
		return embedded, nil
	case "surrounded":
		return surrounded, nil
	case "floating":
		return floating, nil
	}
	return 0, fmt.Errorf("unknown layout: %s", s)
}

func parseInt16(s string) (int16, error) {
	v, err := strconv.Atoi(s)
	return int16(v), err
}

func parseImageParams(s string) (image, error) {
	img := image{}

	for i := 0; i < len(s); {

		for _, c := range s[i:] {
			if !unicode.IsSpace(c) {
				break
			}
			i++
		}

		start := i
		for _, c := range s[i:] {
			if unicode.IsSpace(c) {
				break
			}
			i++
		}

		if start < len(s) {
			const (
				name  = 0
				value = 1
			)

			var err error

			switch p := strings.Split(s[start:i], "="); p[name] {
			case "width":
				img.width, err = parseInt16(p[value])
			case "height":
				img.height, err = parseInt16(p[value])
			case "layout":
				img.layout, err = paseLayout(p[value])
			case "dx":
				img.x, err = parseInt16(p[value])
			case "dy":
				img.y, err = parseInt16(p[value])
			}

			if err != nil {
				return img, err
			}
		}
	}

	return img, nil
}

type point struct {
	x, y int16
}

func (p *point) add(o point) point {
	return point{p.x + o.x, p.y + o.y}
}

type rectangle struct {
	point
	width  int16
	height int16
}

func (r *rectangle) updateHeight(h int16) {
	if r.height < h {
		r.height = h
	}
}

type page struct {
	rectangle
	defaultHeight int16
	curParagraph  paragrath
}

func (p *page) reset(r rectangle) {
	p.rectangle = r
	p.defaultHeight = r.height
	p.curParagraph.reset(rectangle{point{0, 0}, r.width, r.height})
}

func (p *page) arrange(r rectangle, layout layout, intend int16) (point, bool) {

	pnt, ok := p.curParagraph.arrange(r, layout, intend)
	if ok {
		p.updateHeight(p.curParagraph.x + p.curParagraph.height)
	}

	return p.point.add(pnt), true
}

func (p *page) nextParagraph() {
	y := p.curParagraph.y + p.curParagraph.height
	p.curParagraph.reset(rectangle{point{0, y}, p.width, p.defaultHeight})
	p.curParagraph.y = y
}

type paragrath struct {
	rectangle
	defaultHeight     int16
	curLine           line
	surrounded        []rectangle
	surroundedChanged bool
	floatingBase      point
}

func (p *paragrath) reset(r rectangle) {

	if debugEnable {
		log.Printf("paragrath.reset %+v", r)
	}

	p.rectangle = r
	p.defaultHeight = r.height
	p.floatingBase = point{0, 0}
	p.surrounded = p.surrounded[:0]
	p.surroundedChanged = false
	p.curLine.reset(rectangle{point{0, 0}, r.width, r.height})
}

func (p *paragrath) arrange(r rectangle, layout layout, intend int16) (point, bool) {

	if layout == floating {
		return p.arrangeFloating(r)
	}

	for {
		pnt, ok := p.curLine.arrange(r, layout, intend)
		if ok {
			if layout == surrounded {
				p.surrounded = append(p.surrounded, rectangle{
					point:  pnt,
					width:  r.width,
					height: r.height,
				})
				p.surroundedChanged = true
			}

			p.floatingBase = pnt.add(point{r.width, 0})
			p.updateHeight(p.curLine.y + p.curLine.height)

			return p.point.add(pnt), ok
		}

		p.nextLine()
	}
}

func (p *paragrath) arrangeFloating(r rectangle) (point, bool) {
	pnt := p.floatingBase.add(r.point)

	if pnt.x < 0 {
		pnt.x = 0
	} else if x2 := pnt.x + r.width; x2 > p.width {
		pnt.x -= x2 - p.width
	}

	p.floatingBase = pnt.add(point{r.width, 0})
	p.updateHeight(pnt.y + r.height)

	return p.point.add(pnt), true
}

func (p *paragrath) nextLine() {
	y := p.curLine.y + p.curLine.height
	p.curLine.reset(rectangle{point{0, y}, p.width, p.defaultHeight})

	if p.surroundedChanged {
		sort.Slice(p.surrounded, func(i, j int) bool {
			return p.surrounded[i].x < p.surrounded[j].x
		})
		p.surroundedChanged = false
	}

	for i := range p.surrounded {
		r := &p.surrounded[i]
		if r.y+r.height > y {
			p.curLine.fragments.splitLast(r.x-p.x, r.width)
		}
	}
}

type line struct {
	rectangle
	fragments fragments
}

func (l *line) reset(r rectangle) {
	if debugEnable {
		log.Printf("line.reset %+v", r)
	}

	l.rectangle = r
	l.fragments.reset(rectangle{point{0, 0}, r.width, r.height})
}

func (l *line) arrange(r rectangle, layout layout, intend int16) (point, bool) {

	if r.width > l.width {
		// Нам гарантируют, что такое не произойдет
		panic(fmt.Sprintf("*** logical error ***\n"+
			"Can't arrange %+v in line %+v. Width too long", r, l.rectangle))
	}

	pnt, ok := l.fragments.arrange(r, layout, intend)
	if ok {
		l.updateHeight(l.fragments.height)
	}

	return l.point.add(pnt), ok
}

type fragments struct {
	rectangle
	items       []fragment
	curFragment int16
}

func (fs *fragments) reset(r rectangle) {
	fs.rectangle = r
	fs.items = fs.items[:0]

	fs.items = append(fs.items, fragment{
		rectangle: rectangle{
			point:  point{0, 0},
			width:  r.width,
			height: r.height,
		},
	})

	fs.curFragment = 0
}

func (fs *fragments) splitLast(x, width int16) {
	n := len(fs.items)
	last := &fs.items[n-1] // panic if n == 0

	w1 := x - last.x
	w2 := (last.x + last.width) - (x + width)

	if w1 > last.width || w1 < 0 || w2 < 0 {
		// Это логическая ошибка. Такое никогда не должно происходить
		panic(fmt.Sprintf("*** logical error ***\n"+
			"Can't split fragment %+v using segment {x:%d width:%d}.\n"+
			"One of the ends of the segment outside the fragment", *last, x, width))
	}

	last.width = w1

	fs.items = append(fs.items, fragment{
		rectangle: rectangle{
			point:  point{x + width, 0},
			width:  w2,
			height: last.height,
		},
	})
}

func (fs *fragments) arrange(r rectangle, layout layout, intend int16) (point, bool) {

	for ; int(fs.curFragment) < len(fs.items); fs.curFragment++ {
		f := &fs.items[fs.curFragment]
		pnt, ok := f.arrange(r, layout, intend)

		if ok {
			switch layout {
			case embedded:
				fs.updateHeight(f.height)
			case surrounded:
				f.trimLeft()
			}

			return fs.point.add(pnt), true
		}
	}

	return point{}, false
}

type fragment struct {
	rectangle
	curX int16
}

func (f *fragment) arrange(r rectangle, layout layout, intend int16) (point, bool) {

	x := f.curX
	if x != 0 {
		x += intend
	}

	if x+r.width > f.width {
		return point{}, false
	}

	if layout == embedded {
		f.height = r.height
	}

	f.curX = x + r.width

	return point{f.x + x, f.y}, true
}

func (f *fragment) trimLeft() {
	f.x += f.curX
	f.width -= f.curX
	f.curX = 0
}

func solution(pageWidth, lineHeight, charWidth int16, buf []byte) ([]point, error) {
	t := unsafeString(buf)
	var points []point

	page := new(page)
	page.reset(rectangle{point{0, 0}, pageWidth, lineHeight})

	for i := 0; i < len(t); {

		eol := 0
		for _, c := range t[i:] {
			if !unicode.IsSpace(c) {
				break
			}
			if c == '\n' {
				eol++
			}
			i++
		}

		if i == len(t) {
			break
		}

		if eol > 1 {
			page.nextParagraph()
		}

		wordStart := i
		for _, c := range t[i:] {
			if unicode.IsSpace(rune(c)) {
				break
			}
			i++
		}

		switch t[wordStart:i] {
		case imageTag:
			imageStart := i
			j := strings.IndexByte(t[i:], ')')
			if j == -1 {
				return nil, errors.New("closing parenthesis was not found")
			}
			i += j + 1

			img, err := parseImageParams(t[imageStart : i-1])
			if err != nil {
				return nil, err
			}

			var indent int16
			if img.layout == embedded {
				indent = charWidth
			}

			pnt, ok := page.arrange(img.rectangle, img.layout, indent)
			if !ok {
				return nil, fmt.Errorf("can't arrange image %+v", img)
			}

			points = append(points, pnt)

			if debugEnable {
				log.Printf("image %+v at %v", img, pnt)
			}

		default:
			word := t[wordStart:i]
			r := rectangle{width: int16(len(word)) * charWidth, height: lineHeight}

			pnt, ok := page.arrange(r, embedded, charWidth)
			if !ok {
				return nil, fmt.Errorf("can't arrange word %q %+v", word, r)
			}

			if debugEnable {
				log.Printf("word %q %+v at %v", word, r, pnt)
			}
		}
	}

	return points, nil
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func run(in io.Reader, out io.Writer) error {

	// Размер входного файла не превышает 1000 байт.
	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	// Первая строка входного файла содержит три целых числа...
	var w, h, c int16
	if _, err := fmt.Fscan(bytes.NewReader(b), &w, &h, &c); err != nil {
		return err
	}

	// skip first line
	b = b[bytes.IndexByte(b, '\n')+1:]

	images, err := solution(w, h, c, b)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(b[:0]) // reuse the buffer for writing
	for i := range images {
		fmt.Fprintln(buf, images[i].x, images[i].y)
	}

	buf.WriteTo(out)
	return nil
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
