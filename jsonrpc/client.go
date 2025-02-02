package jsonrpc

import (
	"io/ioutil"
	"time"

	"github.com/c4ei/c4exd/domainmessage"

	"github.com/c4ei/c4exapi/config"
	"github.com/pkg/errors"

	"github.com/c4ei/c4exd/util/daghash"

	rpcclient "github.com/c4ei/c4exd/rpc/client"
	"github.com/c4ei/c4exd/util"
)

// Client represents a connection to the JSON-RPC API of a full node
type Client struct {
	*rpcclient.Client
	OnBlockAdded   chan *BlockAddedMsg
	OnChainChanged chan *ChainChangedMsg
}

var client *Client

// GetClient returns an instance of the JSON-RPC client, in case we have an active connection
func GetClient() (*Client, error) {
	if client == nil {
		return nil, errors.New("JSON-RPC is not connected")
	}

	return client, nil
}

// BlockAddedMsg defines the message received in onBlockAdded
type BlockAddedMsg struct {
	ChainHeight uint64
	Header      *domainmessage.BlockHeader
}

// ChainChangedMsg defines the message received in onChainChanged
type ChainChangedMsg struct {
	RemovedChainBlockHashes []*daghash.Hash
	AddedChainBlocks        []*rpcclient.ChainBlock
}

// Close closes the connection to the JSON-RPC API server
func Close() {
	if client == nil {
		return
	}

	client.Disconnect()
	client = nil
}

// Connect initiates a connection to the JSON-RPC API Server
func Connect(cfg *config.KasparovFlags, subscribeToNotifications bool) error {
	var cert []byte
	if !cfg.DisableTLS {
		var err error
		cert, err = ioutil.ReadFile(cfg.RPCCert)
		if err != nil {
			return errors.Errorf("Error reading certificates file: %s", err)
		}
	}

	rpcServerAddress, err := cfg.NetParams().NormalizeRPCServerAddress(cfg.RPCServer)
	if err != nil {
		return err
	}

	connCfg := &rpcclient.ConnConfig{
		Host:           rpcServerAddress,
		Endpoint:       "ws",
		User:           cfg.RPCUser,
		Pass:           cfg.RPCPassword,
		DisableTLS:     cfg.DisableTLS,
		RequestTimeout: time.Minute * 2,
	}

	if !cfg.DisableTLS {
		connCfg.Certificates = cert
	}

	client, err = newClient(connCfg, subscribeToNotifications)
	if err != nil {
		return errors.Errorf("Error connecting to address %s: %s", rpcServerAddress, err)
	}

	return nil
}

func newClient(connCfg *rpcclient.ConnConfig, subscribeToNotifications bool) (*Client, error) {
	client = &Client{
		OnBlockAdded:   make(chan *BlockAddedMsg),
		OnChainChanged: make(chan *ChainChangedMsg),
	}

	var notificationHandlers *rpcclient.NotificationHandlers
	if subscribeToNotifications {
		notificationHandlers = &rpcclient.NotificationHandlers{
			OnFilteredBlockAdded: func(height uint64, header *domainmessage.BlockHeader,
				txs []*util.Tx) {
				client.OnBlockAdded <- &BlockAddedMsg{
					ChainHeight: height,
					Header:      header,
				}
			},
			OnChainChanged: func(removedChainBlockHashes []*daghash.Hash,
				addedChainBlocks []*rpcclient.ChainBlock) {
				client.OnChainChanged <- &ChainChangedMsg{
					RemovedChainBlockHashes: removedChainBlockHashes,
					AddedChainBlocks:        addedChainBlocks,
				}
			},
		}
	}

	var err error
	client.Client, err = rpcclient.New(connCfg, notificationHandlers)
	if err != nil {
		return nil, errors.Errorf("Error connecting to address %s: %s", connCfg.Host, err)
	}

	if subscribeToNotifications {
		if err = client.NotifyBlocks(); err != nil {
			return nil, errors.Errorf("Error while registering client %s for block notifications: %s", client.Host(), err)
		}
		if err = client.NotifyChainChanges(); err != nil {
			return nil, errors.Errorf("Error while registering client %s for chain changes notifications: %s", client.Host(), err)
		}
	}

	return client, nil
}
