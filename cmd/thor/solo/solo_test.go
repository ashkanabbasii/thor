// Copyright (c) 2024 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package solo

import (
	"context"
	"testing"
	"time"

	"github.com/ashkanabbasii/thor/builtin"
	"github.com/ashkanabbasii/thor/chain"
	"github.com/ashkanabbasii/thor/genesis"
	"github.com/ashkanabbasii/thor/logdb"
	"github.com/ashkanabbasii/thor/muxdb"
	"github.com/ashkanabbasii/thor/state"
	"github.com/ashkanabbasii/thor/thor"
	"github.com/ashkanabbasii/thor/txpool"
	"github.com/stretchr/testify/assert"
)

func newSolo() *Solo {
	db := muxdb.NewMem()
	stater := state.NewStater(db)
	gene := genesis.NewDevnet()
	logDb, _ := logdb.NewMem()
	b, _, _, _ := gene.Build(stater)
	repo, _ := chain.NewRepository(db, b)
	mempool := txpool.New(repo, stater, txpool.Options{Limit: 10000, LimitPerAccount: 16, MaxLifetime: 10 * time.Minute})

	return New(repo, stater, logDb, mempool, 0, true, false, thor.BlockInterval, thor.ForkConfig{})
}

func TestInitSolo(t *testing.T) {
	solo := newSolo()

	// init solo -> this should mine a block with the gas price tx
	err := solo.init(context.Background())
	assert.Nil(t, err)

	// check the gas price
	best := solo.repo.BestBlockSummary()
	newState := solo.stater.NewState(best.Header.StateRoot(), best.Header.Number(), best.Conflicts, best.SteadyNum)
	currentBGP, err := builtin.Params.Native(newState).Get(thor.KeyBaseGasPrice)
	assert.Nil(t, err)
	assert.Equal(t, baseGasPrice, currentBGP)
}
