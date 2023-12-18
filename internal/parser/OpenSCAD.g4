/*
 * NOTE: This ANTLR grammer was adapted from the original OpenSCAD grammer downloaded from
 * https://github.com/openscad/openscad/blob/944b83cbce81a63a53ce3c615c006e5eeab27f04/src/parser.y
 *
 * The copyright on the OpenSCAD parser.y file is reproduced here:
 *
 */

/*
 *  OpenSCAD (www.openscad.org)
 *  Copyright (C) 2009-2011 Clifford Wolf <clifford@clifford.at> and
 *                          Marius Kintel <marius@kintel.net>
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  As a special exception, you have permission to link this program
 *  with the CGAL library and distribute executables, as long as you
 *  follow the requirements of the GNU GPL in regard to all of the
 *  software in the executable aside from CGAL.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program; if not, write to the Free Software
 *  Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 *
*/

grammar OpenSCAD;

start
:
    input
;

input
:
    (
        use
        | include
        | statement
    )*
;

include
:
    INCLUDE FILE
;

use
:
    USE FILE
;

statement
:
    nullStatement // is semicolon
    | assignment // ends in semicolon
    | statements
    | moduleDefinition // ends in optional semicolon
    | functionDefinition // ends in semicolon
    | moduleInstantiation
    | ifElseStatement
;

statements
:
    L_CURLY statement* R_CURLY
;

nullStatement
:
    SEMICOLON
;

assignmentExpression
:
    ID EQUALS expr
;

assignment
:
    assignmentExpression SEMICOLON
;

modifierCharacter
:
    '!'
    | '#'
    | '%'
    | '*'
;

moduleDefinition
:
    MODULE ID L_PAREN argumentsDecl R_PAREN statement SEMICOLON?
;

functionDefinition
:
    FUNCTION ID L_PAREN argumentsDecl R_PAREN EQUALS expr SEMICOLON
;

moduleInstantiation
:
    modifierCharacter? moduleId L_PAREN callArguments R_PAREN childStatement
    SEMICOLON?
;

ifElseStatement
:
    ifStatement childStatement
    (
        elseStatement childStatement
    )?
;

ifStatement
:
    IF L_PAREN expr R_PAREN
;

elseStatement
:
    ELSE
;

childStatement
:
    nullStatement
    | moduleInstantiation
    | ifElseStatement
    | childStatements
;

childStatements
:
    L_CURLY childStatementOrAssignment* R_CURLY
;

childStatementOrAssignment
:
    childStatement
    | assignment
;

moduleId
:
    ID
    | FOR
    | LET
    | ASSERT
    | ECHO
    | EACH
;

expr
:
    TRUE leftExpr
    | FALSE leftExpr
    | UNDEF leftExpr
    | ID leftExpr
    | STRING leftExpr
    | NUMBER leftExpr
    | functionCall leftExpr
    | letExpr leftExpr
    | assertion leftExpr
    | echo leftExpr
    | parenthetical leftExpr
    | range leftExpr
    | vector leftExpr
    | unary leftExpr
;

leftExpr
:
/* empty */
    | memberAccess leftExpr
    | binary leftExpr
    | ternary leftExpr
    | arrayAccess leftExpr
;

binary
:
    binaryOperator expr
;

binaryOperator
:
    '*'
    | '/'
    | '%'
    | PLUS
    | MINUS
    | '<'
    | LE
    | EQ
    | NE
    | GE
    | '>'
    | AND
    | OR
;

memberAccess
:
    '.' ID
;

arrayAccess
:
    L_BRACKET expr R_BRACKET
;

ternary
:
    QUESTION_MARK expr COLON expr
;

unary
:
    (
        '!'
        | MINUS
        | PLUS
    ) expr
;

parenthetical
:
    L_PAREN expr R_PAREN
;

vector
:
    L_BRACKET vectorArguments? R_BRACKET
;

vectorArguments
:
    (
        expr
        | listComprehension
    )
    (
        commas childListComprehensionOrExpr
    )*
    optionalCommas
;

functionCall
:
    ID L_PAREN callArguments R_PAREN
;

letStatement
:
    LET L_PAREN callArguments R_PAREN
;

letExpr
:
    letStatement expr
;

assertion
:
    ASSERT L_PAREN callArguments R_PAREN expr?
;

echo
:
    ECHO L_PAREN callArguments R_PAREN expr?
;

range
:
    L_BRACKET expr COLON expr
    (
        COLON expr
    )? R_BRACKET
;

eachStatement
:
    EACH
;

forStatement
:
    FOR L_PAREN callArguments
    (
        SEMICOLON expr SEMICOLON callArguments
    )? R_PAREN
;

ifElseStatementComprehension
:
    ifStatement childListComprehensionOrExpr
    (
        elseStatement childListComprehensionOrExpr
    )?
;

listComprehension
:
    letStatement childListComprehension
    | eachStatement childListComprehensionOrExpr
    | forStatement childListComprehensionOrExpr
    | ifElseStatementComprehension
;

childListComprehension
:
    listComprehension
    | L_PAREN listComprehension R_PAREN
;

childListComprehensionOrExpr
:
    childListComprehension
    | expr
;

argumentsDecl
:
/* empty */
    | optionalCommas argumentDecl
    (
        commas argumentDecl
    )* optionalCommas
;

argumentDecl
:
    ID
    | assignmentExpression
;

callArguments
:
/* empty */
    | optionalCommas callArgument
    (
        commas callArgument
    )*
;

callArgument
:
    expr
    | assignmentExpression
;

commas
:
    COMMA+
;

optionalCommas
:
    COMMA*
;

EQUALS
:
    '='
;

SEMICOLON
:
    ';'
;

COLON
:
    ':'
;

COMMA
:
    ','
;

L_CURLY
:
    '{'
;

R_CURLY
:
    '}'
;

L_PAREN
:
    '('
;

R_PAREN
:
    ')'
;

L_BRACKET
:
    '['
;

R_BRACKET
:
    ']'
;

QUESTION_MARK
:
    '?'
;

GE
:
    '>='
;

EQ
:
    '=='
;

NE
:
    '!='
;

LE
:
    '<='
;

AND
:
    '&&'
;

OR
:
    '||'
;

LET
:
    'let'
;

FOR
:
    'for'
;

IF
:
    'if'
;

INCLUDE
:
    'include'
;

USE
:
    'use'
;

TRUE
:
    'true'
;

FALSE
:
    'false'
;

UNDEF
:
    'undef'
;

ELSE
:
    'else'
;

ASSERT
:
    'assert'
;

ECHO
:
    'echo'
;

EACH
:
    'each'
;

FUNCTION
:
    'function'
;

MODULE
:
    'module'
;

ID
:
    '$'?
    (
        LETTER
        | DIGIT
        | UNDERSCORE
    )+
;

NUMBER
:
    (
        FLOAT
        | INTEGER
    )
;

FLOAT
:
    DIGIT+ '.' DIGIT+
;

INTEGER
:
    DIGIT+
;

PLUS
:
    '+'
;

MINUS
:
    '-'
;

STRING
:
    '"' STRING_CHAR* '"'
;

FILE
:
    '<' ~[\t\r\n>]* '>'
;

// Note: All comments are sent to channel #2
SINGLE_LINE_COMMENT
:
    EOL? SPACES* '//' ~[\r\n]* -> channel ( 2 )
;

// Multiple newlines are preserved as "comments"
MULTI_NEWLINE
:
    SPACES_EOL_SPACES EOL_SPACES+ -> channel ( 2 )
;

MULTILINE_COMMENT
:
     '/*' .*? '*/' -> channel ( 2 )
;

WHITESPACE
:
    [ \t\r\n] -> skip
;

fragment SPACES_EOL_SPACES
:
    [ \t]* EOL [ \t]*
;

fragment EOL_SPACES
:
    EOL [ \t]*
;

fragment SPACES
:
    [ \t]
;


fragment EOL
:
    '\r' '\n'
    | '\n'
;

fragment
STRING_CHAR
:
    ~["\\]
    | ESCAPE_SEQUENCE
;

fragment
ESCAPE_SEQUENCE
:
    '\\' ["\\]
    | OCTAL_ESCAPE_SEQUENCE
    | UNICODE_ESCAPE_SEQUENCE
;

fragment
OCTAL_ESCAPE_SEQUENCE
:
    '\\' OCTAL_DIGIT+
;

fragment
UNICODE_ESCAPE_SEQUENCE
:
    '\\' 'u' HEX_DIGIT+
;

fragment
FILE_CHAR
:
    [a-zA-Z./]
;

fragment
LETTER
:
    [a-zA-Z]
;

fragment
UPPERCASE_LETTER
:
    [A-Z]
;

fragment
UNDERSCORE
:
    '_'
;

fragment
DIGIT
:
    [0-9]
;

fragment
OCTAL_DIGIT
:
    [0-7]
;

fragment
HEX_DIGIT
:
    [0-9a-fA-F]
;



