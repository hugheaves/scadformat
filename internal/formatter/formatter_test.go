// SCADFormat - Formatter / beautifier for OpenSCAD source code
//
// Copyright (C) 2023  Hugh Eaves
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

package formatter

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

const (
	validInputDir   = "testdata/valid"
	invalidInputDir = "testdata/invalid"
	expectedDir     = "testdata/expected"
)

func TestFormat(t *testing.T) {
	inputFiles, err := filepath.Glob(filepath.Join(validInputDir, "*.scad"))
	if err != nil {
		t.Fatal(err)
	}

	for _, inputFile := range inputFiles {
		t.Run(filepath.Base(inputFile), func(t *testing.T) {
			testData := readTestData(t, validInputDir)

			formatter := NewFormatter("")

			output, err := formatter.formatBytes(testData)
			if err != nil {
				t.Fatal("error formatting:", err)
			}

			expected := readTestData(t, expectedDir)

			validateOutput(t, expected, output)
		})
	}
}

// Test that formatting is idempotent
func TestReformat(t *testing.T) {
	testData, err := filepath.Glob(filepath.Join(validInputDir, "*.scad"))
	if err != nil {
		t.Fatal(err)
	}

	for _, inputFile := range testData {
		t.Run(filepath.Base(inputFile), func(t *testing.T) {
			validInput := readTestData(t, validInputDir)

			formatter := NewFormatter("")

			output, err := formatter.formatBytes(validInput)
			if err != nil {
				t.Fatal("error formatting:", err)
			}

			expected := readTestData(t, expectedDir)

			validateOutput(t, expected, output)

			output, err = formatter.formatBytes(output)
			if err != nil {
				t.Fatal("error reformatting:", err)
			}

			validateOutput(t, expected, output)
		})
	}
}

func readTestData(t *testing.T, dir string) []byte {
	filename := filepath.Join(dir, filepath.Base(t.Name()))
	validInput, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal("error reading test file:", err)
	}
	return validInput
}

func validateOutput(t *testing.T, expected []byte, output []byte) {
	edits := myers.ComputeEdits(span.URIFromPath("foo"), string(expected), string(output))
	if len(edits) > 0 {
		diff := fmt.Sprint(gotextdiff.ToUnified("expected", "output", string(expected), edits))
		t.Errorf("Formatted output different than expected:\n" + diff)
	}
}
