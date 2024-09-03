package cards

import (
	"errors"

	"github.com/mrkovshik/memento/internal/crypto"
	"github.com/mrkovshik/memento/internal/validation"
)

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

func (c *CardData) Validate() error {
	if !validation.ValidateCardNumber(c.Number) {
		return errors.New("invalid card number")
	}
	if err := validation.ValidateExpirationDate(c.Expiry); err != nil {
		return err
	}
	if !validation.ValidateCVV(c.CVV) {
		return errors.New("invalid CVV")
	}
	return nil
}
