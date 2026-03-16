package main

import (
	"logCheckLinter/external"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(external.GetLogAnalizer())
}
