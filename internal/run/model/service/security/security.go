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

package security

// Security represents the security settings of a service.
type Security struct {
	InvokerIAMDisabled      bool   `json:"invokerIamDisabled"`
	ServiceAccount          string `json:"serviceAccount,omitempty"`
	EncryptionKey           string `json:"encryptionKey,omitempty"`
	BinaryAuthorization     string `json:"binaryAuthorization,omitempty"`
	BreakglassJustification string `json:"breakglassJustification,omitempty"`
}
