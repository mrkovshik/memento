package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/crypto"
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

func (c *Credential) Encrypt(passphrase string) error {
	var err error
	c.Login, err = crypto.EncryptString(c.Login, passphrase)
	if err != nil {
		return err
	}
	c.Password, err = crypto.EncryptString(c.Password, passphrase)
	if err != nil {
		return err
	}
	return nil
}

func (c *Credential) Decrypt(passphrase string) error {
	var err error
	c.Login, err = crypto.DecryptString(c.Login, passphrase)
	if err != nil {
		return err
	}
	c.Password, err = crypto.DecryptString(c.Password, passphrase)
	if err != nil {
		return err
	}
	return nil
}
