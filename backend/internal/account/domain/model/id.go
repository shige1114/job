package model

import "github.com/google/uuid"

// ID は申請の一意識別子です
type ID struct {
	value uuid.UUID
}

// NewID は UUID から ID を生成します
func NewID(v uuid.UUID) ID {
	return ID{value: v}
}

// ParseID は文字列から ID を再構成します
func ParseID(v string) (ID, error) {
	uid, err := uuid.Parse(v)
	if err != nil {
		return ID{}, err
	}
	return ID{value: uid}, nil
}

func (i ID) String() string {
	return i.value.String()
}

func (i ID) UUID() uuid.UUID {
	return i.value
}
