package sync

import (
	"github.com/c4ei/c4exapi/database"
	"github.com/c4ei/c4exapi/dbaccess"
	"github.com/c4ei/c4exapi/dbmodels"

	"github.com/pkg/errors"
)

func insertTransactionBlocks(dbTx *database.TxContext, blocks []*rawAndVerboseBlock,
	blockHashesToIDs map[string]uint64, transactionHashesToTxsWithMetadata map[string]*txWithMetadata) error {

	transactionBlocksToAdd := make([]interface{}, 0)
	for _, block := range blocks {
		blockID, ok := blockHashesToIDs[block.hash()]
		if !ok {
			return errors.Errorf("couldn't find block ID for block %s", block)
		}
		for i, tx := range block.Verbose.RawTx {
			transactionBlocksToAdd = append(transactionBlocksToAdd, &dbmodels.TransactionBlock{
				TransactionID: transactionHashesToTxsWithMetadata[tx.Hash].id,
				BlockID:       blockID,
				Index:         uint32(i),
			})
		}
	}
	return dbaccess.BulkInsert(dbTx, transactionBlocksToAdd)
}
