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
	"strings"
	"testing"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/hugheaves/scadformat/internal/logutil"
)

func TestFormat(t *testing.T) {
	logutil.ConfigureLogging("info")
	testDataPaths, err := filepath.Glob(filepath.Join("testdata", "*.scad"))
	if err != nil {
		t.Fatal(err)
	}

	for _, testDataPath := range testDataPaths {
		t.Run(testDataPath, execFormat)
	}
}

func execFormat(t *testing.T) {
	paths := strings.Split(t.Name(), string(os.PathSeparator))
	filename := filepath.Join(paths[1:]...)
	source, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal("error reading test file:", err)
	}

	formatter := NewFormatter("")

	output, err := formatter.formatBytes(source)
	if err != nil {
		t.Fatal("error formatting:", err)
	}

	expected := readExpected(t, filename)

	validateOutput(t, expected, output)
}

func testReformat(t *testing.T, filePath string) {
	_, fileName := filepath.Split(filePath)
	t.Run("testExpected("+fileName+")", func(t *testing.T) {
		source, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatal("error reading test file:", err)
		}

		expected := readExpected(t, filePath)

		formatter := NewFormatter("")

		output, err := formatter.formatBytes(source)
		if err != nil {
			t.Fatal("error formatting:", err)
		}

		validateOutput(t, expected, output)

		output, err = formatter.formatBytes(output)
		if err != nil {
			t.Fatal("error reformatting:", err)
		}

		validateOutput(t, expected, output)
	})
}

func readExpected(t *testing.T, filePath string) []byte {
	expectedResultPath := filepath.Join(filePath + ".expected")
	expected, err := os.ReadFile(expectedResultPath)
	if err != nil {
		t.Fatal("error reading expected file:", err)
	}
	return expected
}

func validateOutput(t *testing.T, expected []byte, output []byte) {
	edits := myers.ComputeEdits(span.URIFromPath("expected"), string(expected), string(output))

	// fmt.Println("===============")
	// fmt.Println(string(output))
	// fmt.Println("===============")
	if len(edits) > 0 {
		diff := fmt.Sprint(gotextdiff.ToUnified("expected", "output", string(expected), edits))
		t.Errorf("Formatted output different than expected:\n" + diff)
	}
}
