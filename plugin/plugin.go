package main

import (
	"logCheckLinter/external"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"golang.org/x/tools/go/analysis"
)

func New(conf *config.Config) ([]*linter.Config, error) {
	analyzer := external.GetLogAnalizer()
	linterConfig := goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)

	return []*linter.Config{
		linter.NewConfig(linterConfig),
	}, nil
}
