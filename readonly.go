package readonly

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "readonly",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

const Doc = `check for possible assigning package variables`

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	cms := make([]ast.CommentMap, len(pass.Files))
	for i, f := range pass.Files {
		cms[i] = ast.NewCommentMap(pass.Fset, f, f.Comments)
	}

	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		switch n := n.(type) {
		case *ast.AssignStmt:
			if inInit(pass, stack) {
				return true
			}

			if hasComment(pass, cms, n) {
				return true
			}

			for _, l := range n.Lhs {
				ident, ok := l.(*ast.Ident)
				if ok && isPkgIdent(pass, ident) {
					pass.Reportf(ident.Pos(), "%s shoud not be assigned twice", ident.Name)
				}
			}
		}
		return false
	})

	return nil, nil
}

func hasComment(pass *analysis.Pass, cms []ast.CommentMap, n ast.Node) bool {
	for _, cm := range cms {
		for _, cg := range cm[n] {
			if strings.HasPrefix(strings.TrimSpace(cg.Text()), "not-readonly") {
				return true
			}
		}
	}
	return false
}

func inInit(pass *analysis.Pass, stack []ast.Node) bool {
	if len(stack) == 0 {
		return false
	}
	f, ok := stack[0].(*ast.FuncDecl)
	isInit := ok && f.Name.Name == "init" && isPkgIdent(pass, f.Name)
	return isInit || inInit(pass, stack[1:])
}

func isPkgIdent(pass *analysis.Pass, ident *ast.Ident) bool {
	obj := pass.TypesInfo.Defs[ident]
	if obj == nil {
		obj = pass.TypesInfo.Uses[ident]
	}
	return obj != nil && pass.Pkg.Scope() == obj.Parent()
}
