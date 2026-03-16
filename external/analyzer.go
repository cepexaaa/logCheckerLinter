package external

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"logCheckLinter/domain"
	"logCheckLinter/rules"
)

var (
	config = domain.DefaultConfig()
	flSet  flag.FlagSet
)

func GetLogAnalizer() *analysis.Analyzer {
	Analyzer := &analysis.Analyzer{
		Name:     "loglint",
		Doc:      "checks log calls for compliance with the rules",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    flags(),
	}
	return Analyzer
}

func flags() flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.BoolVar(&config.CheckLowercase, "check-lowercase", true, "Check for lowercase letter")
	fs.BoolVar(&config.CheckEnglish, "check-english", true, "Check only English (ASCII)")
	fs.BoolVar(&config.CheckNoSpecials, "check-no-specials", true, "Check for special characters")
	fs.BoolVar(&config.CheckSensitive, "check-sensitive", true, "Check for sensitive data")
	fs.BoolVar(&config.AutoFix, "autofix", true, "Auto-correction (lowercase letter)")
	fs.StringVar(&configFile, "config", "", "Path to the JSON configuration file")
	fs.StringVar(&sensitiveWords, "sensitive-words", strings.Join(config.SensitivePatterns, ","), "List of sensitive words separated by commas")
	return *fs
}

var (
	configFile     string
	sensitiveWords string
)

func run(pass *analysis.Pass) (any, error) {
	if configFile != "" {
		fileCfg, err := LoadConfig(configFile)
		if err != nil {
			return nil, fmt.Errorf("configuration loading error: %v", err)
		}
		config = fileCfg
	} else {

		if sensitiveWords != "" {
			config.SensitivePatterns = strings.Split(sensitiveWords, ",")
		}
	}

	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	sensitiveChecker := rules.NewSensitiveChecker(config.SensitivePatterns)

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		if !isLoggerCall(call) {
			return
		}

		msg, strParts, varNames, err := extractMessageWithConcat(call)
		if err != nil {
			pass.Reportf(call.Pos(), "couldn't extract the log message, skip the check")
			return
		}

		if config.CheckLowercase {
			if len(call.Args) > 0 {
				if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					ok, fix, errMsg := rules.CheckLowercase(msg, lit.Pos())
					if !ok {
						diag := analysis.Diagnostic{
							Pos:     lit.Pos(),
							Message: errMsg,
						}
						if config.AutoFix && fix != "" {
							replacement := fmt.Sprintf(`"%s"`, fix)
							diag.SuggestedFixes = []analysis.SuggestedFix{{
								Message: "replace with a lowercase letter",
								TextEdits: []analysis.TextEdit{{
									Pos:     lit.Pos(),
									End:     lit.End(),
									NewText: []byte(replacement),
								}},
							}}
						}
						pass.Report(diag)
					}
				}
			}
		}

		if config.CheckEnglish {
			if ok, errMsg := rules.CheckEnglish(msg); !ok {
				pos := call.Args[0].Pos()
				pass.Reportf(pos, "%s", errMsg)
			}
		}

		if config.CheckNoSpecials {
			if ok, errMsg := rules.CheckNoSpecials(msg); !ok {
				pos := call.Args[0].Pos()
				pass.Reportf(pos, "%s", errMsg)
			}
		}

		if config.CheckSensitive {
			if ok, errMsg := sensitiveChecker.CheckSensitiveConcat(strParts, varNames); !ok {
				pos := call.Args[0].Pos()
				pass.Reportf(pos, "%s", errMsg)
			}
		}
	})

	return nil, nil
}

func isLoggerCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	baseNames := []string{"Info", "Error", "Debug", "Warn", "Print", "Fatal", "Panic"}
	suffixes := []string{"", "f", "ln"}

	loggerFuncs := make(map[string]bool)
	for _, base := range baseNames {
		for _, suffix := range suffixes {
			loggerFuncs[base+suffix] = true
		}
	}

	return loggerFuncs[sel.Sel.Name]
}

func extractMessage(call *ast.CallExpr) (string, *ast.BasicLit, error) {
	if len(call.Args) == 0 {
		return "", nil, fmt.Errorf("no arguments")
	}
	first := call.Args[0]
	lit, ok := first.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", nil, fmt.Errorf("the first argument is not a string literal")
	}

	str := strings.Trim(lit.Value, `"`)
	return str, lit, nil
}

func extractMessageWithConcat(call *ast.CallExpr) (string, []string, []string, error) {
	if len(call.Args) == 0 {
		return "", nil, nil, fmt.Errorf("there are no arguments")
	}

	first := call.Args[0]

	if lit, ok := first.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		str := strings.Trim(lit.Value, `"`)
		return str, []string{str}, nil, nil
	}

	if binary, ok := first.(*ast.BinaryExpr); ok && binary.Op == token.ADD {
		return extractFromConcat(binary)
	}

	return "", nil, nil, fmt.Errorf("unsupported message format")
}

func extractFromConcat(expr *ast.BinaryExpr) (string, []string, []string, error) {
	var strParts []string
	var varNames []string
	var fullMsg strings.Builder

	var walk func(e ast.Expr)
	walk = func(e ast.Expr) {
		switch v := e.(type) {
		case *ast.BasicLit:
			if v.Kind == token.STRING {
				str := strings.Trim(v.Value, `"`)
				strParts = append(strParts, str)
				fullMsg.WriteString(str)
			}
		case *ast.Ident:
			varNames = append(varNames, v.Name)
			fullMsg.WriteString("{VAR}")
		case *ast.BinaryExpr:
			if v.Op == token.ADD {
				walk(v.X)
				walk(v.Y)
			}
		}
	}

	walk(expr)
	return fullMsg.String(), strParts, varNames, nil
}
