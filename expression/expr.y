%{
package expression
import (
)
%}

%union {
    exprVal Expression
    sVal string
    iVal int64
    funcArgs []Expression
    structArgs map[string]Expression
}

%token NUM, STRING
%token IDENTIFIER, VARIABLE
%token FUNC, RETURN, IF, ELSE, SWITCH, CASE, DEFAULT
%token NEWLINE
%token COMPARE

%token '=' '+' '-' '*' '/' '%' '(' ')'

%%
result: expr                    { currentExecutor.Push($1.exprVal) }
      | result newline

expr:   expr0                   { 
}

expr0:  expr1                   { }
    |   expr0 COMPARE expr1     {
                                    $$.exprVal = NewMathStatement([]Expression{$1.exprVal, $3.exprVal}, $2.sVal)
                                }

expr1:  expr2                   { }
    |   expr1 '+' expr2         {
                                    $$.exprVal = NewMathStatement([]Expression{$1.exprVal, $3.exprVal}, "+")
                                }
    |   expr1 '-' expr2         {
                                    $$.exprVal = NewMathStatement([]Expression{$1.exprVal, $3.exprVal}, "-")
                                }

expr2:  expr3                   { }
    |   expr2 '*' expr3         {
                                    $$.exprVal = NewMathStatement([]Expression{$1.exprVal, $3.exprVal}, "*")
                                }
    |   expr2 '/' expr3         {
                                    $$.exprVal = NewMathStatement([]Expression{$1.exprVal, $3.exprVal}, "/")
                                }
    |   expr2 '%' expr3         {
                                    $$.exprVal = NewMathStatement([]Expression{$1.exprVal, $3.exprVal}, "%")
                                }

expr3:  number                  { $$.exprVal = $1.exprVal }
    |   STRING                  { $$.exprVal = NewString($1.sVal) }
    |   variable_expr           { $$.exprVal = $1.exprVal }
    |   call_function           { $$.exprVal = $1.exprVal }
    |   declare_function        { $$.exprVal = $1.exprVal }
    |   exprStruct              { $$.exprVal = $1.exprVal }
    |   '(' expr ')'            { $$.exprVal = $2.exprVal }

exprStruct: '{' structItems '}' {
                                    $$.exprVal = NewStructStatement($2.structArgs)
                                }

structItems : IDENTIFIER ':' expr
                                {
                                    $$.structArgs = make(map[string]Expression)
                                    $$.structArgs[$1.sVal] = $3.exprVal
                                }
    | structItems ',' IDENTIFIER ':' expr
                                {
                                    $$.structArgs[$3.sVal] = $5.exprVal
                                }

number: NUM                     { $$.exprVal = NewInt($1.iVal) }
    | '+' NUM                   { $$.exprVal = NewInt($1.iVal) }
    | '-' NUM                   {
                                    $$.exprVal = NewMathStatement([]Expression{NewInt($1.iVal)}, "neg")
                                }

variable_expr: VARIABLE         { 
                                    $$.exprVal = NewVariableStatement($1.sVal)
                                }

call_function: '$' '(' function ')'
                                {
                                    $$.exprVal = $3.exprVal
                                }
function: identifier '(' argument_list ')'  
                                {
                                    $$.exprVal = NewCallStatement($3.funcArgs, $1.sVal)
                                }

argument_list:  /* empty */     { $$.funcArgs = nil }
    |   expr                    { 
                                    $$.funcArgs = append($$.funcArgs, $1.exprVal)
                                }
    |   argument_list ',' expr  { 
                                    $$.funcArgs = append($1.funcArgs, $3.exprVal)
                                }

declare_function: '$' FUNC '('')' '{' statement_list '}'
                                {
                                    $$.exprVal = NewFunctionStatement($6.funcArgs)
                                }

statement: return_statement     { $$.exprVal = $1.exprVal }
    | switch_statement          { $$.exprVal = $1.exprVal }

statement_list:                 { $$.funcArgs = []Expression{}}
    | statement_list statement  { $$.funcArgs = append($1.funcArgs, $2.exprVal)}

case_expr : expr                { $$.exprVal = NewCaseExpression("", $1.exprVal) }
    | COMPARE expr              { $$.exprVal = NewCaseExpression($1.sVal, $2.exprVal)}

default_statement: DEFAULT ':' '{' statement_list '}'
                                { $$.exprVal = NewDefaultStatement($4.funcArgs) }

case_statement: CASE case_expr ':' '{' statement_list '}'
                                { $$.exprVal = NewCaseStatement($2.exprVal, $5.funcArgs)}

case_statement_list:            { $$.funcArgs = []Expression{}}
    | case_statement_list case_statement
                                { $$.funcArgs = append($1.funcArgs, $2.exprVal) }

switch_statement: SWITCH '(' expr ')' '{' case_statement_list '}'
                                { $$.exprVal = NewSwitchStatement($3.exprVal, $6.funcArgs)}
    | SWITCH '(' expr ')' '{' case_statement_list default_statement '}'
                                { 
                                    tempArgs := append($6.funcArgs, $7.exprVal)
                                    $$.exprVal = NewSwitchStatement($3.exprVal, tempArgs)
                                }

return_statement: RETURN expr ';' 
                                {
                                    $$.exprVal = NewReturnStatement($2.exprVal)
                                }

identifier : IDENTIFIER         {}
newline: NEWLINE                {}

%%
