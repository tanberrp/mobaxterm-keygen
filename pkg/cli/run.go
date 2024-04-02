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

package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/rand"
	cliflag "k8s.io/component-base/cli/flag"
)

func Run(cmd *cobra.Command) int {
	if err := run(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

func run(cmd *cobra.Command) (err error) {
	rand.Seed(time.Now().UnixNano())

	cmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)

	// When error printing is enabled for the Cobra command, a flag parse
	// error gets printed first, then optionally the often long usage
	// text. This is very unreadable in a console because the last few
	// lines that will be visible on screen don't include the error.
	//
	// The recommendation from #sig-cli was to print the usage text, then
	// the error. We implement this consistently for all commands here.
	// However, we don't want to print the usage text when command
	// execution fails for other reasons than parsing. We detect this via
	// the FlagParseError callback.
	//
	// Some commands, like kubectl, already deal with this themselves.
	// We don't change the behavior for those.
	if !cmd.SilenceUsage {
		cmd.SilenceUsage = true
		cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
			// Re-enable usage printing.
			c.SilenceUsage = false
			return err
		})
	}

	// In all cases error printing is done below.
	cmd.SilenceErrors = true

	err = cmd.Execute()
	return
}
