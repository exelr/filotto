package filotto

import (
	"container/list"
	"fmt"
	"log"
	"sync"

	"github.com/Pallinder/go-randomdata"
	"github.com/exelr/eddwise"
	"github.com/exelr/filotto/gen/filotto"
)

const (
	Columns = 7
	Rows    = 6
)

type MatchStatus uint

const (
	MatchStatusPlaying MatchStatus = iota
	MatchStatusWaitPlayer1
	MatchStatusWaitPlayer2
	MatchStatusWinPlayer1
	MatchStatusWinPlayer2
)

type WinReason string

const (
	WinReasonLine             WinReason = "line"
	WinReasonRemotePlayerLeft WinReason = "player_left"
)

type Match struct {
	Player1, Player2 *PlayerMatch
	Status           MatchStatus
	WinReason
	Board [Columns][Rows]uint64
}

func (m *Match) GetPlayerNumber(pm *PlayerMatch) uint {
	if pm == m.Player1 {
		return 1
	}
	return 2
}

func (m *Match) GetAdversary(pm *PlayerMatch) *PlayerMatch {
	if pm != m.Player1 {
		return m.Player1
	}
	return m.Player2
}

type PlayerMatch struct {
	Player *filotto.Player
	Match  *Match
	Client *eddwise.Client
	Queue  *list.Element
}

type FilottoChannel struct {
	filotto.Filotto
	Players          map[uint64]*PlayerMatch
	PlayersMx        sync.RWMutex
	Matches          []*Match
	WaitingPlayers   *list.List
	WaitingPlayersMx sync.Mutex
}

func NewFilotto() *FilottoChannel {
	return &FilottoChannel{Players: map[uint64]*PlayerMatch{}, WaitingPlayers: list.New()}
}

func (ch *FilottoChannel) AddPlayerToWaitingFront(p *PlayerMatch, unsafe ...bool) {
	if len(unsafe) == 0 || !unsafe[0] {
		ch.WaitingPlayersMx.Lock()
		defer ch.WaitingPlayersMx.Unlock()
	}
	p.Queue = ch.WaitingPlayers.PushBack(p)
}
func (ch *FilottoChannel) AddPlayerToWaitingBack(p *PlayerMatch, unsafe ...bool) {
	if len(unsafe) == 0 || !unsafe[0] {
		ch.WaitingPlayersMx.Lock()
		defer ch.WaitingPlayersMx.Unlock()
	}
	p.Queue = ch.WaitingPlayers.PushFront(p)
}

func (ch *FilottoChannel) RemovePlayerFromQueue(p *PlayerMatch) {
	if p.Queue == nil {
		return
	}
	ch.WaitingPlayersMx.Lock()
	defer ch.WaitingPlayersMx.Unlock()
	ch.WaitingPlayers.Remove(p.Queue)
}

// PickFirstAvailablePlayer is not thread-safe
func (ch *FilottoChannel) PickFirstAvailablePlayer() *PlayerMatch {
	for ch.WaitingPlayers.Len() > 0 {
		var p = ch.WaitingPlayers.Front().Value.(*PlayerMatch)
		ch.WaitingPlayers.Remove(ch.WaitingPlayers.Front())
		p.Queue = nil
		if !p.Client.Closed {
			return p
		}
	}
	return nil

}

func (ch *FilottoChannel) CheckStartMatch() error {
	var p1, p2 = func() (*PlayerMatch, *PlayerMatch) {
		ch.WaitingPlayersMx.Lock()
		defer ch.WaitingPlayersMx.Unlock()
		if ch.WaitingPlayers.Len() < 2 {
			return nil, nil
		}
		var p1 = ch.PickFirstAvailablePlayer()
		if p1 == nil {
			return nil, nil
		}
		var p2 = ch.PickFirstAvailablePlayer()
		if p2 == nil {
			ch.AddPlayerToWaitingFront(p1, true)
			return nil, nil
		}
		return p1, p2
	}()
	if p1 == nil || p2 == nil {
		return nil
	}
	var match = &Match{
		Player1: p1,
		Player2: p2,
		Board:   [Columns][Rows]uint64{},
		Status:  MatchStatusWaitPlayer1,
	}

	p1.Match = match
	p2.Match = match

	if err := ch.SendMatchStarts(p1.Client, &filotto.MatchStarts{Adversary: p2.Player, Columns: Columns, Rows: Rows, FirstMove: true}); err != nil {
		//return err
	}

	if err := ch.SendMatchStarts(p2.Client, &filotto.MatchStarts{Adversary: p1.Player, Columns: Columns, Rows: Rows}); err != nil {
		//return err
	}

	return nil

}

func (ch *FilottoChannel) Connected(c *eddwise.Client) error {

	log.Println("User connected", c.GetId())

	var name = randomdata.SillyName()
	var id = c.GetId()
	var player = &PlayerMatch{
		Player: &filotto.Player{
			Id:   id,
			Name: name,
		},
		Client: c,
		Match:  nil,
		Queue:  nil,
	}
	ch.PlayersMx.Lock()
	ch.Players[id] = player
	ch.PlayersMx.Unlock()
	//ch.AddPlayerToWaitingBack(player)
	return ch.SendWelcome(c, &filotto.Welcome{You: *player.Player})
}

func (ch *FilottoChannel) Disconnected(c *eddwise.Client) error {
	var id = c.GetId()
	ch.PlayersMx.Lock()
	var pm = ch.Players[id]
	delete(ch.Players, id)
	ch.PlayersMx.Unlock()
	if pm.Match != nil {
		var status = MatchStatusWinPlayer1
		var adversary = pm.Match.Player1
		if adversary == pm {
			status = MatchStatusWinPlayer2
			adversary = pm.Match.Player2
		}
		pm.Match.Status = status
		pm.Match.WinReason = WinReasonRemotePlayerLeft
		_ = ch.SendMatchEnds(adversary.Client, &filotto.MatchEnds{
			Winner:  adversary.Player,
			WinLine: nil,
			Reason:  string(WinReasonRemotePlayerLeft),
		})

		adversary.Match = nil
		ch.RemovePlayerFromQueue(pm)
		//ch.AddPlayerToWaitingBack(adversary)
	}
	log.Println("User disconnected", c.GetId())
	return nil
}

func PerformMove(board *[Columns][Rows]uint64, id uint64, column uint) (uint, error) {
	for row := uint(0); row < Rows; row++ {
		if board[column][row] == 0 {
			board[column][row] = id
			return row, nil
		}
	}
	return 0, fmt.Errorf("column %d is full", column)
}

func CheckWin(board *[Columns][Rows]uint64, column, row uint) []filotto.Point {
	var id = board[column][row]
	if id == 0 {
		return nil
	}
	//check vertical - no need to search above
	if row >= 3 {
		if board[column][row-1] == id && board[column][row-2] == id && board[column][row-3] == id {
			return []filotto.Point{
				{Column: column, Row: row},
				{Column: column, Row: row - 1},
				{Column: column, Row: row - 2},
				{Column: column, Row: row - 3},
			}
		}
	}

	//search horizontal
	var coll = column
	for coll > 0 {
		if id != board[coll-1][row] {
			break
		}
		coll--
	}
	var colr = column
	for colr < Columns-1 {
		if id != board[colr+1][row] {
			break
		}
		colr++
	}

	if colr-coll >= 3 {
		return []filotto.Point{
			{Column: coll, Row: row},
			{Column: coll + 1, Row: row},
			{Column: coll + 2, Row: row},
			{Column: coll + 3, Row: row},
		}
	}

	//search diagonally bl-tr
	var dl1 uint
	for {
		if dl1 == row || dl1 == column {
			break
		}
		if id != board[column-dl1-1][row-dl1-1] {
			break
		}
		dl1++
	}

	var dr1 uint
	for {
		if dr1+row == Rows-1 || dr1+column == Columns-1 {
			break
		}
		if id != board[column+dr1+1][row+dr1+1] {
			break
		}
		dr1++
	}

	if dr1+dl1 >= 3 {
		return []filotto.Point{
			{Column: column - dl1, Row: row - dl1},
			{Column: column - dl1 + 1, Row: row - dl1 + 1},
			{Column: column - dl1 + 2, Row: row - dl1 + 2},
			{Column: column - dl1 + 3, Row: row - dl1 + 3},
		}
	}

	//search diagonally tl-br
	var dl2 uint
	for {
		if dl2+row == Rows-1 || dl2 == column {
			break
		}
		if id != board[column-dl2-1][row+dl2+1] {
			break
		}
		dl2++
	}

	var dr2 uint
	for {
		if dr2 == row || dr2+column == Columns-1 {
			break
		}
		if id != board[column+dr2][row-dr2] {
			break
		}
		dr2++
	}

	if dr2+dl2 >= 3 {
		return []filotto.Point{
			{Column: column - dl2, Row: row + dl2},
			{Column: column - dl2 + 1, Row: row + dl2 - 1},
			{Column: column - dl2 + 2, Row: row + dl2 - 2},
			{Column: column - dl2 + 3, Row: row + dl2 - 3},
		}
	}

	return nil

}

func (ch *FilottoChannel) OnPlayerMove(ctx filotto.FilottoContext, playermove *filotto.PlayerMove) error {

	//log.Println("received event PlayerMove:", playermove, "from", ctx.GetClient().GetId())

	ch.PlayersMx.RLock()
	var pm = ch.Players[ctx.GetClient().GetId()]
	ch.PlayersMx.RUnlock()
	if pm == nil || pm.Match == nil {
		return fmt.Errorf("you cannot move if you are not in a match")
	}

	var playerN = pm.Match.GetPlayerNumber(pm)
	if !(playerN == 1 && pm.Match.Status == MatchStatusWaitPlayer1) && !(playerN == 2 && pm.Match.Status == MatchStatusWaitPlayer2) {
		return fmt.Errorf("you cannot move right now")
	}

	if playermove.Column >= Columns {
		return fmt.Errorf("invalid column")
	}

	var row, err = PerformMove(&pm.Match.Board, pm.Player.Id, playermove.Column)
	if err != nil {
		return err
	}

	var adv = pm.Match.GetAdversary(pm)
	_ = ch.SendPlayerMove(adv.Client, &filotto.PlayerMove{
		Player: pm.Player,
		Column: playermove.Column,
		Row:    &row,
	})

	_ = ch.SendPlayerMove(pm.Client, &filotto.PlayerMove{
		Player: pm.Player,
		Column: playermove.Column,
		Row:    &row,
	})

	var winLine = CheckWin(&pm.Match.Board, playermove.Column, row)
	if winLine == nil {
		if playerN == 1 {
			pm.Match.Status = MatchStatusWaitPlayer2
		} else {
			pm.Match.Status = MatchStatusWaitPlayer1
		}
		return nil
	}

	pm.Match.WinReason = WinReasonLine
	if playerN == 1 {
		pm.Match.Status = MatchStatusWinPlayer1
	} else {
		pm.Match.Status = MatchStatusWinPlayer2
	}

	_ = ch.SendMatchEnds(adv.Client, &filotto.MatchEnds{
		Winner:  pm.Player,
		WinLine: winLine,
		Reason:  string(WinReasonLine),
	})

	_ = ch.SendMatchEnds(pm.Client, &filotto.MatchEnds{
		Winner:  pm.Player,
		WinLine: winLine,
		Reason:  string(WinReasonLine),
	})

	pm.Match = nil
	adv.Match = nil

	//ch.AddPlayerToWaitingBack(adv) // precedence to the loser!
	//ch.AddPlayerToWaitingBack(pm)

	return nil
}

func (ch *FilottoChannel) OnQueueRequest(ctx filotto.FilottoContext, _ *filotto.QueueRequest) error {
	ch.PlayersMx.RLock()
	var pm = ch.Players[ctx.GetClient().GetId()]
	ch.PlayersMx.RUnlock()
	if pm.Match != nil {
		return fmt.Errorf("cannot queue while in a match")
	}
	if pm.Queue != nil {
		return fmt.Errorf("you are already in queue")
	}

	ch.AddPlayerToWaitingBack(pm)

	_ = ch.CheckStartMatch()

	return nil
}
