/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

// import-boss enforces import restrictions in a given repository.
//
// When a directory is verified, import-boss looks for a file called
// ".import-restrictions". If this file is not found, parent directories will be
// recursively searched.
//
// If an ".import-restrictions" file is found, then all imports of the package
// are checked against each "rule" in the file. A rule consists of three parts:
// * A SelectorRegexp, to select the import paths that the rule applies to.
// * A list of AllowedPrefixes
// * A list of ForbiddenPrefixes
// An import is allowed if it matches at least one allowed prefix and does not
// match any forbidden prefix. An example file looks like this:
//
// {
//   "Rules": [
//     {
//       "SelectorRegexp": "k8s[.]io",
//       "AllowedPrefixes": [
//         "k8s.io/kubernetes/cmd/libs/go2idl",
//         "k8s.io/kubernetes/third_party"
//       ],
//       "ForbiddenPrefixes": [
//         "k8s.io/kubernetes/pkg/third_party/deprecated"
//       ]
//     },
//     {
//       "SelectorRegexp": "^unsafe$",
//       "AllowedPrefixes": [
//       ],
//       "ForbiddenPrefixes": [
//         ""
//       ]
//     }
//   ]
// }
//
// Note the secound block explicitly matches the unsafe package, and forbids it
// ("" is a prefix of everything).
package main

import (
	"os"

	"k8s.io/kubernetes/cmd/libs/go2idl/args"
	"k8s.io/kubernetes/cmd/libs/go2idl/import-boss/generators"

	"github.com/golang/glog"
)

func main() {
	arguments := args.Default()

	// Override defaults. These are Kubernetes specific input and output
	// locations.
	arguments.InputDirs = []string{
		"k8s.io/kubernetes/pkg/",
		"k8s.io/kubernetes/cmd/",
		"k8s.io/kubernetes/plugin/",
	}
	arguments.OutputBase = ""
	arguments.Recursive = true
	// arguments.VerifyOnly = true

	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		glog.Errorf("Error: %v", err)
		os.Exit(1)
	}
	glog.Info("Completed successfully.")
}
