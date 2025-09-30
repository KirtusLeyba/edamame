package app

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	gui "github.com/gen2brain/raylib-go/raygui"
	"strconv"
)

type UIState int

const (
	UIMain UIState = iota
	UILoad
)

type UILayer struct {
	currentState UIState
	currentFPS int
	origin Vec2Df32
	size Vec2Df32
}

func (u *UILayer) Attach() {
	u.currentState = UIMain
	u.currentFPS = int(rl.GetFPS())
}
func (u *UILayer) Detach() {}
func (u *UILayer) OnEvent() {}
func (u *UILayer) OnUpdate() {
	u.currentFPS = int(rl.GetFPS())
}
func (u *UILayer) OnRender() {
	u.drawStats()
}

func (u *UILayer) drawStats(){

	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	pixelOrigin := Vec2Di{ int(u.origin.x*screenWidth) , int(u.origin.y*screenHeight) }
	pixelSize := Vec2Di{ int(u.size.x*screenWidth) , int(u.size.y*screenHeight) }

	infoBoxOrigin := Vec2Df32{float32(pixelOrigin.x) + 0.025*float32(pixelSize.x),
							  float32(pixelOrigin.y) + 0.05*float32(pixelSize.y)}
	infoBoxSize := Vec2Df32{0.15*float32(pixelSize.x),
						  0.9*float32(pixelSize.y)}

	gui.GroupBox(rl.Rectangle{infoBoxOrigin.x, infoBoxOrigin.y, infoBoxSize.x, infoBoxSize.y}, "Info")
	fpsStr := strconv.Itoa(u.currentFPS)
	rl.DrawText("FPS: " + fpsStr, int32(infoBoxOrigin.x + 8), int32(infoBoxOrigin.y + 8), 16, rl.White)
}

func (u *UILayer) SetTransform(origin, size Vec2Df32){
	u.origin = origin
	u.size = size
}
