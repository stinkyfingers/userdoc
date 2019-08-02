package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type visitor struct {
	fset       *token.FileSet
	file       *ast.File
	commentMap ast.CommentMap
	selectors  []string
	userdocs   []userdoc
}

type userdoc struct {
	selector string   `markdownTable:"selector"`
	Params   []string `markdownTable:"params"`
	Comments []string `markdownTable:"comments"`
}

var stdOutFunctions = []string{
	"os.Stdout.Write",
	"fmt.Printf",
	"fmt.Println",
	"funk.Funk",
	"Foo",
}

// parseFlie parses a go file into an ast and calls Walk, populating userdocs on the visitor object.
func parseFile(filename string, selectors []string) ([]userdoc, error) {
	fmt.Println("parsing file: ", filename)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	commentMap := ast.NewCommentMap(fset, file, file.Comments)
	v := &visitor{
		fset:       fset,
		file:       file,
		commentMap: commentMap,
		selectors:  append(stdOutFunctions, selectors...),
	}
	ast.Walk(v, file)
	return v.userdocs, nil
}

// Visit is visitor's implementation to satisfy the Visitor interface.
// For each node, if there is a call expression or return-call expression, it calls getCallExprArgsAndComments().
func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {

	// TODO - remove Ident?
	case *ast.ReturnStmt:
		for _, res := range d.Results {
			switch e := res.(type) {
			case *ast.CallExpr:
				v.getCallExprArgsAndComments(e, &n)
				// case *ast.Ident:
				// 	log.Print("E", e)
				// v.getIdentExprArgsAndComments(e, &n)
			}
		}

	case *ast.ExprStmt:
		switch e := d.X.(type) {
		case *ast.CallExpr:
			v.getCallExprArgsAndComments(e, &n)
		}

	}
	return v
}

// getCallExprArgsAndComments determines if the call is required (in the list of selectors) and
// calls recursivelyGetCallExprArgsAndComments if it is.
func (v *visitor) getCallExprArgsAndComments(call *ast.CallExpr, n *ast.Node) {
	selector, err := v.requiredCall(call)
	if err != nil {
		return
	}
	u := &userdoc{selector: selector}
	v.recursivelyGetCallExprArgsAndComments(call, n, u)
}

// recursivelyGetCallExprArgsAndComments gets the call parameter names (recursively) and comments and
// appends a userdoc to the visitor containing them.
func (v *visitor) recursivelyGetCallExprArgsAndComments(call *ast.CallExpr, n *ast.Node, u *userdoc) {
	appendUserDoc := true

	// args
	for _, arg := range call.Args {
		switch d := arg.(type) {
		case *ast.BasicLit:
			u.Params = append(u.Params, d.Value)

		// case *ast.Ident:
		// 	u.Params = append(u.Params, d.String())
		// 	v.getComments(u, n)
		// 	v.userdocs = append(v.userdocs, *u)
		// 	log.Print("ident ")

		case *ast.CallExpr:
			// get func calls' calls recursively
			appendUserDoc = false
			v.recursivelyGetCallExprArgsAndComments(d, n, u)
		}
	}
	if appendUserDoc {
		v.getComments(u, n)
		v.userdocs = append(v.userdocs, *u)
	}
}

// getComments add a node's comments to the userdoc
func (v *visitor) getComments(u *userdoc, n *ast.Node) {
	if commentGroups, ok := v.commentMap[*n]; ok && len(commentGroups) > 0 {
		u.Comments = append(u.Comments, strings.Replace(commentGroups[0].Text(), "\n", "", -1)) // NOTE - inline comments only; first comment
	}
}

// requiredCall returns the selector string if the call is a required selector or identifier
func (v *visitor) requiredCall(sel *ast.CallExpr) (string, error) {
	switch fun := sel.Fun.(type) {
	case *ast.SelectorExpr:
		return v.requiredSelector(fun)

	case *ast.Ident:
		return v.requiredIdentifier(fun)
	}
	return "", errors.New("call not required")
}

// requiredSelector calls recursivelyEvalIfSelectorIsRequired against the visitor's selector strings
func (v *visitor) requiredSelector(sel *ast.SelectorExpr) (string, error) {
	if sel == nil {
		return "", errors.New("nil selector")
	}
	for _, selector := range v.selectors {
		ok := recursivelyEvalIfSelectorIsRequired(sel, selector)
		if ok {
			return selector, nil
		}
	}
	return "", errors.New("selector not required")
}

// recursivelyEvalIfSelectorIsRequired determines if an ast.SelectorExpr is required by recursively evaluating nested selectors and comparing
// them against a selector string
func recursivelyEvalIfSelectorIsRequired(sel *ast.SelectorExpr, str string) bool {
	arr := strings.Split(str, ".")
	if len(arr) < 1 {
		return false
	}
	if sel.Sel.Name != arr[len(arr)-1] {
		return false
	}

	if iden, ok := sel.X.(*ast.Ident); ok {
		if arr[0] != iden.String() {
			return false
		}
	}

	if nextSelector, ok := sel.X.(*ast.SelectorExpr); ok {
		return recursivelyEvalIfSelectorIsRequired(nextSelector, strings.Join(arr[:len(arr)-1], "."))
	}
	return true
}

// requiredIdentifier returns the selector string if a simple Ident matches a visitor's selector
func (v *visitor) requiredIdentifier(iden *ast.Ident) (string, error) {
	for _, selector := range v.selectors {
		if selector == iden.Name {
			return selector, nil
		}
	}
	return "", errors.New("selector not required")
}
