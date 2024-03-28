package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"unsafe"
)

const (
	eps      = 1e-6  // must be <= 1e-4
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

// xxx
type peak struct {
	left  *basin
	right *basin
}

// xxx
var peaks = map[int]*peak{}

type basin struct {
	width          int // ширина водосбора
	top            int // верхний уровень
	bottom         int // нижняя точка
	watershedLevel int
	left           *basin  // левый бассейн
	right          *basin  // правый бассейн
	pnt            []point // рельеф от вернего уровня
	capacity       float64 // максимальный объем выше водораздела
	volume         float64
	isRightSrc     bool
}

func (b *basin) hasWatershed() bool {
	return b.left != nil
}

// возвращает координату x для точки пересечения отрезка с горизонтальной прямой y
func intersection(p1 point, p2 point, y float64) float64 {
	return p1.x + (p2.x-p1.x)/(p2.y-p1.y)*(y-p1.y)
}

// возвращает интеграл для ломанной заданной списком точек
func integral(pnt []point) float64 {
	s := 0.0
	for i := 1; i < len(pnt); i++ {
		s += (pnt[i-1].y + pnt[i].y) * (pnt[i].x - pnt[i-1].x)
	}
	return s / 2
}

// возвращает координаты x для двух крайних точек пересечения ломанной с горизонтальной прямой
// и индексы начала и конца части ломанной лежащей между этими точками
func findBordersByLevel(pnt []point, y float64) (x1, x2 float64, l, r int) {
	l, r = 0, len(pnt)-1

	if y1, y2 := pnt[l].y, pnt[r].y; y-min(y1, y2) > eps {
		panic(fmt.Errorf("findBordersByLevel: level must be less or equal both borders. got %g, want <= %g", y, min(y1, y2)))
	}

	for pnt[l+1].y-y > eps {
		l++
	}

	for pnt[r-1].y-y > eps {
		r--
	}

	x1 = intersection(pnt[l], pnt[l+1], y)
	x2 = intersection(pnt[r], pnt[r-1], y)

	return x1, x2, l + 1, r - 1
}

// xxx
func iFindBordersByLevel(pnt []ipoint, y int) (x1, x2 float64, l, r int) {
	l, r = 0, len(pnt)-1

	if y1, y2 := pnt[l].y, pnt[r].y; y-min(y1, y2) > 0 {
		panic(fmt.Errorf("iFindBordersByLevel: level must be less or equal both borders. got %d, want <= %d", y, min(y1, y2)))
	}
	for pnt[l+1].y-y > 0 {
		l++
	}

	for pnt[r-1].y-y > 0 {
		r--
	}

	x1 = intersection(Point(pnt[l]), Point(pnt[l+1]), float64(y))
	x2 = intersection(Point(pnt[r]), Point(pnt[r-1]), float64(y))

	return x1, x2, l + 1, r - 1
}

func newBasin(pnt []ipoint) *basin {
	if debugEnable {
		log.Println("newBasin:", pnt)
	}

	n := len(pnt)

	// ширина водосбора
	width := pnt[n-1].x - pnt[0].x

	// верхний уровень (по нижнему краю)
	top := min(pnt[0].y, pnt[n-1].y)

	// поиск водораздела и нижней точки
	watershed := -1
	watershedLevel := -infinity
	bottom := infinity
	{
		for i := 1; i < n-1; i++ {
			y := pnt[i].y

			if y > pnt[i-1].y && y > pnt[i+1].y && (watershed == -1 || y > watershedLevel) {
				watershed = i
				watershedLevel = y
			}

			bottom = min(bottom, y)
		}
	}

	if watershed != -1 {
		peaks[watershedLevel] = &peak{}
	}

	// поиск верхней границы (по нижнему краю)
	x1, x2, l, r := iFindBordersByLevel(pnt, top)
	if debugEnable {
		log.Printf("top: %d -> x1: %0.3f x2: %0.3f l: %d r: %d", top, x1, x2, l, r)
	}

	// поиск нижней границы (по уровню водораздела)
	var x3, x4 float64
	if watershed != -1 {
		x3, x4, _, _ = iFindBordersByLevel(pnt, watershedLevel)
		if debugEnable {
			log.Printf("watershed: %d -> x3: %0.3f x4: %0.3f", watershedLevel, x3, x4)
		}
	}

	// собираем рельеф без учета левого и правого бассейна
	pnt2 := make([]point, 0, len(pnt))
	{
		pnt2 = append(pnt2, point{x1, float64(top)})

		for i := l; i <= r; i++ {

			if watershed != -1 && pnt[i].y < watershedLevel {
				pnt2 = append(pnt2, point{x3, float64(watershedLevel)})

				for pnt[i].y < watershedLevel || i == watershed {
					i++
				}

				pnt2 = append(pnt2, point{x4, float64(watershedLevel)})
			}

			pnt2 = append(pnt2, Point(pnt[i]))
		}

		pnt2 = append(pnt2, point{x2, float64(top)})
	}

	// считаем объем (без учета левого и правого бассейна)
	capacity := float64(top)*(x2-x1) - integral(pnt2)

	b := &basin{
		width:          width,
		top:            top,
		bottom:         bottom,
		watershedLevel: watershedLevel,
		pnt:            pnt2,
		capacity:       capacity,
	}

	if y1, y2 := pnt[0].y, pnt[len(pnt)-1].y; y1 < y2 {
		b.isRightSrc = true
		peaks[y2].left = b
	} else {
		peaks[y1].right = b
	}

	// создаем левый и правый бассейны
	if watershed != -1 {

		b.left = newBasin(pnt[:watershed+1])
		b.right = newBasin(pnt[watershed:])

		if debugEnable {
			log.Printf("left : %+v", b.left)
			log.Printf("right: %+v", b.right)
		}
	}

	return b
}

func (b *basin) distribute(v float64) float64 {
	if debugEnable {
		log.Printf("== filling == b: %+v v: %0.3f", b, v)
	}

	if b.hasWatershed() {
		// делим объем между правым и левым бассейном пропорцианально
		// ширине водозабора каждого
		lv := float64(b.left.width) / float64(b.width) * v
		rv := float64(b.right.width) / float64(b.width) * v
		v = 0

		v += b.left.distribute(lv)
		v += b.right.distribute(rv)
	}

	return b.pour(nil, v)
}

func (b *basin) pour(src *basin, v float64) float64 {

	if b == nil || b == src {
		return v
	}

	available := b.capacity - b.volume
	if available < eps {
		return v
	}

	dv := min(v, available)
	b.volume += dv
	v -= dv

	return b.pourOut(v)
}

func (b *basin) pourOut(v float64) float64 {
	if v < eps {
		return 0
	}
	p := peaks[b.top]
	if b.isRightSrc {
		v = p.left.pour(b, v)
	} else {
		v = p.right.pour(b, v)
	}
	return v
}

func (b *basin) findLevel(v float64) (level float64) {
	top := float64(b.top)

	// xxx
	bottom := float64(b.bottom)
	if b.hasWatershed() {
		bottom = float64(b.watershedLevel)
	}

	var (
		hasPrev      bool
		prevL, prevR int
		prevIntegral float64
	)

	for top-bottom > eps {
		mid := float64(top+bottom) / 2
		x1, x2, l, r := findBordersByLevel(b.pnt, mid)

		if debugEnable {
			log.Printf("mid: %0.3f x1: %0.3f x2: %0.3f l: %d r: %d", mid, x1, x2, l, r)
		}

		if !hasPrev || prevL != l || prevR != r {
			prevIntegral = integral(b.pnt[l : r+1])
			prevL, prevR, hasPrev = l, r, true
		}

		vol := mid*(b.pnt[r].x-b.pnt[l].x) - prevIntegral
		vol += ((mid-b.pnt[l].y)*(b.pnt[l].x-x1) + (mid-b.pnt[r].y)*(x2-b.pnt[r].x)) / 2

		if vol > v {
			top = mid
		} else {
			bottom = mid
		}
	}

	if debugEnable {
		log.Printf("== findLevel == b: %+v v: %0.3f -> %0.3f", b, v, (top+bottom)/2)
	}

	return (top + bottom) / 2
}

func (b *basin) findHeigth() float64 {
	if b == nil {
		return 0
	}

	if b.volume < eps {
		return max(b.left.findHeigth(), b.right.findHeigth())
	}

	return b.findLevel(b.volume) - float64(b.bottom)
}

func solve(h float64, pnt []ipoint) float64 {

	// xxx фиктивные пики
	peaks[pnt[0].y] = &peak{}
	peaks[pnt[len(pnt)-1].y] = &peak{}

	root := newBasin(pnt)

	// if !checkWidth(root) {
	// 	log.Fatal("check width fail")
	// }
	// if !checkBottom(root) {
	// 	log.Fatal("check bottom fail")
	// }
	// if !checkTop(root) {
	// 	log.Fatal("check top fail")
	// }

	if debugEnable {
		log.Printf("root basin: %+v", root)
	}

	w := float64(root.width)
	v := root.distribute(h * w)

	if math.Abs(caclTotalVolume(root) + v - h*w) > eps {
		panic("bad total volume")
	}

	checkVolume(root)

	if v > 0 {
		return float64(root.top-root.bottom) + v/w
	}

	return root.findHeigth()
}

func caclTotalVolume(b *basin) float64 {
	if b == nil {
		return 0
	}
	return b.volume + caclTotalVolume(b.left) + caclTotalVolume(b.right)
}

func checkVolume(b *basin) (full bool) {
	if b == nil {
		return true
	}
	if b.volume > eps && (!checkVolume(b.left) || !checkVolume(b.right)) {
		panic("volume != 0 whit child not full")
	}
	return b.capacity - b.volume < eps
}

// func checkWidth(b *basin) bool {

// 	if !root.hasWatershed() {
// 		return true
// 	}

// 	if (root.pnt[len(root.pnt)-1].x-root.pnt[0].x)-root.width > eps {
// 		return false
// 	}

// 	w := root.width
// 	w -= root.left.width
// 	w -= root.right.width

// 	if w > eps {
// 		log.Printf("chechWidth: %+v -> %v", root, w)
// 		return false
// 	}

// 	return checkWidth(root.left) && checkWidth(root.right)
// }

// func checkBottom(b *basin) bool {

// 	if !b.hasWatershed() {
// 		return true
// 	}

// 	if b.bottom-min(b.left.bottom, b.right.bottom) > eps {
// 		return false
// 	}

// 	return checkBottom(b.left) && checkBottom(b.right)
// }

// func checkTop(b *basin) bool {

// 	if b == nil {
// 		return true
// 	}

// 	if b.top-b.pnt[0].y > eps || b.top-b.pnt[len(b.pnt)-1].y > eps {
// 		return false
// 	}

// 	return checkTop(b.left) && checkTop(b.right)
// }

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
	pnt := make([]ipoint, n+2)

	maxY := -infinity
	for i := 1; i <= n; i++ {
		x, y, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		maxY = max(maxY, y)
		pnt[i] = ipoint{x, y}
	}

	pnt[0] = ipoint{pnt[1].x, maxY + 1}
	pnt[n+1] = ipoint{pnt[n].x, maxY + 2}
	// pnt[0] = ipoint{pnt[1].x, infinity}
	// pnt[n+1] = ipoint{pnt[n].x, infinity}

	res := solve(h, pnt)
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
