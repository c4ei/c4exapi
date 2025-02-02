package sync

import (
	"strconv"
	"time"

	"github.com/c4ei/c4exapi/database"
	"github.com/c4ei/c4exapi/serializer"

	"github.com/c4ei/c4exapi/dbaccess"
	"github.com/c4ei/c4exapi/dbmodels"
	rpcmodel "github.com/c4ei/c4exd/rpc/model"
	"github.com/pkg/errors"
)

func insertBlocks(dbTx *database.TxContext, blocks []*rawAndVerboseBlock, transactionHashesToTxsWithMetadata map[string]*txWithMetadata) error {
	blocksToAdd := make([]interface{}, len(blocks))
	for i, block := range blocks {
		blockMass := uint64(0)
		for _, tx := range block.Verbose.RawTx {
			blockMass += transactionHashesToTxsWithMetadata[tx.Hash].mass
		}
		var err error
		blocksToAdd[i], err = dbBlockFromVerboseBlock(block.Verbose, blockMass)
		if err != nil {
			return err
		}
	}
	return dbaccess.BulkInsert(dbTx, blocksToAdd)
}

func getBlocksWithTheirAcceptedBlocksAndParentIDs(dbTx *database.TxContext, blocks []*rawAndVerboseBlock) (map[string]uint64, error) {
	blockSet := make(map[string]struct{})
	for _, block := range blocks {
		blockSet[block.hash()] = struct{}{}
		for _, parentHash := range block.Verbose.ParentHashes {
			blockSet[parentHash] = struct{}{}
		}
		for _, acceptedBlockHash := range block.Verbose.AcceptedBlockHashes {
			blockSet[acceptedBlockHash] = struct{}{}
		}
	}

	blockHashes := stringsSetToSlice(blockSet)

	dbBlocks, err := dbaccess.BlocksByHashes(dbTx, blockHashes)
	if err != nil {
		return nil, err
	}

	if len(dbBlocks) != len(blockSet) {
		return nil, errors.Errorf("couldn't retrieve all block IDs")
	}

	blockHashesToIDs := make(map[string]uint64)
	for _, dbBlock := range dbBlocks {
		blockHashesToIDs[dbBlock.BlockHash] = dbBlock.ID
	}
	return blockHashesToIDs, nil
}

func dbBlockFromVerboseBlock(verboseBlock *rpcmodel.GetBlockVerboseResult, mass uint64) (*dbmodels.Block, error) {
	bits, err := strconv.ParseUint(verboseBlock.Bits, 16, 32)
	if err != nil {
		return nil, err
	}

	dbBlock := dbmodels.Block{
		BlockHash:            verboseBlock.Hash,
		Version:              verboseBlock.Version,
		HashMerkleRoot:       verboseBlock.HashMerkleRoot,
		AcceptedIDMerkleRoot: verboseBlock.AcceptedIDMerkleRoot,
		UTXOCommitment:       verboseBlock.UTXOCommitment,
		Timestamp:            time.Unix(verboseBlock.Time, 0),
		Bits:                 uint32(bits),
		Nonce:                serializer.Uint64ToBytes(verboseBlock.Nonce),
		BlueScore:            verboseBlock.BlueScore,
		IsChainBlock:         false, // This must be false for updateSelectedParentChain to work properly
		Mass:                 mass,
	}

	// Set genesis block as the initial chain block
	if len(verboseBlock.ParentHashes) == 0 {
		dbBlock.IsChainBlock = true
	}
	return &dbBlock, nil
}
