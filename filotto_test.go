package filotto

import "testing"

func TestCheckWin(t *testing.T) {
	var board = [Columns][Rows]uint64{}
	board[0][0] = 1
	board[0][1] = 1
	board[0][2] = 1
	board[0][3] = 1

	var points = CheckWin(&board, 0, 3)
	if len(points) == 0 {
		t.Fatal("expected vertical win")
	}

	board = [Columns][Rows]uint64{}
	board[0][0] = 1
	board[1][0] = 0
	board[2][0] = 1
	board[3][0] = 1

	points = CheckWin(&board, 0, 3)
	if len(points) != 0 {
		t.Fatal("expected loss")
	}

	board = [Columns][Rows]uint64{}
	board[0][0] = 1
	board[1][0] = 1
	board[2][0] = 1
	board[3][0] = 1

	points = CheckWin(&board, 0, 0)
	if len(points) == 0 {
		t.Fatal("expected horizontal win")
	}

	board = [Columns][Rows]uint64{}
	board[0][0] = 1
	board[1][0] = 1
	board[2][0] = 1
	board[3][0] = 1
	board[4][0] = 1

	points = CheckWin(&board, 3, 0)
	if len(points) == 0 {
		t.Fatal("expected horizontal long win")
	}

	board = [Columns][Rows]uint64{}
	board[0][0] = 1
	board[1][1] = 1
	board[2][2] = 1
	board[3][3] = 1

	points = CheckWin(&board, 0, 0)
	if len(points) == 0 {
		t.Fatal("expected diagonal bl-tr")
	}

	board = [Columns][Rows]uint64{}
	board[3][3] = 1
	board[2][2] = 1
	board[1][1] = 1
	board[0][0] = 1

	points = CheckWin(&board, 1, 1)
	if len(points) == 0 {
		t.Fatal("expected diagonal tl-br")
	}

}
