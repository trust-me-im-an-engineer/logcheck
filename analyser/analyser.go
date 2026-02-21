package analyser

import (
	"encoding/json"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"

	"github.com/trust-me-im-an-engineer/logcheck/analyser/rules"
)

var Analyzer = &analysis.Analyzer{
	Name: "logcheck",
	Doc:  "check log messages for log/slog pkg",
	Run:  run,
}

type WatchCalle struct {
	Functions map[string]int            `json:"functions"`
	Methods   map[string]map[string]int `json:"methods"`
}

// AnalyserConfig holds the pre-processed configuration
type AnalyserConfig struct {
	SensitiveKeywords map[string]struct{}
	WatchedLogs       map[string]WatchCalle
}

var rawConfig struct {
	SensitiveKeywords string
	WatchedLogs       string
}

var (
	parsedConfig *AnalyserConfig
	once         sync.Once
)

func init() {
	Analyzer.Flags.StringVar(
		&rawConfig.SensitiveKeywords,
		"sensitive-keywords",
		"password,token,secret,key",
		"comma-separated list of sensitive keywords",
	)
	Analyzer.Flags.StringVar(
		&rawConfig.WatchedLogs,
		"watched-logs",
		"",
		"JSON string defining custom loggers and argument positions",
	)
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Initialize the configuration exactly once
	once.Do(func() {
		parsedConfig = &AnalyserConfig{
			SensitiveKeywords: parseKeywords(rawConfig.SensitiveKeywords),
			WatchedLogs:       parseWatchedLogs(rawConfig.WatchedLogs),
		}
	})

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			checkNode(n, pass, parsedConfig)
			return true
		})
	}
	return nil, nil
}

func parseKeywords(s string) map[string]struct{} {
	keywordsList := strings.Split(s, ",")
	sensitiveMap := make(map[string]struct{})
	for _, kw := range keywordsList {
		kw = strings.TrimSpace(kw)
		if kw != "" {
			sensitiveMap[kw] = struct{}{}
		}
	}
	return sensitiveMap
}

// checkNode inspects one ast.Node looking for log msg
func checkNode(n ast.Node, pass *analysis.Pass, cfg *AnalyserConfig) {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return
	}

	fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if !ok || fn.Pkg() == nil {
		return
	}

	reg, pkgExists := cfg.WatchedLogs[fn.Pkg().Path()]
	if !pkgExists {
		return
	}

	receiver := fn.Type().(*types.Signature).Recv()
	msgPos := -1
	if receiver == nil { // Function
		if pos, exists := reg.Functions[fn.Name()]; exists {
			msgPos = pos
		}
	} else { // Method
		typeName := ""
		// Handle pointer and value receivers
		t := receiver.Type()
		if ptr, ok := t.Underlying().(*types.Pointer); ok {
			t = ptr.Elem()
		}
		if named, ok := t.(*types.Named); ok {
			typeName = named.Obj().Name()
		}

		if typeMethods, exists := reg.Methods[typeName]; exists {
			if pos, exists := typeMethods[fn.Name()]; exists {
				msgPos = pos
			}
		}
	}

	if msgPos == -1 || len(call.Args) <= msgPos {
		return
	}

	// Check message (rules 1-4)
	checkLogMsg(pass, call.Args[msgPos])

	// Check log arguments for sensitive names (rule 4)
	for _, arg := range call.Args[msgPos:] {
		checkLogArg(pass, arg, cfg.SensitiveKeywords)
	}
}

// checkLogMsg inspects msg applying linter rules 1-3
func checkLogMsg(pass *analysis.Pass, expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			msg, err := strconv.Unquote(e.Value)
			if err != nil {
				return
			}

			checkMessage(pass, e.Pos(), msg)
		}

	case *ast.BinaryExpr:
		// Recursively check concatenations
		if e.Op == token.ADD {
			checkLogMsg(pass, e.X)
			checkLogMsg(pass, e.Y)
		}
	}
}

// checkLogArg inspects log argument for sensitive names (rule 4)
func checkLogArg(pass *analysis.Pass, arg ast.Expr, sensitiveKeywords map[string]struct{}) {
	switch a := arg.(type) {
	case *ast.Ident:
		// Rule 4: Check for sensitive variable names like "password"
		checkIdentSensitiveName(pass, a, sensitiveKeywords)

	case *ast.SelectorExpr:
		// Rule 4: Check for sensitive fields like "user.Token"
		checkIdentSensitiveName(pass, a.Sel, sensitiveKeywords)

	case *ast.BinaryExpr:
		// Recursively check concatenations
		if a.Op == token.ADD {
			checkLogArg(pass, a.X, sensitiveKeywords)
			checkLogArg(pass, a.Y, sensitiveKeywords)
		}
	}
}

// checkIdentSensitiveName inspects ident for sensitive naming
func checkIdentSensitiveName(pass *analysis.Pass, ident *ast.Ident, keywords map[string]struct{}) {
	if i, n := rules.FindSensitiveName(ident.Name, keywords); i != -1 {
		pass.Reportf(ident.Pos()+token.Pos(i), "potential sensitive data leak: argument contains '%s'", n)
	}
}

// checkMessage inspects log msg applying linter rules 1-3.
func checkMessage(pass *analysis.Pass, pos token.Pos, msg string) {
	if len(msg) < 1 {
		return
	}

	// Rules 2 and 3
	if i := rules.IndexIllegalCharacter(msg); i != -1 {
		pass.Reportf(pos+token.Pos(i+1), "log message should only contain english letters, numbers and spaces")
		return
	}

	// Rule 1
	if !rules.StartsWithLowercase(msg) {

		// If first char is Uppercase letter suggest fix, otherwise just report
		if msg[0] >= 'A' && msg[0] <= 'Z' {
			pass.Report(analysis.Diagnostic{
				Pos:     pos + 1,
				Message: "log message should start with a lowercase letter",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "lowercase the first letter",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     pos + 1,
								End:     pos + 2,
								NewText: []byte{msg[0] + 32},
							},
						},
					},
				},
			})
		} else {
			pass.Reportf(pos+token.Pos(1), "log message should start with a lowercase letter")
		}
	}
}

func parseWatchedLogs(s string) map[string]WatchCalle {
	// Default watched logs
	watchedLogs := map[string]WatchCalle{
		"log/slog": {
			Functions: map[string]int{
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
			Methods: map[string]map[string]int{
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
			Methods: map[string]map[string]int{
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

	// Parse the JSON from the flag
	var custom map[string]WatchCalle
	if err := json.Unmarshal([]byte(s), &custom); err != nil {
		return watchedLogs
	}

	// Merge or Override
	for pkg, reg := range custom {
		watchedLogs[pkg] = reg
	}
	return watchedLogs
}
