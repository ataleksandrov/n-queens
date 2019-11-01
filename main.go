package main

import (
	"fmt"
	"math"
	"math/rand"
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

type cells []cell

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
	return rand.Perm(n)
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
	return col - row + n - 1, row + col
}

func solveNqueens(board Board) {
	var (
		col, newRow, d1, d2 int
	)

	for {
		for i := 0; i < MAX_ITER; i++ {
			col = getColWithMaxConf(board)
			newRow = getRowWithMinConf(col)

			rows[board[col]]--
			rows[newRow]++

			d1, d2 = getDiags(col, board[col])
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
		if hasConflict(board) {
			board = createBoard(n)
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
	c := findAllConflicts(board)
	var cs cells
	max := math.MinInt32
	for i := 0; i < n; i++ {
		if c[i].conf > max {
			max = c[i].conf
			cs = cells{}
			cs = append(cs, c[i])
		}
		if c[i].conf == max {
			cs = append(cs, c[i])
		}
	}

	return cs[rand.Intn(len(cs))].col
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

func findAllConflicts(board Board) cells {
	cells := make(cells, n)
	for col := 0; col < n; col++ {
		cells[col] = findConflictsToCell(col, board[col])
	}

	return cells
}

func getRowWithMinConf(col int) int {
	var c cell
	min := math.MaxInt32
	cs := cells{}
	for i := 0; i < n; i++ {
		c = findConflictsToCell(col, i)
		if c.conf < min {
			min = c.conf
			cs = cells{}
			cs = append(cs, c)
		}
		if c.conf == min {
			cs = append(cs, c)
		}
	}

	return cs[rand.Intn(len(cs))].row
}

func findConflictsToCell(col, row int) cell {
	d1, d2 := getDiags(col, row)
	confs := 0

	confs += rows[row]
	confs += primaryDiags[d1]
	confs += secondaryDiags[d2]

	if board[col] == row {
		confs -= 3
	}
	return cell{
		col:      col,
		row:      row,
		conf:     confs,
		rConf:    rows[col],
		primConf: primaryDiags[d1],
		secCOnf:  secondaryDiags[d2],
	}
}

func getRandomCellWithConfs(conf int, cells cells) cell {
	cnt := 0
	for i := 0; i < len(cells); i++ {
		if cells[i].conf == conf {
			cnt++
		} else {
			break
		}
	}
	return cells[rand.Intn(cnt)]
}
