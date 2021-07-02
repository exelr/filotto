package filotto

import (
	"github.com/exelr/eddwise"
	"github.com/exelr/filotto/gen/filotto"
	filottobehave "github.com/exelr/filotto/gen/filotto/behave"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStartMatchAndSinglePlayerMovePlayerMove(t *testing.T) {
	var behave = filottobehave.NewFilottoBehave(t)
	behave.Given("an empty Filotto channel", func() eddwise.ImplChannel {
		return NewFilotto()
	}, func() {
		var ch = behave.Recv().(*FilottoChannel)
		behave.ThenClientJoins(1, func() {
			So(ch.Players, ShouldHaveLength, 1)
			behave.ThenClientShouldReceiveEvent("with silly name", 1, &filotto.Welcome{You: *ch.Players[1].Player})
			behave.OnQueueRequest(1, &filotto.QueueRequest{}, func() {
				So(ch.WaitingPlayers.Len(), ShouldEqual, 1)
				behave.ThenClientJoins(2, func() {
					So(ch.Players, ShouldHaveLength, 2)
					behave.ThenClientShouldReceiveEvent("with silly name", 2, &filotto.Welcome{You: *ch.Players[2].Player})

					behave.OnQueueRequest(2, &filotto.QueueRequest{}, func() {
						So(ch.WaitingPlayers.Len(), ShouldEqual, 0)
						//So(ch.Matches, ShouldHaveLength, 1)
						So(ch.Players[1].Match, ShouldEqual, ch.Players[2].Match)
						var match = ch.Players[1].Match
						So(match.Status, ShouldEqual, MatchStatusWaitPlayer1)
						behave.ThenClientShouldReceiveEvent("with info of player 2", 1, &filotto.MatchStarts{
							Rows:      Rows,
							Adversary: *match.Player2.Player,
							FirstMove: true,
							Columns:   Columns,
						})
						behave.ThenClientShouldReceiveEvent("with info of player 1", 2, &filotto.MatchStarts{
							Rows:      Rows,
							Adversary: *match.Player1.Player,
							FirstMove: false,
							Columns:   Columns,
						})
						behave.OnPlayerMove(1, &filotto.PlayerMove{
							Column: 0,
						}, func() {
							So(match.Status, ShouldEqual, MatchStatusWaitPlayer2)
							So(match.Board[0][0], ShouldEqual, 1)
							var r = uint(0)
							behave.ThenClientShouldReceiveEvent("the move of the player 1", 2, &filotto.PlayerMove{
								Player: ch.Players[1].Player,
								Column: 0,
								Row:    &r,
							})
						})
					})
				})
			})
		})

	})
}

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
