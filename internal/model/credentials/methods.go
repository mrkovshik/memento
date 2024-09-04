package credentials

import "github.com/mrkovshik/memento/internal/crypto"

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
