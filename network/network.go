package network

import (
	"encoding/json"
	"net"
	"tic-tac-go/game"
)

type AppState int
const (
	Menu AppState = iota
	Hosting
	Joining
	Playing
)

func HostGame(port string) (net.Conn, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func JoinGame(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SendGameState(conn net.Conn, g game.Game) error {
	encoder := json.NewEncoder(conn)
	return encoder.Encode(g)
}

func ReceiveGameState(conn net.Conn) (game.Game, error) {
	var g game.Game
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&g)
	return g, err
}
