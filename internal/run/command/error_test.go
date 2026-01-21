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

package command_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/command"
	"github.com/stretchr/testify/assert"
)

func TestPrintError(t *testing.T) {
	w := &bytes.Buffer{}
	testErr := errors.New("something went wrong")

	err := command.PrintError(w, testErr)

	assert.NoError(t, err)
	assert.Equal(t, "something went wrong\n", w.String())
}
