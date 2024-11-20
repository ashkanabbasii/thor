// Copyright (c) 2022 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>
package bft

import (
	"bytes"
	"sort"

	"github.com/ashkanabbasii/thor/block"
	"github.com/ashkanabbasii/thor/thor"
)

// casts stores the master's overall casts, maintaining the map of quality to checkpoint.
type casts map[thor.Bytes32]uint32

func (engine *Engine) newCasts() error {
	return nil
}

// Slice dumps the casts that is after finalized into slice.
func (ca casts) Slice(finalized thor.Bytes32) []struct {
	checkpoint thor.Bytes32
	quality    uint32
} {
	list := make([]struct {
		checkpoint thor.Bytes32
		quality    uint32
	}, 0, len(ca))

	for checkpoint, quality := range ca {
		if block.Number(checkpoint) >= block.Number(finalized) {
			list = append(list, struct {
				checkpoint thor.Bytes32
				quality    uint32
			}{checkpoint: checkpoint, quality: quality})
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return bytes.Compare(list[i].checkpoint.Bytes(), list[j].checkpoint.Bytes()) > 0
	})

	return list
}

// Mark marks the master's cast.
func (ca casts) Mark(checkpoint thor.Bytes32, quality uint32) {
	ca[checkpoint] = quality
}
