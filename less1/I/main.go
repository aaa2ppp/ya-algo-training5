package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func run(in io.Reader, out io.Writer) error {

	sc := bufio.NewScanner(in)

	sc.Scan()
	n, err := strconv.Atoi(sc.Text())
	if err != nil {
		return err
	}

	sc.Scan()
	year, err := strconv.Atoi(sc.Text())
	if err != nil {
		return err
	}

	holidays := make([]bool, 367) // TODO: bitArray
	for i := 0; i < n; i++ {
		sc.Scan()
		d, err := time.Parse("2 January 2006", fmt.Sprintf("%s %d", sc.Text(), year))
		if err != nil {
			return err
		}

		if debugEnable {
			log.Printf("%d %s -> %d", year, sc.Text(), d.YearDay())
		}

		holidays[d.YearDay()] = true
	}

	sc.Scan()
	firstWeekday := sc.Text()

	weekDays := make([]int, 7)

	d := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	if firstWeekday != d.Weekday().String() {
		panic(fmt.Sprintf("oops!.. нас обманули: got %s, want %v", firstWeekday, d.Weekday()))
	}

	for d.Year() == year {
		if !holidays[d.YearDay()] {
			weekDays[d.Weekday()]++
		}
		d = d.AddDate(0, 0, 1)
	}

	if debugEnable {
		log.Printf("weekDays: %v", weekDays)
	}

	// NOTE: В России неделя начинается в понедельник, а не в воскресенье
	minWeekday, maxWeekday := 1, 1 

	for _, i := range []int{2,3,4,5,6,0} {
		if weekDays[i] > weekDays[maxWeekday] {
			maxWeekday = i
		}
		if weekDays[i] < weekDays[minWeekday] {
			minWeekday = i
		}
	}

	fmt.Fprintln(out, time.Weekday(maxWeekday), time.Weekday(minWeekday))

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
