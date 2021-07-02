// Code generated by eddwise, DO NOT EDIT.

package filotto

import (
	"errors"

	"github.com/exelr/eddwise"
)

var _ eddwise.ImplChannel = (*Filotto)(nil)
var _ FilottoRecv = (*Filotto)(nil)

type FilottoRecv interface {
	OnPlayerMove(eddwise.Context, *PlayerMove) error
	OnQueueRequest(eddwise.Context, *QueueRequest) error
}

type Filotto struct {
	server eddwise.Server
	recv   FilottoRecv
}

func (ch *Filotto) Name() string {
	return "Filotto"
}

func (ch *Filotto) Bind(server eddwise.Server) error {
	ch.server = server
	return nil
}

func (ch *Filotto) SetReceiver(chr eddwise.ImplChannel) error {
	if _, ok := chr.(FilottoRecv); !ok {
		return errors.New("unexpected channel type while SetReceiver on 'Filotto' channel")
	}
	ch.recv = chr.(FilottoRecv)
	return nil
}

func (ch *Filotto) GetServer() eddwise.Server {
	return ch.server
}

func (ch *Filotto) Route(ctx eddwise.Context, evt *eddwise.EventMessage) error {
	switch evt.Name {
	default:
		return eddwise.ErrMissingServerHandler(evt.Channel, evt.Name)

	case "PlayerMove":
		var msg = &PlayerMove{}
		if err := ch.server.GetSerializer().Deserialize(evt.Body, msg); err != nil {
			return err
		}
		return ch.recv.OnPlayerMove(ctx, msg)

	case "QueueRequest":
		var msg = &QueueRequest{}
		if err := ch.server.GetSerializer().Deserialize(evt.Body, msg); err != nil {
			return err
		}
		return ch.recv.OnQueueRequest(ctx, msg)

	}
}

func (ch *Filotto) OnPlayerMove(eddwise.Context, *PlayerMove) error {
	return errors.New("event 'PlayerMove' is not handled on server")
}

func (ch *Filotto) OnQueueRequest(eddwise.Context, *QueueRequest) error {
	return errors.New("event 'QueueRequest' is not handled on server")
}

func (ch *Filotto) SendMatchEnds(client eddwise.Client, msg *MatchEnds) error {
	return client.Send(ch.Name(), msg)
}

func (ch *Filotto) SendMatchStarts(client eddwise.Client, msg *MatchStarts) error {
	return client.Send(ch.Name(), msg)
}

func (ch *Filotto) SendPlayerMove(client eddwise.Client, msg *PlayerMove) error {
	return client.Send(ch.Name(), msg)
}

func (ch *Filotto) SendWelcome(client eddwise.Client, msg *Welcome) error {
	return client.Send(ch.Name(), msg)
}

func (ch *Filotto) BroadcastMatchEnds(clients []eddwise.Client, msg *MatchEnds) error {
	return eddwise.Broadcast(ch.Name(), msg, clients)
}

func (ch *Filotto) BroadcastMatchStarts(clients []eddwise.Client, msg *MatchStarts) error {
	return eddwise.Broadcast(ch.Name(), msg, clients)
}

func (ch *Filotto) BroadcastPlayerMove(clients []eddwise.Client, msg *PlayerMove) error {
	return eddwise.Broadcast(ch.Name(), msg, clients)
}

func (ch *Filotto) BroadcastWelcome(clients []eddwise.Client, msg *Welcome) error {
	return eddwise.Broadcast(ch.Name(), msg, clients)
}

// Event structures

// MatchEnds sent from server to both players in a match when the match ends for whatever reason
type MatchEnds struct {
	Winner  Player  `json:"Winner"`
	WinLine []Point `json:"WinLine"`
	// Reason can be "line" or "player_left"
	Reason string `json:"Reason"`
}

func (evt *MatchEnds) GetEventName() string {
	return "MatchEnds"
}

// MatchStarts sent from server to two clients when the match is found for two players
type MatchStarts struct {
	Rows      uint64 `json:"Rows"`
	Adversary Player `json:"Adversary"`
	FirstMove bool   `json:"FirstMove"`
	Columns   uint64 `json:"Columns"`
}

func (evt *MatchStarts) GetEventName() string {
	return "MatchStarts"
}

type Player struct {
	Id   uint64 `json:"Id"`
	Name string `json:"Name"`
}

func (evt *Player) GetEventName() string {
	return "Player"
}

// PlayerMove sent when client performs any move. Server will relay to adversary
type PlayerMove struct {
	Player *Player `json:"Player,omitempty"` // ServerToClient
	Column uint    `json:"Column"`
	Row    *uint   `json:"Row,omitempty"` // ServerToClient
}

func (evt *PlayerMove) GetEventName() string {
	return "PlayerMove"
}

type Point struct {
	Row    uint `json:"Row"`
	Column uint `json:"Column"`
}

func (evt *Point) GetEventName() string {
	return "Point"
}

type QueueRequest struct {
}

func (evt *QueueRequest) GetEventName() string {
	return "QueueRequest"
}

// Welcome is sent from server to client whenever connects
type Welcome struct {
	You Player `json:"You"`
}

func (evt *Welcome) GetEventName() string {
	return "Welcome"
}
