package main

import (
	"edamame/app"
)

func main() {
	var defaultWidth int32 = 800
	var defaultHeight int32 = 600
	app.Execute(defaultWidth, defaultHeight)
}
