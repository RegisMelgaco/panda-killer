package transfer

import "errors"

var ErrInsufficientFundsToMakeTransaction error = errors.New("Insufficient funds to execute transaction.")
