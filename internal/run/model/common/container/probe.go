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

package container

// Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
type Probe struct {
	InitialDelaySeconds int32            `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      int32            `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       int32            `json:"periodSeconds,omitempty"`
	SuccessThreshold    int32            `json:"successThreshold,omitempty"`
	FailureThreshold    int32            `json:"failureThreshold,omitempty"`
	HTTPGet             *HTTPGetAction   `json:"httpGet,omitempty"`
	TCPSocket           *TCPSocketAction `json:"tcpSocket,omitempty"`
	Exec                *ExecAction      `json:"exec,omitempty"`
}

// HTTPGetAction describes an action involving a HTTP GET request.
type HTTPGetAction struct {
	Path        string        `json:"path,omitempty"`
	HTTPHeaders []*HTTPHeader `json:"httpHeaders,omitempty"`
	Port        int32         `json:"port,omitempty"`
}

// HTTPHeader describes a custom header to be used in HTTP probes
type HTTPHeader struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// TCPSocketAction describes an action involving a TCP port.
type TCPSocketAction struct {
	Port int32 `json:"port,omitempty"`
}
