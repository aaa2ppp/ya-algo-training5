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

type point struct {
	x, y int16
}

type image struct {
	point
	width  int16
	height int16
	dx, dy int16
	layout layout
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
				img.dx, err = parseInt16(p[value])
			case "dy":
				img.dy, err = parseInt16(p[value])
			}

			if err != nil {
				return img, err
			}
		}
	}

	return img, nil
}

type fragment struct {
	x     int16
	width int16
	curX  int16 // relative current x position to insert
}

func (f *fragment) trimLeft() {
	f.x += f.curX
	f.width -= f.curX
	f.curX = 0
}

func (f *fragment) arrange(width, prefix int16) int16 {

	x := f.curX // relative
	if x != 0 {
		x += prefix
	}

	if x+width > f.width {
		return -1
	}

	f.curX = x + width
	return f.x + x
}

type fragments []fragment

func (fs *fragments) devideLast(x, width int16) {
	n := len(*fs)
	f := &(*fs)[n-1] // panic if n == 0

	w1 := x - f.x
	w2 := (f.x + f.width) - (x + width)

	if w1 > f.width || w1 < 0 || w2 < 0 {
		panic(fmt.Sprintf("fragment %+v devide by {x:%d width:%d}: bad segment", *f, x, width))
	}

	f.width = w1
	*fs = append(*fs, fragment{x: x + width, width: w2})
}

type page struct {
	width             int16
	defaultLineHeight int16
	defaultCharWidth  int16
	curLineY          int16
	nextLineY         int16
	nextParagraphY    int16
	floatingBase      point // base position to insert floating image
	fragments         fragments
	curFragment       int16
	surroundedImages  []*image
	images            []image
}

func (p *page) updateNextParagraphY(y int16) {
	p.nextParagraphY = max(p.nextParagraphY, y)
}

func (p *page) updateNextLineY(y int16) {
	p.nextLineY = max(p.nextLineY, y)
	p.updateNextParagraphY(p.nextLineY)
}

func (p *page) nextParagraph() {
	p.nextLineY = p.nextParagraphY
	p.floatingBase = point{0, p.nextParagraphY}
	p.surroundedImages = p.surroundedImages[:0]
	p.nextLine()
}

func (p *page) nextLine() {
	p.curLineY = p.nextLineY
	p.updateNextLineY(p.curLineY + p.defaultLineHeight)
	p.calcFragments()
}

// Перестраивает список фрагментов строки. Должно вызыватся только при преходе к новой строке.
func (p *page) calcFragments() {
	p.prepareSurroundedImages()
	p.resetFragments()

	for _, img := range p.surroundedImages {
		p.fragments.devideLast(img.x, img.width)
	}
}

func (p *page) prepareSurroundedImages() {
	list := p.surroundedImages

	// Предварительно очистим список от изображений которые ушли вверх
	y := p.curLineY
	n := len(list) - 1
	for i := n; i >= 0; i-- {
		if img := list[i]; y >= img.y+img.height {
			list[i] = list[n]
			list = list[:n]
			n--
		}
	}

	// TODO: От сортировки можно избавится. Если учесть, что все новые изобращения
	// добавляются по порядку, то здесь можно просто слить два списка старый и новый
	sort.Slice(list, func(i, j int) bool {
		return list[i].x < list[j].x
	})

	p.surroundedImages = list
}

func (p *page) resetFragments() {
	p.fragments = p.fragments[:0]
	p.fragments = append(p.fragments, fragment{width: p.width})
	p.curFragment = 0
}

func (p *page) arrangeImage(image image) {
	p.images = append(p.images, image)
	img := &p.images[len(p.images)-1]

	switch img.layout {
	case embedded:
		x, y := p.arrangeWidth(img.width, p.defaultCharWidth)
		img.x = x
		img.y = y
		p.updateNextLineY(y + img.height)

	case surrounded:
		x, y := p.arrangeWidth(img.width, 0)
		img.x = x
		img.y = y
		p.updateNextParagraphY(y + img.height)
		p.surroundedImages = append(p.surroundedImages, img)
		p.fragments[p.curFragment].trimLeft()

	case floating:
		x := p.floatingBase.x + img.dx
		y := p.floatingBase.y + img.dy

		if x < 0 {
			x = 0
		} else if x+img.width > p.width {
			x = p.width - img.width
		}

		img.x = x
		img.y = y
	}

	p.floatingBase = point{img.x + img.width, img.y}
}

func (p *page) arrangeWord(s string) {
	wordWidth := int16(len(s)) * p.defaultCharWidth

	x, y := p.arrangeWidth(wordWidth, p.defaultCharWidth)
	p.floatingBase = point{x + wordWidth, y}
}

func (p *page) arrangeWidth(width, prefix int16) (x, y int16) {

	x = p.fragments[p.curFragment].arrange(width, prefix)
	for x == -1 {
		p.curFragment++
		if int(p.curFragment) == len(p.fragments) {
			p.nextLine()
		}
		x = p.fragments[p.curFragment].arrange(width, prefix)
	}

	return x, p.curLineY
}

func solution(pageWidth, lineHeight, charWidth int16, buf []byte) ([]image, error) {
	t := unsafeString(buf)

	page := page{
		width:             pageWidth,
		defaultLineHeight: lineHeight,
		defaultCharWidth:  charWidth,
	}
	page.nextParagraph()

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

			page.arrangeImage(img)

		default:
			page.arrangeWord(t[wordStart:i])
		}
	}

	return page.images, nil
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
