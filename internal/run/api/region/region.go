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

package region

// Represents all regions.
const ALL = "all"

// List returns a list of supported Cloud Run regions.
func List() []string {
	return []string{
		"asia-east1", "asia-northeast1", "asia-northeast2", "asia-northeast3",
		"asia-south1", "asia-southeast1", "asia-southeast2", "australia-southeast1",
		"europe-central2", "europe-north1", "europe-west1", "europe-west2",
		"europe-west3", "europe-west4", "europe-west6", "northamerica-northeast1",
		"southamerica-east1", "us-central1", "us-east1", "us-east4",
		"us-west1", "us-west2", "us-west3", "us-west4",
	}
}
