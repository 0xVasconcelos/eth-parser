package parser

import (
	"context"
	"math/big"
)

func (n *Notifier) index(ctx context.Context) {
	defer func() {
		<-n.sem
	}()

	storageLastBlock, err := n.s.GetLastBlock()
	if err != nil {
		n.log.Println("Error getting last block:", err)
	}

	n.log.Println("Storage last block:", storageLastBlock)
	networkLastBlock, err := n.e.GetCurrentBlock(ctx)
	if err != nil {
		n.log.Println("Error getting current block:", err)
	}

	n.log.Println("Network last block:", networkLastBlock)
	if storageLastBlock.Cmp(networkLastBlock) == 0 {
		n.log.Println("No new blocks")
		return
	}

	if storageLastBlock.Cmp(big.NewInt(0)) == 0 {
		n.log.Println("No blocks in storage, set last block to network last block")
		err = n.indexBlocks(ctx, networkLastBlock)
		if err != nil {
			n.log.Println("Error indexing block:", err)
			return
		}
		n.s.SetLastBlock(networkLastBlock)
		return
	}

	for i := storageLastBlock; i.Cmp(networkLastBlock) < 0; i.Add(i, big.NewInt(1)) {
		err = n.indexBlocks(ctx, i)
		if err != nil {
			n.log.Println("Error indexing block:", err)
			continue
		}
		err = n.s.SetLastBlock(i)
		if err != nil {
			n.log.Println("Error setting last block:", err)
			continue
		}

	}

}

func (n *Notifier) indexBlocks(ctx context.Context, blockNumber *big.Int) error {
	n.log.Println("Indexing block:", blockNumber)
	block, err := n.e.GetBlockByNumber(ctx, blockNumber, true)
	if err != nil {
		return err
	}

	for _, tx := range block.Transactions {
		ok, err := n.s.IsSubscribed(tx.From)
		if err != nil {
			return err
		}
		if ok {
			n.log.Println("Indexing transaction from:", tx.From)
			err := n.s.AddTransaction(tx.From, tx)
			if err != nil {
				return err
			}
		}
		ok, err = n.s.IsSubscribed(tx.To)
		if err != nil {
			return err
		}
		if ok {
			n.log.Println("Indexing transaction to:", tx.To)
			err := n.s.AddTransaction(tx.To, tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
