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

// TODO: Сейчас фрагменты строк реализованы как элементы связанного списка. Это сделано для упрощения
//  "разрезания" текущего фрагмента при вставке surrounded image. Но поскольку мы полностью перестраиваем
//  этот список каждый раз, когда переходим к новой строке, кажется, что поддерживать список это чрезмерная
//  сложность. Возможно, будет достаточно держать в руках только текущий фрагмент и пересчитывать его по
//  мере продвижения вдоль строки.

type fragmentFactory struct {
	free *fragment
}

func (ff *fragmentFactory) new(x, width int16, next *fragment) *fragment {
	if ff.free == nil {
		buf := make([]fragment, 1024/unsafe.Sizeof(fragment{})) // XXX много/мало?
		for i := 1; i < len(buf); i++ {
			buf[i-1].next = &buf[i]
		}
		ff.free = &buf[0]
	}

	return &fragment{
		owner: ff,
		x:     x,
		width: width,
		next:  next,
	}
}

type fragment struct {
	owner *fragmentFactory
	next  *fragment
	x     int16
	width int16
	curX  int16 // relative current x position to insert
}

func (f *fragment) free() {
	f.next = f.owner.free
	f.owner.free = f
}

func (f *fragment) devide(x1, x2 int16) *fragment {

	if !(x1 < x2 && f.x <= x1 && x2 <= f.x+f.width) {
		panic(fmt.Sprintf("fragment %+v devide by [%d, %d]: bad segment", *f, x1, x2))
	}

	next := f.owner.new(x2, (f.x+f.width)-x2, f.next)
	f.width = x1 - f.x
	f.next = next

	return next
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

type page struct {
	width             int16
	defaultLineHeight int16
	defaultCharWidth  int16
	curLineY          int16
	nextLineY         int16
	nextParagraphY    int16
	floatingBase      point // base position to insert floating image
	fragmentFactory   fragmentFactory
	curFragment       *fragment
	surroundedImages  []image
	imagePoints       []point
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
	list := p.surroundedImages

	// предварительно очистим список от изображений которые ушли вверх
	y := p.curLineY
	n := len(list) - 1
	for i := n; i >= 0; i-- {
		if img := &list[i]; y >= img.y+img.height {
			list[i] = list[n]
			list = list[:n]
			n--
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].x < list[j].x
	})

	p.surroundedImages = list

	p.freeFragments()
	f := p.fragmentFactory.new(0, p.width, nil)
	p.curFragment = f

	for i := range list {
		img := &list[i]
		f.devide(img.x, img.x+img.width)
		f = f.next
	}
}

func (p *page) freeFragments() {
	for p.curFragment != nil {
		f := p.curFragment
		p.curFragment = f.next
		f.free()
	}
}

func (p *page) arrangeImage(img image) {

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
		p.curFragment = p.curFragment.devide(x, x+img.width)

		p.surroundedImages = append(p.surroundedImages, img)
		p.updateNextParagraphY(y + img.height)

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
	p.imagePoints = append(p.imagePoints, img.point)
}

func (p *page) arrangeWord(s string) {
	wordWidth := int16(len(s)) * p.defaultCharWidth

	x, y := p.arrangeWidth(wordWidth, p.defaultCharWidth)
	p.floatingBase = point{x + wordWidth, y}
}

func (p *page) arrangeWidth(width, prefix int16) (x, y int16) {

	x = p.curFragment.arrange(width, prefix)
	for x == -1 {
		p.curFragment = p.curFragment.next
		if p.curFragment == nil {
			p.nextLine()
		}
		x = p.curFragment.arrange(width, prefix)
	}

	return x, p.curLineY
}

func solution(pageWidth, lineHeight, charWidth int16, buf []byte) ([]point, error) {
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

	return page.imagePoints, nil
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
