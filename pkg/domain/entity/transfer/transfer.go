package transfer

import "time"

type Transfer struct {
	ID, AccountOrigin, AccountDestination int
	Amount                                float64
	CreatedAt                             time.Time
}
