package main

import (
	"edamame/app"
	"flag"
	"os"
)

func isSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	var opt app.EdamameOptions

	//Parse flags
	flag.BoolVar(&opt.Headless, "headless", false, "run with gui or headless")
	flag.StringVar(&opt.NodeFilePath, "nodeFilePath", "", "File path to find nodes in headless mode")
	flag.StringVar(&opt.EdgeFilePath, "edgeFilePath", "", "File path to find edges in headless mode")
	flag.StringVar(&opt.OutputFilePath, "outputFilePath", "", "File path to save result in headless mode")
	flag.IntVar(&opt.MaxWorkers, "maxWorkers", 1, "Number of go routines to use to generate layout")
	flag.IntVar(&opt.MaxIters, "maxIters", 1, "Number of iterations in the layout algorithm")
	flag.Float64Var(&opt.Repulsion, "repulsion", 80, "Repulsive force")
	flag.Parse()

	if !opt.Headless {
		var defaultWidth int32 = 800
		var defaultHeight int32 = 600
		app.Execute(defaultWidth, defaultHeight)
	} else {
		if !isSet("nodeFilePath") || !isSet("edgeFilePath") || !isSet("outputFilePath") {
			flag.PrintDefaults()
			os.Exit(0)
		}
		app.ExecuteHeadless(&opt)
	}
}
