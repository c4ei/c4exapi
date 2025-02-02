package mqtt

import (
	"github.com/c4ei/c4exapi/apimodels"
	rpcmodel "github.com/c4ei/c4exd/rpc/model"
)

// SelectedParentChainTopic is an MQTT topic for changes in the
// selected parent chain
const SelectedParentChainTopic = "dag/selected-parent-chain"

// PublishSelectedParentChainNotifications publishes notifications for changes in the selected parent chain
func PublishSelectedParentChainNotifications(removedChainHashes []string, addedChainBlocks []rpcmodel.ChainBlock) error {
	if !isConnected() {
		return nil
	}

	notificationData := &apimodels.SelectedParentChainNotification{
		AddedChainBlocks:   make([]*apimodels.AddedChainBlock, len(addedChainBlocks)),
		RemovedBlockHashes: make([]string, len(removedChainHashes)),
	}

	for i, block := range addedChainBlocks {
		acceptedBlockHashes := make([]string, len(block.AcceptedBlocks))
		for i, acceptedBlock := range block.AcceptedBlocks {
			acceptedBlockHashes[i] = acceptedBlock.Hash
		}
		notificationData.AddedChainBlocks[i] = &apimodels.AddedChainBlock{
			Hash:                block.Hash,
			AcceptedBlockHashes: acceptedBlockHashes,
		}
	}
	notificationData.RemovedBlockHashes = removedChainHashes

	return publish(SelectedParentChainTopic, notificationData)
}
