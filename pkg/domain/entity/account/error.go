package account

import "errors"

var ErrAccountNameIsObligatory error = errors.New("Account name is Obligatory")
var ErrAccountCPFShouldHaveLength11 error = errors.New("Account cpf should have length 11")
var ErrAccountNotFound error = errors.New("Account not found")
