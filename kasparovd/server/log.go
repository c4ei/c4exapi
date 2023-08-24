package server

import (
	"github.com/c4ei/c4exapi/logger"
	"github.com/c4ei/c4exd/util/panics"
)

var (
	log   = logger.Logger("REST")
	spawn = panics.GoroutineWrapperFunc(log)
)
