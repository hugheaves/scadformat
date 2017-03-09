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

import java.io.PrintStream;

import org.antlr.v4.runtime.CommonTokenStream;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * The Class TokenFormatter.
 */
public class TokenFormatter {

    /**
     * The logger.
     */
    private static Logger logger = LoggerFactory.getLogger(TokenFormatter.class);

    /**
     * The Constant INDENT_SIZE.
     */
    private static final int INDENT_SIZE = 2;

    /**
     * The token stream.
     */
    //
    private CommonTokenStream tokenStream;

    /**
     * The output stream.
     */
    private final PrintStream outputStream;

    /**
     * The indent.
     */
    private int indent = 0;

    /**
     * The line pos.
     */
    private int linePos = 0;

    /**
     * The in line.
     */
    private boolean inLine = false;

    /**
     * The wrapped line.
     */
    private boolean wrappedLine = false;

    /**
     * The Constant LINE_LEN.
     */
    private static final int LINE_LEN = Integer.MAX_VALUE;

    /**
     * Instantiates a new token formatter.
     *
     * @param tokenStream
     *            the token stream
     * @param outputStream
     *            the output stream
     */
    public TokenFormatter(final CommonTokenStream tokenStream, final PrintStream outputStream) {
        // this.tokenStream = tokenStream;
        this.outputStream = outputStream;
    }

    /**
     * Prints the string.
     *
     * @param string
     *            the string
     */
    public void printString(final String string) {
        final String[] lines = string.split("\\R", -1);
        for (int line = 0; line < lines.length - 1; ++line) {
            printWithLineWrap(lines[line]);
            printNewLine();
        }

        if (lines.length > 0) {
            printWithLineWrap(lines[lines.length - 1]);
        }

    }

    /**
     * Prints the with line wrap.
     *
     * @param string
     *            the string
     */
    private void printWithLineWrap(final String string) {
        if (inLine && string.length() > lineRemaining()) {
            printNewLine();
            indent();
            wrappedLine = true;
        }
        appendToLine(string);
    }

    /**
     * Prints the space.
     */
    public void printSpace() {
        if (inLine) {
            printWithLineWrap(" ");
        }
    }

    /**
     * End line.
     */
    void endLine() {
        if (inLine) {
            printNewLine();
        }
    }

    /**
     * Prints the new line.
     */
    void printNewLine() {
        outputStream.print("\n");
        if (wrappedLine) {
            unindent();
            wrappedLine = false;
        }
        inLine = false;
        linePos = 0;
    }

    /**
     * Line remaining.
     *
     * @return the int
     */
    private int lineRemaining() {
        if (!inLine) {
            return LINE_LEN - indent;
        } else {
            return LINE_LEN - linePos;
        }
    }

    /**
     * Append to line.
     *
     * @param string
     *            the string
     */
    private void appendToLine(final String string) {
        if (string == null || string.length() == 0) {
            return;
        }

        if (!inLine) {
            outputIndent();
            inLine = true;
        }
        outputStream.print(string);
        linePos += string.length();
    }

    /**
     * Output indent.
     */
    private void outputIndent() {
        for (int i = 0; i < indent; ++i) {
            outputStream.print(" ");
        }
        linePos += indent;
    }

    /**
     * Indent.
     */
    public void indent() {
        indent += INDENT_SIZE;
        // System.out.print(">");

    }

    /**
     * Unindent.
     */
    public void unindent() {
        indent -= INDENT_SIZE;
        // System.out.print("<");
    }

    /**
     * Gets the token stream.
     *
     * @return the token stream
     */
    public CommonTokenStream getTokenStream() {
        return tokenStream;
    }

}
