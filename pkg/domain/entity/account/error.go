package account

import "errors"

// Is used when the account name is empty.
var ErrAccountNameIsObligatory error = errors.New("Account name is Obligatory")

// Is used when the account cpf doesn't have a length equals to 11.
var ErrAccountCPFShouldHaveLength11 error = errors.New("Account cpf should have length 11")

// Is used when trying to use a non existing account.
var ErrAccountNotFound error = errors.New("Account not found")
