/*
Copyright 2024 ByteDance and/or its affiliates.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import (
	"reflect"
	"testing"
)

func TestParseKernelVersion(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected *KernelVersion
	}{
		{
			name:     "ubuntu distro kernel with abi suffix",
			input:    "Linux version 5.15.0-88-generic (buildd@lcy02-amd64-045) (gcc) #88-Ubuntu SMP Tue Nov 14 00:36:04 UTC 2023",
			expected: &KernelVersion{Major: 5, Minor: 15, Patch: 0},
		},
		{
			name:     "centos/rhel distro kernel with el8 suffix",
			input:    "Linux version 4.18.0-553.el8_9.1.x86_64 (gcc) #1 SMP",
			expected: &KernelVersion{Major: 4, Minor: 18, Patch: 0},
		},
		{
			name:     "debian distro kernel with amd64 suffix",
			input:    "Linux version 5.10.0-26-amd64 (gcc) #1 SMP Debian 5.10.197-1 (2023-09-19)",
			expected: &KernelVersion{Major: 5, Minor: 10, Patch: 0},
		},
		{
			name:     "newer ubuntu 6.x distro kernel",
			input:    "Linux version 6.5.0-35-generic (gcc) #35 SMP PREEMPT_DYNAMIC",
			expected: &KernelVersion{Major: 6, Minor: 5, Patch: 0},
		},
		{
			name:     "vanilla kernel without abi suffix",
			input:    "Linux version 5.15.0 #1 SMP built-from-source",
			expected: &KernelVersion{Major: 5, Minor: 15, Patch: 0},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "too few fields",
			input:    "Linux version",
			expected: nil,
		},
		{
			name:     "too few dot components",
			input:    "Linux version 5.15 (gcc) #1 SMP",
			expected: nil,
		},
		{
			name:     "non numeric version",
			input:    "Linux version foo.bar.baz (gcc) #1 SMP",
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseKernelVersion(tc.input)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("parseKernelVersion(%q) = %+v, want %+v", tc.input, got, tc.expected)
			}
		})
	}
}

// TestParseKernelVersionDistroNotNil guards the regression where the distro ABI
// suffix (e.g. "0-88-generic") made strconv.Atoi fail and the whole version
// parsed as nil, which silently disabled kernel-gated auto-fixes.
func TestParseKernelVersionDistroNotNil(t *testing.T) {
	distro := []string{
		"Linux version 5.15.0-88-generic (buildd) #88-Ubuntu SMP",
		"Linux version 4.18.0-553.el8_9.1.x86_64 (gcc) #1 SMP",
		"Linux version 5.10.0-26-amd64 (gcc) #1 SMP Debian",
		"Linux version 6.5.0-35-generic (gcc) #35 SMP",
	}
	for _, s := range distro {
		if v := parseKernelVersion(s); v == nil {
			t.Fatalf("distro kernel version parsed as nil: %s", s)
		}
	}
}
