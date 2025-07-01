package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func solve(
	regionMap [][]int,
	row int,
	usedColumns map[int]int,
	regions map[int]bool,
) bool {
	if row == len(regionMap) {
		return true
	}

	for col := range len(regionMap) {
		_, isCurrentColumnUsed := usedColumns[col]

		colPrev, ok := usedColumns[col-1]
		hasDiagonalConflictLeft := ok && (colPrev == row-1 || colPrev == row+1)

		colNext, ok := usedColumns[col+1]
		hasDiagonalConflictRight := ok && (colNext == row-1 || colNext == row+1)

		region := regionMap[row][col]
		isRegionUsed := regions[region]

		if isCurrentColumnUsed || hasDiagonalConflictLeft || hasDiagonalConflictRight || isRegionUsed {
			continue
		}

		usedColumns[col] = row
		regions[region] = true

		if solve(regionMap, row+1, usedColumns, regions) {
			return true
		}

		// Backtrack
		delete(usedColumns, col)
		delete(regions, region)
	}

	return false
}

func printSolution(regionMap [][]int, usedColumns map[int]int) {
	n := len(regionMap)
	table := make([][]string, n)

	for i := range n {
		table[i] = make([]string, n)
		for j := range n {
			row, ok := usedColumns[j]
			if ok && row == i {
				table[i][j] = "👑"
			} else {
				table[i][j] = "✘"
			}
		}
	}

	writer := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewBlueprint(
		tw.Rendition{
			Settings: tw.Settings{
				Separators: tw.Separators{
					BetweenRows: tw.On,
				},
			},
		},
	)))

	writer.Bulk(table)
	writer.Render()
}

func readRegionMapFromCSV(path string) ([][]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	regionMap := make([][]int, len(records))
	for i, record := range records {
		regionMap[i] = make([]int, len(record))
		for j, value := range record {
			num, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("error parsing value %s at row %d, column %d: %v", value, i+1, j+1, err)
			}
			regionMap[i][j] = num
		}
	}

	return regionMap, nil
}

func main() {
	source := flag.String("source", "", "Path to the source file containing the region map")
	flag.Parse()

	if source == nil || *source == "" {
		fmt.Println("Please provide a source file using the -source flag.")
		return
	}

	regions := make(map[int]bool)
	usedColumns := make(map[int]int)
	regionMap, err := readRegionMapFromCSV(*source)
	if err != nil {
		panic(fmt.Sprintf("Error reading region map: %v", err))
	}

	start := time.Now()
	found := solve(regionMap, 0, usedColumns, regions)
	elpased := time.Since(start)

	if !found {
		fmt.Println("No valid solution exists for this board.")
		return
	}

	printSolution(regionMap, usedColumns)
	fmt.Printf("Solution found in %s\n", elpased)
}
