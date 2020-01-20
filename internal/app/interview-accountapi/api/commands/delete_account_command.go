package commands

import "github.com/google/uuid"

type DeleteAccountCommand struct {
	AccountId uuid.UUID
	Version       int64
}
