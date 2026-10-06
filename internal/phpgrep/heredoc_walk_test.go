package phpgrep

import (
	"fmt"
	"testing"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/ir/irconv"
	"github.com/VKCOM/noverify/src/phpdoc"
	"github.com/VKCOM/php-parser/pkg/version"
)

func TestComplexEncapsedInterpolationWalk(t *testing.T) {
	ver, err := version.New("8.0")
	if err != nil {
		t.Fatal(err)
	}
	w := &worker{
		irconv:     irconv.NewConverter(phpdoc.NewTypeParser()),
		phpVersion: ver,
	}

	src := "<?php\n" +
		"$x = <<<EOT\n" +
		"hello $name\n" +
		"${name}\n" +
		"{$name}\n" +
		"${'foo'}\n" +
		"${$a . $b}\n" +
		"${foo()}\n" +
		"${$foo?->bar}\n" +
		"${bar[0]}\n" +
		"EOT;\n" +
		"$y = \"a ${'foo'} b\";\n" +
		"$z = `cmd ${'bar'}`;\n"
	root, err := w.parseFile([]byte(src))
	if err != nil {
		t.Fatal(err)
	}

	seen := map[string]int{}
	root.Walk(countVisitor{seen: seen})

	for _, kind := range []string{
		"*ir.SimpleVar",
		"*ir.String",
		"*ir.ConcatExpr",
		"*ir.FunctionCallExpr",
		"*ir.NullsafePropertyFetchExpr",
		"*ir.ArrayDimFetchExpr",
	} {
		if seen[kind] == 0 {
			t.Errorf("missing %s; seen=%v", kind, seen)
		}
	}
}

type countVisitor struct {
	seen map[string]int
}

func (c countVisitor) EnterNode(n ir.Node) bool {
	c.seen[fmt.Sprintf("%T", n)]++
	return true
}

func (c countVisitor) LeaveNode(ir.Node) {}
