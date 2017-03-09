grammar SCAD;


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
    emptyStatement // is semicolon

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

emptyStatement
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
    emptyStatement
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
    L_BRACKET
    (
        expr
        | listComprehension
    )
    (
        commas childListComprehensionOrExpr
    )* optionalCommas R_BRACKET
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

MULTI_NEWLINE
:
    EOL EOL+ -> channel ( 2 )
;

LINE_COMMENT
:
    EOL* '//' .*? EOL+ -> channel ( 2 )
;

COMMENT
:
    '/*' .*? '*/' EOL* -> channel ( 2 )
;

WHITESPACE
:
    [ \t\r\n]+ -> skip
;

EOL
:
    '\r' '\n'   
    |   '\n'       
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
    


