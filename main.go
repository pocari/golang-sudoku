package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func eachLine(r io.Reader, handler func(string)) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		handler(s.Text())
	}
}

func initBoard() [][]int {
	ret := make([][]int, 9)
	for i := range ret {
		ret[i] = make([]int, 9)
	}
	return ret
}

func posToRowCol(pos int) (int, int) {
	return pos / 9, pos % 9
}

func loadQuestion(numbers []string) [][]int {
	board := initBoard()
	var v int
	for i, e := range numbers {
		r, c := posToRowCol(i)
		if e == "." {
			v = 0
		} else {
			var err error
			if v, err = strconv.Atoi(e); err != nil {
				panic("hoge")
			}
		}
		board[r][c] = v
	}
	return board
}

func dumpBoard(board [][]int) {
	for _, row := range board {
		for _, cell := range row {
			if cell == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%v ", cell)
			}
		}
		fmt.Println()
	}
}

func loadQuestions(r io.Reader) [][][]int {
	i := 0
	questions := make([][][]int, 0)
	eachLine(r, func(text string) {
		board := loadQuestion(strings.Split(text, ""))
		questions = append(questions, board)
		i++
	})

	return questions
}

func eachCells(board [][]int, fn func(int, int, int)) {
	for r, row := range board {
		for c, cell := range row {
			fn(r, c, cell)
		}
	}
}

func clearCheck(check []int) {
	for i := range check {
		check[i] = 0
	}
}

func blockID(r, c int) int {
	return (r/3)*3 + c/3
}

func isValid(board [][]int, r, c int) bool {
	check := make([]int, 10)
	for i := 0; i < 9; i++ {
		j := board[r][i]
		check[j]++
		if j != 0 && check[j] > 1 {
			return false
		}
	}

	clearCheck(check)
	for i := 0; i < 9; i++ {
		j := board[i][c]
		check[j]++
		if j != 0 && check[j] > 1 {
			return false
		}
	}

	clearCheck(check)
	b := blockID(r, c)
	for i := 0; i < 9; i++ {
		offsetr := b / 3 * 3
		offsetc := b % 3 * 3
		rx := (i / 3) + offsetr
		cx := (i % 3) + offsetc
		j := board[rx][cx]
		check[j]++
		if j != 0 && check[j] > 1 {
			return false
		}
	}
	return true
}

func solveSudokuHelper(board [][]int, pos int) bool {
	if pos == 81 {
		return true
	}

	r, c := posToRowCol(pos)
	v := board[r][c]
	if v == 0 {
		// 数字が設定されていなかったら
		for i := 1; i <= 9; i++ {
			// どれかを試す
			board[r][c] = i
			if isValid(board, r, c) {
				if solveSudokuHelper(board, pos+1) {
					return true
				}
			}
			board[r][c] = 0
		}
		return false
	}
	return solveSudokuHelper(board, pos+1)
}

// とりあえず全部バックトラックで解く
func solveSudoku(board [][]int) bool {
	return solveSudokuHelper(board, 0)
}

func main() {
	qs := loadQuestions(os.Stdin)
	for i, q := range qs {
		fmt.Printf("No. %v\n", i+1)
		fmt.Println("q ---------------")
		dumpBoard(q)
		fmt.Println("a ---------------")
		if !solveSudoku(q) {
			panic("error solve")
		}
		dumpBoard(q)
		fmt.Println()
	}
}
