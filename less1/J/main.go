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

// XXX это порно! но примеры проходит, посмотрим на тесты...

const imageTag = "(image"

type layout int

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

type image struct {
	x, y   int
	width  int
	height int
	layout layout
	dx, dy int
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func parseImageParams(s string) (image, error) {
	img := image{}

	i := 0
	for i < len(s) {

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
				img.width, err = strconv.Atoi(p[value])
			case "height":
				img.height, err = strconv.Atoi(p[value])
			case "layout":
				img.layout, err = paseLayout(p[value])
			case "dx":
				img.dx, err = strconv.Atoi(p[value])
			case "dy":
				img.dy, err = strconv.Atoi(p[value])
			}

			if err != nil {
				return img, err
			}
		}
	}

	return img, nil
}

type page struct {
	width            int
	lineHeight       int
	charWidth        int
	nextParagraphY   int
	curLine          line
	curFragment      *fragment
	images           []image
	surroundedImages []image
	lastX, lastY     int // позиция от которой осчитываем floating
}

func (p *page) nextParagraph() {
	p.curLine = line{
		y:      p.nextParagraphY,
		height: p.lineHeight,
	}

	p.curFragment = &fragment{
		width: p.width,
	}

	p.lastX = 0
	p.lastY = p.curLine.y

	p.nextParagraphY += p.lineHeight
	p.surroundedImages = p.surroundedImages[:0]
}

func (p *page) nextLine() {
	y := p.curLine.y + p.curLine.height
	p.curLine = line{
		y:      y,
		height: p.lineHeight,
	}

	p.nextParagraphY = max(p.nextParagraphY, p.curLine.y + p.curLine.height) // update position of next paragrath

	p.lastX = 0
	p.lastY = p.curLine.y

	p.curFragment = &fragment{
		width: p.width,
	}

	var images []*image
	for i := range p.surroundedImages {
		img := &p.surroundedImages[i]

		if img.layout != surrounded || !(img.y <= y && y < img.y+img.height) {
			continue
		}

		images = append(images, img)
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].x < images[j].x
	})

	f := p.curFragment
	for _, img := range images {
		f.devide(img.x, img.x+img.width)
		f = f.next
	}
}

func (p *page) arrangeImage(img image) {
	switch img.layout {
	case embedded:
		x, y := p.arrangeWidth(img.width, p.charWidth)
		img.x = x
		img.y = y
		p.curLine.height = max(p.curLine.height, img.height)
		p.nextParagraphY = max(p.nextParagraphY, img.y+img.height)

	case surrounded:
		x, y := p.arrangeWidth(img.width, 0)
		img.x = x
		img.y = y
		p.curFragment.devide(x, x+img.width)
		p.curFragment = p.curFragment.next
		p.nextParagraphY = max(p.nextParagraphY, img.y+img.height)
		p.surroundedImages = append(p.surroundedImages, img)

	case floating:
		x := p.lastX + img.dx
		y := p.lastY + img.dy
		if x < 0 {
			x = 0
		} else if x+img.width > p.width {
			x = p.width - img.width
		}
		img.x = x
		img.y = y
	}

	if debugEnable {
		log.Printf("arrangeImage: %+v", img)
	}
	p.lastX = img.x + img.width
	p.lastY = img.y
	p.images = append(p.images, img)
}

func (p *page) arrangeWord(s string) {
	w := len(s) * p.charWidth
	x, y := p.arrangeWidth(w, p.charWidth)
	p.lastX = x + w
	p.lastY = y
	if debugEnable {
		log.Printf("arrangeWord: %s %d %d %d", s, x, y, w)
	}
}

func (p *page) arrangeWidth(width int, prefix int) (int, int) {

	x := p.curFragment.arrange(width, prefix)

	for x == -1 {
		p.curFragment = p.curFragment.next
		if p.curFragment == nil {
			p.nextLine()
		}
		x = p.curFragment.arrange(width, p.charWidth)
	}

	return x, p.curLine.y
}

type line struct {
	y      int
	height int
}

type fragment struct {
	next    *fragment
	x       int
	width   int
	curWord word
}

type word struct {
	x     int
	width int
}

func (f *fragment) devide(x1, x2 int) {

	if !(x1 < x2 && f.x <= x1 && x2 <= f.x+f.width) {
		// TODO: more info
		panic("fragment.devide: bad segment")
	}

	next := &fragment{
		next:    f.next,
		x:       x2,
		width:   (f.x + f.width) - x2,
		curWord: word{x: x2},
	}
	if debugEnable {
		log.Printf("fragment.devide: next: %+v", *next)
	}

	f.width = x1 - f.x
	f.next = next
}

func (f *fragment) arrange(width int, prefix int) int {

	x := f.curWord.x + f.curWord.width

	if f.curWord.width > 0 {
		x += prefix
	}

	if x+width > f.x+f.width {
		return -1
	}

	f.curWord = word{
		x:     x,
		width: width,
	}

	return x
}

func solution(pageWidth, lineHeight, charWidth int, buf []byte) ([]image, error) {
	t := unsafeString(buf)

	page := page{
		width:      pageWidth,
		lineHeight: lineHeight,
		charWidth:  charWidth,
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

func run(in io.Reader, out io.Writer) error {

	// Размер входного файла не превышает 1000 байт.
	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	// Первая строка входного файла содержит три целых числа
	var w, h, c int
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
