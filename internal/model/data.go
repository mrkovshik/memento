package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	DataStatusSuccess    DataStatus = "success"
	DataStatusFailure    DataStatus = "failure"
	DataStatusProcessing DataStatus = "processing"
	DataStatusUndefined  DataStatus = "undefined"
)

type DataStatus string

type VariousData struct {
	ID        uint
	UserID    uint `db:"user_id"`
	UUID      uuid.UUID
	Title     string
	DataType  int `db:"data_type"`
	Status    DataStatus
	Meta      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
