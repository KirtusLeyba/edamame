package app


import (
	rl "github.com/gen2brain/raylib-go/raylib"
	gui "github.com/gen2brain/raylib-go/raygui"
	"strconv"
)

func drawBackground(){
	rl.ClearBackground(rl.White)
}

func drawStats(){
	gui.GroupBox(rl.Rectangle{10, 10, 100, 100}, "Info")
	fpsStr := strconv.Itoa(int(rl.GetFPS()))
	rl.DrawText("FPS: " + fpsStr, 12, 14, 16, rl.Black)
}
