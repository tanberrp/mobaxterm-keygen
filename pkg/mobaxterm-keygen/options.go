// Copyright 2024 The mobaxterm-keygen Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mobaxterm_keygen

import (
	"os"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Options struct {
	UserName  string
	Version   string
	OutputDir string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) validate() field.ErrorList {
	errs := field.ErrorList{}
	if o.UserName == "" {
		errs = append(errs, field.Required(field.NewPath("--username"), ""))
	}
	if err := o.validateVersion(); err != nil {
		errs = append(errs, err)
	}

	if err := o.validateOutputDir(); err != nil {
		errs = append(errs, err)
	}
	return errs
}

func (o *Options) validateVersion() *field.Error {
	if o.Version == "" {
		return field.Required(field.NewPath("--version"), "")
	}
	v := strings.Split(o.Version, ".")
	if len(v) != 2 {
		return field.Invalid(field.NewPath("--version"), o.Version, "valid version like: 23.5")
	}
	return nil
}
func (o *Options) validateOutputDir() *field.Error {
	if o.OutputDir == "" {
		return field.Required(field.NewPath("--output-dir"), "")
	}
	dir, err := os.Stat(o.OutputDir)
	if err != nil {
		return field.Invalid(field.NewPath("--output-dir"), o.OutputDir, err.Error())
	}
	if !dir.IsDir() {
		return field.Invalid(field.NewPath("--output-dir"), o.OutputDir, "is not dir")
	}
	return nil
}

func (o *Options) getVersion() (major, minor int64, err error) {
	v := strings.Split(o.Version, ".")
	major, err = strconv.ParseInt(v[0], 10, 64)
	if err != nil {
		return
	}
	minor, err = strconv.ParseInt(v[1], 10, 64)
	if err != nil {
		return
	}
	return
}
