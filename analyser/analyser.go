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
		},
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			checkNode(n, pass)
			return true
		})
	}
	return nil, nil
}

// checkNode inspects one ast.Node looking for log msg
func checkNode(n ast.Node, pass *analysis.Pass) {
	// Filter Call Expression
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return
	}

	// Get fn
	fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if !ok || fn.Pkg() == nil {
		return
	}

	// Filter by package
	reg, pkgExists := watchedLogs[fn.Pkg().Path()]
	if !pkgExists {
		return
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
		return
	}

	arg := call.Args[msgPos]
	checkLogArg(pass, arg)
}

// checkLogArg inspects one ast Expression (presumably a log msg), passing through string literals
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

// validateMessage inspects log msg applying linter rules
func validateMessage(pass *analysis.Pass, pos token.Pos, msg string) {
	if len(msg) < 1 {
		return
	}

	// Rules 2 and 3
	for i := range msg {
		b := msg[i]
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == ' ') {
			pass.Reportf(pos+token.Pos(i), "log message should only contain english letters, numbers and spaces")
			return
		}
	}

	// Rule 1
	if !(msg[0] >= 'a' && msg[0] <= 'z') {
		pass.Reportf(pos, "log message should start with a lowercase letter")
	}
}
