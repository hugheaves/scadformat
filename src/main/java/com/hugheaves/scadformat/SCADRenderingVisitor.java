/*
 * Copyright (C) 2017  Hugh Eaves
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package com.hugheaves.scadformat;

import java.util.List;

import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.RuleContext;
import org.antlr.v4.runtime.Token;
import org.antlr.v4.runtime.tree.ParseTree;
import org.antlr.v4.runtime.tree.TerminalNode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.hugheaves.scadformat.antlr.SCADBaseVisitor;
import com.hugheaves.scadformat.antlr.SCADLexer;
import com.hugheaves.scadformat.antlr.SCADParser.AssignmentContext;
import com.hugheaves.scadformat.antlr.SCADParser.AssignmentExpressionContext;
import com.hugheaves.scadformat.antlr.SCADParser.BinaryContext;
import com.hugheaves.scadformat.antlr.SCADParser.ChildStatementsContext;
import com.hugheaves.scadformat.antlr.SCADParser.CommasContext;
import com.hugheaves.scadformat.antlr.SCADParser.EmptyStatementContext;
import com.hugheaves.scadformat.antlr.SCADParser.FunctionDefinitionContext;
import com.hugheaves.scadformat.antlr.SCADParser.IncludeContext;
import com.hugheaves.scadformat.antlr.SCADParser.ModuleDefinitionContext;
import com.hugheaves.scadformat.antlr.SCADParser.ModuleInstantiationContext;
import com.hugheaves.scadformat.antlr.SCADParser.OptionalCommasContext;
import com.hugheaves.scadformat.antlr.SCADParser.StartContext;
import com.hugheaves.scadformat.antlr.SCADParser.StatementsContext;
import com.hugheaves.scadformat.antlr.SCADParser.TernaryContext;
import com.hugheaves.scadformat.antlr.SCADParser.UseContext;

/**
 * The Class SCADRenderingVisitor.
 */
public class SCADRenderingVisitor extends SCADBaseVisitor<Boolean> {

    /**
     * The logger.
     */
    private static Logger logger = LoggerFactory.getLogger(SCADRenderingVisitor.class);

    /**
     * The token stream.
     */
    private final CommonTokenStream tokenStream;

    /**
     * The formatter.
     */
    private final TokenFormatter formatter;

    /**
     * The deferred comment index.
     */
    private int deferredCommentIndex = -1;

    /**
     * Instantiates a new SCAD rendering visitor.
     *
     * @param tokenStream
     *            the token stream
     * @param formatter
     *            the formatter
     */
    public SCADRenderingVisitor(final CommonTokenStream tokenStream, final TokenFormatter formatter) {
        this.tokenStream = tokenStream;
        this.formatter = formatter;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitStart(scadformat.antlr.SCADParser.
     * StartContext)
     */
    @Override
    public Boolean visitStart(final StartContext ctx) {
        formatComments(0);
        visitChildren(ctx);
        if (deferredCommentIndex != -1) {
            formatComments(deferredCommentIndex);
        }
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see scadformat.antlr.SCADBaseVisitor#visitAssignment(scadformat.antlr.
     * SCADParser.AssignmentContext)
     */
    @Override
    public Boolean visitAssignment(final AssignmentContext ctx) {
        visit(ctx.assignmentExpression());
        visit(ctx.SEMICOLON());
        formatter.endLine();
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitAssignmentExpression(scadformat.
     * antlr.SCADParser.AssignmentExpressionContext)
     */
    @Override
    public Boolean visitAssignmentExpression(final AssignmentExpressionContext ctx) {
        visit(ctx.ID());
        formatter.printSpace();
        visit(ctx.EQUALS());
        formatter.printSpace();
        visit(ctx.expr());
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitBinary(scadformat.antlr.SCADParser.
     * BinaryContext)
     */
    @Override
    public Boolean visitBinary(final BinaryContext ctx) {
        formatter.printSpace();
        visit(ctx.binaryOperator());
        formatter.printSpace();
        visit(ctx.expr());
        return null;

    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitTernary(scadformat.antlr.SCADParser
     * .TernaryContext)
     */
    @Override
    public Boolean visitTernary(final TernaryContext ctx) {
        formatter.printSpace();
        visit(ctx.QUESTION_MARK());
        formatter.printSpace();
        visit(ctx.expr(0));
        formatter.printSpace();
        visit(ctx.COLON());
        formatter.printSpace();
        visit(ctx.expr(1));
        return null;

    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitChildStatements(scadformat.antlr.
     * SCADParser.ChildStatementsContext)
     */
    @Override
    public Boolean visitChildStatements(final ChildStatementsContext ctx) {
        formatter.printSpace();
        visit(ctx.L_CURLY());
        formatter.endLine();
        formatter.indent();
        listVisit(ctx.childStatementOrAssignment());
        formatter.unindent();
        visit(ctx.R_CURLY());
        formatter.endLine();
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitCommas(scadformat.antlr.SCADParser.
     * CommasContext)
     */
    @Override
    public Boolean visitCommas(final CommasContext ctx) {
        for (final TerminalNode node : ctx.COMMA()) {
            visit(node);
            formatter.printSpace();
        }
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitEmptyStatement(scadformat.antlr.
     * SCADParser.EmptyStatementContext)
     */
    @Override
    public Boolean visitEmptyStatement(final EmptyStatementContext ctx) {
        visitChildren(ctx);
        formatter.endLine();
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitFunctionDefinition(scadformat.antlr
     * .SCADParser.FunctionDefinitionContext)
     */
    @Override
    public Boolean visitFunctionDefinition(final FunctionDefinitionContext ctx) {
        visit(ctx.FUNCTION());
        formatter.printSpace();
        visit(ctx.ID());
        visit(ctx.L_PAREN());
        visit(ctx.argumentsDecl());
        visit(ctx.R_PAREN());
        formatter.printSpace();
        visit(ctx.EQUALS());
        formatter.printSpace();
        visit(ctx.expr());
        visit(ctx.SEMICOLON());
        formatter.endLine();
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitInclude(scadformat.antlr.SCADParser
     * .IncludeContext)
     */
    @Override
    public Boolean visitInclude(final IncludeContext ctx) {
        visitChildren(ctx);
        formatter.endLine();
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitModuleDefinition(scadformat.antlr.
     * SCADParser.ModuleDefinitionContext)
     */
    @Override
    public Boolean visitModuleDefinition(final ModuleDefinitionContext ctx) {
        visit(ctx.MODULE());
        formatter.printSpace();
        visit(ctx.ID());
        visit(ctx.L_PAREN());
        visit(ctx.argumentsDecl());
        visit(ctx.R_PAREN());
        visit(ctx.statement());
        visit(ctx.SEMICOLON());
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitModuleInstantiation(scadformat.
     * antlr.SCADParser.ModuleInstantiationContext)
     */
    @Override
    public Boolean visitModuleInstantiation(final ModuleInstantiationContext ctx) {
        formatter.printSpace();
        visit(ctx.modifierCharacter());
        visit(ctx.moduleId());
        visit(ctx.L_PAREN());
        visit(ctx.callArguments());
        visit(ctx.R_PAREN());
        if (ctx.childStatement() != null
                && (ctx.childStatement().emptyStatement() == null && ctx.childStatement().childStatements() == null)) {
            logger.info("CS : {}", ctx.childStatement().getText());
            formatter.indent();
            formatter.printNewLine();
            visit(ctx.childStatement());
            visit(ctx.SEMICOLON());
            formatter.unindent();
        } else {
            visit(ctx.childStatement());
            visit(ctx.SEMICOLON());
        }
        if (ctx.SEMICOLON() != null) {
            logger.info("Gotta semicolon");
        }

        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitOptionalCommas(scadformat.antlr.
     * SCADParser.OptionalCommasContext)
     */
    @Override
    public Boolean visitOptionalCommas(final OptionalCommasContext ctx) {
        for (final TerminalNode node : ctx.COMMA()) {
            visit(node);
            formatter.printSpace();
        }
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see scadformat.antlr.SCADBaseVisitor#visitStatements(scadformat.antlr.
     * SCADParser.StatementsContext)
     */
    @Override
    public Boolean visitStatements(final StatementsContext ctx) {
        formatter.printSpace();
        visit(ctx.L_CURLY());
        formatter.endLine();
        formatter.indent();
        listVisit(ctx.statement());
        formatter.unindent();
        visit(ctx.R_CURLY());
        formatter.endLine();
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * scadformat.antlr.SCADBaseVisitor#visitUse(scadformat.antlr.SCADParser.
     * UseContext)
     */
    @Override
    public Boolean visitUse(final UseContext ctx) {
        visitChildren(ctx);
        formatter.endLine();
        return null;
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * org.antlr.v4.runtime.tree.AbstractParseTreeVisitor#visitTerminal(org.
     * antlr.v4.runtime.tree.TerminalNode)
     */
    @Override
    public Boolean visitTerminal(final TerminalNode node) {
        if (deferredCommentIndex != -1) {
            formatComments(deferredCommentIndex);
            deferredCommentIndex = -1;
        }
        formatter.printString(node.getText());
        formatComments(node.getSymbol().getTokenIndex() + 1);
        return null;
    }

    /**
     * Format comments.
     *
     * @param startIndex
     *            the start index
     */
    public void formatComments(final int startIndex) {
        int i = startIndex;
        for (Token token = tokenStream.get(i); token != null; token = tokenStream.get(++i)) {
            if (token.getType() == SCADLexer.LINE_COMMENT) {
                final String text = token.getText();
                /*
                 * If line comment starts with a newline, then we need to print
                 * it at the level of indention used for the token _after_ the
                 * comment. So, instead of printing it now (which would use the
                 * indention level for the current token) we wait to print it
                 * until the next token is about to be renderered.
                 */
                if ((text.startsWith("\n") || text.startsWith("\r")) && deferredCommentIndex == -1) {
                    deferredCommentIndex = startIndex;
                    break;
                } else if (deferredCommentIndex != -1) {
                    formatter.printString(token.getText().substring(1));
                    deferredCommentIndex = -1;
                } else {
                    formatter.printSpace();
                    formatter.printString(token.getText());
                }
            } else if (token.getType() == SCADLexer.COMMENT) {
                formatter.printString(token.getText());
            } else if (token.getType() == SCADLexer.MULTI_NEWLINE) {
                formatter.printString(token.getText());
            } else {
                break;
            }
        }

    }

    /**
     * List visit.
     *
     * @param ctxList
     *            the ctx list
     */
    public void listVisit(final List<? extends RuleContext> ctxList) {
        for (final RuleContext ctx : ctxList) {
            visit(ctx);
        }
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * org.antlr.v4.runtime.tree.AbstractParseTreeVisitor#visit(org.antlr.v4.
     * runtime.tree.ParseTree)
     */
    @Override
    public Boolean visit(final ParseTree tree) {
        if (tree != null) {
            return tree.accept(this);
        }
        return null;
    }
}
