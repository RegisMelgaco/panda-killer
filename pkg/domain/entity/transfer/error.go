package transfer

import "errors"

var (
	ErrInsufficientFundsToMakeTransaction           error = errors.New("insufficient funds to execute transaction")
	ErrTransferAmountShouldBeGreatterThanZero       error = errors.New("transfer amount should be greatter than zero")
	ErrTransferOriginAndDestinationNeedToBeDiffrent error = errors.New("transfer origin and destination need to be diffrent")
)
