package app

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	gui "github.com/gen2brain/raylib-go/raygui"
)
//

// type Layer interface {
// 	Attach(*[]Layer)
// 	Detach()
// 	OnEvent()
// 	OnUpdate()
// 	OnRender()
// 	SetTransform(origin, size Vec2Df32)
// }

type FileLoadLayer struct {
	origin Vec2Df32
	size Vec2Df32
	scrollVal int32
	layerStack *[]Layer
}

func (f *FileLoadLayer) Attach(layerStack *[]Layer){
	f.layerStack = layerStack
	f.scrollVal = 0
}

func (f *FileLoadLayer) Detach(){}
func (f *FileLoadLayer) OnEvent(){}
func (f *FileLoadLayer) OnUpdate(){}
func (f *FileLoadLayer) OnRender(){

	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	fileDialogRect := rl.Rectangle{( f.origin.x * screenWidth  ),
								   ( f.origin.y * screenHeight ),
								   ( f.size.x   * screenWidth  ),
								   ( f.size.y   * screenHeight )}
	fileDialogColor := rl.NewColor(77, 77, 77, 200)
	rl.DrawRectangle(int32(fileDialogRect.X),
					 int32(fileDialogRect.Y),
					 int32(fileDialogRect.Width),
					 int32(fileDialogRect.Height), fileDialogColor)

	var sliderOrigin, sliderSize Vec2Df32
	sliderOrigin.x =  (f.origin.x*screenWidth)+(f.size.x*screenWidth)
	sliderOrigin.y  = (f.origin.y*screenHeight)
	sliderSize.x = 0.02*fileDialogRect.Width
	sliderSize.y = fileDialogRect.Height

	scrollBarBounds := rl.Rectangle{sliderOrigin.x, sliderOrigin.y, sliderSize.x, sliderSize.y}

	f.scrollVal = gui.ScrollBar(scrollBarBounds, f.scrollVal, 0, 100)

}
func (f* FileLoadLayer) SetTransform(origin, size Vec2Df32){
	f.origin = origin
	f.size = size
}
