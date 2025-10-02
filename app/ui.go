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
	layerStack *[]Layer
}

func (u *UILayer) Attach(layerStack *[]Layer) {
	u.layerStack = layerStack
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
	loadFilePressed := u.drawButtons()
	if(loadFilePressed && u.currentState == UIMain){
		var fileLoadLayer FileLoadLayer
		fileLoadLayer.SetTransform(Vec2Df32{0.1, 0.1}, Vec2Df32{0.8, 0.8})
		PushLayer(u.layerStack, &fileLoadLayer)
		u.currentState = UILoad
	}
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

func (u *UILayer) drawButtons() bool {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	pixelOrigin := Vec2Di{ int(u.origin.x*screenWidth) , int(u.origin.y*screenHeight) }
	pixelSize := Vec2Di{ int(u.size.x*screenWidth) , int(u.size.y*screenHeight) }
	infoBoxOrigin := Vec2Df32{float32(pixelOrigin.x) + 0.025*float32(pixelSize.x),
		float32(pixelOrigin.y) + 0.05*float32(pixelSize.y)}
	infoBoxSize := Vec2Df32{0.15*float32(pixelSize.x),
		0.9*float32(pixelSize.y)}

	buttonOrigin := Vec2Df32{x : infoBoxOrigin.x + 0.1*infoBoxSize.x,
							 y : infoBoxOrigin.y + 0.05*infoBoxSize.y}
	buttonSize := Vec2Df32{  x : 0.8*infoBoxSize.x,
							 y : 0.05*infoBoxSize.y}

	loadFilePressed := gui.Button( rl.Rectangle{buttonOrigin.x, buttonOrigin.y, buttonSize.x, buttonSize.y}, "Load Network")
	return loadFilePressed
}

func (u *UILayer) SetTransform(origin, size Vec2Df32){
	u.origin = origin
	u.size = size
}
