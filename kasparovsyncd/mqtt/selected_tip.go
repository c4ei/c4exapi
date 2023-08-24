package mqtt

import (
	"github.com/c4ei/c4exapi/apimodels"
	"github.com/c4ei/c4exapi/database"
	"github.com/c4ei/c4exapi/dbaccess"
	"github.com/c4ei/c4exapi/dbmodels"
)

const (
	// SelectedTipTopic is an MQTT topic for DAG selected tips
	SelectedTipTopic = "dag/selected-tip"
)

// PublishSelectedTipNotification publishes notification for a new selected tip
func PublishSelectedTipNotification(selectedTipHash string) error {
	if !isConnected() {
		return nil
	}
	dbBlock, err := dbaccess.BlockByHash(database.NoTx(), selectedTipHash, dbmodels.BlockRecommendedPreloadedFields...)
	if err != nil {
		return err
	}

	block := apimodels.ConvertBlockModelToBlockResponse(dbBlock, dbBlock.BlueScore)
	return publish(SelectedTipTopic, block)
}
