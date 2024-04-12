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
	"archive/zip"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
)

// NewMobaXtermKeygenCommand creates the `mobaxterm-keygen` command with default arguments
func NewMobaXtermKeygenCommand() *cobra.Command {
	o := NewOptions()
	cmd := &cobra.Command{
		Use:   "mobaxterm-keygen",
		Short: "generate mobaxterm key",
		Long:  "generate mobaxterm key",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if errs := o.validate(); len(errs) != 0 {
				return errs.ToAggregate()
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&o.UserName, "username", "", "The user name to licensed, like: tanber")
	flags.StringVar(&o.Version, "version", "", "The version of MobaXterm, like: 23.5")
	flags.StringVar(&o.MobaxtermDir, "mobaxterm-dir", "", "The dir of mobaxterm has installed, like: ~/MobaXterm_Portable_V23.5")
	return cmd
}

func run(o *Options) error {
	major, minor, err := o.getVersion()
	if err != nil {
		return err
	}
	license := generateLicense(1, 1, o.UserName, major, minor)
	fileName := path.Join(o.MobaxtermDir, "Custom.mxtpro")
	_ = os.Remove(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	zipFile := zip.NewWriter(f)
	defer zipFile.Close()
	header := &zip.FileHeader{
		Name:               "Pro.key",
		Method:             zip.Store,
		CompressedSize64:   38,
		UncompressedSize64: 38,
	}
	header.SetModTime(time.Now())
	proFile, err := zipFile.CreateRaw(header)
	if err != nil {
		return err
	}
	_, err = proFile.Write(license)
	if err != nil {
		return err
	}
	return nil
}

func generateLicense(userType, count int, username string, major, minor int64) []byte {
	licenseString := fmt.Sprintf("%d#%s|%d%d#%d#%d3%d6%d#%d#%d#%d#", userType, username, major, minor, count, major, minor, minor, 0, 0, 0)
	return variantBase64Encode(encryptBytes(0x787, []byte(licenseString)))
}

func encryptBytes(key int, bs []byte) []byte {
	var result []byte
	for _, b := range bs {
		encryptedByte := b ^ byte((key>>8)&0xff)
		result = append(result, encryptedByte)
		key = (int(result[len(result)-1]) & key) | 0x482D
	}
	return result
}

var (
	variantBase64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	variantBase64Map   = func() map[int]byte {
		result := make(map[int]byte)
		for i, v := range variantBase64Table {
			result[i] = byte(v)
		}
		return result
	}()
)

func variantBase64Encode(bs []byte) []byte {
	var result []byte
	blocksCount := len(bs) / 3
	leftBytes := len(bs) % 3
	for i := 0; i < blocksCount; i++ {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*i : 3*i+3])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>12)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>18)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
	}
	if leftBytes == 0 {
		return result
	} else if leftBytes == 1 {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*blocksCount:])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
		return result
	} else {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*blocksCount:])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>12)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
		return result
	}
}

func littleEndianBytes(bs []byte) int {
	var result = int(bs[0])
	for i := 1; i < len(bs); i++ {
		result = result | int(bs[i])<<(8*i)
	}
	return result
}
