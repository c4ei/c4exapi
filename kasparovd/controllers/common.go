package controllers

import (
	"net/http"

	"github.com/c4ei/c4exapi/httpserverutils"
	"github.com/c4ei/c4exapi/kasparovd/config"
	"github.com/c4ei/c4exd/util"
	"github.com/pkg/errors"
)

func validateAddress(address string) error {
	_, err := util.DecodeAddress(address, config.ActiveConfig().ActiveNetParams.Prefix)
	if err != nil {
		return httpserverutils.NewHandlerErrorWithCustomClientMessage(http.StatusUnprocessableEntity,
			errors.Wrap(err, "error decoding address"),
			"The given address is not a well-formatted P2PKH or P2SH address")
	}

	return nil
}
