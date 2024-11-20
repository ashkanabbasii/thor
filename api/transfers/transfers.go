// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package transfers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ashkanabbasii/thor/api/utils"
	"github.com/ashkanabbasii/thor/chain"
	"github.com/ashkanabbasii/thor/logdb"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Transfers struct {
	repo  *chain.Repository
	db    *logdb.LogDB
	limit uint64
}

func New(repo *chain.Repository, db *logdb.LogDB, logsLimit uint64) *Transfers {
	return &Transfers{
		repo,
		db,
		logsLimit,
	}
}

// Filter query logs with option
func (t *Transfers) filter(ctx context.Context, filter *TransferFilter) ([]*FilteredTransfer, error) {
	return nil, nil
}

func (t *Transfers) handleFilterTransferLogs(w http.ResponseWriter, req *http.Request) error {
	var filter TransferFilter
	if err := utils.ParseJSON(req.Body, &filter); err != nil {
		return utils.BadRequest(errors.WithMessage(err, "body"))
	}
	if filter.Options != nil && filter.Options.Limit > t.limit {
		return utils.Forbidden(fmt.Errorf("options.limit exceeds the maximum allowed value of %d", t.limit))
	}
	if filter.Options == nil {
		// if filter.Options is nil, set to the default limit +1
		// to detect whether there are more logs than the default limit
		filter.Options = &logdb.Options{
			Offset: 0,
			Limit:  t.limit + 1,
		}
	}

	tLogs, err := t.filter(req.Context(), &filter)
	if err != nil {
		return err
	}

	// ensure the result size is less than the configured limit
	if len(tLogs) > int(t.limit) {
		return utils.Forbidden(fmt.Errorf("the number of filtered logs exceeds the maximum allowed value of %d, please use pagination", t.limit))
	}

	return utils.WriteJSON(w, tLogs)
}

func (t *Transfers) Mount(root *mux.Router, pathPrefix string) {
	sub := root.PathPrefix(pathPrefix).Subrouter()

	sub.Path("").
		Methods(http.MethodPost).
		Name("logs_filter_transfer").
		HandlerFunc(utils.WrapHandlerFunc(t.handleFilterTransferLogs))
}
