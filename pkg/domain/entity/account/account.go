package account

import "time"

type Account struct {
	ID                int
	Balance           int
	Name, CPF, Secret string
	CreatedAt         time.Time
}
