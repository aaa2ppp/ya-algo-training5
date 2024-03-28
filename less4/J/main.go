package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

const (
	eps      = 1e-5  // must be <= 1e-4
	infinity = 10001 // must be > max(abs(y))
)

type point struct {
	x, y float64
}

func Point(p ipoint) point {
	return point{float64(p.x), float64(p.y)}
}

type ipoint struct {
	x, y int
}

// возвращает интеграл для ломанной заданной списком точек
func integral(first, last *list.Element) float64 {
	s := 0.0

	p1 := first.Value.(point)
	for first != last {
		next := first.Next()
		p2 := next.Value.(point)

		s += (p1.y + p2.y) * (p2.x - p1.x)

		first = next
		p1 = p2
	}

	return s / 2
}

// возвращает координату x для точки пересечения отрезка с горизонтальной прямой y
func intersection(p1 point, p2 point, y float64) float64 {
	return p1.x + (p2.x-p1.x)/(p2.y-p1.y)*(y-p1.y)
}

// находит и создает следующую точку пересечения ломаной с горизонталью
func nextIntersec(points *list.List, el *list.Element, y float64) *list.Element {

	p1 := el.Value.(point)
	if p1.y < y {
		panic("nextIntersec: p1.y < y")
	}

	for el.Next() != nil {
		p2 := el.Next().Value.(point)
		// if p1.y >= y && y > p2.y {
		if y > p2.y {
			x := intersection(p1, p2, y)
			return points.InsertAfter(point{x, y}, el)
		}
		p1 = p2
		el = el.Next()
	}

	panic("nextIntersec: not found")
}

// находит и создает предыдущую точку пересечения ломаной с горизонталью
func prevIntersec(points *list.List, el *list.Element, y float64) *list.Element {

	p1 := el.Value.(point)
	if p1.y < y {
		panic("prevIntersec: p1.y < y")
	}

	for el.Prev() != nil {
		p2 := el.Prev().Value.(point)
		// if p1.y >= y && y > p2.y {
		if y > p2.y {
			x := intersection(p1, p2, y)
			return points.InsertBefore(point{x, y}, el)
		}
		p1 = p2
		el = el.Prev()
	}

	panic("prevIntersec: not found")
}

type basin struct {
	points     *list.List
	prev       *list.Element
	post       *list.Element
	first      *list.Element
	last       *list.Element
	needRemove *list.Element
	width      float64
	top        float64
	bottom     float64
	capacity   float64
	volume     float64
	isRight    bool
}

func newBasin(points *list.List, prev, post *list.Element) *basin {

	prevP := prev.Value.(point)
	postP := post.Value.(point)

	width := postP.x - prevP.x

	cur := prev.Next()
	bottom := float64(infinity)
	for cur != post {
		bottom = min(bottom, cur.Value.(point).y)
		cur = cur.Next()
	}

	return &basin{
		points:  points,
		prev:    prev,
		post:    post,
		width:   width,
		top:     min(prevP.y, postP.y),
		bottom:  bottom,
		isRight: prevP.y < postP.y,
	}
}

// xxx
func (b *basin) makeTopPoints() {
	p1 := b.prev.Value.(point)
	p2 := b.post.Value.(point)
	if p1.y < p2.y {
		b.top = p1.y
		b.isRight = true
		b.first = b.prev
		b.last = prevIntersec(b.points, b.post, b.top)
		b.needRemove = b.last // xxx
	} else {
		b.top = p2.y
		b.isRight = false
		b.first = nextIntersec(b.points, b.prev, b.top)
		b.needRemove = b.first // xxx
		b.last = b.post
	}
	b.updateCapacity()
}

func (b *basin) updateCapacity() {
	p1 := b.first.Value.(point)
	p2 := b.last.Value.(point)
	b.capacity = b.top*(p2.x-p1.x) - integral(b.first, b.last)
}

// возвращает высоту осадков, которая должна выпасть, чтобы заполнить бассейн до краев
func (b *basin) precipHiToFill() float64 {
	return (b.capacity - b.volume) / b.width
}

func (b *basin) removeRelief() {
	pnt := b.first.Next()
	for pnt != b.last {
		next := pnt.Next()
		b.points.Remove(pnt)
		pnt = next
	}
}

func (b *basin) level() float64 {
	top := b.top
	bottom := b.bottom

	for top-bottom > eps {
		mid := (top + bottom) / 2

		el1 := nextIntersec(b.points, b.first, mid)
		el2 := prevIntersec(b.points, b.last, mid)

		p1 := el1.Value.(point)
		p2 := el2.Value.(point)
		s := mid*(p2.x-p1.x) - integral(el1, el2)

		b.points.Remove(el1)
		b.points.Remove(el2)

		if b.volume < s {
			top = mid
		} else {
			bottom = mid
		}
	}

	return (top + bottom) / 2
}

func newPointList(pnt []ipoint) *list.List {
	maxY := -infinity

	for _, p := range pnt {
		maxY = max(maxY, p.y)
	}

	l := list.New()
	l.PushBack(Point(ipoint{pnt[0].x, maxY + 1}))

	for _, p := range pnt {
		l.PushBack(Point(p))
	}

	l.PushBack(Point(ipoint{pnt[len(pnt)-1].x, maxY + 1}))

	return l
}

type basinList struct {
	*list.List
	points *list.List
}

func newBasinList(points *list.List) basinList {
	l := list.New()

	{
		prev := points.Front()
		cur := points.Front().Next()

		p0 := prev.Value.(point)
		p1 := cur.Value.(point)

		for cur.Next() != nil {
			p2 := cur.Next().Value.(point)
			if p1.y > p0.y && p1.y > p2.y {
				l.PushBack(newBasin(points, prev, cur))
				prev = cur
			}
			cur = cur.Next()
			p0, p1 = p1, p2
		}

		l.PushBack(newBasin(points, prev, points.Back()))
	}

	for cur := l.Front(); cur != nil; cur = cur.Next() {
		b := cur.Value.(*basin)
		b.makeTopPoints()

	}

	return basinList{
		List:   l,
		points: points,
	}
}

// возвращает минимальную высоту осадков для заполнения покрайней мере одного бассейна
func (l basinList) findMinPrecipHiToFill() float64 {
	minimum := float64(infinity)

	cur := l.Front()
	for cur != nil {
		h := cur.Value.(*basin).precipHiToFill()
		minimum = min(minimum, h)
		cur = cur.Next()
	}

	return minimum
}

func (l basinList) pour(h float64) {
	for cur := l.Front(); cur != nil; cur = cur.Next() {
		b := cur.Value.(*basin)
		b.volume += b.width * h
	}
}

func (l basinList) removeFullBasins() {
	cur := l.Front()
	for cur != nil {
		next := cur.Next()
		b := cur.Value.(*basin)
		if b.capacity-b.volume < eps {
			l._remove(cur)
		}
		cur = next
	}
}

func (l basinList) _remove(el *list.Element) {

	// объединяем бассейн с соседним
	b := el.Value.(*basin)
	var b2 *basin
	if b.isRight {
		b2 = el.Prev().Value.(*basin)
		b2.post = b.post
	} else {
		b2 = el.Next().Value.(*basin)
		b2.prev = b.prev
	}
	b2.width += b.width
	b2.bottom = min(b2.bottom, b.bottom)
	l.points.Remove(b2.needRemove)

	b.removeRelief()
	b2.makeTopPoints()
	b2.updateCapacity()
	l.List.Remove(el)
}

func (l basinList) findMaxHeight() float64 {
	maximum := float64(-infinity)
	cur := l.List.Front()
	for cur != nil {
		b := cur.Value.(*basin)
		maximum = max(maximum, b.level()-b.bottom)
		cur = cur.Next()
	}
	return maximum
}

// O(N^2)
func solve(h float64, ipoints []ipoint) float64 {
	points := newPointList(ipoints)
	basins := newBasinList(points)

	// чтобы не заморачиватся с переливами, в цикле заполняем басейны начиная
	// от самого мелкого (по потреблению). заполненые удаляем из системы O(N)
	for {

		// ищем минимальную высоту осадков для наполнения по крайней мере одного бассейна
		// O(N), можно сделать O(logN). но зачем, при N <= 100?
		minH := basins.findMinPrecipHiToFill()

		if minH > h {
			// наполняем бассейны и расчитываем уровень неполного заполнения
			// O(N)+O(N*log(Y')) Y'=max(Y*2)*<точность>, где <точность> кол-во знаков после запятой
			basins.pour(h)
			return basins.findMaxHeight() 
		}

		// наполняем бассейны
		// O(N) и ничего не сделать :(, наполнить нужно каждый
		h -= minH
		basins.pour(minH) 
		if basins.Len() == 1 {
			break
		}

		basins.removeFullBasins() // O(N), можно сделать O(1)
	}

	b := basins.Front().Value.(*basin)
	return b.top - b.bottom + h
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	h, err := scanFloat(sc)
	if err != nil {
		return err
	}

	n++ // "В последующих N + 1 строках..."  что за подстава!!!
	points := make([]ipoint, n)

	for i := 0; i < n; i++ {
		x, y, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		points[i] = ipoint{x, y}
	}

	res := solve(h, points)
	fmt.Fprintf(bw, "%0.5f", res)

	return nil
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(unsafeString(sc.Bytes()))
}

func scanFloat(sc *bufio.Scanner) (float64, error) {
	sc.Scan()
	return strconv.ParseFloat(unsafeString(sc.Bytes()), 64)
}

func scanTwoInt(sc *bufio.Scanner) (v1, v2 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	return v1, v2, err
}

func scanThreeInt(sc *bufio.Scanner) (v1, v2, v3 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	return v1, v2, v3, err
}

func scanFourInt(sc *bufio.Scanner) (v1, v2, v3, v4 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	if err == nil {
		v4, err = scanInt(sc)
	}
	return v1, v2, v3, v4, err
}

func scanInts(sc *bufio.Scanner, a []int) error {
	for i := range a {
		v, err := scanInt(sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type writeOpts struct {
	sep byte
	end byte
}

func writeInt[I Int](bw *bufio.Writer, v I, opts writeOpts) error {
	var buf [32]byte

	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))

	if err == nil && opts.end != 0 {
		bw.WriteByte(opts.end)
	}

	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, opts writeOpts) error {
	var err error

	if len(a) != 0 {
		var buf [32]byte

		if opts.sep == 0 {
			opts.sep = ' '
		}

		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))

		for i := 1; err == nil && i < len(a); i++ {
			err = bw.WriteByte(opts.sep)
			if err == nil {
				_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
			}
		}
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
