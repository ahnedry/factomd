package databaseOverlay

import (
	"fmt"

	"github.com/FactomProject/factomd/common/interfaces"
)

func (db *Overlay) FetchFactoidTransactionByHash(hash interfaces.IHash) (interfaces.ITransaction, error) {
	in, err := db.FetchIncludedIn(hash)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	block, err := db.FetchFBlockByKeyMR(in)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("Block not found, should not happen")
	}
	txs := block.GetTransactions()
	for _, tx := range txs {
		if tx.GetHash().IsSameAs(hash) {
			return tx, nil
		}
		if tx.GetSigHash().IsSameAs(hash) {
			return tx, nil
		}
	}
	return nil, fmt.Errorf("Transaction not found in block, should not happen")
}

func (db *Overlay) FetchECTransactionByHash(hash interfaces.IHash) (interfaces.IECBlockEntry, error) {
	in, err := db.FetchIncludedIn(hash)
	if err != nil {
		return nil, err
	}
	if in == nil {
		return nil, nil
	}
	block, err := db.FetchECBlockByHash(in)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("Block not found, should not happen")
	}
	tx := block.GetEntryByHash(hash)
	if tx != nil {
		return tx, nil
	}
	return nil, fmt.Errorf("Transaction not found in block, should not happen")
}
