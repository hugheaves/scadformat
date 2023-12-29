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

start: input;

input: ( includeOrUseFile | statement)*;

includeOrUseFile: INCLUDE_OR_USE_FILE;

/*
Equivalent from parser.y:
------------------
statement
        : 
        ';'
        | '{' inner_input '}'
        | module_instantiation
        | assignment
        | TOK_MODULE TOK_ID '(' parameters ')' statement
        | TOK_FUNCTION TOK_ID '(' parameters ')' '=' expr ';'
        | TOK_EOT
        ;
------------------
*/
statement:
	semicolon // is semicolon
	| innerInput
	| moduleInstantiation
	| assignment // ends in semicolon
	| moduleDefinition // no semicolon
	| functionDefinition; // ends in semicolon

moduleDefinition:
	MODULE ID L_PAREN parameters R_PAREN statement;

functionDefinition:
	FUNCTION ID L_PAREN parameters R_PAREN EQUALS expr SEMICOLON;

/*
Equivalent from parser.y:
------------------
inner_input
        : /x empty x/
        | inner_input statement
        ;
------------------
*/
innerInput: L_CURLY statement* R_CURLY;

/*
Equivalent from parser.y:
------------------
assignment
        : TOK_ID '=' expr ';'
        ;
------------------
*/
assignment: ID EQUALS expr SEMICOLON;

/*
Equivalent from parser.y:
------------------
module_instantiation
        : '!' module_instantiation
        | '#' module_instantiation
        | '%' module_instantiation
        | '*' module_instantiation
        | single_module_instantiation child_statement
        | ifelse_statement
        ;
------------------
*/
moduleInstantiation:
	modifierCharacters (
		singleModuleInstantiation // Note: childStatement moved to singleModuleInstantiation
		| ifElseStatement
	);

/*
Equivalent from parser.y:
------------------
ifelse_statement
        : if_statement %prec NO_ELSE
        | if_statement TOK_ELSE child_statement
        ;
------------------
*/
ifElseStatement: ifStatement | ifStatement ELSE childStatement;

/*
Equivalent from parser.y:
------------------
if_statement
        : TOK_IF '(' expr ')' child_statement
        ;
------------------
*/
ifStatement: IF L_PAREN expr R_PAREN childStatement;

/*
Equivalent from parser.y:
------------------
child_statement
        : ';'
        | '{' child_statements '}'
        | module_instantiation
        ;
------------------
*/
childStatement:
	semicolon
	| childStatements
	| moduleInstantiation;

/*
Equivalent from parser.y:
------------------
child_statements
        : /x empty x/
        | child_statements child_statement
        | child_statements assignment
        ;
------------------
*/
childStatements: L_CURLY childStatementOrAssignment* R_CURLY;

childStatementOrAssignment: (childStatement | assignment);

/*
Equivalent from parser.y:
------------------
// "for", "let" and "each" are valid module identifiers
module_id
        : TOK_ID 
        | TOK_FOR
        | TOK_LET
        | TOK_ASSERT
        | TOK_ECHO
        | TOK_EACH
        ;
------------------
*/
moduleId: ID | FOR | LET | ASSERT | ECHO | EACH;

/*
Equivalent from parser.y:
------------------
single_module_instantiation
        : module_id '(' arguments ')'
        ;
 ------------------
*/
singleModuleInstantiation:
	moduleId L_PAREN arguments R_PAREN childStatement;

/*
Equivalent from parser.y:
------------------
expr
        : logic_or
        | TOK_FUNCTION '(' parameters ')' expr %prec NO_ELSE
        | logic_or '?' expr ':' expr
        | TOK_LET '(' arguments ')' expr
        | TOK_ASSERT '(' arguments ')' expr_or_empty
        | TOK_ECHO '(' arguments ')' expr_or_empty
        ;
------------------
*/
expr:
	call										# callExpr
	| expr binaryOperator expr					# binaryExpr
	| ('!' | MINUS | PLUS) expr					# unary
	| FUNCTION '(' parameters ')' expr			# functionDef
	| expr QUESTION_MARK expr COLON expr		# ternaryExpr
	| LET L_PAREN arguments R_PAREN expr		# letExpr
	| ASSERT L_PAREN arguments R_PAREN expr?	# assertion
	| ECHO L_PAREN arguments R_PAREN expr?		# echo;

/*
Equivalent from parser.y:
------------------
call
        : primary
        | call '(' arguments ')'
        | call '[' expr ']'
        | call '.' TOK_ID
------------------
*/
call:
	primary								# primaryElem
	| call L_PAREN arguments R_PAREN	# functionCall
	| call L_BRACKET expr R_BRACKET		# arrayAccess
	| call '.' ID						# memberAccess;

/*
Equivalent from parser.y:
------------------
primary
        : TOK_TRUE
        | TOK_FALSE
        | TOK_UNDEF
        | TOK_NUMBER
        | TOK_STRING
        | TOK_ID
        | '(' expr ')'
        | '[' expr ':' expr ']'
        | '[' expr ':' expr ':' expr ']'
        | '[' ']'
        | '[' vector_elements optional_trailing_comma ']'
		;
------------------
*/
primary:
	(TRUE | FALSE | UNDEF | NUMBER | STRING | ID)				# literalOrId
	| L_PAREN expr R_PAREN										# parenthetical
	| L_BRACKET expr COLON expr (COLON expr)? R_BRACKET			# range
	| L_BRACKET R_BRACKET										# emptyVector
	| L_BRACKET vectorElements optionalTrailingComma R_BRACKET	# vector;

binaryOperator:
	'*'
	| '/'
	| '%'
	| PLUS
	| MINUS
	| POW
	| '<'
	| LE
	| EQ
	| NE
	| GE
	| '>'
	| AND
	| OR;

/*
Equivalent from parser.y:
------------------
vector_elements
        : vector_element
        | vector_elements ',' vector_element
        ;
------------------
*/
vectorElements: vectorElement (comma vectorElement)*;

/*
Equivalent from parser.y:
------------------
list_comprehension_elements
        : TOK_LET '(' arguments ')' list_comprehension_elements_p
        | TOK_EACH vector_element
        | TOK_FOR '(' arguments ')' vector_element
        | TOK_FOR '(' arguments ';' expr ';' arguments ')' vector_element
        | TOK_IF '(' expr ')' vector_element %prec NO_ELSE
        | TOK_IF '(' expr ')' vector_element TOK_ELSE vector_element
        ;
------------------
*/
listComprehensionElements:
	LET L_PAREN arguments R_PAREN listComprehensionElementsP							# letStatement
	| EACH vectorElement																# eachStatement
	| FOR L_PAREN arguments (SEMICOLON expr SEMICOLON arguments)? R_PAREN vectorElement	# forStatement
	| IF L_PAREN expr R_PAREN vectorElement (ELSE vectorElement)?                       # ifStatementComprehension;

/*
Equivalent from parser.y:
------------------
// list_comprehension_elements with optional parenthesis
list_comprehension_elements_p
        : list_comprehension_elements
        | '(' list_comprehension_elements ')'
        ;
------------------
*/
listComprehensionElementsP:
	listComprehensionElements
	| L_PAREN listComprehensionElements R_PAREN;

/*
Equivalent from parser.y:
------------------
vector_element
        : list_comprehension_elements_p
        | expr
        ;
------------------
*/
vectorElement: listComprehensionElementsP | expr;

/*
Equivalent from parser.y:
------------------
parameters
        : /x empty x/
        | parameter_list optional_trailing_comma
        ;
parameter_list
        : parameter
        | parameter_list ',' parameter
        ;
parameter
        : TOK_ID
        | TOK_ID '=' expr
        ;
------------------
*/
parameters:
	// empty
	| parameter (comma parameter)* optionalTrailingComma;

parameter: ID | assignmentExpression;

/*
Equivalent from parser.y:
------------------
arguments
        : /x empty x/
        | argument_list optional_trailing_comma
        ;

argument_list
        : argument
        | argument_list ',' argument
        ;

argument
        : expr
        | TOK_ID '=' expr
        ;
------------------
*/
arguments:
	// empty
	| argument ( comma argument)* optionalTrailingComma;

argument: expr | assignmentExpression;

modifierCharacters: ('!' | '#' | '%' | '*')*;

optionalTrailingComma: comma?;

comma: COMMA;

semicolon: SEMICOLON;

assignmentExpression: ID EQUALS expr;

EQUALS: '=';

SEMICOLON: ';';

COLON: ':';

COMMA: ',';

L_CURLY: '{';

R_CURLY: '}';

L_PAREN: '(';

R_PAREN: ')';

L_BRACKET: '[';

R_BRACKET: ']';

QUESTION_MARK: '?';

GE: '>=';

EQ: '==';

NE: '!=';

LE: '<=';

AND: '&&';

OR: '||';

LET: 'let';

FOR: 'for';

IF: 'if';

TRUE: 'true';

FALSE: 'false';

UNDEF: 'undef';

ELSE: 'else';

ASSERT: 'assert';

ECHO: 'echo';

EACH: 'each';

FUNCTION: 'function';

MODULE: 'module';

ID: '$'? ( LETTER | DIGIT | UNDERSCORE)+;

NUMBER: ( FLOAT | INTEGER);

FLOAT:
	DIGIT+ FLOAT_EXPONENT?
	| DIGIT* '.' DIGIT+ FLOAT_EXPONENT?
	| DIGIT+ '.' DIGIT* FLOAT_EXPONENT?;

INTEGER: DIGIT+;

PLUS: '+';

MINUS: '-';

POW: '^';

STRING: '"' STRING_CHAR* '"';

// Note: All comments are sent to channel #2
SINGLE_LINE_COMMENT: EOL? SPACES '//' ~[\r\n]* -> channel ( 2 );

// Multiple newlines are handled as "comments"
MULTI_NEWLINE: SPACES_EOL_SPACES EOL_SPACES+ -> channel ( 2 );

MULTILINE_COMMENT: '/*' .*? '*/' -> channel ( 2 );

WHITESPACE: [ \t\r\n] -> skip;

INCLUDE_OR_USE_FILE: (INCLUDE | USE) SPACES FILE;

fragment INCLUDE: 'include';

fragment USE: 'use';

fragment FILE: '<' ~[\t\r\n>]* '>';

fragment SPACES_EOL_SPACES: [ \t]* EOL [ \t]*;

fragment EOL_SPACES: EOL [ \t]*;

fragment SPACES: [ \t]*;

fragment EOL: '\r' '\n' | '\n';

fragment FLOAT_EXPONENT: [eE] [+-]? DIGIT+;

fragment STRING_CHAR: ~["\\] | ESCAPE_SEQUENCE;

fragment ESCAPE_SEQUENCE:
	'\\' ["\\rnt]
	| OCTAL_ESCAPE_SEQUENCE
	| UNICODE_ESCAPE_SEQUENCE;

fragment OCTAL_ESCAPE_SEQUENCE: '\\' OCTAL_DIGIT+;

fragment UNICODE_ESCAPE_SEQUENCE: '\\' 'u' HEX_DIGIT+;

fragment FILE_CHAR: [a-zA-Z./];

fragment LETTER: [a-zA-Z];

fragment UPPERCASE_LETTER: [A-Z];

fragment UNDERSCORE: '_';

fragment DIGIT: [0-9];

fragment OCTAL_DIGIT: [0-7];

fragment HEX_DIGIT: [0-9a-fA-F];