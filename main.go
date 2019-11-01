package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
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
	start := time.Now()

	solveNqueens(board, 3*n)
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
	for i := 0; i < n; i++ {
		primaryDiags[i-board[i]+n-1]++
		secondaryDiags[i+board[i]]++
		rows[board[i]]++
	}
}

func solveNqueens(board Board, maxIter int) {
	var (
		col, newRow int
	)

	for {
		for i := 0; i < maxIter; i++ {
			col = getColWithMaxConf(board)
			newRow = getRowWithMinConf(col)

			rows[board[col]]--
			rows[newRow]++

			primaryDiags[col-board[col]+n-1]--
			secondaryDiags[col+board[col]]--

			primaryDiags[col-newRow+n-1]++
			secondaryDiags[col+newRow]++

			board[col] = newRow

			if !hasConflict(board) {
				// printSolution(board)
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
				fmt.Print("*")
			} else {
				fmt.Print("_")
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func getColWithMaxConf(board Board) int {
	var cs []cell
	var c cell
	max := math.MinInt32
	for i := 0; i < n; i++ {
		c = findConflictsToCell(i, board[i])
		if c.conf > max {
			max = c.conf
			cs = cs[:0]
			cs = append(cs, c)
		}
		if c.conf == max {
			cs = append(cs, c)
		}
	}

	return cs[rand.Intn(len(cs))].col
}

func hasConflict(board Board) bool {
	var c cell
	for col := 0; col < n; col++ {
		c = findConflictsToCell(col, board[col])
		if c.rConf > 1 || c.primConf > 1 || c.secCOnf > 1 {
			return true
		}
	}

	return false
}

func getRowWithMinConf(col int) int {
	var c cell
	min := math.MaxInt32
	cs := make([]cell, 0)
	for i := 0; i < n; i++ {
		c = findConflictsToCell(col, i)
		if c.conf < min {
			min = c.conf
			cs = cs[:0]
			cs = append(cs, c)
		}
		if c.conf == min {
			cs = append(cs, c)
		}
	}

	return cs[rand.Intn(len(cs))].row
}

func findConflictsToCell(col, row int) cell {
	d1, d2 := col-row+n-1, col+row
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
