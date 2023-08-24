package mqtt

import (
	"github.com/c4ei/c4exapi/apimodels"
	"github.com/c4ei/c4exapi/database"
	"github.com/c4ei/c4exapi/dbaccess"
	"github.com/c4ei/c4exapi/dbmodels"
)

// BlocksTopic is an MQTT topic for new blocks
const BlocksTopic = "dag/blocks"

// PublishBlockAddedNotifications publishes notifications for the block
// that was added, and notifications for its transactions.
func PublishBlockAddedNotifications(hash string) error {
	if !isConnected() {
		return nil
	}

	preloadedFields := dbmodels.PrefixFieldNames(dbmodels.BlockFieldNames.Transactions, dbmodels.TransactionRecommendedPreloadedFields)
	preloadedFields = append(preloadedFields, dbmodels.BlockFieldNames.ParentBlocks)

	dbBlock, err := dbaccess.BlockByHash(database.NoTx(), hash, preloadedFields...)
	if err != nil {
		return err
	}

	selectedTipBlueScore, err := dbaccess.SelectedTipBlueScore(database.NoTx())
	if err != nil {
		return err
	}

	err = publish(BlocksTopic, apimodels.ConvertBlockModelToBlockResponse(dbBlock, selectedTipBlueScore))
	if err != nil {
		return err
	}

	return publishTransactionsNotifications(TransactionsTopic, dbBlock.Transactions, selectedTipBlueScore)
}
