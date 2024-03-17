package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	Name  = 1
	Score = 12
)

func readFile(fileName string, scores map[string]int) error {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(buf))
	names := make(map[string]struct{}, 1000)

	for i := 1; ; i++ {
		row, err := r.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("row%d: %w", i, err)
		}

		if i == 1 && row[Score] == "Очки" {
			// это заголовок
			continue
		}

		name := uniqueName(names, row[Name])

		score, err := strconv.Atoi(row[Score])
		if err != nil {
			return fmt.Errorf("row%d: %w", i, err)
		}

		scores[name] += score
	}
}

func uniqueName(names map[string]struct{}, name string) string {
	i := 1
	unique := name
	for _, ok := names[unique]; ok; _, ok = names[unique] {
		i++
		unique = name + strconv.Itoa(i)
	}
	names[unique] = struct{}{}
	return unique
}

func addScores(stats map[int]int, scores map[string]int) {
	for _, v := range scores {
		stats[v]++
	}
}

func printStats(stats map[int]int) {
	for k := 20; k > 0; k-- {
		if v, ok := stats[k]; ok {
			fmt.Printf("%2d %4d\n", k, v)
		}
	}
}

func main() {
	scores := make(map[string]int, 1000)
	stats := make(map[int]int, 20)

	for i := 1; i <= 4; i++ {
		fileName := fmt.Sprintf("less%d.csv", i)

		if err := readFile(fileName, scores); err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("%s: %v", fileName, err)
			}
			log.Printf("%s: not found - skip", fileName)
			continue
		}

		fmt.Println(fileName)
		addScores(stats, scores)
		printStats(stats)
	}
}
