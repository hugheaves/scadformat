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
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hugheaves/scadformat/internal/parser"
	"go.uber.org/zap"
)

type Formatter struct {
	settings *FormatSettings
}

func NewFormatter(fileName string) *Formatter {
	return &Formatter{
		settings: DefaultFormatSettings(fileName),
	}
}

func (f *Formatter) Format() error {
	if f.settings.fileName != "" {
		return f.formatFile()
	} else {
		return f.formatStdio()
	}
}

func (f *Formatter) formatFile() error {
	zap.S().Infof("formatting file %s", f.settings.fileName)
	err := checkFile(f.settings.fileName)
	if err != nil {
		return err
	}

	input, err := os.ReadFile(f.settings.fileName)
	if err != nil {
		zap.S().Errorf("failed to read file %s: %s", f.settings.fileName, err)
		return err
	}

	output, err := f.formatBytes(input)
	if err != nil {
		zap.S().Errorf("failed to format file %s: %s", f.settings.fileName, err)
		return err
	}

	timeStamp := time.Now().Format("2006-01-02_15-04-05")
	backupFileName := strings.TrimSuffix(f.settings.fileName, filepath.Ext(f.settings.fileName)) + "_" + timeStamp + ".scadbak"
	err = os.WriteFile(backupFileName, input, 0666)
	if err != nil {
		zap.S().Errorf("failed to write file %s: %s", backupFileName, err)
		return err
	}

	err = os.WriteFile(f.settings.fileName, output, 0666)
	if err != nil {
		zap.S().Errorf("failed to write file %s: %s", f.settings.fileName, err)
		return err
	}
	return nil
}

func (f *Formatter) formatStdio() error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		zap.S().Errorf("failed to read data: %s", err)
		return err
	}

	output, err := f.formatBytes(input)
	if err != nil {
		zap.S().Errorf("failed to format data: %s", err)
		return err
	}

	os.Stdout.Write(output)
	if err != nil {
		zap.S().Errorf("failed to write data: %s", err)
		return err
	}

	return nil
}

func (f *Formatter) formatBytes(input []byte) ([]byte, error) {
	antlrStream := antlr.NewIoStream(bytes.NewBuffer(input))
	lexer := parser.NewOpenSCADLexer(antlrStream)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewOpenSCADParser(tokens)
	outputBuffer := &bytes.Buffer{}
	formatter := NewTokenFormatter(f.settings, outputBuffer)
	v := NewFormattingVisitor(tokens, formatter)
	e := &ErrorListener{}
	p.AddErrorListener(e)
	p.Start_().Accept(v)

	return outputBuffer.Bytes(), e.lastErr
}

func checkFile(sourceFile string) error {
	sourceFileStat, err := os.Stat(sourceFile)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourceFile)
	}
	return nil
}
