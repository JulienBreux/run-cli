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

package traffic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTrafficParams(t *testing.T) {
	tests := []struct {
		name    string
		params  []string
		wantSum int64
		wantErr bool
	}{
		{
			name:    "Valid sum 100",
			params:  []string{"50", "50"},
			wantSum: 100,
			wantErr: false,
		},
		{
			name:    "Invalid sum != 100",
			params:  []string{"50", "40"},
			wantSum: 90,
			wantErr: true,
		},
		{
			name:    "Invalid input",
			params:  []string{"50", "abc"},
			wantErr: true,
		},
		{
			name:    "Single 100",
			params:  []string{"100"},
			wantSum: 100,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum, err := validateTrafficParams(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSum, sum)
			}
		})
	}
}
