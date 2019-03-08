package readonly

import (
	"go/ast"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "readonly",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
		commentmap.Analyzer,
	},
}

// flags
var annotation = "assign"

func init() {
	Analyzer.Flags.StringVar(&annotation, "annotation", annotation, "annotation for explicit assignment")
}

const Doc = `check for assignment package variables`

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)
	pkg := pass.ResultOf[buildssa.Analyzer].(buildssa.SSA).Package

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
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

			if cmaps.Annotated(n, annotation) {
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
