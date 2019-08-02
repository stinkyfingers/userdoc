package main

import (
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	filename := filepath.Join("test_files", "foobar.go")
	userdocs, err := parseFile(filename, []string{})
	if err != nil {
		t.Error(err)
	}
	t.Log(userdocs)
}

// func TestParseFunctionSelectorString(t *testing.T) {
// 	str := "os.Stdout.Write"
// 	sel := &ast.SelectorExpr{}
// 	parseFunctionSelectorString(str, sel)
// 	t.Log(sel)
// }

// func TestRequiredSelector(t *testing.T) {
// 	v := &visitor{}
// 	sel := &ast.SelectorExpr{
// 		Sel: &ast.Ident{
// 			Name: "Poo",
// 		},
// 		X: &ast.SelectorExpr{
// 			Sel: &ast.Ident{
// 				Name: "Foo",
// 			},
// 			X: &ast.SelectorExpr{
// 				Sel: &ast.Ident{
// 					Name: "foobar",
// 				},
// 			},
// 		},
// 	}
//
// 	tests := []struct {
// 		pkg      string
// 		expected bool
// 	}{
// 		// {
// 		// 	"foobar.Foo",
// 		// 	false,
// 		// },
// 		{
// 			"foobar.Bar",
// 			false,
// 		},
// 		// {
// 		// 	"foobar.Foo.Poo",
// 		// 	true,
// 		// },
// 	}
//
// 	for _, test := range tests {
// 		ok := v.requiredSelector(sel, test.pkg)
// 		if ok != test.expected {
// 			t.Errorf("expected %t, got %t on pkg %s", test.expected, ok, test.pkg)
// 		}
// 	}
// }
