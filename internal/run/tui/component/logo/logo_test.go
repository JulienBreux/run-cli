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

package logo_test

import (
	"strings"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/tui/component/logo"
)

func TestString(t *testing.T) {
	s := logo.String()

	if !strings.Contains(s, "___ _   _ _  _") {
		t.Errorf("logo.String() should contain the logo art")
	}
	if !strings.Contains(s, "Julien Breux") {
		t.Errorf("logo.String() should contain the attribution")
	}
}

func TestStringSimple(t *testing.T) {
	s := logo.StringSimple()

	if !strings.Contains(s, "___ _   _ _  _") {
		t.Errorf("logo.StringSimple() should contain the logo art")
	}
	if strings.Contains(s, "Julien Breux") {
		t.Errorf("logo.StringSimple() should NOT contain the attribution")
	}
}

func TestNew(t *testing.T) {
	l := logo.New()

	if l == nil {
		t.Error("logo.New() should return a non-nil TextView")
	}

	// We can't easily extract text from TextView directly without drawing,
	// but we can check if it was initialized (not crashing).
	// In a real TUI test we might use a screen simulation, but for unit test ensuring it returns expected type is often enough.
	// However, we can assert properties if we want.
}
