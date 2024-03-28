package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
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

type basin struct {
	width          float64 // ширина водосбора
	top            float64 // верхний уровень
	bottom         float64 // нижняя точка
	watershedLevel float64
	left           *basin  // левый бассейн
	right          *basin  // правый бассейн
	pnt            []point // рельеф от вернего уровня
	volume         float64 // максимальный объем
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
func findBordersByLevel(points []point, y float64) (x1, x2 float64, l, r int) {
	l, r = 0, len(points)-1

	if y1, y2 := points[l].y, points[r].y; y > min(y1, y2) {
		panic(fmt.Errorf("findBordersByLevel: level must be less or equal both borders. got %g, want <= %g", y, min(y1, y2)))
	}

	// defer func() {
	// 	if p := recover(); p != nil {
	// 		log.Printf("panic recovere: %v", p)
	// 		log.Printf("y: %0.3f points: %v", y, points)
	// 		panic(p)
	// 	}
	// }()

	for points[l+1].y > y {
		l++
	}

	for points[r-1].y > y {
		r--
	}

	x1 = intersection(points[l], points[l+1], y)
	x2 = intersection(points[r], points[r-1], y)

	return x1, x2, l + 1, r - 1
}

func newBasin(pnt []point) *basin {
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
	watershedLevel := float64(-infinity)
	bottom := float64(infinity)
	{
		for i := 1; i < n-1; i++ {
			y := pnt[i].y

			if y > pnt[i-1].y && y > pnt[i+1].y && (watershed == -1 || y > pnt[watershed].y) {
				watershed = i
				watershedLevel = y
			}

			bottom = min(bottom, y)
		}
	}

	if watershed == -1 {
		watershedLevel = bottom
	}

	// поиск верхней границы (по нижнему краю)
	x1, x2, l, r := findBordersByLevel(pnt, top)
	if debugEnable {
		log.Printf("top: %0.3f -> x1: %0.3f x2: %0.3f l: %d r: %d", top, x1, x2, l, r)
	}

	// поиск нижней границы (по уровню водораздела)
	var x3, x4 float64
	if watershed != -1 {
		x3, x4, _, _ = findBordersByLevel(pnt, watershedLevel)
		if debugEnable {
			log.Printf("watershed: %0.3f -> x3: %0.3f x4: %0.3f", watershedLevel, x3, x4)
		}
	}

	// собираем рельеф без учета левого и правого бассейна
	pnt2 := make([]point, 0, len(pnt))
	{
		pnt2 = append(pnt2, point{x1, top})

		for i := l; i <= r; i++ {

			if watershed != -1 && pnt[i].y < watershedLevel {
				pnt2 = append(pnt2, point{x3, watershedLevel})

				for pnt[i].y < watershedLevel || i == watershed {
					i++
				}

				pnt2 = append(pnt2, point{x4, watershedLevel})
			}

			pnt2 = append(pnt2, pnt[i])
		}

		pnt2 = append(pnt2, point{x2, top})
	}

	// считаем объем (без учета левого и правого бассейна)
	volume := top*(x2-x1) - integral(pnt2)

	// создаем левый и правый бассейны
	var left, right *basin
	if watershed != -1 {
		left = newBasin(pnt[:watershed+1])
		right = newBasin(pnt[watershed:])
		volume += left.volume + right.volume

		if debugEnable {
			log.Printf("left : %+v", left)
			log.Printf("right: %+v", right)
		}
	}

	return &basin{
		width:          width,
		top:            top,
		bottom:         bottom,
		watershedLevel: watershedLevel,
		left:           left,
		right:          right,
		pnt:            pnt2,
		volume:         volume,
	}
}

func (b *basin) filling(v float64) (depth float64) {
	if debugEnable {
		log.Printf("== filling == b: %+v v: %0.3f", b, v)
	}

	if v > b.volume {
		panic(fmt.Errorf("gorge.filling: volume too large. got %g, want <= %g", v, b.volume))
	}

	if v < eps {
		// не стоит внимания
		return 0
	}

	var lv, rv float64
	if b.hasWatershed() {

		// делим объем между правым и левым бассейном пропорцианально
		// ширине водозабора каждого
		lv = b.left.width / b.width * v
		rv = b.right.width / b.width * v
		v = 0

		// преливаем избытки

		if lv > b.left.volume {
			rv += lv - b.left.volume
			lv = b.left.volume
		}

		if rv > b.right.volume {
			lv += rv - b.right.volume
			rv = b.right.volume
		}

		if lv > b.left.volume {
			v = lv - b.left.volume
			lv = b.left.volume
		}
	}

	if v > 0 {
		// бассейнов нет или оба заполнены до краев
		return b.findLevel(v) - b.bottom
	}

	var ld float64
	if lv < b.left.volume {
		ld = b.left.filling(lv)
	} else {
		ld = b.left.top - b.left.bottom
	}

	var rd float64
	if rv < b.right.volume {
		rd = b.right.filling(rv)
	} else {
		rd = b.right.top - b.right.bottom
	}

	return max(ld, rd)
}

func (b *basin) findLevel(v float64) (level float64) {
	top := b.top
	bottom := b.watershedLevel

	var (
		hasPrev      bool
		prevL, prevR int
		prevIntegral float64
	)

	for top-bottom > eps {
		mid := (top + bottom) / 2
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

func solve(h float64, pnt []point) float64 {
	b := newBasin(pnt)
	if debugEnable {
		log.Printf("root basin: %+v", b)
	}
	v := h * b.width
	if v > b.volume {
		v -= b.volume
		return b.top-b.bottom + v / b.width
	}
	return b.filling(v)
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
	pnt := make([]point, n+2)

	maxY := -infinity
	for i := 1; i <= n; i++ {
		x, y, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		maxY = max(maxY, y)
		pnt[i] = point{float64(x), float64(y)}
	}

	pnt[0] = point{pnt[1].x, float64(maxY + 1)}
	pnt[n+1] = point{pnt[n].x, float64(maxY + 1)}

	res := solve(h, pnt)
	fmt.Fprintf(bw, "%0.6f", res)

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
