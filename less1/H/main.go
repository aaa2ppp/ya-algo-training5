package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

type F struct {
	t0 float64
	l0 float64
	v  float64
}

func (f F) String() string {
	return fmt.Sprintf("(%v,%v) %v", f.t0, f.l0, f.v)
}

func (f F) a() float64 {
	return -f.v*f.t0 + f.l0
}

func eq(v1, v2 float64) bool {
	return math.Abs(v1-v2)/math.Max(1, math.Abs(v2)) <= 1e-9
}

func lq(v1, v2 float64) bool {
	return v1 < v2 || eq(v1, v2)
}

func intersection(f1, f2 F) (float64, float64, bool) {
	if eq(f1.v, f2.v) {
		return 0, 0, false
	}

	a1 := f1.a()
	a2 := f2.a()

	t := (a2 - a1) / (f1.v - f2.v)
	l := f1.v*t + a1

	return t, l, true
}

func solution(l, x1, v1, x2, v2 int) (float64, bool) {

	top := float64(l) / 2

	init := func(x, v int) F {
		if float64(x) < top {
			return F{
				l0: float64(x),
				v:  float64(v),
			}
		} else {
			return F{
				l0: float64(l - x),
				v:  float64(-v),
			}
		}
	}

	next := func(f F) F {
		if !lq(0, f.v) { // f.v < 0
			t, _, ok := intersection(f, F{l0: 0})
			if !ok {
				panic(fmt.Sprintf("can't next for %v", f))
			}
			return F{t0: t, l0: 0, v: -f.v}
		} else if !lq(f.v, 0) { // f.v > 0
			t, _, ok := intersection(f, F{l0: top})
			if !ok {
				panic(fmt.Sprintf("can't next for %v", f))
			}
			return F{t0: t, l0: top, v: -f.v}
		}
		return f // f.v == 0
	}

	f1 := init(x1, v1)
	f2 := init(x2, v2)
	if debugEnable {
		log.Printf("f1: %v", f1)
		log.Printf("f2: %v", f2)
	}

	for i := 0; i < 3; i++ {

		next1 := next(f1)
		next2 := next(f2)
		if debugEnable {
			log.Printf("%d: next1: %v", i, next1)
			log.Printf("%d: next2: %v", i, next2)
		}

		t0 := math.Max(f1.t0, f2.t0)

		var t1 float64
		if eq(next1.v, 0) {
			t1 = next2.t0
		} else if eq(next2.v, 0) {
			t1 = next1.t0
		} else {
			t1 = math.Min(next1.t0, next2.t0)
		}

		t, l, ok := intersection(f1, f2)
		if debugEnable {
			log.Printf("%d: interval: [%g, %g]", i, t0, t1)
			log.Printf("%d: intersection: (%g,%g) %v", i, t, l, ok)
		}

		if ok && lq(t0, t) && lq(t, t1) && lq(0, l) && lq(l, top) {
			return t, true
		}

		if eq(next1.t0, t1) {
			f1 = next1
		}

		if eq(next2.t0, t1) {
			f2 = next2
		}
	}

	return 0, false
}

func run(in io.Reader, out io.Writer) error {

	var l, x1, v1, x2, v2 int
	if _, err := fmt.Fscan(in, &l, &x1, &v1, &x2, &v2); err != nil {
		return err
	}

	if x1 == x2 || x1 == l-x2 {
		fmt.Fprintln(out, "YES\n0")
		return nil
	}

	if v1 == 0 && v2 == 0 {
		fmt.Fprintln(out, "NO")
		return nil
	}

	t, ok := solution(l, x1, v1, x2, v2)
	if !ok {
		panic("oops!.. can't solve")
	}

	fmt.Fprintf(out, "YES\n%0.10f\n", t)

	return nil
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
