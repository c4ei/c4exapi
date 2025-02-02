package sync

import (
	"github.com/c4ei/c4exapi/database"
	"github.com/c4ei/c4exapi/dbaccess"
	"github.com/c4ei/c4exapi/dbmodels"
	rpcmodel "github.com/c4ei/c4exd/rpc/model"

	"github.com/pkg/errors"
)

func insertBlockParents(dbTx *database.TxContext, blocks []*rawAndVerboseBlock, blockHashesToIDs map[string]uint64) error {
	parentsToAdd := make([]interface{}, 0)
	for _, block := range blocks {
		dbBlockParents, err := dbParentBlocksFromVerboseBlock(blockHashesToIDs, block.Verbose)
		if err != nil {
			return err
		}
		for _, dbBlockParent := range dbBlockParents {
			parentsToAdd = append(parentsToAdd, dbBlockParent)
		}
	}
	err := dbaccess.BulkInsert(dbTx, parentsToAdd)
	if err != nil {
		return err
	}
	return nil
}

func dbParentBlocksFromVerboseBlock(blockHashesToIDs map[string]uint64, verboseBlock *rpcmodel.GetBlockVerboseResult) ([]*dbmodels.ParentBlock, error) {
	// Exit early if this is the genesis block
	if len(verboseBlock.ParentHashes) == 0 {
		return nil, nil
	}

	blockID, ok := blockHashesToIDs[verboseBlock.Hash]
	if !ok {
		return nil, errors.Errorf("couldn't find block ID for block %s", verboseBlock.Hash)
	}
	dbParentBlocks := make([]*dbmodels.ParentBlock, len(verboseBlock.ParentHashes))
	for i, parentHash := range verboseBlock.ParentHashes {
		parentID, ok := blockHashesToIDs[parentHash]
		if !ok {
			return nil, errors.Errorf("missing parent %s for block %s", parentHash, verboseBlock.Hash)
		}
		dbParentBlocks[i] = &dbmodels.ParentBlock{
			BlockID:       blockID,
			ParentBlockID: parentID,
		}
	}
	return dbParentBlocks, nil
}
