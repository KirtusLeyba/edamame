package app

import (
	// "edamame/core/networks"
	// "strconv"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Execute(defaultWidth, defaultHeight int32) {

	initWindow(defaultWidth, defaultHeight)
	defer rl.CloseWindow()

	mainLoop()

}

func updateStep(){

}

func drawStep(){
	rl.BeginDrawing()

	drawBackground()
	drawStats()

	rl.EndDrawing()
}

func mainLoop(){
	for !rl.WindowShouldClose() {
		updateStep()
		drawStep()
	}
}

func initWindow(width, height int32){
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(width, height, "edamame")
}
