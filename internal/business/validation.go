package business

import (
	"errors"

	"github.com/Toflex/directory_v2/pkg/util"
)

func (b *CreateBusinessRequest) Validate() error {
	if b.UserID == "" {
		return errors.New("user not found!")
	}

	if b.Name == "" {
		return errors.New("business name cannot be empty!")
	}

	if b.Category == "" {
		return errors.New("business category is empty!")
	}

	err := util.ValidateEmail(b.Email)
	if err != nil {
		return err
	}

	return nil
}
