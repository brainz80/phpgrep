package phpgrep

import (
	"reflect"

	"github.com/VKCOM/php-parser/pkg/ast"
)

// rewriteComplexEncapsedVars rewrites `${expr}` string parts so noverify's
// IR converter keeps the inner expression.
//
// irconv turns ScalarEncapsedStringVar into a node only when the name is a
// plain variable or identifier. Any other expression becomes a nil part, and
// ir.Heredoc.Walk (and the encapsed-string walker) panic on that nil.
// `{$expr}` is converted as the expression itself, so represent the complex
// form the same way before conversion.
func rewriteComplexEncapsedVars(n ast.Vertex) {
	if n == nil {
		return
	}
	rv := reflect.ValueOf(n)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return
	}

	vertexType := reflect.TypeOf((*ast.Vertex)(nil)).Elem()
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		switch f.Kind() {
		case reflect.Interface:
			if child, ok := vertexFrom(f); ok {
				rewriteComplexEncapsedVars(child)
			}
		case reflect.Ptr:
			if child, ok := vertexFrom(f); ok {
				rewriteComplexEncapsedVars(child)
			}
		case reflect.Slice:
			if f.Type().Elem() != vertexType || !f.CanSet() {
				continue
			}
			for j := 0; j < f.Len(); j++ {
				el := f.Index(j)
				child, ok := vertexFrom(el)
				if !ok {
					continue
				}
				if replaced := complexEncapsedAsBrackets(child); replaced != nil {
					el.Set(reflect.ValueOf(replaced))
					child = replaced
				}
				rewriteComplexEncapsedVars(child)
			}
		}
	}
}

func vertexFrom(f reflect.Value) (ast.Vertex, bool) {
	if !f.IsValid() || (f.Kind() != reflect.Interface && f.Kind() != reflect.Ptr) || f.IsNil() || !f.CanInterface() {
		return nil, false
	}
	child, ok := f.Interface().(ast.Vertex)
	if !ok || child == nil {
		return nil, false
	}
	return child, true
}

// complexEncapsedAsBrackets returns a {$expr} node for `${expr}` when expr is
// not a plain variable or identifier. Those two forms are already converted
// without producing a nil part.
func complexEncapsedAsBrackets(n ast.Vertex) ast.Vertex {
	sv, ok := n.(*ast.ScalarEncapsedStringVar)
	if !ok || sv == nil || sv.Dim != nil || sv.Name == nil {
		return nil
	}
	switch sv.Name.(type) {
	case *ast.Identifier, *ast.ExprVariable:
		return nil
	default:
		return &ast.ScalarEncapsedStringBrackets{
			Position:             sv.Position,
			OpenCurlyBracketTkn:  sv.DollarOpenCurlyBracketTkn,
			Var:                  sv.Name,
			CloseCurlyBracketTkn: sv.CloseCurlyBracketTkn,
		}
	}
}
