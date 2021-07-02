package transfer

import "errors"

var (
	ErrInsufficientFundsToMakeTransaction     error = errors.New("Insufficient funds to execute transaction.")
	ErrTransferAmountShouldBeGreatterThanZero error = errors.New("Transfer amount should be greatter than zero.")
)
