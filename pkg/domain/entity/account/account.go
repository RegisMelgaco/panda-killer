package account

import "time"

type Account struct {
	ID                int
	Balance           float64
	Name, CPF, Secret string
	CreatedAt         time.Time
}
