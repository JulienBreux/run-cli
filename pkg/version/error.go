package version

// DateParseError represents the date parse error
type DateParseError struct {
	Date string
	Err  error
}

// Error returns human readable error
func (e *DateParseError) Error() string {
	return "unable to parse date: " + e.Date
}

// Unwrap unwraps the original error
func (e *DateParseError) Unwrap() error {
	return e.Err
}
