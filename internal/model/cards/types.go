package cards

import (
	"time"

	"github.com/google/uuid"
)

type CardData struct {
	ID        uint
	UserID    uint `db:"user_id"`
	UUID      uuid.UUID
	Number    string
	Expiry    string
	CVV       string
	Name      string
	Meta      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
