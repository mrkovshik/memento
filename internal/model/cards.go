package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrkovshik/memento/internal/crypto"
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

func (c *CardData) Encrypt(passphrase string) error {
	var err error
	c.CVV, err = crypto.EncryptString(c.CVV, passphrase)
	if err != nil {
		return err
	}
	c.Number, err = crypto.EncryptString(c.Number, passphrase)
	if err != nil {
		return err
	}
	c.Name, err = crypto.EncryptString(c.Name, passphrase)
	if err != nil {
		return err
	}
	c.Expiry, err = crypto.EncryptString(c.Expiry, passphrase)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardData) Decrypt(passphrase string) error {
	var err error
	c.CVV, err = crypto.DecryptString(c.CVV, passphrase)
	if err != nil {
		return err
	}
	c.Number, err = crypto.DecryptString(c.Number, passphrase)
	if err != nil {
		return err
	}
	c.Name, err = crypto.DecryptString(c.Name, passphrase)
	if err != nil {
		return err
	}
	c.Expiry, err = crypto.DecryptString(c.Expiry, passphrase)
	if err != nil {
		return err
	}
	return nil
}
