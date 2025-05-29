package id

import (
	"errors"

	"github.com/google/uuid"
)

type Id uuid.UUID

func NewId() *Id {
	newId := Id(uuid.New())
	return &newId
}

func FromString(from string) (*Id, error) {
	newUUID, err := uuid.Parse(from)
	if err != nil {
		return nil, err
	}

	newId := Id(newUUID)

	return &newId, nil
}

func FromUUID(from uuid.UUID) (*Id, error) {
	if from == uuid.Nil {
		return nil, errors.New("from uuid cannot bem nil")
	}

	newId := Id(from)

	return &newId, nil
}

func (i *Id) UUID() uuid.UUID {
	return uuid.UUID(*i)
}

func (i *Id) String() string {
	return i.UUID().String()
}

func (i *Id) Equal(src Id) bool {
	if i.UUID() == src.UUID() {
		return true
	}

	return false
}
