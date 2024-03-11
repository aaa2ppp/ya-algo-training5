package main

// version without time package

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var daysInMonth = []struct {
	name string
	days int
}{
	{"January", 31},
	{"February", 28}, // or 29 if it is a leap year
	{"March", 31},
	{"April", 30},
	{"May", 31},
	{"June", 30},
	{"July", 31},
	{"August", 31},
	{"September", 30},
	{"October", 31},
	{"November", 30},
	{"December", 31},
}

var weekdays = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func isLeapYear(year int) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}

func daysInYear(year int) int {
	if isLeapYear(year) {
		return 366
	}
	return 365
}

func yearDay(year int, day int, monthName string) (int, error) {
	n := 0

	for i := range daysInMonth {
		if strings.EqualFold(monthName, daysInMonth[i].name) { // case insensitive
			if i >= 2 && isLeapYear(year) {
				n++
			}
			return n + day, nil
		}
		n += daysInMonth[i].days
	}

	return 0, errors.New(monthName + ": unknown month")
}

func parseWeekday(weekdayName string) (int, error) {
	i := 6
	for !strings.EqualFold(weekdayName, weekdays[i]) { // case insensitive
		i--
	}

	if i == -1 {
		return -1, errors.New(weekdayName + ": unknown weekday")
	}

	return i, nil
}

func run(in io.Reader, out io.Writer) error {

	sc := bufio.NewScanner(in)

	sc.Scan()
	holidayCount, err := strconv.Atoi(sc.Text())
	if err != nil {
		return err
	}

	sc.Scan()
	year, _ := strconv.Atoi(sc.Text())
	if err != nil {
		return err
	}

	holidays := make([]bool, 367) // TODO: bitArray
	for i := 0; i < holidayCount; i++ {
		sc.Scan()
		p := strings.Split(sc.Text(), " ")

		day, err := strconv.Atoi(p[0])
		if err != nil {
			return err
		}

		monthName := p[1]
		d, err := yearDay(year, day, monthName)
		if err != nil {
			return err
		}

		if debugEnable {
			log.Printf("%d %s -> %d", year, sc.Text(), d)
		}
		holidays[d] = true
	}

	sc.Scan()
	weekday, err := parseWeekday(sc.Text())
	if err != nil {
		return err
	}

	workings := make([]int, 7)

	for i, n := 1, daysInYear(year); i <= n; i++ {
		if !holidays[i] {
			workings[weekday]++
		}
		weekday++
		if weekday == 7 {
			weekday = 0
		}
	}

	if debugEnable {
		log.Printf("weekDays: %v", workings)
	}

	// NOTE: В России неделя начинается в понедельник, а не в воскресенье
	minWeekday, maxWeekday := 1, 1

	for _, i := range []int{2, 3, 4, 5, 6, 0} {
		if workings[i] > workings[maxWeekday] {
			maxWeekday = i
		}
		if workings[i] < workings[minWeekday] {
			minWeekday = i
		}
	}

	fmt.Fprintln(out, weekdays[maxWeekday], weekdays[minWeekday])

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
