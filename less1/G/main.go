package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func check(i, x, z int) int {
	for ; ; i++ {
		if debugEnable {
			log.Printf("%d.1 %3d --- %3d", i, x, z)
		}
		if z <= 0 {
			return i
		}

		x -= min(x, z)
		if debugEnable {
			log.Printf("%d.2 %3d --- %3d", i, x, z)
			log.Println("---")
		}
		if x <= 0 {
			return -1
		}

		z -= min(x, z)
	}
}

func solution(x, y, z, p int) int {

	if debugEnable {
		log.Printf("=== %3d %3d %3d", x, y, z)
		log.Println("---")
	}

	res := -1

	// XXX У меня нет доказательств. Это эмпирическое решение теста 9, 13 и 20
	for i := 1; res == -1 || i < res; i++ {
		if x >= y {
			i2 := check(i, x, z-(x-y))
			if i2 == i {
				if debugEnable {
					log.Printf("bingo! %d", i2)
				}
				return i2
			}
			if i2 != -1 && (res == -1 || i2 < res) {
				if debugEnable {
					log.Printf("found: %d", i2)
				}
				res = i2
			}
		}

		{
			x2 := x

			// гарантируем уменьшение y
			if y > 0 {
				y--
				x2--
			}

			// максимально уменьшаем z
			d := min(x2, z)
			z -= d
			x2 -= d

			y -= min(x2, y)
		}
		if debugEnable {
			log.Printf("%d.1 %3d %3d %3d", i, x, y, z)
		}
		if y <= 0 && z <= 0 {
			if debugEnable {
				log.Printf("bingo! %d", i)
			}
			return i
		}

		x -= min(x, z)
		if debugEnable {
			log.Printf("%d.2 %3d %3d %3d", i, x, y, z)
		}
		if x <= 0 {
			break
		}

		if y > 0 {
			z += p
		}
		if debugEnable {
			log.Printf("%d.3 %3d %3d %3d", i, x, y, z)
			log.Println("---")
		}
	}

	if debugEnable {
		log.Printf("result %d", res)
	}
	return res
}

func run(in io.Reader, out io.Writer) error {

	var x, y, p int
	if _, err := fmt.Fscan(in, &x, &y, &p); err != nil {
		return err
	}

	res := solution(x, y, 0, p)

	fmt.Fprintln(out, res)

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
