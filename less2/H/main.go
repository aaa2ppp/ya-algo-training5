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

// возвращает индекс масимального значения и второе максимальное значение
func idx1max2(nums []pair) (int, pair) {
	idx1, max1 := 0, nums[0]
	idx2, max2 := 1, nums[1]

	if max1.less(max2) {
		max1, max2 = max2, max1
		idx1, idx2 = idx2, idx1
	}

	for i, v := range nums[2:] {

		if max1.less(v) {
			idx2, max2 = idx1, max1
			idx1, max1 = i+2, v
			continue
		}

		if max2.less(v) {
			idx2, max2 = i, v
		}
	}

	if debugEnable {
		log.Println("idx1:", idx1, "max1:", max1, "idx2:", idx2, "max2:", max2)
	}

	_ = idx2

	return idx1, max2
}

type pair [2]int

func (p1 pair) less(p2 pair) bool {
	return p1[0] < p2[0] || p1[0] == p2[0] && p1[1] < p2[1]
}

// находит максимальные значения для каждой строки и колонки
func findRowColMax(matrix [][]int, rows, cols []pair) {
	for i := range rows {
		rows[i] = [2]int{-1, -1}
	}
	for i := range cols {
		cols[i] = [2]int{-1, -1}
	}

	for i, row := range matrix {
		for j, v := range row {
			if v > rows[i][0] {
				rows[i][1] = rows[i][0]
				rows[i][0] = v
			} else if v > rows[i][1] {
				rows[i][1] = v
			}
			if v > cols[j][0] {
				cols[j][1] = cols[j][0]
				cols[j][0] = v
			} else if v > cols[j][1] {
				cols[j][1] = v
			}
		}
	}
}

func clearRow(matrix [][]int, idx int) {
	row := matrix[idx]
	for j := range row {
		row[j] = 0
	}
}

func clearCol(matrix [][]int, idx int) {
	for i := range matrix {
		matrix[i][idx] = 0
	}
}

func solution(matrix [][]int) (int, int) {
	n := len(matrix)
	m := len(matrix[0])

	rows := make([]pair, n)
	cols := make([]pair, m)

	findRowColMax(matrix, rows, cols)

	if debugEnable {
		log.Println("rows:", rows, "cols:", cols)
	}

	row, rowMax2 := idx1max2(rows)
	if debugEnable {
		log.Printf("row:%d max2:%d", row, rowMax2)
	}

	col, colMax2 := idx1max2(cols)
	if debugEnable {
		log.Printf("col:%d max2:%d", col, colMax2)
	}

	var resRow, resCol int

	if rowMax2[0] < colMax2[1] {
		resRow = row
		clearRow(matrix, row)
		if debugEnable {
			logMatrix(matrix)
		}
		findRowColMax(matrix, rows, cols)
		resCol, _ = idx1max2(cols)

	} else if colMax2[0] < rowMax2[0] {
		resCol = col
		clearCol(matrix, col)
		if debugEnable {
			logMatrix(matrix)
		}
		findRowColMax(matrix, rows, cols)
		resRow, _ = idx1max2(rows)

	} else {

		if debugEnable {
			log.Println("try remove row first:")
		}
		bak := matrix[row]
		matrix[row] = make([]int, m)
		if debugEnable {
			logMatrix(matrix)
		}
		findRowColMax(matrix, rows, cols)
		col2, colMax2 := idx1max2(cols)
		if debugEnable {
			log.Printf("col2: %d max2: %d", col2, colMax2)
		}

		matrix[row] = bak // restore row

		if debugEnable {
			log.Println("try remove col first:")
		}
		clearCol(matrix, col)
		if debugEnable {
			logMatrix(matrix)
		}
		findRowColMax(matrix, rows, cols)
		row2, rowMax2 := idx1max2(rows)
		if debugEnable {
			log.Printf("row2: %d max2: %d", row2, rowMax2)
		}

		if rowMax2.less(colMax2) {
			resRow, resCol = row2, col
		} else {
			resRow, resCol = row, col2
		}
	}

	return resRow + 1, resCol + 1 // индексация с 1
}

func logMatrix(matrix [][]int) {
	for _, row := range matrix {
		log.Printf("%2d", row)
	}
}

func run(in io.Reader, out io.Writer) error {

	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)

	n, m, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	matrix := makeMatrix(n, m)
	scanMatrix(sc, matrix)

	if debugEnable {
		logMatrix(matrix)
	}

	row, col := solution(matrix)
	fmt.Fprintln(out, row, col)

	return nil
}

func makeMatrix(n, m int) [][]int {
	buf := make([]int, n*m)
	matrix := make([][]int, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}

func scanMatrix(sc *bufio.Scanner, matrix [][]int) error {
	for _, row := range matrix {
		for j := range row {
			v, err := scanInt(sc)
			if err != nil {
				return err
			}
			row[j] = v
		}
	}
	return nil
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(unsafeString(sc.Bytes()))
}

func scanTwoInt(sc *bufio.Scanner) (v1, v2 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	return v1, v2, err
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
