package filotto

// Filotto is the edd channel on which the server and the players communicate
type Filotto interface {
	Enable(
		Welcome,
		MatchStarts,
		MatchEnds,
		PlayerMove,
		QueueRequest,
	)
	ServerToClient(Welcome, MatchStarts, MatchEnds) // Prevent clients to send unwanted events (prevent docs to be generated in client library)
	ClientToServer(QueueRequest)                    // Prevent server to send players a QueueRequest
}

// Welcome is sent from server to client whenever connects
type Welcome struct {
	You Player
}

type Player struct {
	Id   uint64
	Name string
}

// QueueRequest sent from client to server when the player wants to queue for a match
type QueueRequest struct {
}

// MatchStarts sent from server to two clients when the match is found for two players
type MatchStarts struct {
	Adversary *Player
	FirstMove bool
	Columns   uint64
	Rows      uint64
}

// MatchEnds sent from server to both players in a match when the match ends for whatever reason
type MatchEnds struct {
	Winner  *Player
	WinLine []Point
	Reason  string //"line" or "player_left"
}

type Point struct {
	Row    uint
	Column uint
}

// PlayerMove sent when client performs any move. Server will relay to adversary
type PlayerMove struct {
	Player *Player
	Column uint
	Row    *uint //sent only from server to client
}
