package parser

import "github.com/0xVasconcelos/ethparser/pkg/ethereum"

type Transaction struct {
	BlockNumber uint64 `json:"blockNumber"`
	From        string `json:"from"`
	Hash        string `json:"hash"`
	To          string `json:"to"`
	Value       uint64 `json:"value"`
	Gas         uint64 `json:"gas"`
}

func (p *Parser) FromEthereumTransaction(ethereumTransaction ethereum.Transaction) Transaction {
	return Transaction{
		BlockNumber: ethereumTransaction.BlockNumber.Uint64(),
		From:        ethereumTransaction.From,
		Hash:        ethereumTransaction.Hash,
		To:          ethereumTransaction.To,
		Value:       ethereumTransaction.Value.Uint64(),
		Gas:         ethereumTransaction.Gas.Uint64(),
	}
}
