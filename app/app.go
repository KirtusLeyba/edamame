package app

import (
	ednet "edamame/core/networks"
	// "strconv"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func Execute(defaultWidth, defaultHeight int32) {

	initWindow(defaultWidth, defaultHeight)
	defer rl.CloseWindow()

	//set up layer tree
	//with a ui layer as root
	var ui UILayer
	ui.SetTransform(Vec2Df32{0.0, 0.0}, Vec2Df32{1.0, 1.0})
	root := NewRootLayerTreeNode(&ui)

	//add the netviz layer
	var netLayer  NetworkLayer
	netLayer.SetTransform(Vec2Df32{0.1, 0.1}, Vec2Df32{0.8, 0.8})
	netLayer.EquilCon = 50.0
	netLayer.EquilDis = 4.0
	netLayer.SpringConstant = 3.0
	netLayer.TimeStep = 0.1
	netLayer.Friction = 0.001
	netLayer.Net = ednet.NewNetwork()
	root.AddChild(&netLayer)

	mainLoop(root)
	exitApp(root)
}

func initWindow(width, height int32){
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(width, height, "edamame")
}

func mainLoop(root *LayerTreeNode){
	for !rl.WindowShouldClose() {
		root.UpdateTree()

		backgroundColor := rl.Color{57, 56, 84, 255}
		rl.BeginDrawing()
		rl.ClearBackground(backgroundColor)
		root.RenderTree()
		rl.EndDrawing()
	}
}

func exitApp(root *LayerTreeNode){
	root.Remove()
}
