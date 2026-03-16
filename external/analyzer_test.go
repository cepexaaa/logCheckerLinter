package external_test

import (
	"testing"

	"logCheckLinter/external"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, external.GetLogAnalizer(), "slog", "zap")
}
