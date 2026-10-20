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
	"os"
	"strconv"
	"strings"
)

type KernelVersion struct {
	Major int
	Minor int
	Patch int
}

func GetKernelVersion() *KernelVersion {

	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return nil
	}

	return parseKernelVersion(string(data))
}

// parseKernelVersion extracts the major.minor.patch triplet from the content of
// /proc/version, e.g.:
//
//	Linux version 5.15.0-88-generic (buildd@...) (gcc ...) #88-Ubuntu SMP ...
//
// The third whitespace field is the kernel release. Distro kernels (Ubuntu,
// Debian, CentOS, ...) append an ABI/build suffix after a '-', such as
// "0-88-generic" or "0-553.el8_9.1.x86_64". Only the numeric part before the
// '-' is the patch level; the suffix must be stripped before converting it to
// an integer, otherwise the whole version fails to parse.
func parseKernelVersion(versionString string) *KernelVersion {

	versionFields := strings.Fields(versionString)
	if len(versionFields) < 3 {
		return nil
	}

	version := strings.Split(versionFields[2], ".")
	if len(version) < 3 {
		return nil
	}

	major, err := strconv.Atoi(version[0])
	if err != nil {
		return nil
	}

	minor, err := strconv.Atoi(version[1])
	if err != nil {
		return nil
	}

	// Strip the distro ABI/build suffix, e.g. "0-88-generic" -> "0".
	patch, err := strconv.Atoi(strings.SplitN(version[2], "-", 2)[0])
	if err != nil {
		return nil
	}

	ver := &KernelVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}

	return ver
}
