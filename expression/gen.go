//go:generate nex -e -o expr.lex.go expr.lex
//go:generate goyacc -o expr.y.go expr.y
package expression
