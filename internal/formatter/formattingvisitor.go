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
	v.VisitChildren(ctx)
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
	v.formatter.indent()
	for _, child := range ctx.AllChildStatementOrAssignment() {
		v.Visit(child)
	}
	v.formatter.unindent()
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
	v.Visit(ctx.Statement())
	return nil
}

func (v *FormattingVisitor) VisitAssertExpr(ctx *parser.AssertExprContext) interface{} {
	v.Visit(ctx.ASSERT())
	v.formatter.printSpace()
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Arguments())
	v.Visit(ctx.R_PAREN())
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
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Arguments())
	v.Visit(ctx.R_PAREN())
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
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Arguments())
	v.Visit(ctx.R_PAREN())
	v.formatter.endLine()
	v.formatter.indent()
	v.Visit(ctx.Expr())
	v.formatter.unindent()
	return nil
}

func (v *FormattingVisitor) VisitSingleModuleInstantiation(ctx *parser.SingleModuleInstantiationContext) interface{} {
	v.Visit(ctx.ModuleId())
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Arguments())
	v.Visit(ctx.R_PAREN())
	v.formatChildStatement(ctx.ChildStatement(), false)
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitInnerInput(ctx *parser.InnerInputContext) interface{} {
	v.formatter.printSpace()
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
	nested := false
	for _, ve := range ctx.AllVectorElement() {
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
	for i, ve := range ctx.AllVectorElement() {
		v.Visit(ve)
		if ctx.Comma(i) != nil {
			v.Visit(ctx.Comma(i))
		}
		if nested {
			v.formatter.endLine()
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
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Arguments())
	v.Visit(ctx.R_PAREN())
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
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Expr())
	v.Visit(ctx.R_PAREN())
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
		v.formatChildStatement(ctx.ChildStatement(), true)
	}
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitIfStatement(ctx *parser.IfStatementContext) interface{} {
	v.Visit(ctx.IF())
	v.formatter.printSpace()
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Expr())
	v.Visit(ctx.R_PAREN())
	v.formatChildStatement(ctx.ChildStatement(), false)
	return nil
}

func (v *FormattingVisitor) VisitComma(ctx *parser.CommaContext) interface{} {
	v.Visit(ctx.COMMA())
	v.formatter.printSpace()
	return nil
}

// Formats a childStatement, which can be a semicolon, a singleModuleInstantiation, or a childStatements node.
// isElse is set to true if this childStatement is the else part of an "if else" statement.
func (v *FormattingVisitor) formatChildStatement(ctx parser.IChildStatementContext, isElse bool) interface{} {
	if ctx.Semicolon() != nil {
		v.Visit(ctx.Semicolon())
	} else if ctx.ChildStatements() != nil {
		v.formatter.printSpace()
		v.Visit(ctx.ChildStatements())
	} else if ctx.ModuleInstantiation() != nil {
		continueLine := isElse && ctx.ModuleInstantiation().IfElseStatement() != nil
		if continueLine {
			v.formatter.printSpace()
		} else {
			v.formatter.endLine()
			v.formatter.indent()
		}
		v.Visit(ctx.ModuleInstantiation())
		if !continueLine {
			v.formatter.unindent()
		}
	} else {
		// not possible to hit this unless there's a parser bug
		zap.S().Fatalf("Invalid child statement state")
	}
	return nil
}

func (v *FormattingVisitor) printEndOfLineCommentAfter(tokenIndex int) {
	tokenIndex = tokenIndex + 1
	if tokenIndex >= v.tokenStream.Size() {
		return
	}
	token := v.tokenStream.Get(tokenIndex)
	tokenType := token.GetTokenType()
	if tokenType == parser.OpenSCADLexerEND_OF_LINE_COMMENT {
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
	case parser.OpenSCADLexerMULTILINE_COMMENT:
		zap.S().Debugf("Printing MULTILINE_COMMENT, token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printMultilineComment(token)
	case parser.OpenSCADLexerMULTI_NEWLINE:
		zap.S().Debugf("Printing rMULTI_NEWLINE, token index = %d, text=[%s]", token.GetTokenIndex(), token.GetText())
		v.printMultiNewlineComment(token.GetText())
	default:
		zap.S().Debugf("skipping non-comment token, token index = %d", token.GetTokenIndex())
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

func (v *FormattingVisitor) printMultilineComment(token antlr.Token) error {
	v.formatter.endLine()
	strVal := strings.TrimSpace(token.GetText())
	lines := strings.Split(strVal, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		// if i > 0 && i < len(lines)-1 {
		// 	if !strings.HasPrefix(line, "*") {
		// 		v.formatter.printString(" ")
		// 	}
		// }
		err := v.formatter.printWithLineWrap(line)
		if err != nil {
			return err
		}
		if i < len(lines)-1 {
			err = v.formatter.printNewLine()
			if err != nil {
				return err
			}
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
