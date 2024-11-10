/[ \t]/             { /* Skip blanks and tabs. */ }
/[\n]/              { /* return NEWLINE */ }
/[0-9]*/            { lval.iVal, _ = strconv.ParseInt(yylex.Text(), 10, 64); return NUM }
/func/              { return FUNC }
/return/            { return RETURN }
/if/                { return IF }
/else/              { return ELSE }
/switch/            { return SWITCH }
/case/              { return CASE }
/default/           { return DEFAULT }

/[a-zA-Z_][a-zA-Z_0-9]*/ 
                    { lval.sVal = yylex.Text(); return IDENTIFIER }
/[$][{][^}]+[}]/    {
                        str := yylex.Text()
                        lval.sVal = str[2:len(str) - 1]
                        return VARIABLE
                    }
/["]([^"]|\\")*["]/
                    {
                        str := yylex.Text()
                        str = str[1:len(str) - 1]
                        str = strings.ReplaceAll(str, "\\\"", "\"")
                        lval.sVal = str
                        return STRING
                    }
/[>]/               { lval.sVal = yylex.Text(); return COMPARE }
/[<]/               { lval.sVal = yylex.Text(); return COMPARE }
/[>][=]/            { lval.sVal = yylex.Text(); return COMPARE }
/[<][=]/            { lval.sVal = yylex.Text(); return COMPARE }
/[=][=]/            { lval.sVal = yylex.Text(); return COMPARE }
/[!][=]/            { lval.sVal = yylex.Text(); return COMPARE }

/./                 { return int(yylex.Text()[0]) }

%%
//
package expression
import (
    "strconv"
)

func (yylex Lexer) Error(e string) {
    panic(ErrSyntaxError)
}

