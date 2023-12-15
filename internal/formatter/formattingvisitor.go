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
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hugheaves/scadformat/internal/parser"
	"go.uber.org/zap"
)

// check that Visitor implements OpenSCADVisitor
var _ parser.OpenSCADVisitor = &FormattingVisitor{}

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
		return tree.Accept(v)
	}
	return nil
}

func (v *FormattingVisitor) VisitChildren(tree antlr.RuleNode) interface{} {
	zap.S().Debug("entering VisitChildren")
	for i, child := range tree.GetChildren() {
		zap.S().Debug("visiting child", "i", i)
		val := child.(antlr.ParseTree)
		_ = v.Visit(val)
	}
	zap.L().Debug("exiting VisitChildren")
	return nil
}

func (v *FormattingVisitor) VisitTerminal(node antlr.TerminalNode) interface{} {
	zap.S().Debug("Visiting TerminalNode", "node", node)
	if v.deferredCommentIndex != -1 {
		zap.L().Debug("Formatting deferred comment")
		v.processCommentTokens(v.deferredCommentIndex, true)
		v.deferredCommentIndex = -1
	}
	v.formatter.printString(node.GetText())
	v.processCommentTokens(node.GetSymbol().GetTokenIndex()+1, false)
	return nil
}

func (v *FormattingVisitor) VisitErrorNode(_ antlr.ErrorNode) interface{} {
	return nil
}

func (v *FormattingVisitor) VisitStart(ctx *parser.StartContext) interface{} {
	zap.L().Debug("entering VisitStart")
	v.processCommentTokens(0, false)
	v.VisitChildren(ctx)
	if v.deferredCommentIndex != -1 {
		v.processCommentTokens(v.deferredCommentIndex, true)
	}
	zap.L().Debug("exiting VisitStart")
	return nil
}

func (v *FormattingVisitor) VisitAssignment(ctx *parser.AssignmentContext) interface{} {
	zap.L().Debug("VisitAssignment")
	v.Visit(ctx.AssignmentExpression())
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

func (v *FormattingVisitor) VisitBinary(ctx *parser.BinaryContext) interface{} {
	zap.L().Debug("visitBinary")
	v.formatter.printSpace()
	v.Visit(ctx.BinaryOperator())
	v.formatter.printSpace()
	v.Visit(ctx.Expr())
	return nil

}

func (v *FormattingVisitor) VisitTernary(ctx *parser.TernaryContext) interface{} {
	zap.L().Debug("visitTernary")
	v.formatter.printSpace()
	v.Visit(ctx.QUESTION_MARK())
	v.formatter.printSpace()
	v.Visit(ctx.Expr(0))
	v.formatter.printSpace()
	v.Visit(ctx.COLON())
	v.formatter.printSpace()
	v.Visit(ctx.Expr(1))
	return nil

}

func (v *FormattingVisitor) VisitChildStatements(ctx *parser.ChildStatementsContext) interface{} {
	zap.L().Debug("visitChildStatements")
	v.formatter.printSpace()
	v.Visit(ctx.L_CURLY())
	v.formatter.endLine()
	v.formatter.indent()
	for _, childCtx := range ctx.GetChildren() {
		if _, ok := childCtx.(parser.IChildStatementOrAssignmentContext); ok {
			v.Visit(childCtx.(antlr.RuleContext))
		}
	}
	v.formatter.unindent()
	v.Visit(ctx.R_CURLY())
	nextToken := v.tokenStream.Get(ctx.R_CURLY().GetSymbol().GetTokenIndex() + 1)
	if nextToken.GetTokenType() != parser.OpenSCADLexerELSE {
		v.formatter.endLine()
	} else {
		v.formatter.printSpace()
	}

	return nil
}

func (v *FormattingVisitor) VisitCommas(ctx *parser.CommasContext) interface{} {
	zap.L().Debug("visitCommas")
	for _, comma := range ctx.AllCOMMA() {
		v.Visit(comma)
		v.formatter.printSpace()
	}
	return nil
}

func (v *FormattingVisitor) VisitNullStatement(ctx *parser.NullStatementContext) interface{} {
	zap.L().Debug("VisitNullStatement")
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
	v.Visit(ctx.ArgumentsDecl())
	v.Visit(ctx.R_PAREN())
	v.formatter.printSpace()
	v.Visit(ctx.EQUALS())
	v.formatter.printSpace()
	v.Visit(ctx.Expr())
	v.Visit(ctx.SEMICOLON())
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitInclude(ctx *parser.IncludeContext) interface{} {
	zap.L().Debug("visitInclude")
	v.Visit(ctx.INCLUDE())
	v.formatter.printSpace()
	v.Visit(ctx.FILE())
	v.formatter.endLine()
	return nil
}

func (v *FormattingVisitor) VisitModuleDefinition(ctx *parser.ModuleDefinitionContext) interface{} {
	zap.L().Debug("visitModuleDefinition")
	v.Visit(ctx.MODULE())
	v.formatter.printSpace()
	v.Visit(ctx.ID())
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.ArgumentsDecl())
	v.Visit(ctx.R_PAREN())
	v.Visit(ctx.Statement())
	v.Visit(ctx.SEMICOLON())
	return nil
}

func (v *FormattingVisitor) VisitModuleInstantiation(ctx *parser.ModuleInstantiationContext) interface{} {
	v.formatter.printSpace()
	v.Visit(ctx.ModifierCharacter())
	v.Visit(ctx.ModuleId())
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.CallArguments())
	v.Visit(ctx.R_PAREN())
	if ctx.ChildStatement() != nil && (ctx.ChildStatement().NullStatement() == nil && ctx.ChildStatement().ChildStatements() == nil) {
		v.formatter.indent()
		v.formatter.printNewLine()
		v.Visit(ctx.ChildStatement())
		v.Visit(ctx.SEMICOLON())
		v.formatter.unindent()
	} else {
		v.Visit(ctx.ChildStatement())
		v.Visit(ctx.SEMICOLON())
	}

	return nil
}

func (v *FormattingVisitor) VisitOptionalCommas(ctx *parser.OptionalCommasContext) interface{} {
	for _, comma := range ctx.AllCOMMA() {
		v.Visit(comma)
		v.formatter.printSpace()
	}
	return nil
}

func (v *FormattingVisitor) VisitStatements(ctx *parser.StatementsContext) interface{} {
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

func (v *FormattingVisitor) VisitUse(ctx *parser.UseContext) interface{} {
	v.VisitChildren(ctx)
	v.formatter.endLine()
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
	v.Visit(ctx.VectorArguments())
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
	return v.VisitChildren(ctx)
}

func (v *FormattingVisitor) VisitIfStatement(ctx *parser.IfStatementContext) interface{} {
	v.Visit(ctx.IF())
	v.formatter.printSpace()
	v.Visit(ctx.L_PAREN())
	v.Visit(ctx.Expr())
	v.Visit(ctx.R_PAREN())
	return nil
}

func (v *FormattingVisitor) VisitElseStatement(ctx *parser.ElseStatementContext) interface{} {
	v.VisitChildren(ctx)
	v.formatter.printSpace()
	return nil
}

// processCommentTokens processes any comment tokens in the token stream starting at
// index "startIndex".
func (v *FormattingVisitor) processCommentTokens(startIndex int, printDeferred bool) {
	zap.S().Debug("format comments", "startIndex", startIndex)

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
				zap.S().Debug("deferring comment print", "i", i)
				v.deferredCommentIndex = i
				break loop
			} else {
				v.printComment(token)
			}
		case parser.OpenSCADLexerMULTILINE_COMMENT:
			v.printComment(token)
		case parser.OpenSCADLexerMULTI_NEWLINE:
			v.printMultiNewlineComment(token.GetText())
		default:
			zap.S().Debug("Exiting printComments", "i", i, "text", token.GetText())
			break loop
		}
	}
}

// func (v *FormattingVisitor) printDeferredComment() {
// 	if v.deferredCommentIndex != -1 {
// 		token := v.tokenStream.Get(v.deferredCommentIndex)
// 		tokenType := token.GetTokenType()
// 		// this should never happen
// 		if tokenType != parser.OpenSCADLexerSINGLE_LINE_COMMENT {
// 			panic("comment was not a single line comment")
// 		}
// 		v.printComment(token)
// 		v.deferredCommentIndex = -1
// 	}
// }

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
