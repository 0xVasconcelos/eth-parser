package ethereum

import (
	"context"
	"errors"
	"math/big"
)

// GetBlockByNumber fetches the current block with transactions
func (c *Client) GetBlockByNumber(ctx context.Context, blockNumber *big.Int, includeTxs bool) (*Block, error) {
	params := &GetByBlockNumberParams{
		BlockNumber:         blockNumber,
		IncludeTransactions: includeTxs,
	}

	req := &RPCRequest{
		Jsonrpc: "2.0",
		Method:  MethodGetBlockByNumber,
		Params:  params,
		Id:      1,
	}

	var r GetByBlockNumberResponse
	err := c.post(ctx, req, &r)
	if err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, errors.New(r.Error.Error())
	}

	return &r.Result, nil
}

// GetCurrentBlock gets the current block number
func (c *Client) GetCurrentBlock(ctx context.Context) (*big.Int, error) {
	req := &RPCRequest{
		Jsonrpc: "2.0",
		Method:  MethodBlockNumber,
		Id:      1,
		Params:  nil,
	}

	var b GetCurrentBlockResponse

	err := c.post(ctx, req, &b)
	if err != nil {
		return nil, err
	}

	var n big.Int
	_, ok := n.SetString(b.Result, 0)
	if !ok {
		return nil, errors.New("failed to parse block number")
	}

	return &n, nil
}
