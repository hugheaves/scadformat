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
	currentIndent int
	linePos       int
	inLine        bool
	wrappedLine   bool
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
	for i := 0; i < len(lines)-1; i++ {
		err := tokenFormatter.printWithLineWrap(lines[i])
		if err != nil {
			return err
		}
		err = tokenFormatter.printNewLine()
		if err != nil {
			return err
		}
	}
	if len(lines) > 0 {
		err := tokenFormatter.printWithLineWrap(lines[len(lines)-1])
		if err != nil {
			return err
		}
	}
	return nil
}

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

func (tokenFormatter *TokenFormatter) printSpace() error {
	if tokenFormatter.inLine {
		err := tokenFormatter.printWithLineWrap(" ")
		if err != nil {
			return err
		}
	}
	return nil
}

func (tokenFormatter *TokenFormatter) endLine() error {
	if tokenFormatter.inLine {
		err := tokenFormatter.printNewLine()
		if err != nil {
			return err
		}
	}
	return nil
}

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
	//	zap.L().Debug("Indenting to " + tokenFormatter.indentSize)

}

func (tokenFormatter *TokenFormatter) unindent() {
	tokenFormatter.currentIndent -= tokenFormatter.settings.indentSize
	// zap.L().Debug("Unindenting to " + tokenFormatter.indentSize)
}
