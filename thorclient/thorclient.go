// Copyright (c) 2024 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

// Package thorclient provides a client for interacting with the VeChainThor blockchain.
// It offers a set of methods to interact with accounts, transactions, blocks, events, and other
// features via HTTP and WebSocket connections.

package thorclient

import (
	"fmt"

	"github.com/ashkanabbasii/thor/api/blocks"
	"github.com/ashkanabbasii/thor/api/events"
	"github.com/ashkanabbasii/thor/api/transactions"
	"github.com/ashkanabbasii/thor/api/transfers"
	"github.com/ashkanabbasii/thor/thor"
	"github.com/ashkanabbasii/thor/thorclient/httpclient"
	"github.com/ashkanabbasii/thor/tx"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"

	tccommon "github.com/ashkanabbasii/thor/thorclient/common"
)

// Client represents the VeChainThor client, allowing communication over HTTP and WebSocket.
type Client struct {
	httpConn *httpclient.Client
}

// New creates a new Client using the provided HTTP URL.
func New(url string) *Client {
	return &Client{
		httpConn: httpclient.New(url),
	}
}

// Option represents a functional option for customizing client requests.
type Option func(*getOptions)

// getOptions holds configuration options for client requests.
type getOptions struct {
	revision string
	pending  bool
}

// applyOptions applies the given functional options to the default options.
func applyOptions(opts []Option) *getOptions {
	options := &getOptions{
		revision: tccommon.BestRevision,
		pending:  false,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// Revision returns an Option to specify the revision for requests.
func Revision(revision string) Option {
	return func(o *getOptions) {
		o.revision = revision
	}
}

// Pending returns an Option to specify that the client should fetch pending results.
func Pending() Option {
	return func(o *getOptions) {
		o.pending = true
	}
}

// RawHTTPClient returns the underlying HTTP client.
func (c *Client) RawHTTPClient() *httpclient.Client {
	return c.httpConn
}

// Transaction retrieves a transaction by its ID.
func (c *Client) Transaction(id *thor.Bytes32, opts ...Option) (*transactions.Transaction, error) {
	options := applyOptions(opts)
	return c.httpConn.GetTransaction(id, options.revision, options.pending)
}

// RawTransaction retrieves the raw transaction data by its ID.
func (c *Client) RawTransaction(id *thor.Bytes32, opts ...Option) (*transactions.RawTransaction, error) {
	options := applyOptions(opts)
	return c.httpConn.GetRawTransaction(id, options.revision, options.pending)
}

// TransactionReceipt retrieves the receipt for a transaction by its ID.
func (c *Client) TransactionReceipt(id *thor.Bytes32, opts ...Option) (*transactions.Receipt, error) {
	options := applyOptions(opts)
	return c.httpConn.GetTransactionReceipt(id, options.revision)
}

// SendTransaction sends a signed transaction to the blockchain.
func (c *Client) SendTransaction(tx *tx.Transaction) (*transactions.SendTxResult, error) {
	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, fmt.Errorf("unable to encode transaction - %w", err)
	}

	return c.SendRawTransaction(rlpTx)
}

// SendRawTransaction sends a raw RLP-encoded transaction to the blockchain.
func (c *Client) SendRawTransaction(rlpTx []byte) (*transactions.SendTxResult, error) {
	return c.httpConn.SendTransaction(&transactions.RawTx{Raw: hexutil.Encode(rlpTx)})
}

// Block retrieves a block by its revision.
func (c *Client) Block(revision string) (blocks *blocks.JSONCollapsedBlock, err error) {
	return c.httpConn.GetBlock(revision)
}

// ExpandedBlock retrieves an expanded block by its revision.
func (c *Client) ExpandedBlock(revision string) (blocks *blocks.JSONExpandedBlock, err error) {
	return c.httpConn.GetExpandedBlock(revision)
}

// FilterEvents filters events based on the provided filter request.
func (c *Client) FilterEvents(req *events.EventFilter) ([]events.FilteredEvent, error) {
	return c.httpConn.FilterEvents(req)
}

// FilterTransfers filters transfers based on the provided filter request.
func (c *Client) FilterTransfers(req *transfers.TransferFilter) ([]*transfers.FilteredTransfer, error) {
	return c.httpConn.FilterTransfers(req)
}

// ChainTag retrieves the chain tag from the genesis block.
func (c *Client) ChainTag() (byte, error) {
	genesisBlock, err := c.Block("0")
	if err != nil {
		return 0, err
	}
	return genesisBlock.ID[31], nil
}
