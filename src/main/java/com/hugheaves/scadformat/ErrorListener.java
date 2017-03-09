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

import org.antlr.v4.runtime.BaseErrorListener;
import org.antlr.v4.runtime.RecognitionException;
import org.antlr.v4.runtime.Recognizer;
import org.antlr.v4.runtime.misc.ParseCancellationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * The listener interface for receiving error events.
 * The class that is interested in processing a error
 * event implements this interface, and the object created
 * with that class is registered with a component using the
 * component's <code>addErrorListener<code> method. When
 * the error event occurs, that object's appropriate
 * method is invoked.
 *
 * @see ErrorEvent
 */
public class ErrorListener extends BaseErrorListener {
    
    /**
     * The logger.
     */
    private static Logger logger = LoggerFactory.getLogger(ErrorListener.class);

    /* (non-Javadoc)
     * @see org.antlr.v4.runtime.BaseErrorListener#syntaxError(org.antlr.v4.runtime.Recognizer, java.lang.Object, int, int, java.lang.String, org.antlr.v4.runtime.RecognitionException)
     */
    @Override
    public void syntaxError(final Recognizer<?, ?> recognizer, final Object offendingSymbol, final int line,
            final int charPositionInLine, final String msg, final RecognitionException e) {
        throw new ParseCancellationException(String.format("line %d:%d %s", line, charPositionInLine, msg));
    }

}
