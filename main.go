package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var (
	MAX_ITER       int
	n              int
	rows           []int
	primaryDiags   []int
	secondaryDiags []int
	board          Board
)

type Board []int

type cell struct {
	col, row, conf, rConf, primConf, secCOnf int
}

func main() {
	//fmt.Println("Enter N:")
	//if _, e := fmt.Scanf("%d", &n); e != nil {
	//	fmt.Println("N must be an integer!")
	//	os.Exit(1)
	//}
	n = 10000
	board = createBoard(n)
	initRowDiags(board)

	MAX_ITER = 3 * n
	start := time.Now()
	solveNqueens(board)
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func createBoard(n int) Board {
	board := make([]int, n)
	for i := 0; i < n; i++ {
		board[i] = i
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { board[i], board[j] = board[j], board[i] })
	return board
}

func initRowDiags(board Board) {
	rows = make([]int, n)
	primaryDiags = make([]int, 2*n-1)
	secondaryDiags = make([]int, 2*n-1)
	for i, q := range board {
		d1, d2 := getDiags(i, q)
		primaryDiags[d1]++
		secondaryDiags[d2]++
		rows[q]++
	}
}

func getDiags(col, row int) (int, int) {
	d1 := col - row + n - 1
	d2 := row + col
	return d1, d2
}

func solveNqueens(board Board) {
	var (
		col, newRow int
	)

	for {
		for i := 0; i < MAX_ITER; i++ {
			col = getColWithMaxConf(board)
			newRow = getRowWithMinConf(col, board[col])

			if findConflictsToCell(col, board[col]).conf >= findConflictsToCell(col, newRow).conf {
				rows[board[col]]--
				rows[newRow]++

				d1, d2 := getDiags(col, board[col])
				primaryDiags[d1]--
				secondaryDiags[d2]--

				d1, d2 = getDiags(col, newRow)
				primaryDiags[d1]++
				secondaryDiags[d2]++

				board[col] = newRow

				if !hasConflict(board) {
					//printSolution(board)
					return
				}
			}

		}
		if hasConflict(board) {
			board = createBoard(n)
			//fmt.Println("Reset to", board)
			initRowDiags(board)
		}
	}
}

func printSolution(board Board) {
	for j := 0; j < n; j++ {
		for _, v := range board {
			if v == j {
				fmt.Print("Q")
			} else {
				fmt.Print("_")
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func getColWithMaxConf(board Board) int {
	cells := findAllConflicts(board)
	sort.SliceStable(cells, func(i, j int) bool {
		return cells[i].conf > cells[j].conf
	})

	return getRandomCellWithConfs(cells[0].conf, cells).col
}

func hasConflict(board Board) bool {
	cells := findAllConflicts(board)
	for _, c := range cells {
		if c.rConf > 1 || c.primConf > 1 || c.secCOnf > 1 {
			return true
		}
	}

	return false
}

func findAllConflicts(board Board) []*cell {
	cells := make([]*cell, n)
	for col := 0; col < n; col++ {
		cells[col] = findConflictsToCell(col, board[col])
	}

	return cells
}

func getRowWithMinConf(col, row int) int {
	cells := make([]*cell, 0)
	for i := 0; i < n; i++ {
		if i == row {
			continue
		}
		cells = append(cells, findConflictsToCell(col, i))
	}

	sort.SliceStable(cells, func(i, j int) bool {
		return cells[i].conf < cells[j].conf
	})

	return getRandomCellWithConfs(cells[0].conf, cells).row
}

func findConflictsToCell(col, row int) *cell {
	d1, d2 := getDiags(col, row)
	confs := 0
	isQueen := board[col] == row

	if (isQueen && rows[row] > 1) || (!isQueen && rows[row] > 0) {
		confs++
	}
	if (isQueen && primaryDiags[d1] > 1) || (!isQueen && primaryDiags[d1] > 0) {
		confs++
	}
	if (isQueen && secondaryDiags[d2] > 1) || (!isQueen && secondaryDiags[d2] > 0) {
		confs++
	}
	return &cell{
		col:      col,
		row:      row,
		conf:     confs,
		rConf:    rows[col],
		primConf: primaryDiags[d1],
		secCOnf:  secondaryDiags[d2],
	}
}

func getRandomCellWithConfs(conf int, cells []*cell) *cell {
	cnt := 0
	for i := 0; i < len(cells); i++ {
		if cells[i].conf == conf {
			cnt++
		} else {
			break
		}
	}
	index := rand.Intn(cnt)
	return cells[index]
}
