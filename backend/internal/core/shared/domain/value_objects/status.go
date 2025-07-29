package value_objects

// Status represents a generic status value object
type Status string

// String returns string representation
func (s Status) String() string {
	return string(s)
}

// IsEmpty checks if status is empty
func (s Status) IsEmpty() bool {
	return string(s) == ""
}

// Equals compares two statuses
func (s Status) Equals(other Status) bool {
	return string(s) == string(other)
}
