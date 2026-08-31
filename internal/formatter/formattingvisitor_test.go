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
	"errors"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hugheaves/scadformat/internal/parser"
)

// These tests exercise the visitor code paths that used to call
// zap.S().Fatalf/Fatal (which terminates the process via os.Exit) and now
// instead record the error on the visitor via setErr. None of these states
// are reachable through the normal parse->visit pipeline in formatBytes,
// since a parser-reported syntax error already prevents visitation (see
// TestInvalid), so they're exercised directly here against a visitor backed
// by a real token stream.

// newTestVisitor builds a FormattingVisitor backed by a real token stream
// (lexed from input) so that visitor methods relying on the token stream
// (e.g. printCommentsBefore) behave correctly without a full parse tree.
func newTestVisitor(input string) *FormattingVisitor {
	lexer := parser.NewOpenSCADLexer(antlr.NewInputStream(input))
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	tokens.Fill()
	formatter := NewTokenFormatter(DefaultFormatSettings(""), &bytes.Buffer{})
	return NewFormattingVisitor(tokens, formatter)
}

func TestVisitErrorNodeRecordsError(t *testing.T) {
	v := newTestVisitor("garbage")

	errorNode := antlr.NewErrorNodeImpl(v.tokenStream.Get(0))

	v.VisitErrorNode(errorNode)

	if v.err == nil {
		t.Fatal("expected VisitErrorNode to record an error instead of exiting the process")
	}
}

func TestVisitIncludeOrUseFileRecordsErrorOnRegexMismatch(t *testing.T) {
	input := "include <foo.scad>;"
	v := newTestVisitor(input)

	ctx := parser.NewIncludeOrUseFileContext(nil, nil, 0)
	// The lexer only ever produces INCLUDE_OR_USE_FILE tokens matching
	// includeOrUseRegex, so mangle the token's text to force a mismatch -
	// this simulates the "should never happen" defensive branch.
	token := v.tokenStream.Get(0)
	token.SetText("not a valid include statement")
	ctx.AddTokenNode(token)
	ctx.SetStart(token)

	v.VisitIncludeOrUseFile(ctx)

	if v.err == nil {
		t.Fatal("expected VisitIncludeOrUseFile to record an error instead of exiting the process")
	}
}

func TestVisitChildStatementRecordsErrorOnInvalidState(t *testing.T) {
	v := newTestVisitor("")

	ctx := parser.NewChildStatementContext(nil, nil, 0)

	v.VisitChildStatement(ctx)

	if v.err == nil {
		t.Fatal("expected VisitChildStatement to record an error instead of exiting the process")
	}
}

func TestSetErrKeepsFirstError(t *testing.T) {
	v := newTestVisitor("")

	errFirst := errors.New("first error")
	errSecond := errors.New("second error")

	v.setErr(errFirst)
	v.setErr(errSecond)

	if v.err != errFirst {
		t.Fatalf("expected first error to be retained, got: %v", v.err)
	}
}
