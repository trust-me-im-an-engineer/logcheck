package analyser

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "logcheck",
	Doc:  "check log messages for log/slog pkg",
	Run:  run,
}

type logRegistry struct {
	functions map[string]int            // MethodName -> ArgPos
	methods   map[string]map[string]int // TypeName -> MethodName -> ArgPos
}

var watchedLogs = map[string]logRegistry{
	"log/slog": {
		functions: map[string]int{
			"Debug":        0,
			"DebugContext": 1,
			"Info":         0,
			"InfoContext":  1,
			"Warn":         0,
			"WarnContext":  1,
			"Error":        0,
			"ErrorContext": 1,
			"Log":          2,
		},
		methods: map[string]map[string]int{
			"Logger": {
				"Debug":        0,
				"DebugContext": 1,
				"Info":         0,
				"InfoContext":  1,
				"Warn":         0,
				"WarnContext":  1,
				"Error":        0,
				"ErrorContext": 1,
				"Log":          2,
			},
		},
	},
	"go.uber.org/zap": {
		methods: map[string]map[string]int{
			"Logger": {
				"Debug":  0,
				"Info":   0,
				"Warn":   0,
				"Error":  0,
				"DPanic": 0,
				"Panic":  0,
				"Fatal":  0,
				"Log":    1,
			},
			"SugaredLogger": {
				"Debugw":  0,
				"Infow":   0,
				"Warnw":   0,
				"Errorw":  0,
				"DPanicw": 0,
				"Panicw":  0,
				"Fatalw":  0,
				"Logw":    1,
				"Debugf":  0,
				"Infof":   0,
				"Warnf":   0,
				"Errorf":  0,
				"DPanicf": 0,
				"Panicf":  0,
				"Fatalf":  0,
				"Logf":    1,
			},
		},
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Filter Call Expression
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Get fn
			fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
			if !ok || fn.Pkg() == nil {
				return true
			}

			// Filter by package
			reg, pkgExists := watchedLogs[fn.Pkg().Path()]
			if !pkgExists {
				return true
			}

			// Filter by function or method
			// Differ between function and method by checking receiver
			receiver := fn.Type().(*types.Signature).Recv() // it's safe to cast, because fn is *types.Func
			msgPos := -1
			if receiver == nil { // Function
				if pos, exists := reg.functions[fn.Name()]; exists {
					msgPos = pos
				}
			} else { // Method
				typeName := ""
				if named, ok := receiver.Type().Underlying().(*types.Pointer); ok {
					// Pointer receiver
					if t, ok := named.Elem().(*types.Named); ok {
						typeName = t.Obj().Name()
					}
				} else if t, ok := receiver.Type().(*types.Named); ok {
					// Value receiver
					typeName = t.Obj().Name()
				}

				if typeMethods, exists := reg.methods[typeName]; exists {
					if pos, exists := typeMethods[fn.Name()]; exists {
						msgPos = pos
					}
				}
			}

			if msgPos == -1 || len(call.Args) <= msgPos {
				return true
			}

			arg := call.Args[msgPos]
			checkLogArg(pass, arg)

			return true
		})
	}
	return nil, nil
}

func checkLogArg(pass *analysis.Pass, expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			msg, err := strconv.Unquote(e.Value)
			if err != nil {
				return
			}

			validateMessage(pass, e.Pos(), msg)
		}

	case *ast.BinaryExpr:
		// Recursively check concatenations
		if e.Op == token.ADD {
			checkLogArg(pass, e.X)
			checkLogArg(pass, e.Y)
		}
	}
}

func validateMessage(pass *analysis.Pass, pos token.Pos, value string) {
	pass.Reportf(pos, "msg")
}
