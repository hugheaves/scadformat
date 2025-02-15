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
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/hugheaves/scadformat/internal/logutil"
)

const (
	validInputDir   = "testdata/valid"
	invalidInputDir = "testdata/invalid"
	expectedDir     = "testdata/expected"
)

func TestMain(m *testing.M) {
	err := logutil.ConfigureLogging("")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestFormat(t *testing.T) {
	runTestOnDir(t, validInputDir, func(t *testing.T) {
		testData := readTestData(t, validInputDir)

		formatter := NewFormatter("")

		output, err := formatter.formatBytes(testData)
		if err != nil {
			t.Fatal("error formatting:", err)
		}

		expected := readTestData(t, expectedDir)

		err = validateOutput(t, expected, output)
		if err != nil {
			t.Fatal(err)
		}
	})
}

// Test that formatting is idempotent
func TestReformat(t *testing.T) {
	runTestOnDir(t, validInputDir, func(t *testing.T) {
		validInput := readTestData(t, validInputDir)

		formatter := NewFormatter("")

		output, err := formatter.formatBytes(validInput)
		if err != nil {
			t.Fatal("error formatting:", err)
		}

		expected := readTestData(t, expectedDir)

		err = validateOutput(t, expected, output)
		if err != nil {
			t.Fatal(err)
		}

		output, err = formatter.formatBytes(output)
		if err != nil {
			t.Fatal("error reformatting:", err)
		}

		err = validateOutput(t, expected, output)
		if err != nil {
			t.Fatal(err)
		}
	})
}

// This is not actually a test - it updates the contents of the "expected" testdata
// with the output of the formatter.
func TestUpdate(t *testing.T) {
	//t.Skip()

	runTestOnDir(t, validInputDir, func(t *testing.T) {
		validInput := readTestData(t, validInputDir)

		formatter := NewFormatter("")

		output, err := formatter.formatBytes(validInput)
		if err != nil {
			t.Fatal("error formatting:", err)
		}

		writeTestData(t, expectedDir, output)

	})
}

func runTestOnDir(t *testing.T, dir string, testFunc func(t *testing.T)) {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		path = strings.Join(strings.Split(path, string(os.PathSeparator))[2:], string(os.PathSeparator))
		var matched bool
		matched, err = filepath.Match("*.scad", filepath.Base(path))
		if err != nil {
			t.Fatal(err)
		}
		if !matched {
			return nil
		}

		t.Run(path, testFunc)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func readTestData(t *testing.T, dir string) []byte {
	filename := filepath.Join(dir, strings.Join(strings.Split(t.Name(), string(os.PathSeparator))[1:], string(os.PathSeparator)))
	validInput, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal("error reading test file:", err)
	}
	return validInput
}

func writeTestData(t *testing.T, dir string, data []byte) {
	filename := filepath.Join(dir, strings.Join(strings.Split(t.Name(), string(os.PathSeparator))[1:], string(os.PathSeparator)))
	err := os.WriteFile(filename, data, 0666)
	if err != nil {
		t.Fatal("error writing test file:", err)
	}
}

func validateOutput(t *testing.T, expected []byte, output []byte) error {
	edits := myers.ComputeEdits(span.URIFromPath("foo"), string(expected), string(output))
	if len(edits) > 0 {
		diff := fmt.Sprint(gotextdiff.ToUnified("expected", "output", string(expected), edits))
		return fmt.Errorf("Formatted output different than expected:\n%s", diff)
	}
	return nil
}
