package api

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

func parseUUID(ID string) (uuid.UUID, error) {
	log.Printf("Trying to parse id %s", ID)
	var result uuid.UUID
	result, err := uuid.Parse(ID)
	if err != nil {
		return result, errors.New("Invalid ID")
	}
	return result, nil
}
