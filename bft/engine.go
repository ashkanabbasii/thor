// Copyright (c) 2022 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>
package bft

import (
	"github.com/ashkanabbasii/thor/block"
	"github.com/ashkanabbasii/thor/cache"
	"github.com/ashkanabbasii/thor/chain"
	"github.com/ashkanabbasii/thor/kv"
	"github.com/ashkanabbasii/thor/thor"
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
)

const dataStoreName = "bft.engine"

var finalizedKey = []byte("finalized")

type Committer interface {
	Finalized() thor.Bytes32
	Justified() (thor.Bytes32, error)
}

type justified struct {
	search thor.Bytes32
	value  thor.Bytes32
}

// Engine tracks all votes of blocks, computes the finalized checkpoint.
// Not thread-safe!
type Engine struct {
	//repo *chain.Repository
	data kv.Store
	//stater     *state.Stater
	forkConfig thor.ForkConfig
	master     thor.Address
	casts      casts
	finalized  atomic.Value
	justified  atomic.Value
	caches     struct {
		state     *lru.Cache
		quality   *lru.Cache
		justifier *cache.PrioCache
	}
}

// NewEngine creates a new bft engine.
//func NewEngine(repo *chain.Repository, mainDB *muxdb.MuxDB, forkConfig thor.ForkConfig, master thor.Address) (*Engine, error) {
//	engine := Engine{
//		repo:       repo,
//		//data:       mainDB.NewStore(dataStoreName),
//		//stater:     state.NewStater(mainDB),
//		forkConfig: forkConfig,
//		master:     master,
//	}
//
//	engine.caches.state, _ = lru.New(256)
//	engine.caches.quality, _ = lru.New(16)
//	engine.caches.justifier = cache.NewPrioCache(16)
//
//	// Restore finalized block, if any
//	if val, err := engine.data.Get(finalizedKey); err != nil {
//		if !engine.data.IsNotFound(err) {
//			return nil, err
//		}
//		engine.finalized.Store(engine.repo.GenesisBlock().Header().ID())
//	} else {
//		engine.finalized.Store(thor.BytesToBytes32(val))
//	}
//
//	return &engine, nil
//}

// Finalized returns the finalized checkpoint.
func (engine *Engine) Finalized() thor.Bytes32 {
	return engine.finalized.Load().(thor.Bytes32)
}

// Justified returns the justified checkpoint.
func (engine *Engine) Justified() (thor.Bytes32, error) {
	return thor.Bytes32{}, nil
}

// Accepts checks if the given block is on the same branch of finalized checkpoint.
//func (engine *Engine) Accepts(parentID thor.Bytes32) (bool, error) {
//	finalized := engine.Finalized()
//
//	if block.Number(finalized) != 0 {
//		return engine.repo.NewChain(parentID).HasBlock(finalized)
//	}
//
//	return true, nil
//}

// Select selects between the new block and the current best, return true if new one is better.
//func (engine *Engine) Select(header *block.Header) (bool, error) {
//	newSt, err := engine.computeState(header)
//	if err != nil {
//		return false, err
//	}
//
//	best := engine.repo.BestBlockSummary().Header
//	bestSt, err := engine.computeState(best)
//	if err != nil {
//		return false, err
//	}
//
//	if newSt.Quality != bestSt.Quality {
//		return newSt.Quality > bestSt.Quality, nil
//	}
//
//	return header.BetterThan(best), nil
//}

// CommitBlock commits bft state to storage.
func (engine *Engine) CommitBlock(header *block.Header, isPacking bool) error {
	// save quality and finalized at the end of each round

	return nil
}

// ShouldVote decides if vote COM for a given parent block ID.
// Packer only.
//func (engine *Engine) ShouldVote(parentID thor.Bytes32) (bool, error) {
//	// laze init casts
//	if engine.casts == nil {
//		if err := engine.newCasts(); err != nil {
//			return false, err
//		}
//	}
//
//	// do not vote COM at the first round
//	if absRound := (block.Number(parentID)+1)/thor.CheckpointInterval - engine.forkConfig.FINALITY/thor.CheckpointInterval; absRound == 0 {
//		return false, nil
//	}
//
//	sum, err := engine.repo.GetBlockSummary(parentID)
//	if err != nil {
//		return false, err
//	}
//	st, err := engine.computeState(sum.Header)
//	if err != nil {
//		return false, err
//	}
//	if st.Quality == 0 {
//		return false, nil
//	}
//
//	headQuality := st.Quality
//	finalized := engine.Finalized()
//	chain := engine.repo.NewChain(parentID)
//	// most recent justified checkpoint
//	var recentJC thor.Bytes32
//	if st.Justified {
//		// if justified in this round, use this round's checkpoint
//		checkpoint, err := chain.GetBlockID(getCheckPoint(block.Number(parentID)))
//		if err != nil {
//			return false, err
//		}
//		recentJC = checkpoint
//	} else {
//		// if current round is not justified, find the most recent justified checkpoint
//		prev, err := chain.GetBlockID(getStorePoint(block.Number(parentID) - thor.CheckpointInterval))
//		if err != nil {
//			return false, err
//		}
//		checkpoint, err := engine.findCheckpointByQuality(headQuality, finalized, prev)
//		if err != nil {
//			return false, err
//		}
//		recentJC = checkpoint
//	}
//
//	// see https://github.com/vechain/VIPs/blob/master/vips/VIP-220.md
//	for _, cast := range engine.casts.Slice(finalized) {
//		if cast.quality >= headQuality-1 {
//			x, y := recentJC, cast.checkpoint
//			if block.Number(cast.checkpoint) > block.Number(recentJC) {
//				x, y = cast.checkpoint, recentJC
//			}
//			// checks if the voted checkpoint belongs to the head chain
//			includes, err := engine.repo.NewChain(x).HasBlock(y)
//			if err != nil {
//				return false, err
//			}
//
//			// if one votes a checkpoint was within [headQuality-1, +∞) and conflict with head
//			// should not vote COM
//			if !includes {
//				return false, nil
//			}
//		}
//	}
//
//	return true, nil
//}

// computeState computes the bft state regarding the given block header to the closest checkpoint.
//func (engine *Engine) computeState(header *block.Header) (*bftState, error) {
//	if cached, ok := engine.caches.state.Get(header.ID()); ok {
//		return cached.(*bftState), nil
//	}
//
//	if header.Number() == 0 || header.Number() < engine.forkConfig.FINALITY {
//		return &bftState{}, nil
//	}
//
//	var (
//		js  *justifier
//		end uint32
//	)
//
//	if entry := engine.caches.justifier.Remove(header.ParentID()); !isCheckPoint(header.Number()) && entry != nil {
//		js = (entry.Entry.Value).(*justifier)
//		end = header.Number()
//	} else {
//		// create a new vote set if cache missed or new block is checkpoint
//		var err error
//		js, err = engine.newJustifier(header.ParentID())
//		if err != nil {
//			return nil, errors.Wrap(err, "failed to create vote set")
//		}
//		end = js.checkpoint
//	}
//
//	h := header
//	for {
//		if h.Number() < engine.forkConfig.FINALITY {
//			break
//		}
//
//		signer, _ := h.Signer()
//		js.AddBlock(signer, h.COM())
//
//		if h.Number() <= end {
//			break
//		}
//
//		sum, err := engine.repo.GetBlockSummary(h.ParentID())
//		if err != nil {
//			return nil, err
//		}
//		h = sum.Header
//	}
//
//	st := js.Summarize()
//	engine.caches.state.Add(header.ID(), st)
//	engine.caches.justifier.Set(header.ID(), js, float64(header.Number()))
//	return st, nil
//}

// findCheckpointByQuality finds the first checkpoint reaches the given quality.
// It is caller's responsibility to ensure the epoch that headID belongs to is concluded.
//func (engine *Engine) findCheckpointByQuality(target uint32, finalized, headID thor.Bytes32) (blockID thor.Bytes32, err error) {
//	defer func() {
//		if e := recover(); e != nil {
//			err = e.(error)
//			return
//		}
//	}()
//
//	searchStart := block.Number(finalized)
//	if searchStart == 0 {
//		searchStart = getCheckPoint(engine.forkConfig.FINALITY)
//	}
//
//	c := engine.repo.NewChain(headID)
//	get := func(i int) (uint32, error) {
//		id, err := c.GetBlockID(getStorePoint(searchStart + uint32(i)*thor.CheckpointInterval))
//		if err != nil {
//			return 0, err
//		}
//		return engine.getQuality(id)
//	}
//
//	// sort.Search searches from [0, n)
//	n := int((block.Number(headID)-searchStart)/thor.CheckpointInterval) + 1
//	num := sort.Search(n, func(i int) bool {
//		quality, err := get(i)
//		if err != nil {
//			panic(err)
//		}
//
//		return quality >= target
//	})
//
//	// n means not found for sort.Search
//	if num == n {
//		return thor.Bytes32{}, errors.New("failed find the block by quality")
//	}
//
//	quality, err := get(num)
//	if err != nil {
//		return thor.Bytes32{}, err
//	}
//
//	if quality != target {
//		return thor.Bytes32{}, errors.New("failed to find the block by quality")
//	}
//
//	return c.GetBlockID(searchStart + uint32(num)*thor.CheckpointInterval)
//}

func (engine *Engine) getMaxBlockProposers(sum *chain.BlockSummary) (uint64, error) {
	//state := engine.stater.NewState(sum.Header.StateRoot(), sum.Header.Number(), sum.Conflicts, sum.SteadyNum)
	return 2, nil
}

func (engine *Engine) getQuality(id thor.Bytes32) (quality uint32, err error) {
	if cached, ok := engine.caches.quality.Get(id); ok {
		return cached.(uint32), nil
	}

	defer func() {
		if err == nil {
			engine.caches.quality.Add(id, quality)
		}
	}()

	quality, err = loadQuality(engine.data, id)
	// no quality saved yet
	if engine.data.IsNotFound(err) {
		return 0, nil
	}
	return
}

func getCheckPoint(blockNum uint32) uint32 {
	return blockNum / thor.CheckpointInterval * thor.CheckpointInterval
}

func isCheckPoint(blockNum uint32) bool {
	return getCheckPoint(blockNum) == blockNum
}

// save quality at the end of round
func getStorePoint(blockNum uint32) uint32 {
	return getCheckPoint(blockNum) + thor.CheckpointInterval - 1
}
