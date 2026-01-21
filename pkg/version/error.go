/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
