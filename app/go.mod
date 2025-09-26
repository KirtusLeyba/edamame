module edamame/app

go 1.25.1

replace edamame/core => ../core

require (
	edamame/core v0.0.0-00010101000000-000000000000
	github.com/ebitengine/purego v0.7.1
	github.com/gen2brain/raylib-go/raylib v0.55.1
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842
	golang.org/x/sys v0.20.0
)
