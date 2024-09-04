package credentials

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID        uint
	UserID    uint `db:"user_id"`
	UUID      uuid.UUID
	Login     string
	Password  string
	Meta      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
