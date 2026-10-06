package phpgrep

import (
	"strings"
	"testing"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/ir/irconv"
	"github.com/VKCOM/noverify/src/phpdoc"
	"github.com/VKCOM/php-parser/pkg/version"
)

func TestAmpersandAtEOFDoesNotPanic(t *testing.T) {
	ver, err := version.New("8.0")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{name: "bitwise without newline", src: "<?php $x = $a & 1", wantErr: false},
		{name: "bitwise with semicolon", src: "<?php $x = $a & 1;", wantErr: false},
		{name: "reference parameter", src: "<?php function f(&$a) { return $a; }", wantErr: false},
		{name: "ampersand then newline", src: "<?php function f(&\n", wantErr: true},
		{name: "ampersand then dollar", src: "<?php $x = $a & $", wantErr: true},
		{name: "ampersand then ellipsis", src: "<?php $x = $a & ...", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &worker{
				irconv:     irconv.NewConverter(phpdoc.NewTypeParser()),
				phpVersion: ver,
			}
			root, err := w.parseFile([]byte(test.src))
			if test.wantErr {
				if err == nil {
					t.Fatal("expected parse error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			root.Walk(nopWalk{})
			if test.name == "bitwise without newline" && !strings.Contains(string(w.data), ";") {
				t.Fatal("expected salvaged source to include a trailing semicolon")
			}
		})
	}
}

type nopWalk struct{}

func (nopWalk) EnterNode(ir.Node) bool { return true }
func (nopWalk) LeaveNode(ir.Node)      {}
