package app

import (
	// "edamame/core/networks"
	// "strconv"
	rl "github.com/gen2brain/raylib-go/raylib"
	"errors"
	"fmt"
)

func Execute(defaultWidth, defaultHeight int32) {

	initWindow(defaultWidth, defaultHeight)
	defer rl.CloseWindow()

	//set up default layers
	var layerStack []Layer

	//ui layer
	var uiLayer UILayer
	uiLayer.SetTransform(Vec2Df32{0.0, 0.0}, Vec2Df32{1.0, 1.0})
	PushLayer(&layerStack, &uiLayer)

	mainLoop(&layerStack)
	exitApp(&layerStack)
}

func initWindow(width, height int32){
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(width, height, "edamame")
}

func mainLoop(layerStack *[]Layer){
	for !rl.WindowShouldClose() {
		update(layerStack)
		render(layerStack)
	}
}

func PushLayer(layerStack *[]Layer, lay Layer){
	*layerStack = append(*layerStack, lay)
	lay.Attach(layerStack)
}

func PopLayer(layerStack *[]Layer) (*[]Layer, Layer, error){

	if(len(*layerStack) == 0){
		return layerStack, nil, errors.New("cannot pop layer from an empty stack")
	}

	lastLayer := (*layerStack)[len(*layerStack) - 1]
	*layerStack = (*layerStack)[:len(*layerStack) - 1]
	lastLayer.Detach()
	return layerStack, lastLayer, nil
}

func update(layerStack *[]Layer){
	for _, layer := range *layerStack {
		layer.OnUpdate()
	}
}

func render(layerStack *[]Layer){
	backgroundColor := rl.Color{57, 56, 84, 255}


	rl.BeginDrawing()
	rl.ClearBackground(backgroundColor)

	for _, layer := range *layerStack {
		layer.OnRender()
	}

	rl.EndDrawing()
}

func exitApp(layerStack *[]Layer){
	for len(*layerStack) > 0 {
		var err error
		layerStack, _, err = PopLayer(layerStack)
		if(err != nil){
			panic(fmt.Sprintf("%v\n", err))
		}
	}
}
