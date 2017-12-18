package falcona

import (
	"fmt"
	"sort"
	"time"
)

type SearchInfo struct {
	startTime int64
	endTime   int64
	nodes     int64
	stopped   bool
}

func (info *SearchInfo) check() {
	now := time.Now().Unix()
	if now >= info.endTime {
		info.stopped = true
	}
}

func (board *Board) isRepetition() bool {
	pos := &board.pos[board.ply]
	for ply := 0; ply < board.ply; ply++ {
		if pos.poskey == board.pos[ply].poskey {
			return true
		}
	}
	return false
}

func (board *Board) quiescence(info *SearchInfo, alpha, beta int) int {
	pos := &board.pos[board.ply]
	if (info.nodes & 4095) == 0 {
		info.check()
	}
	info.nodes++
	if board.isRepetition() || pos.fiftymove >= 100 {
		return 0
	}
	if board.ply > MaxPly {
		return board.evaluate()
	}

	score := board.evaluate()
	if score >= beta {
		return beta
	}
	if score > alpha {
		alpha = score
	}

	ml := pos.generateCaptures()
	score = -100000

	sort.Slice(ml.moves[:], func(i, j int) bool {
		return ml.moves[i].score > ml.moves[j].score
	})

	for i := 0; i < ml.count; i++ {
		move := ml.moves[i].move
		if !board.makeMove(move) {
			continue
		}
		score = -1 * board.quiescence(info, -beta, -alpha)
		board.unmakeMove()

		if info.stopped {
			return 0
		}
		if score > alpha {
			if score >= beta {
				return beta
			}
			alpha = score
		}
	}
	return alpha
}

func (board *Board) alphaBeta(info *SearchInfo, alpha, beta, depth int) int {
	pos := &board.pos[board.ply]
	if (info.nodes & 4095) == 0 {
		info.check()
	}
	if depth == 0 {
		return board.quiescence(info, alpha, beta)
	}
	info.nodes++
	if board.isRepetition() || pos.fiftymove >= 100 {
		return 0
	}
	if board.ply > MaxPly {
		return board.evaluate()
	}

	pvMove := uint32(0)
	score := -100000

	found, pvMove, score := board.probeEntry(alpha, beta, depth)
	if found {
		table.cut++
		return score
	}

	ml := pos.generateMoves()
	legal := 0
	oldAlpha := alpha
	bestMove := uint32(0)
	bestScore := -100000

	if pvMove > 0 {
		for i := 0; i < ml.count; i++ {
			if ml.moves[i].move == pvMove {
				ml.moves[i].score = 20000000
				break
			}
		}
	}

	sort.Slice(ml.moves[:], func(i, j int) bool {
		return ml.moves[i].score > ml.moves[j].score
	})

	for i := 0; i < ml.count; i++ {
		move := ml.moves[i].move
		if !board.makeMove(move) {
			continue
		}
		legal++
		score = -1 * board.alphaBeta(info, -beta, -alpha, depth-1)
		board.unmakeMove()

		if info.stopped {
			return 0
		}
		if score > bestScore {
			bestScore = score
			bestMove = move
			if score > alpha {
				if score >= beta {
					board.storeEntry(bestMove, beta, depth, FlagBeta)
					return beta
				}
				alpha = score
				bestMove = move
			}
		}
	}
	if legal == 0 {
		if pos.attackedBy(int(pos.kings[pos.side]), pos.side^1) {
			return -1000000
		}
		return 0
	}
	if oldAlpha != alpha {
		board.storeEntry(bestMove, bestScore, depth, FlagExact)
	} else {
		board.storeEntry(bestMove, alpha, depth, FlagAlpha)
	}
	return alpha
}

func (board *Board) searchPosition(seconds int64) {
	bestMove := uint32(0)
	bestScore := -10000000
	info := &SearchInfo{time.Now().Unix(), time.Now().Unix() + seconds, 0, false}
	table.cut, table.hit, table.newWrite, table.overWrite = 0, 0, 0, 0
	for depth := 1; depth < MaxDepth; depth++ {
		bestScore = board.alphaBeta(info, -100000, 100000, depth)
		if info.stopped {
			break
		}
		numPV, pvMoves := board.getPV(depth)
		bestMove = pvMoves[0]

		fmt.Print("info score cp ", bestScore, " depth ", depth, " nodes ", info.nodes, " time ", (time.Now().Unix()-info.startTime)*1000, " pv")
		for i := 0; i < numPV; i++ {
			fmt.Print(" " + moveToString(pvMoves[i]))
		}
		fmt.Println()
		fmt.Println("hit", table.hit, " newwrite", table.newWrite, " overwrite", table.overWrite)
	}
	fmt.Println("bestmove", moveToString(bestMove))
}

func moveToString(move uint32) string {
	from := from(move)
	to := to(move)
	promoted := promoted(move)
	pieceChars := [12]string{"p", "p", "n", "n", "b", "b", "r", "r", "q", "q", "k", "k"}
	fr, fc := toRowCol(int(from))
	tr, tc := toRowCol(int(to))
	movestr := string('a'+fc) + string('1'+fr) + string('a'+tc) + string('1'+tr)
	if promoted != 0xF {
		movestr += pieceChars[promoted]
	}
	return movestr
}
