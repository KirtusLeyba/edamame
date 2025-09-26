module edamame

go 1.25.1

replace edamame/core => ./core

replace edamame/app => ./app

require (
	edamame/app v0.0.0-00010101000000-000000000000 // indirect
	edamame/core v0.0.0-00010101000000-000000000000 // indirect
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/gen2brain/raylib-go/raylib v0.55.1 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
