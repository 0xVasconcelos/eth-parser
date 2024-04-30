package parser

import (
	"context"
	"math/big"
)

func (p *Parser) index(ctx context.Context) {
	defer func() {
		<-p.sem
	}()

	storageLastBlock, err := p.s.GetLastBlock()
	if err != nil {
		p.log.Println("Error getting last block:", err)
		return
	}

	p.log.Println("Storage last block:", storageLastBlock)
	networkLastBlock, err := p.e.GetCurrentBlock(ctx)
	if err != nil {
		p.log.Println("Error getting current block:", err)
		return
	}

	p.log.Println("Network last block:", networkLastBlock)
	if storageLastBlock.Cmp(networkLastBlock) == 0 {
		p.log.Println("No new blocks")
		return
	}

	if storageLastBlock.Cmp(big.NewInt(0)) == 0 {
		p.log.Println("No blocks in storage, set last block to network last block")
		err = p.indexBlocks(ctx, networkLastBlock)
		if err != nil {
			p.log.Println("Error indexing block:", err)
			return
		}
		p.s.SetLastBlock(networkLastBlock)
		return
	}

	for i := storageLastBlock; i.Cmp(networkLastBlock) < 0; i.Add(i, big.NewInt(1)) {
		err = p.indexBlocks(ctx, i)
		if err != nil {
			p.log.Println("Error indexing block:", err)
			continue
		}
		err = p.s.SetLastBlock(i)
		if err != nil {
			p.log.Println("Error setting last block:", err)
			continue
		}

	}

}

func (p *Parser) indexBlocks(ctx context.Context, blockNumber *big.Int) error {
	p.log.Println("Indexing block:", blockNumber)
	block, err := p.e.GetBlockByNumber(ctx, blockNumber, true)
	if err != nil {
		return err
	}

	for _, tx := range block.Transactions {
		ok, err := p.s.IsSubscribed(tx.From)
		if err != nil {
			return err
		}
		if ok {
			p.log.Println("Indexing transaction from:", tx.From)
			nTx := p.FromEthereumTransaction(tx)
			err := p.s.AddTransaction(tx.From, nTx)
			if err != nil {
				return err
			}
		}
		ok, err = p.s.IsSubscribed(tx.To)
		if err != nil {
			return err
		}
		if ok {
			p.log.Println("Indexing transaction to:", tx.To)
			nTx := p.FromEthereumTransaction(tx)
			err := p.s.AddTransaction(tx.To, nTx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
