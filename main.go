package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	maxn     = uint(9)
	allBitOn = (1 << maxn) - 1
	numBits  = []int{
		0x000, // 0 000000000
		0x001, // 1 000000001
		0x002, // 2 000000010
		0x004, // 3 000000100
		0x008, // 4 000001000
		0x010, // 5 000010000
		0x020, // 6 000100000
		0x040, // 7 001000000
		0x080, // 8 010000000
		0x100, // 9 100000000
	}
)

type question struct {
	board [][]int
	uf    *usedFlags
}

func (q *question) candidates(r, c int) int {
	b := blockID(r, c)
	used := q.uf.row[r] | q.uf.col[c] | q.uf.blk[b]
	return allBitOn ^ used
}

func (q *question) enableNumbers(r, c int) []int {
	unused := q.candidates(r, c)
	enableNumbers := make([]int, 0)
	for i := 1; i <= int(maxn); i++ {
		if unused&numBits[i] != 0 {
			enableNumbers = append(enableNumbers, i)
		}
	}
	return enableNumbers
}

func (q *question) canPut(number, r, c int) bool {
	return q.uf.canPut(number, r, c)
}

func (q *question) putNumber(number, r, c int) {
	q.uf.putNumber(number, r, c)
	q.board[r][c] = number
}

func (q *question) removeNumber(r, c int) {
	n := q.board[r][c]
	q.board[r][c] = 0
	q.uf.removeNumber(n, r, c)
}

type usedFlags struct {
	row []int
	col []int
	blk []int
}

func newUsedFlags(board [][]int) *usedFlags {
	uf := new(usedFlags)
	uf.row = make([]int, maxn)
	uf.col = make([]int, maxn)
	uf.blk = make([]int, maxn)

	eachCells(board, func(r, c, cell int) {
		uf.putNumber(board[r][c], r, c)
	})

	return uf
}

func (uf *usedFlags) canPut(number, r, c int) bool {
	// uf.dumpFlags()
	// fmt.Printf("(number, r, c) ... (%v, %v, %v)\n", number, r, c)
	if uf.row[r]&numBits[number] != 0 {
		// fmt.Printf("(number, r, c) ... (%v, %v, %v, ng_row)\n", number, r, c)
		return false
	}

	if uf.col[c]&numBits[number] != 0 {
		// fmt.Printf("(number, r, c) ... (%v, %v, %v, ng_col)\n", number, r, c)
		return false
	}

	b := blockID(r, c)
	if uf.blk[b]&numBits[number] != 0 {
		// fmt.Printf("(number, r, c) ... (%v, %v, %v, ng_blk)\n", number, r, c)
		return false
	}
	// fmt.Printf("(number, r, c) ... (%v, %v, %v, ok)\n", number, r, c)
	return true
}

func (uf *usedFlags) putNumber(number, r, c int) {
	if uf == nil {
		panic("uf is nil")
	}
	if number == 0 {
		return
	}
	b := blockID(r, c)
	uf.row[r] |= numBits[number]
	uf.col[c] |= numBits[number]
	uf.blk[b] |= numBits[number]
}

func (uf *usedFlags) removeNumber(number, r, c int) {
	if uf == nil {
		panic("uf is nil")
	}
	if number == 0 {
		return
	}
	xor := allBitOn ^ numBits[number]
	// fmt.Printf("remove (number, r, c) ... (%v, %v, %v)\n", number, r, c)
	// uf.dumpFlagHelper(xor)
	b := blockID(r, c)
	uf.row[r] = uf.row[r] & xor
	uf.col[c] = uf.col[c] & xor
	uf.blk[b] = uf.blk[b] & xor
}

func (*usedFlags) dumpFlagHelper(flags int) {
	for i := 0; i < int(maxn); i++ {
		if (flags & numBits[int(maxn)-i]) != 0 {
			fmt.Print("1 ")
		} else {
			fmt.Print("0 ")
		}
	}
	fmt.Println()
}

func (uf *usedFlags) dumpFlags() {
	for i := 0; i < int(maxn); i++ {
		fmt.Printf("%v ", int(maxn)-i)
	}
	fmt.Println()
	fmt.Println("row:")
	for _, r := range uf.row {
		uf.dumpFlagHelper(r)
	}
	fmt.Println("col")
	for _, r := range uf.col {
		uf.dumpFlagHelper(r)
	}
	fmt.Println("blk")
	for _, r := range uf.blk {
		uf.dumpFlagHelper(r)
	}
}

func eachLine(r io.Reader, handler func(string)) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		handler(s.Text())
	}
}

func blankBoard() [][]int {
	ret := make([][]int, maxn)
	for i := range ret {
		ret[i] = make([]int, 9)
	}
	return ret
}

func posToRowCol(pos int) (int, int) {
	return pos / int(maxn), pos % int(maxn)
}

func loadQuestion(numbers []string) *question {
	board := blankBoard()
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

	return &question{
		board, newUsedFlags(board),
	}
}

func (q *question) dumpBoard() {
	for r, row := range q.board {
		if r != 0 && r%3 == 0 {
			fmt.Println("------+------+------")
		}
		for c, cell := range row {
			if c != 0 && c%3 == 0 {
				fmt.Print("|")
			}
			if cell == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%v ", cell)
			}
		}
		fmt.Println()
	}
}

func loadQuestions(r io.Reader) []*question {
	i := 0
	questions := make([]*question, 0)
	eachLine(r, func(text string) {
		q := loadQuestion(strings.Split(text, ""))
		questions = append(questions, q)
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

func blockID(r, c int) int {
	return (r/3)*3 + c/3
}

func solveSudokuHelper(q *question, pos int) bool {
	if pos == int(maxn*maxn) {
		return true
	}

	r, c := posToRowCol(pos)
	v := q.board[r][c]
	if v == 0 {
		nums := q.enableNumbers(r, c)
		for _, i := range nums {
			// どれかを試す
			if q.canPut(i, r, c) {
				q.putNumber(i, r, c)
				if solveSudoku(q, true) {
					return true
				}
				q.removeNumber(r, c)
			}
		}
		return false
	}
	return solveSudokuHelper(q, pos+1)
}

type journal struct {
	number int
	row    int
	col    int
}

func heuristicA(q *question) []journal {
	journals := make([]journal, 0)
	for {
		changed := false
		eachCells(q.board, func(r, c, num int) {
			if num == 0 {
				candidates := q.candidates(r, c)
				for i := 1; i <= int(maxn); i++ {
					if candidates == numBits[i] {
						// 候補の数字が一つだけならそのマスはその数字で確定
						q.putNumber(i, r, c)
						journals = append(journals, journal{number: i, row: r, col: c})
						//fmt.Printf("(num, r, c) ... (%v, %v, %v)\n", i, r, c)
						//q.dumpBoard(q.board)
						changed = true
					}
				}
			}
		})

		// 一回でも数字が確定していたら、その数字を元に確定する数字があるかもしれないので、
		// もう一度確定サーチ実行
		if !changed {
			break
		} else {
			//fmt.Println("do method A")
		}
	}
	return journals
}

// 数独を解く
func solveSudoku(q *question, tryMode bool) bool {
	//fmt.Printf("solveSudoku tryMode: %v\n", tryMode)
	//確定サーチA 候補が一つの場合はその数字で確定させる
	journals := heuristicA(q)
	if !solveSudokuHelper(q, 0) && tryMode {
		for _, j := range journals {
			q.removeNumber(j.row, j.col)
		}
		return false
	}
	return true
}

func main() {
	qs := loadQuestions(os.Stdin)
	for i, q := range qs {
		fmt.Printf("No. %v\n", i+1)
		fmt.Println("q")
		q.dumpBoard()
		fmt.Println("\nans")
		start := time.Now()
		if !solveSudoku(q, false) {
			panic("error solve")
		}
		end := time.Now()

		q.dumpBoard()
		fmt.Printf("%f sec\n", (end.Sub(start)).Seconds())
		fmt.Println()

	}
}
