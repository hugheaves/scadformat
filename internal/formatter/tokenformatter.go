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
	"io"
	"strings"

	"go.uber.org/zap"
)

type TokenFormatter struct {
	settings      *FormatSettings
	writer        io.Writer
	currentIndent int  // current indent size for new lines
	linePos       int  // position that next character will be written to the line
	inLine        bool // true if the current line contains text
	wrappedLine   bool // true if the previous print statement caused the text to wrap to the next line
}

func NewTokenFormatter(settings *FormatSettings, writer io.Writer) *TokenFormatter {
	return &TokenFormatter{
		settings:      settings,
		writer:        writer,
		currentIndent: 0,
		linePos:       0,
		inLine:        false,
		wrappedLine:   false,
	}
}

func (tokenFormatter *TokenFormatter) printString(strVal string) error {
	zap.L().Debug("printString |" + strVal + "|")
	lines := strings.Split(strVal, "\n")
	for i, line := range lines {
		err := tokenFormatter.printWithLineWrap(line)
		if err != nil {
			return err
		}
		if i < len(lines)-1 {
			err = tokenFormatter.printNewLine()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// prints a string on the current line. If the length of the string exceeds
// the remaining space on the line, the text is printed on a new line, with
// additional indention applied to indicate that the line was wrapped.
func (tokenFormatter *TokenFormatter) printWithLineWrap(strVal string) error {
	if tokenFormatter.inLine && len(strVal) > tokenFormatter.lineRemaining() {
		err := tokenFormatter.printNewLine()
		if err != nil {
			return err
		}
		tokenFormatter.indent()
		tokenFormatter.wrappedLine = true
	}
	err := tokenFormatter.appendToLine(strVal)
	return err
}

// printSpace adds a space to the line, if the current line contains is not empty. Otherwise, this function does nothing.
func (tokenFormatter *TokenFormatter) printSpace() error {
	if tokenFormatter.inLine {
		err := tokenFormatter.printWithLineWrap(" ")
		if err != nil {
			return err
		}
	}
	return nil
}

// endLine calls printNewLine if the current line is not empty. Otherwise, it does nothing.
func (tokenFormatter *TokenFormatter) endLine() error {
	if tokenFormatter.inLine {
		err := tokenFormatter.printNewLine()
		if err != nil {
			return err
		}
	}
	return nil
}

// printNewLine prints a new line character, and removes any
// indentation applied by printWithLineWrap.
func (tokenFormatter *TokenFormatter) printNewLine() error {
	_, err := fmt.Fprintln(tokenFormatter.writer)
	if err != nil {
		return err
	}
	if tokenFormatter.wrappedLine {
		tokenFormatter.unindent()
		tokenFormatter.wrappedLine = false
	}
	tokenFormatter.inLine = false
	tokenFormatter.linePos = 0
	return nil
}

func (tokenFormatter *TokenFormatter) lineRemaining() int {
	if !tokenFormatter.inLine {
		return tokenFormatter.settings.maxLineLen - tokenFormatter.currentIndent
	} else {
		return tokenFormatter.settings.maxLineLen - tokenFormatter.linePos
	}
}

// appendToLine appends a string to the current line. If the line is empty,
// outputIndent is called to indent the line before adding the string.
func (tokenFormatter *TokenFormatter) appendToLine(strVal string) error {
	if len(strVal) == 0 {
		return nil
	}
	if !tokenFormatter.inLine {
		tokenFormatter.outputIndent()
		tokenFormatter.inLine = true
	}
	_, err := fmt.Fprint(tokenFormatter.writer, strVal)
	if err != nil {
		return err
	}
	tokenFormatter.linePos += len(strVal)
	return nil
}

func (tokenFormatter *TokenFormatter) outputIndent() error {
	for i := 0; i < tokenFormatter.currentIndent; i++ {
		_, err := fmt.Fprint(tokenFormatter.writer, " ")
		if err != nil {
			return err
		}
	}
	tokenFormatter.linePos += tokenFormatter.currentIndent
	return nil
}

func (tokenFormatter *TokenFormatter) indent() {
	tokenFormatter.currentIndent += tokenFormatter.settings.indentSize
}

func (tokenFormatter *TokenFormatter) unindent() {
	tokenFormatter.currentIndent -= tokenFormatter.settings.indentSize
}
