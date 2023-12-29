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
	tokenStream          antlr.TokenStream
	formatter            *TokenFormatter
	deferredCommentIndex int
	arrayDepth           int // stores nested array element "depth" during array formatting
}

func NewFormattingVisitor(tokenStream antlr.TokenStream, formatter *TokenFormatter) *FormattingVisitor {
	visitor := &FormattingVisitor{
		tokenStream:          tokenStream,
		formatter:            formatter,
		deferredCommentIndex: -1,
		arrayDepth:           0,
	}

	// Override VisitChildren in BaseOpenClassVisitor
	// as part of workaround for https://github.com/antlr/antlr4/issues/2504
	visitor.BaseOpenSCADVisitor.VisitChildren = visitor.VisitChildren

	return visitor
}

func (v *FormattingVisitor) Visit(tree antlr.ParseTree) interface{} {
	if tree != nil {
		zap.S().Debugf("Visiting: %s", reflect.TypeOf(tree).String())
		val := tree.Accept(v)
		zap.S().Debugf("Visited: %s", reflect.TypeOf(tree).String())
		return val
	}
	return nil
}

func (v *FormattingVisitor) VisitChildren(tree antlr.RuleNode) interface{} {
	for i, child := range tree.GetChildren() {
		zap.S().Debugf("VisitChildren: visiting child %d of %s", i, reflect.TypeOf(tree).String())
		val := child.(antlr.ParseTree)
		_ = v.Visit(val)
	}
	return nil
}

func (v *FormattingVisitor) VisitTerminal(node antlr.TerminalNode) interface{} {
	zap.S().Debugf("Visiting TerminalNode: %s", node)
	v.processCommentTokens(0, true)
	v.formatter.printString(node.GetText())
	v.processCommentTokens(node.GetSymbol().GetTokenIndex()+1, false)
	return nil
}

func (v *FormattingVisitor) VisitErrorNode(_ antlr.ErrorNode) interface{} {
	panic("visited ErrorNode")
}

func (v *FormattingVisitor) VisitStart(ctx *parser.StartContext) interface{} {
	zap.L().Debug("entering VisitStart")
	v.processCommentTokens(0, false)
	v.VisitChildren(ctx)
	v.processCommentTokens(0, true)
	zap.L().Debug("exiting VisitStart")
	return nil
}

func (v *FormattingVisitor) VisitAssignment(ctx *parser.AssignmentContext) interface{} {
	zap.L().Debug("VisitAssignment")
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
	zap.L().Debug("visitAssignmentExpression")
	v.Visit(ctx.ID())
	v.formatter.printSpace()
	v.Visit(ctx.EQUALS())
	v.formatter.printSpace()
	v.Visit(ctx.Expr())
	return nil
}

func (v *FormattingVisitor) VisitBinaryExpr(ctx *parser.BinaryExprContext) interface{} {
	zap.L().Debug("visitBinary")
	v.Visit(ctx.Expr(0))
	v.formatter.printSpace()
	v.Visit(ctx.BinaryOperator())
	v.formatter.printSpace()
	v.Visit(ctx.Expr(1))
	return nil
}

func (v *FormattingVisitor) VisitTernaryExpr(ctx *parser.TernaryExprContext) interface{} {
	zap.L().Debug("visitTernaryExprContext")
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

func (v *FormattingVisitor) formatChildStatement(ctx parser.IChildStatementContext, afterElse bool) interface{} {
	if ctx.Semicolon() != nil {
		v.Visit(ctx.Semicolon())
	} else if ctx.ChildStatements() != nil {
		v.formatter.printSpace()
		v.Visit(ctx.ChildStatements())
	} else if ctx.ModuleInstantiation() != nil {
		continueLine := afterElse && ctx.ModuleInstantiation().IfElseStatement() != nil
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
		panic("unexpected state")
	}
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
	zap.L().Debug("VisitSemicolon")
	v.VisitChildren(ctx)
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitFunctionDefinition(ctx *parser.FunctionDefinitionContext) interface{} {
	zap.L().Debug("visitFunctionDefinition")
	v.Visit(ctx.FUNCTION())
	v.formatter.printSpace()
	v.Visit(ctx.ID())
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Parameters())
	v.Visit(ctx.R_PAREN())
	v.formatter.printSpace()
	v.Visit(ctx.EQUALS())
	v.formatter.printSpace()
	v.Visit(ctx.Expr())
	v.Visit(ctx.SEMICOLON())
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitModuleDefinition(ctx *parser.ModuleDefinitionContext) interface{} {
	zap.L().Debug("visitModuleDefinition")
	v.Visit(ctx.MODULE())
	v.formatter.printSpace()
	v.Visit(ctx.ID())
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Parameters())
	v.Visit(ctx.R_PAREN())
	v.Visit(ctx.Statement())
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

func (v *FormattingVisitor) VisitIncludeOrUseFile(ctx *parser.IncludeOrUseFileContext) interface{} {
	matches := includeOrUseRegex.FindStringSubmatch(ctx.INCLUDE_OR_USE_FILE().GetText())
	if len(matches) != 3 {
		panic("uh oh")
	}
	v.processCommentTokens(0, true)
	v.formatter.printString(matches[1])
	v.formatter.printSpace()
	v.formatter.printString(matches[2])
	v.formatter.endLine()
	v.processCommentTokens(ctx.GetStart().GetTokenIndex()+1, false)
	return nil
}

func (v *FormattingVisitor) VisitArrayAccess(ctx *parser.ArrayAccessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *FormattingVisitor) VisitVector(ctx *parser.VectorContext) interface{} {
	zap.L().Debug("Visiting vector")

	v.arrayDepth++

	if v.arrayDepth > 1 {
		v.formatter.endLine()
		v.formatter.indent()
	}
	v.Visit(ctx.L_BRACKET())
	v.Visit(ctx.VectorElements())
	v.Visit(ctx.R_BRACKET())
	if v.arrayDepth > 1 {
		v.formatter.unindent()
		nextToken := v.tokenStream.Get(ctx.R_BRACKET().GetSymbol().GetTokenIndex() + 1)
		if nextToken.GetTokenType() != parser.OpenSCADLexerCOMMA {
			v.formatter.endLine()
		}
	}

	v.arrayDepth--

	return nil
}

func (v *FormattingVisitor) VisitForStatement(ctx *parser.ForStatementContext) interface{} {
	v.formatter.printSpace()
	return v.VisitChildren(ctx)
}

func (v *FormattingVisitor) VisitLetStatement(ctx *parser.LetStatementContext) interface{} {
	v.formatter.printSpace()
	return v.VisitChildren(ctx)
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

// processCommentTokens processes any comment tokens in the token stream starting at
// index "startIndex".
func (v *FormattingVisitor) processCommentTokens(startIndex int, printDeferred bool) {

	if printDeferred {
		startIndex = v.deferredCommentIndex
		v.deferredCommentIndex = -1
	}

	if startIndex == -1 {
		return
	}
loop:
	for i := startIndex; ; i++ {
		token := v.tokenStream.Get(i)
		if token == nil {
			break loop
		}
		switch tokenType := token.GetTokenType(); tokenType {
		case parser.OpenSCADLexerSINGLE_LINE_COMMENT:
			text := token.GetText()
			if (strings.Contains(text, "\n") || strings.Contains(text, "\r")) && !printDeferred {
				zap.S().Debugf("Deferring print of single line comment, token index = %d, text=[%s]", i, token.GetText())
				v.deferredCommentIndex = i
				break loop
			} else {
				zap.S().Debugf("Printing single line comment, token index = %d, text=[%s]", i, token.GetText())
				v.printComment(token)
			}
		case parser.OpenSCADLexerMULTILINE_COMMENT:
			zap.S().Debugf("Printing multiline comment, token index = %d, text=[%s]", i, token.GetText())
			v.printComment(token)
		case parser.OpenSCADLexerMULTI_NEWLINE:
			zap.S().Debugf("Printing multinewline comment, token index = %d, text=[%s]", i, token.GetText())
			v.printMultiNewlineComment(token.GetText())
		default:
			break loop
		}
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

func (v *FormattingVisitor) printComment(token antlr.Token) {
	var commentText string
	if token.GetTokenType() == parser.OpenSCADLexerSINGLE_LINE_COMMENT {
		commentText = strings.TrimSpace(token.GetText())
	} else {
		commentText = token.GetText()
	}
	v.formatter.printString(commentText)
	if token.GetTokenType() == parser.OpenSCADLexerSINGLE_LINE_COMMENT || token.GetTokenType() == parser.OpenSCADLexerMULTILINE_COMMENT {
		v.formatter.endLine()
	}
}
