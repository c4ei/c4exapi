package jsonrpc

import (
	"github.com/c4ei/c4exapi/logger"
	"github.com/c4ei/c4exd/logs"
	rpcclient "github.com/c4ei/c4exd/rpc/client"
)

func init() {
	rpcclient.UseLogger(logger.BackendLog, logs.LevelInfo)
}
