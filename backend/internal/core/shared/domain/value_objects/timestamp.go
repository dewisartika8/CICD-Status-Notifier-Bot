package value_objects

import (
	"time"
)

// Timestamp represents a time value object
type Timestamp struct {
	value time.Time
}

// NewTimestamp creates a new timestamp with current time
func NewTimestamp() Timestamp {
	return Timestamp{value: time.Now()}
}

// NewTimestampFromTime creates a timestamp from existing time
func NewTimestampFromTime(t time.Time) Timestamp {
	return Timestamp{value: t}
}

// Value returns the underlying time
func (ts Timestamp) Value() time.Time {
	return ts.value
}

// ToTime returns the underlying time (alias for Value)
func (ts Timestamp) ToTime() time.Time {
	return ts.value
}

// String returns string representation
func (ts Timestamp) String() string {
	return ts.value.String()
}

// Unix returns unix timestamp
func (ts Timestamp) Unix() int64 {
	return ts.value.Unix()
}

// IsZero checks if timestamp is zero
func (ts Timestamp) IsZero() bool {
	return ts.value.IsZero()
}

// Before checks if timestamp is before another
func (ts Timestamp) Before(other Timestamp) bool {
	return ts.value.Before(other.value)
}

// After checks if timestamp is after another
func (ts Timestamp) After(other Timestamp) bool {
	return ts.value.After(other.value)
}
