package app

import (
	ednet "github.com/KirtusLeyba/edamame/core/networks"
	// "strconv"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type EdamameOptions struct{
	Headless bool
	NodeFilePath, EdgeFilePath, OutputFilePath string
	MaxWorkers, MaxIters int
	Repulsion float64
}

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
	netLayer.Equilibrium = 8.0
	netLayer.Repulsion = 80.0
	netLayer.SpringConstant = 0.1
	netLayer.StepSize = 0.1
	netLayer.Friction = 0.125
	netLayer.MaxIters = 100
	netLayer.MaxWorkers = 10
	netLayer.Net = ednet.NewSpatialNet()
	root.AddChild(&netLayer)

	mainLoop(root)
	exitApp(root)
}

func ExecuteHeadless(opt *EdamameOptions){
	var headless HeadlessLayer
	headless.opt = opt
	headless.Equilibrium = 8.0
	headless.Repulsion = float32(opt.Repulsion)
	headless.SpringConstant = 0.1
	headless.StepSize = 0.1
	headless.Friction = 0.125
	headless.MaxIters = opt.MaxIters
	headless.MaxWorkers = uint(opt.MaxWorkers)
	headless.Net = ednet.NewSpatialNet()

	root := NewRootLayerTreeNode(&headless)
	mainLoopHeadless(root)
}

func initWindow(width, height int32){
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(width, height, "edamame")
}

func mainLoop(root *LayerTreeNode){

	edamameGray := rl.Color(94, 94, 94, 94)

	for !rl.WindowShouldClose() {
		root.UpdateTree()

		backgroundColor := edamameGray
		rl.BeginDrawing()
		rl.ClearBackground(backgroundColor)
		root.RenderTree()
		rl.EndDrawing()
	}
}

func exitApp(root *LayerTreeNode){
	root.Remove()
}

func mainLoopHeadless(root *LayerTreeNode){
	for !root.Removed {
		root.UpdateTree()
		root.RenderTree()
	}
}
