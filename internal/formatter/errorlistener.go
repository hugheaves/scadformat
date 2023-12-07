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
	"errors"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"go.uber.org/zap"
)

type ErrorListener struct {
	antlr.DefaultErrorListener
	lastErr error
}

func (e *ErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line int, column int, msg string, _ antlr.RecognitionException) {
	errMsg := fmt.Sprintf("syntax error on line %d:%d - %s", line, column, msg)
	zap.L().Error(errMsg)
	e.lastErr = errors.New(errMsg)
}
