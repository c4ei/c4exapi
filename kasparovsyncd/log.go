package main

import (
	"github.com/c4ei/c4exapi/logger"
	"github.com/c4ei/c4exd/util/panics"
)

var (
	log   = logger.Logger("KVSD")
	spawn = panics.GoroutineWrapperFunc(log)
)
