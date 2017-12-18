package falcona

import (
	"sync"
)

const (
	FlagNone = iota
	FlagAlpha
	FlagBeta
	FlagExact
)

type HashEntry struct {
	poskey uint64
	move   uint32
	score  int
	depth  int
	flag   int
}

type HashTable struct {
	sync.RWMutex
	entries   [8000000]HashEntry
	count     int
	newWrite  int
	overWrite int
	hit       int
	cut       int
}

func (board *Board) probePV() uint32 {
	pos := &board.pos[board.ply]
	index := pos.poskey % uint64(table.count)
	table.RLock()
	entry := table.entries[index]
	table.RUnlock()
	if entry.poskey == pos.poskey {
		return entry.move
	}
	return 0
}

func (board *Board) getPV(depth int) (numPv int, pv [MaxMoves]uint32) {
	for i := 0; i < depth; i++ {
		move := board.probePV()
		if move <= 0 {
			board.ply = 0
			numPv = i
			return
		}
		pv[i] = move
		board.makeMove(move)
	}
	board.ply = 0
	numPv = depth
	return
}

func (board *Board) storeEntry(move uint32, score, depth, flag int) {
	pos := &board.pos[board.ply]
	entry := HashEntry{pos.poskey, move, score, depth, flag}
	index := pos.poskey % uint64(table.count)
	if table.entries[index].poskey == 0 {
		table.newWrite++
	} else {
		table.overWrite++
	}
	table.Lock()
	table.entries[index] = entry
	table.Unlock()
}

func (board *Board) probeEntry(alpha, beta, depth int) (hit bool, move uint32, score int) {
	hit = true
	pos := &board.pos[board.ply]
	index := pos.poskey % uint64(table.count)
	table.RLock()
	entry := table.entries[index]
	table.RUnlock()
	if entry.poskey == pos.poskey {
		move = entry.move
		if entry.depth >= depth {
			table.hit++
			score = entry.score
			switch entry.flag {
			case FlagAlpha:
				if score <= alpha {
					score = alpha
					return
				}
			case FlagBeta:
				if score >= beta {
					score = beta
					return
				}
			case FlagExact:
				return
			}
		}
	}
	hit = false
	return
}
