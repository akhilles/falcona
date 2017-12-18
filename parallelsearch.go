package falcona

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

const NumThreads = 4
const CSSize = 32768
const CSWays = 4
const DeferDepth = 3

var currentlySearching [CSSize][CSWays]uint64

func deferMove(moveHash uint64, depth int) bool {
	if depth < DeferDepth {
		return false
	}
	index := moveHash & (CSSize - 1)
	for i := 0; i < CSWays; i++ {
		if currentlySearching[index][i] == moveHash {
			return true
		}
	}
	return false
}

func startingSearch(moveHash uint64, depth int) {
	if depth < DeferDepth {
		return
	}
	index := moveHash & (CSSize - 1)
	for i := 0; i < CSWays; i++ {
		if currentlySearching[index][i] == 0 {
			currentlySearching[index][i] = moveHash
			return
		}
		if currentlySearching[index][i] == moveHash {
			return
		}
	}
	currentlySearching[index][0] = moveHash
}

func finishedSearch(moveHash uint64, depth int) {
	if depth < DeferDepth {
		return
	}
	index := moveHash & (CSSize - 1)
	for i := 0; i < CSWays; i++ {
		if currentlySearching[index][i] == moveHash {
			currentlySearching[index][i] = 0
		}
	}
}

func (board *Board) parallelAlphaBeta(info *SearchInfo, alpha, beta, depth int) int {
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
	mlDefer := &MoveList{[128]Move{}, 0}
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

		if i == 0 {
			if !board.makeMove(move) {
				continue
			}
			legal++
			score = -1 * board.parallelAlphaBeta(info, -beta, -alpha, depth-1)
			board.unmakeMove()
		} else {
			moveHash := pos.poskey ^ ((uint64(move) * 1664525) + 1013904223)
			if deferMove(moveHash, depth) {
				mlDefer.moves[mlDefer.count] = Move{move, 0}
				mlDefer.count++
				continue
			}
			if !board.makeMove(move) {
				continue
			}
			legal++
			startingSearch(moveHash, depth)
			score = -1 * board.parallelAlphaBeta(info, -beta, -alpha, depth-1)
			finishedSearch(moveHash, depth)
			board.unmakeMove()
		}

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
	for i := 0; i < mlDefer.count; i++ {
		move := mlDefer.moves[i].move
		if !board.makeMove(move) {
			continue
		}
		legal++
		score = -1 * board.parallelAlphaBeta(info, -beta, -alpha, depth-1)
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

func (board *Board) parallelSearchPosition(seconds int64) {
	bestMove := uint32(0)
	bestScore := -10000000
	info := &SearchInfo{time.Now().Unix(), time.Now().Unix() + seconds, 0, false}
	var startiter, elapsed int64
	table.cut, table.hit, table.newWrite, table.overWrite = 0, 0, 0, 0
	for depth := 1; depth < MaxDepth; depth++ {
		startiter = time.Now().Unix()
		var wg sync.WaitGroup
		wg.Add(NumThreads)
		for g := 0; g < NumThreads; g++ {
			go func() {
				local := &Board{}
				*local = *board
				defer wg.Done()
				bestScore = local.parallelAlphaBeta(info, -100000, 100000, depth)
			}()
		}
		wg.Wait()

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

		elapsed = time.Now().Unix() - startiter
		if (info.endTime - time.Now().Unix()) < 4*elapsed {
			break
		}
	}
	fmt.Println("bestmove", moveToString(bestMove))
}
