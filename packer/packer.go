// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package packer

import (
	"github.com/ashkanabbasii/thor/block"
	"github.com/ashkanabbasii/thor/chain"
	"github.com/ashkanabbasii/thor/poa"
	"github.com/ashkanabbasii/thor/thor"
)

// Packer to pack txs and build new blocks.
type Packer struct {
	repo           *chain.Repository
	nodeMaster     thor.Address
	beneficiary    *thor.Address
	targetGasLimit uint64
	forkConfig     thor.ForkConfig
	seeder         *poa.Seeder
}

// New create a new Packer instance.
// The beneficiary is optional, it defaults to endorsor if not set.
//func New(
//	repo *chain.Repository,
//	nodeMaster thor.Address,
//	beneficiary *thor.Address,
//	forkConfig thor.ForkConfig,
//) *Packer {
//	return &Packer{
//		repo,
//		stater,
//		nodeMaster,
//		beneficiary,
//		0,
//		forkConfig,
//		poa.NewSeeder(repo),
//	}
//}

// Schedule schedule a packing flow to pack new block upon given parent and clock time.
func (p *Packer) Schedule(parent *chain.BlockSummary, nowTimestamp uint64) (flow *Flow, err error) {
	return nil, nil
}

// Mock create a packing flow upon given parent, but with a designated timestamp.
// It will skip the PoA verification and scheduling, and the block produced by
// the returned flow is not in consensus.
func (p *Packer) Mock(parent *chain.BlockSummary, targetTime uint64, gasLimit uint64) (*Flow, error) {
	return nil, nil
}

func (p *Packer) gasLimit(parentGasLimit uint64) uint64 {
	if p.targetGasLimit != 0 {
		return block.GasLimit(p.targetGasLimit).Qualify(parentGasLimit)
	}
	return parentGasLimit
}

// SetTargetGasLimit set target gas limit, the Packer will adjust block gas limit close to
// it as it can.
func (p *Packer) SetTargetGasLimit(gl uint64) {
	p.targetGasLimit = gl
}
