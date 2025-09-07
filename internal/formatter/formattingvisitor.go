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
	"reflect"
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hugheaves/scadformat/internal/parser"
	"go.uber.org/zap"
)

// check that Visitor implements OpenSCADVisitor
var _ parser.OpenSCADVisitor = &FormattingVisitor{}

var includeOrUseRegex = regexp.MustCompile("(include|use)[ \t]*(<[^\t\r\n>]*>)")

type FormattingVisitor struct {
	parser.BaseOpenSCADVisitor
	tokenStream             antlr.TokenStream
	formatter               *TokenFormatter
	lastPrintedCommentIndex int
	endLineAfterComma       bool
}

func NewFormattingVisitor(tokenStream antlr.TokenStream, formatter *TokenFormatter) *FormattingVisitor {
	visitor := &FormattingVisitor{
		tokenStream:             tokenStream,
		formatter:               formatter,
		lastPrintedCommentIndex: 0,
	}

	// Override VisitChildren in BaseOpenClassVisitor
	// as part of workaround for https://github.com/antlr/antlr4/issues/2504
	visitor.BaseOpenSCADVisitor.VisitChildren = visitor.VisitChildren

	return visitor
}

func (v *FormattingVisitor) Visit(tree antlr.ParseTree) interface{} {
	if tree != nil {
		zap.S().Debugf("Visiting: %s", reflect.TypeOf(tree).String())
		tree.Accept(v)
		zap.S().Debugf("Visited: %s", reflect.TypeOf(tree).String())
	}
	return nil
}

func (v *FormattingVisitor) VisitChildren(tree antlr.RuleNode) interface{} {
	for _, child := range tree.GetChildren() {
		val := child.(antlr.ParseTree)
		v.Visit(val)
	}
	return nil
}

func (v *FormattingVisitor) VisitTerminal(node antlr.TerminalNode) interface{} {

	v.printCommentsBefore(node.GetSymbol().GetTokenIndex())
	v.formatter.printString(node.GetText())
	v.printEndOfLineCommentAfter(node.GetSymbol().GetTokenIndex())
	return nil
}

func (v *FormattingVisitor) VisitErrorNode(errorNode antlr.ErrorNode) interface{} {
	zap.S().Fatalf("Unable to resolve parsing error: %s", errorNode.GetText())
	return nil
}

func (v *FormattingVisitor) VisitStart(ctx *parser.StartContext) interface{} {
	// Visit only "input", not the EOF token
	v.Visit(ctx.Input())
	v.printCommentsBefore(v.tokenStream.Size())
	return nil
}

func (v *FormattingVisitor) VisitAssignment(ctx *parser.AssignmentContext) interface{} {
	v.Visit(ctx.ID())
	v.formatter.printSpace()
	v.Visit(ctx.EQUALS())
	v.formatter.printSpace()
	v.Visit(ctx.Expr())
	v.Visit(ctx.SEMICOLON())
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitAssignmentExpression(ctx *parser.AssignmentExpressionContext) interface{} {
	v.Visit(ctx.ID())
	v.formatter.printSpace()
	v.Visit(ctx.EQUALS())
	v.formatter.printSpace()
	v.Visit(ctx.Expr())
	return nil
}

func (v *FormattingVisitor) VisitBinaryExpr(ctx *parser.BinaryExprContext) interface{} {
	v.Visit(ctx.Expr(0))
	v.formatter.printSpace()
	v.Visit(ctx.BinaryOperator())
	v.formatter.printSpace()
	v.Visit(ctx.Expr(1))
	return nil
}

func (v *FormattingVisitor) VisitTernaryExpr(ctx *parser.TernaryExprContext) interface{} {
	v.Visit(ctx.Expr(0))
	v.formatter.printSpace()
	v.Visit(ctx.QUESTION_MARK())
	v.formatter.printSpace()
	v.Visit(ctx.Expr(1))
	v.formatter.printSpace()
	v.Visit(ctx.COLON())
	v.formatter.printSpace()
	v.Visit(ctx.Expr(2))
	return nil

}

func (v *FormattingVisitor) VisitChildStatements(ctx *parser.ChildStatementsContext) interface{} {
	v.Visit(ctx.L_CURLY())
	v.formatter.endLine()
	for _, child := range ctx.AllChildStatementOrAssignment() {
		if child.Assignment() != nil {
			v.formatter.indent()
		}
		v.Visit(child)
		if child.Assignment() != nil {
			v.formatter.unindent()
		}
	}
	v.Visit(ctx.R_CURLY())
	return nil
}

func (v *FormattingVisitor) VisitSemicolon(ctx *parser.SemicolonContext) interface{} {
	v.Visit(ctx.SEMICOLON())
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitFunctionDefinition(ctx *parser.FunctionDefinitionContext) interface{} {
	v.Visit(ctx.FUNCTION())
	v.formatter.printSpace()
	v.Visit(ctx.ID())
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Parameters())
	v.Visit(ctx.R_PAREN())
	v.formatter.printSpace()
	v.Visit(ctx.EQUALS())
	v.formatter.endLine()
	v.formatter.indent()
	v.Visit(ctx.Expr())
	v.Visit(ctx.SEMICOLON())
	v.formatter.endLine()
	v.formatter.unindent()
	return nil
}

func (v *FormattingVisitor) VisitModuleDefinition(ctx *parser.ModuleDefinitionContext) interface{} {
	v.Visit(ctx.MODULE())
	v.formatter.printSpace()
	v.Visit(ctx.ID())
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Parameters())
	v.Visit(ctx.R_PAREN())
	v.formatter.printSpace()
	v.Visit(ctx.Statement())
	return nil
}

func (v *FormattingVisitor) VisitAssertExpr(ctx *parser.AssertExprContext) interface{} {
	v.Visit(ctx.ASSERT())
	v.formatter.printSpace()
	v.Visit(ctx.ParenArgs())
	if ctx.Expr() != nil {
		v.formatter.endLine()
		v.formatter.indent()
		v.Visit(ctx.Expr())
		v.formatter.unindent()
	}
	return nil
}

func (v *FormattingVisitor) VisitEchoExpr(ctx *parser.EchoExprContext) interface{} {
	v.Visit(ctx.ECHO())
	v.formatter.printSpace()
	v.Visit(ctx.ParenArgs())
	if ctx.Expr() != nil {
		v.formatter.endLine()
		v.formatter.indent()
		v.Visit(ctx.Expr())
		v.formatter.unindent()
	}
	return nil
}

func (v *FormattingVisitor) VisitLetExpr(ctx *parser.LetExprContext) interface{} {
	v.Visit(ctx.LET())
	v.formatter.printSpace()
	v.Visit(ctx.ParenArgs())
	v.formatter.endLine()
	v.formatter.indent()
	v.Visit(ctx.Expr())
	v.formatter.unindent()
	return nil
}

func (v *FormattingVisitor) VisitSingleModuleInstantiation(ctx *parser.SingleModuleInstantiationContext) interface{} {
	v.VisitChildren(ctx)
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitStatements(ctx *parser.StatementsContext) interface{} {
	v.Visit(ctx.L_CURLY())
	v.formatter.endLine()
	v.formatter.indent()
	for _, childCtx := range ctx.GetChildren() {
		if _, ok := childCtx.(parser.IStatementContext); ok {
			v.Visit(childCtx.(antlr.RuleContext))
		}
	}
	v.formatter.unindent()
	v.Visit(ctx.R_CURLY())
	v.formatter.endLine()
	return nil
}

// Unlike all other statements in the grammar, the grammar parses the entire "IncludeOrUseFile" statement
// as a single token. Therefore we can't Visit each part of the statement to apply formatting rules.
// Instead, we break apart the token using a regex, and format the individual pieces.
func (v *FormattingVisitor) VisitIncludeOrUseFile(ctx *parser.IncludeOrUseFileContext) interface{} {
	// printCommentsBefore normally gets called by "Visit", but we're not calling Visit here.
	v.printCommentsBefore(ctx.GetStart().GetTokenIndex())
	matches := includeOrUseRegex.FindStringSubmatch(ctx.INCLUDE_OR_USE_FILE().GetText())
	if len(matches) != 3 {
		zap.S().Fatal("Failed to parse the include or use statement: %s", ctx.GetText())
	}
	v.formatter.printString(matches[1])
	v.formatter.printSpace()
	v.formatter.printString(matches[2])
	v.formatter.endLine()
	return nil
}

// func (v *FormattingVisitor) VisitArrayAccess(ctx *parser.ArrayAccessContext) interface{} {
// 	return v.VisitChildren(ctx)
// }

func (v *FormattingVisitor) VisitVector(ctx *parser.VectorContext) interface{} {
	// Determine if this vector has other vectors
	// or list comprehension elements nested inside
	var allVectorElements = ctx.AllVectorElement()
	nested := false
	for _, ve := range allVectorElements {
		if ve.ListComprehensionElementsP() != nil {
			nested = true
		} else {
			expr := ve.Expr()
			if expr != nil {
				ec, ok := expr.(*parser.CallExprContext)
				if ok {
					if ec.Call().Primary().Vector() != nil {
						nested = true
					}
				}
			}
		}
	}
	v.Visit(ctx.L_BRACKET())
	if nested {
		v.formatter.endLine()
		v.formatter.indent()
	}

	for i, ve := range allVectorElements {
		v.Visit(ve)
		if ctx.Comma(i) != nil {
			v.Visit(ctx.Comma(i))
		}
		if nested {
			v.formatter.endLine()
		} else if i < len(allVectorElements)-1 {
			v.formatter.printSpace()
		}
	}
	if nested {
		v.formatter.unindent()
	}
	v.Visit(ctx.R_BRACKET())
	return nil
}

func (v *FormattingVisitor) VisitForStatementComprehension(ctx *parser.ForStatementComprehensionContext) interface{} {
	v.Visit(ctx.FOR())
	v.formatter.printSpace()
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Arguments(0))
	if ctx.SEMICOLON(0) != nil {
		v.Visit(ctx.SEMICOLON(0))
		v.Visit(ctx.Expr())
		v.Visit(ctx.SEMICOLON(1))
		v.Visit(ctx.Arguments(1))
	}
	v.Visit(ctx.R_PAREN())
	v.formatter.endLine()
	v.formatter.indent()
	v.Visit(ctx.VectorElement())
	v.formatter.unindent()
	return nil
}

func (v *FormattingVisitor) VisitLetStatementComprehension(ctx *parser.LetStatementComprehensionContext) interface{} {
	v.Visit(ctx.LET())
	v.formatter.printSpace()
	v.Visit(ctx.ParenArgs())
	v.formatter.endLine()
	v.formatter.indent()
	v.Visit(ctx.ListComprehensionElementsP())
	v.formatter.unindent()
	return nil
}

func (v *FormattingVisitor) VisitEachStatementComprehension(ctx *parser.EachStatementComprehensionContext) interface{} {
	v.Visit(ctx.EACH())
	v.formatter.printSpace()
	v.Visit(ctx.VectorElement())
	return nil
}

func (v *FormattingVisitor) VisitIfStatementComprehension(ctx *parser.IfStatementComprehensionContext) interface{} {
	v.Visit(ctx.IF())
	v.formatter.printSpace()
	v.Visit(ctx.ParenExpr())
	v.formatter.endLine()
	v.formatter.indent()
	v.Visit(ctx.VectorElement(0))
	v.formatter.unindent()
	if ctx.ELSE() != nil {
		v.formatter.endLine()
		v.Visit(ctx.ELSE())
		v.formatter.endLine()
		v.formatter.indent()
		v.Visit(ctx.VectorElement(1))
		v.formatter.unindent()
	}
	return nil
}

func (v *FormattingVisitor) VisitIfElseStatement(ctx *parser.IfElseStatementContext) interface{} {
	v.Visit(ctx.IfStatement())
	if ctx.ELSE() != nil {
		v.formatter.printSpace()
		v.Visit(ctx.ELSE())
		v.Visit(ctx.ChildStatement())
	}
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitIfStatement(ctx *parser.IfStatementContext) interface{} {
	v.Visit(ctx.IF())
	v.formatter.printSpace()
	v.Visit(ctx.ParenExpr())
	v.Visit(ctx.ChildStatement())
	return nil
}

// func (v *FormattingVisitor) VisitComma(ctx *parser.CommaContext) interface{} {
// 	v.Visit(ctx.COMMA())
// 	v.formatter.printSpace()
// 	return nil
// }

// Formats a childStatement, which can be a semicolon, a moduleInstantiation, or a childStatements node.
func (v *FormattingVisitor) VisitChildStatement(ctx *parser.ChildStatementContext) interface{} {
	zap.S().Debugf("formatChildStatement: parent=%s", reflect.TypeOf(ctx.GetParent()))

	if ctx.Semicolon() != nil {
		v.Visit(ctx.Semicolon())
	} else if ctx.ChildStatements() != nil {
		v.formatter.printSpace()
		v.Visit(ctx.ChildStatements())
	} else if ctx.ModuleInstantiation() != nil {
		_, parentIsIfElse := ctx.GetParent().(*parser.IfElseStatementContext)
		// determine if this statement is part of an else-if
		elseIf := ctx.ModuleInstantiation().IfElseStatement() != nil && parentIsIfElse
		// special formatting for else-if statements: don't start a new line before the "if"
		if !elseIf {
			v.formatter.endLine()
			v.formatter.indent()
		} else {
			v.formatter.printSpace()
		}
		v.Visit(ctx.ModuleInstantiation())
		if !elseIf {
			v.formatter.unindent()
		}
	} else {
		// not possible to hit this unless there's a parser bug
		zap.S().Fatalf("Invalid child statement state")
	}
	return nil
}

func (v *FormattingVisitor) VisitArguments(ctx *parser.ArgumentsContext) interface{} {
	allArgs := ctx.AllArgument()
	for i, arg := range allArgs {
		v.Visit(arg)
		if ctx.Comma(i) != nil {
			v.Visit(ctx.Comma(i))
		}
		if i < len(allArgs)-1 {
			v.formatter.printSpace()
		}
	}
	v.Visit(ctx.OptionalTrailingComma())
	return nil
}

func (v *FormattingVisitor) VisitParameters(ctx *parser.ParametersContext) interface{} {
	allParams := ctx.AllParameter()
	for i, arg := range allParams {
		v.Visit(arg)
		if ctx.Comma(i) != nil {
			v.Visit(ctx.Comma(i))
		}
		if i < len(allParams)-1 {
			v.formatter.printSpace()
		}
	}
	v.Visit(ctx.OptionalTrailingComma())
	return nil
}

func (v *FormattingVisitor) printEndOfLineCommentAfter(tokenIndex int) {
	tokenIndex = tokenIndex + 1
	if tokenIndex >= v.tokenStream.Size() {
		return
	}
	token := v.tokenStream.Get(tokenIndex)
	tokenType := token.GetTokenType()
	if tokenType == parser.OpenSCADLexerEND_OF_LINE_COMMENT ||
		tokenType == parser.OpenSCADLexerEND_OF_LINE_COMMENT_BLOCK ||
		tokenType == parser.OpenSCADLexerMULTILINE_COMMENT_BLOCK {
		v.printCommentsBefore(tokenIndex)
	}
}

func (v *FormattingVisitor) printCommentsBefore(tokenIndex int) {
	for ; v.lastPrintedCommentIndex <= tokenIndex && v.lastPrintedCommentIndex < v.tokenStream.Size(); v.lastPrintedCommentIndex++ {
		token := v.tokenStream.Get(v.lastPrintedCommentIndex)
		v.printCommentToken(token)
	}
}

func (v *FormattingVisitor) printCommentToken(token antlr.Token) {
	switch tokenType := token.GetTokenType(); tokenType {
	case parser.OpenSCADLexerEND_OF_LINE_COMMENT:
		zap.S().Debugf("Printing END_OF_LINE_COMMENT: token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printEndOfLineComment(token)
	case parser.OpenSCADLexerSINGLE_LINE_COMMENT:
		zap.S().Debugf("Printing SINGLE_LINE_COMMENT, token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printSingleLineComment(token)
	case parser.OpenSCADLexerEND_OF_LINE_COMMENT_BLOCK:
		zap.S().Debugf("Printing END_OF_LINE_COMMENT_BLOCK: token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printEndOfLineComment(token)
	case parser.OpenSCADLexerSINGLE_LINE_COMMENT_BLOCK:
		zap.S().Debugf("Printing SINGLE_LINE_COMMENT_BLOCK, token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printSingleLineComment(token)
	case parser.OpenSCADLexerMULTILINE_COMMENT_BLOCK:
		zap.S().Debugf("Printing MULTILINE_COMMENT_BLOCK, token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printMultilineComment(token)
	case parser.OpenSCADLexerMULTI_NEWLINE:
		zap.S().Debugf("Printing MULTI_NEWLINE, token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printMultiNewlineComment(token.GetText())
	default:
		zap.S().Debugf("skipping non-comment token, token index = %d [%s]", token.GetTokenIndex(), token.GetText())
	}
}

func (v *FormattingVisitor) printMultiNewlineComment(text string) {
	newLineCount := strings.Count(text, "\n")
	if newLineCount > 0 {
		v.formatter.endLine()
	}
	for i := 1; i < newLineCount; i++ {
		v.formatter.printNewLine()
	}
}

// TO DO Implement re-formatting options for multiline comments.
func (v *FormattingVisitor) printMultilineComment(token antlr.Token) error {
	lines := strings.Split(token.GetText(), "\n")
	for _, line := range lines {
		err := v.formatter.appendToLine(line, false)
		if err != nil {
			return err
		}
		err = v.formatter.printNewLine()
		if err != nil {
			return err
		}
	}
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) printSingleLineComment(token antlr.Token) {
	v.formatter.endLine()
	v.formatter.printString(strings.TrimSpace(token.GetText()))
	v.formatter.endLine()
}

func (v *FormattingVisitor) printEndOfLineComment(token antlr.Token) {
	if v.formatter.inLine {
		v.formatter.printSpace()
	}
	v.formatter.printString(strings.TrimSpace(token.GetText()))
	v.formatter.endLine()
}
