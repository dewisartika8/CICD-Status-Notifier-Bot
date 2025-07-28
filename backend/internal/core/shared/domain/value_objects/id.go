package value_objects

import (
	"github.com/google/uuid"
)

// ID represents a universal identifier value object
type ID struct {
	value uuid.UUID
}

// NewID creates a new ID with generated UUID
func NewID() ID {
	return ID{value: uuid.New()}
}

// NewIDFromString creates an ID from string
func NewIDFromString(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID{value: id}, nil
}

// NewIDFromUUID creates an ID from existing UUID
func NewIDFromUUID(id uuid.UUID) ID {
	return ID{value: id}
}

// Value returns the underlying UUID
func (id ID) Value() uuid.UUID {
	return id.value
}

// String returns string representation
func (id ID) String() string {
	return id.value.String()
}

// IsNil checks if ID is nil
func (id ID) IsNil() bool {
	return id.value == uuid.Nil
}

// Equals compares two IDs
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}
