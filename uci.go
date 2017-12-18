package falcona

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func uci() {
	var board *Board

	doIsReady := func(args []string) {
		fmt.Println("readyok")
	}

	doUci := func(args []string) {
		fmt.Println("id name Falcona")
		fmt.Println("id author Akhil Velagapudi")
		fmt.Println("uciok")
	}

	doPosition := func(args []string) {
		board = &Board{}
		switch args[0] {
		case "startpos":
			board.initStandard()
			args = args[1:]
		case "fen":
			fen := []string{}
			for _, token := range args[1:] {
				args = args[1:]
				if token == "moves" {
					break
				}
				fen = append(fen, token)
			}
			board.initFEN(strings.Join(fen, " "))
		default:
			return
		}

		if len(args) > 0 && args[0] == "moves" {
			for _, move := range args[1:] {
				args = args[1:]
				board.makeMove(board.pos[board.ply].moveFromString(move))
			}
		}
		board.pos[0] = board.pos[board.ply]
		board.ply = 0
		board.pos[board.ply].print()
	}

	doGo := func(args []string) {
		for i, token := range args {
			if len(args) > i+1 {
				switch token {
				case `wtime`:
					if board.pos[board.ply].side == White {
						if n, err := strconv.Atoi(args[i+1]); err == nil {
							board.parallelSearchPosition((int64(n) / 20000) + 1)
						}
					}
				case `btime`:
					if board.pos[board.ply].side == Black {
						if n, err := strconv.Atoi(args[i+1]); err == nil {
							board.parallelSearchPosition((int64(n) / 20000) + 1)
						}
					}
				}
			}
		}
	}

	var commands = map[string]func([]string){
		"isready":  doIsReady,
		"uci":      doUci,
		"position": doPosition,
		"go":       doGo,
	}

	in := bufio.NewReader(os.Stdin)
	for {
		command, err := in.ReadString('\n')
		if err != io.EOF && len(command) > 0 {
			args := strings.Split(strings.Trim(command, " \t\r\n"), " ")
			if args[0] == "quit" {
				break
			}
			if handler, ok := commands[args[0]]; ok {
				handler(args[1:])
			}
		}
	}
}
