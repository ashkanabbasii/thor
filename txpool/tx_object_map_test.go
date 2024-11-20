// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package txpool

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ashkanabbasii/thor/genesis"
	"github.com/ashkanabbasii/thor/muxdb"
	"github.com/ashkanabbasii/thor/state"
	"github.com/ashkanabbasii/thor/thor"
	"github.com/ashkanabbasii/thor/tx"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	db := muxdb.NewMem()
	repo := newChainRepo(db)

	// Creating transactions
	tx1 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[0])
	tx2 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[1])

	// Resolving transactions into txObjects
	txObj1, _ := resolveTx(tx1, false)
	txObj2, _ := resolveTx(tx2, false)

	// Creating a new txObjectMap and adding transactions
	m := newTxObjectMap()
	assert.Nil(t, m.Add(txObj1, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
	assert.Nil(t, m.Add(txObj2, 1, func(_ thor.Address, _ *big.Int) error { return nil }))

	// Testing GetByID
	retrievedTxObj1 := m.GetByID(txObj1.ID())
	assert.Equal(t, txObj1, retrievedTxObj1, "The retrieved transaction object should match the original for tx1")

	retrievedTxObj2 := m.GetByID(txObj2.ID())
	assert.Equal(t, txObj2, retrievedTxObj2, "The retrieved transaction object should match the original for tx2")

	// Testing retrieval of a non-existing transaction
	nonExistingTxID := thor.Bytes32{} // An arbitrary non-existing ID
	retrievedTxObj3 := m.GetByID(nonExistingTxID)
	assert.Nil(t, retrievedTxObj3, "Retrieving a non-existing transaction should return nil")
}

func TestFill(t *testing.T) {
	db := muxdb.NewMem()
	repo := newChainRepo(db)

	// Creating transactions
	tx1 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[0])
	tx2 := newDelegatedTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, genesis.DevAccounts()[1], genesis.DevAccounts()[2])

	// Resolving transactions into txObjects
	txObj1, _ := resolveTx(tx1, false)
	txObj2, _ := resolveTx(tx2, false)

	// Creating a new txObjectMap
	m := newTxObjectMap()

	// Filling the map with transactions
	m.Fill([]*txObject{txObj1, txObj2, txObj1})

	// Asserting the length of the map
	assert.Equal(t, 2, m.Len(), "Map should contain only 2 unique transactions")

	// Asserting the transactions are correctly added
	assert.True(t, m.ContainsHash(txObj1.Hash()), "Map should contain txObj1")
	assert.True(t, m.ContainsHash(txObj2.Hash()), "Map should contain txObj2")

	// Asserting duplicate handling
	assert.Equal(t, m.GetByID(txObj1.ID()), txObj1, "Duplicate tx1 should not be added again")
	assert.Equal(t, m.GetByID(txObj2.ID()), txObj2, "txObj2 should be retrievable by ID")

	assert.Equal(t, 1, m.quota[genesis.DevAccounts()[0].Address], "Account quota should be 1 for account 0")
	assert.Equal(t, 1, m.quota[genesis.DevAccounts()[1].Address], "Account quota should be 1 for account 1")
	assert.Equal(t, 1, m.quota[genesis.DevAccounts()[2].Address], "Delegator quota should be 1 for account 2")
}

func TestTxObjMap(t *testing.T) {
	db := muxdb.NewMem()
	repo := newChainRepo(db)

	tx1 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[0])
	tx2 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[0])
	tx3 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[1])

	txObj1, _ := resolveTx(tx1, false)
	txObj2, _ := resolveTx(tx2, false)
	txObj3, _ := resolveTx(tx3, false)

	m := newTxObjectMap()
	assert.Zero(t, m.Len())

	assert.Nil(t, m.Add(txObj1, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
	assert.Nil(t, m.Add(txObj1, 1, func(_ thor.Address, _ *big.Int) error { return nil }), "should no error if exists")
	assert.Equal(t, 1, m.Len())

	assert.Equal(t, errors.New("account quota exceeded"), m.Add(txObj2, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
	assert.Equal(t, 1, m.Len())

	assert.Nil(t, m.Add(txObj3, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
	assert.Equal(t, 2, m.Len())

	assert.True(t, m.ContainsHash(tx1.Hash()))
	assert.False(t, m.ContainsHash(tx2.Hash()))
	assert.True(t, m.ContainsHash(tx3.Hash()))

	assert.True(t, m.RemoveByHash(tx1.Hash()))
	assert.False(t, m.ContainsHash(tx1.Hash()))
	assert.False(t, m.RemoveByHash(tx2.Hash()))

	assert.Equal(t, []*txObject{txObj3}, m.ToTxObjects())
	assert.Equal(t, tx.Transactions{tx3}, m.ToTxs())
}

func TestLimitByDelegator(t *testing.T) {
	db := muxdb.NewMem()
	repo := newChainRepo(db)

	tx1 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[0])
	tx2 := newDelegatedTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, genesis.DevAccounts()[0], genesis.DevAccounts()[1])
	tx3 := newDelegatedTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, genesis.DevAccounts()[2], genesis.DevAccounts()[1])

	txObj1, _ := resolveTx(tx1, false)
	txObj2, _ := resolveTx(tx2, false)
	txObj3, _ := resolveTx(tx3, false)

	m := newTxObjectMap()
	assert.Nil(t, m.Add(txObj1, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
	assert.Nil(t, m.Add(txObj3, 1, func(_ thor.Address, _ *big.Int) error { return nil }))

	m = newTxObjectMap()
	assert.Nil(t, m.Add(txObj2, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
	assert.Equal(t, errors.New("delegator quota exceeded"), m.Add(txObj3, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
	assert.Equal(t, errors.New("account quota exceeded"), m.Add(txObj1, 1, func(_ thor.Address, _ *big.Int) error { return nil }))
}

func TestPendingCost(t *testing.T) {
	db := muxdb.NewMem()
	repo := newChainRepo(db)
	stater := state.NewStater(db)

	// Creating transactions
	tx1 := newTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, tx.Features(0), genesis.DevAccounts()[0])
	tx2 := newDelegatedTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, genesis.DevAccounts()[1], genesis.DevAccounts()[2])
	tx3 := newDelegatedTx(repo.ChainTag(), nil, 21000, tx.BlockRef{}, 100, nil, genesis.DevAccounts()[1], genesis.DevAccounts()[2])

	// Resolving transactions into txObjects
	txObj1, _ := resolveTx(tx1, false)
	txObj2, _ := resolveTx(tx2, false)
	txObj3, _ := resolveTx(tx3, false)

	chain := repo.NewBestChain()
	best := repo.BestBlockSummary()
	state := stater.NewState(best.Header.StateRoot(), best.Header.Number(), best.Conflicts, best.SteadyNum)

	var err error
	txObj1.executable, err = txObj1.Executable(chain, state, best.Header)
	assert.Nil(t, err)
	assert.True(t, txObj1.executable)

	txObj2.executable, err = txObj2.Executable(chain, state, best.Header)
	assert.Nil(t, err)
	assert.True(t, txObj2.executable)

	txObj3.executable, err = txObj3.Executable(chain, state, best.Header)
	assert.Nil(t, err)
	assert.True(t, txObj3.executable)

	// Creating a new txObjectMap
	m := newTxObjectMap()

	m.Add(txObj1, 10, func(_ thor.Address, _ *big.Int) error { return nil })
	m.Add(txObj2, 10, func(_ thor.Address, _ *big.Int) error { return nil })
	m.Add(txObj3, 10, func(_ thor.Address, _ *big.Int) error { return nil })

	assert.Equal(t, txObj1.Cost(), m.cost[genesis.DevAccounts()[0].Address])
	// No cost for txObj2's origin, should be counted on the delegator
	assert.Nil(t, m.cost[genesis.DevAccounts()[1].Address])
	assert.Equal(t, new(big.Int).Add(txObj2.Cost(), txObj3.Cost()), m.cost[genesis.DevAccounts()[2].Address])

	m.RemoveByHash(txObj1.Hash())
	assert.Nil(t, m.cost[genesis.DevAccounts()[0].Address])
	m.RemoveByHash(txObj2.Hash())
	assert.Equal(t, txObj3.Cost(), m.cost[genesis.DevAccounts()[2].Address])
	m.RemoveByHash(txObj2.Hash())
	assert.Equal(t, txObj3.Cost(), m.cost[genesis.DevAccounts()[2].Address])
	m.RemoveByHash(txObj3.Hash())
	assert.Nil(t, m.cost[genesis.DevAccounts()[2].Address])
}
