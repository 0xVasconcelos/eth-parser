package ethereum

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type Block struct {
	Number       BigInt        `json:"number"`
	Hash         string        `json:"hash"`
	Timestamp    BigInt        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	BlockNumber BigInt `json:"blockNumber"`
	From        string `json:"from"`
	Hash        string `json:"hash"`
	To          string `json:"to"`
	Value       BigInt `json:"value"`
	Gas         BigInt `json:"gas"`
}

type BigInt struct {
	*big.Int
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	b.Int = new(big.Int)
	if _, ok := b.Int.SetString(str[2:], 16); !ok {
		return fmt.Errorf("invalid big integer: %s", str)
	}

	return nil
}

type GetByBlockNumberParams struct {
	BlockNumber         *big.Int
	IncludeTransactions bool
}

func (g *GetByBlockNumberParams) MarshalJSON() ([]byte, error) {
	params := []interface{}{fmt.Sprintf("0x%s", g.BlockNumber.Text(16)), g.IncludeTransactions}
	return json.Marshal(params)
}

type GetByBlockNumberResponse struct {
	RPCResponse
	Result Block `json:"result"`
}

type GetCurrentBlockResponse struct {
	RPCResponse
	Result string `json:"result"`
}
