package app

import (
	ednet "edamame/core/networks"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
)

type NetworkLayer struct {
	origin                                                     Vec2Df32
	size                                                       Vec2Df32
	ltNode                                                     *LayerTreeNode
	Net                                                        *ednet.SpatialNet
	SpringConstant, StepSize, Equilibrium, Repulsion, Friction float32
	NodeTexture                                                rl.RenderTexture2D
	StartLayout                                                bool
	RunningLayout                                              bool
	MaxIters, MaxWorkers                                       uint
}

func (nl *NetworkLayer) OnCreate() {

	//Create a render texture for nodes since drawing circles is expensive
	nl.NodeTexture = rl.LoadRenderTexture(32, 32)
	rl.BeginTextureMode(nl.NodeTexture)
	rl.DrawCircle(16, 16, 8.0, rl.White)
	rl.EndTextureMode()

	log.Printf("NetworkLayer Layer created with unique id: %v\n", nl.ltNode.UniqueID)
}
func (nl *NetworkLayer) OnRemove() {
	log.Printf("NetworkLayer Layer removed with unique id: %v\n", nl.ltNode.UniqueID)
}
func (nl *NetworkLayer) OnEvent() {}

func (nl *NetworkLayer) OnUpdate() {
	if nl.StartLayout {
		if !nl.RunningLayout {
			nl.RunningLayout = true
			go func() {
				for range nl.MaxIters {
					nl.Net.SpringUpdateParallel(
						nl.SpringConstant,
						nl.StepSize,
						nl.Equilibrium,
						nl.Repulsion,
						nl.Friction,
						nl.MaxWorkers)
					if !nl.RunningLayout {
						break
					}
				}
				nl.RunningLayout = false
			}()
		} else {
			nl.RunningLayout = false
		}
		nl.StartLayout = false
	}
}
func (nl *NetworkLayer) OnRender() {
	nl.drawEdges()
	nl.drawNodes()
}
func (nl *NetworkLayer) SetLTNode(ltNode *LayerTreeNode) {
	nl.ltNode = ltNode
}
func (nl *NetworkLayer) SetTransform(origin, size Vec2Df32) {
	nl.origin = origin
	nl.size = size
}
func (nl *NetworkLayer) GetTransform() (Vec2Df32, Vec2Df32) {
	return nl.origin, nl.size
}

func (nl *NetworkLayer) drawEdges() {
	frame := nl.ltNode.GetFrame()
	cameraCenter := Vec2Df32{frame.X + frame.Width/2, frame.Y + frame.Height/2}
	cx, cy := nl.Net.GetCOM()
	com := Vec2Df32{cx, cy}
	for sourceNodeName, targetNodeSet := range nl.Net.Adjacencies {
		for targetNodeName, _ := range targetNodeSet {
			nodeA := nl.Net.NodeSlice[nl.Net.NodeIndeces[sourceNodeName]]
			nodeB := nl.Net.NodeSlice[nl.Net.NodeIndeces[targetNodeName]]
			posRealA := Vec2Df32{nodeA.X, nodeA.Y}
			posAdjustedA := Vec2Df32{posRealA.X - com.X,
				posRealA.Y - com.Y}
			posAdjustedA.X = cameraCenter.X + posAdjustedA.X
			posAdjustedA.Y = cameraCenter.Y + posAdjustedA.Y
			posRealB := Vec2Df32{nodeB.X, nodeB.Y}
			posAdjustedB := Vec2Df32{posRealB.X - com.X,
				posRealB.Y - com.Y}
			posAdjustedB.X = cameraCenter.X + posAdjustedB.X
			posAdjustedB.Y = cameraCenter.Y + posAdjustedB.Y
			//TODO: don't hardcode size of circle texture
			rl.DrawLineEx(rl.Vector2{X: posAdjustedA.X + 16, Y: posAdjustedA.Y + 16}, rl.Vector2{X: posAdjustedB.X + 16, Y: posAdjustedB.Y + 16}, 1.0, rl.Black)
		}
	}
}

func (nl *NetworkLayer) drawNodes() {
	frame := nl.ltNode.GetFrame()
	cx, cy := nl.Net.GetCOM()
	com := Vec2Df32{cx, cy}
	for _, n := range nl.Net.NodeSlice {
		posReal := Vec2Df32{n.X, n.Y}
		posAdjusted := Vec2Df32{posReal.X - com.X,
			posReal.Y - com.Y}
		cameraCenter := Vec2Df32{frame.X + frame.Width/2, frame.Y + frame.Height/2}
		posAdjusted.X = cameraCenter.X + posAdjusted.X
		posAdjusted.Y = cameraCenter.Y + posAdjusted.Y
		nodeColor := rl.NewColor(0, 0, 255, 255)
		rl.DrawTexture(nl.NodeTexture.Texture, int32(posAdjusted.X), int32(posAdjusted.Y), nodeColor)
	}
}

func (nl *NetworkLayer) DrawEdgesImage(img *rl.Image, width, height uint, edgeScale, spaceScale float32) {
	frame := rl.Rectangle{0.0, 0.0, float32(width), float32(height)}
	cameraCenter := Vec2Df32{frame.X + frame.Width/2, frame.Y + frame.Height/2}
	cx, cy := nl.Net.GetCOM()
	com := Vec2Df32{cx, cy}
	for sourceNodeName, targetNodeSet := range nl.Net.Adjacencies {
		for targetNodeName, _ := range targetNodeSet {
			nodeA := nl.Net.NodeSlice[nl.Net.NodeIndeces[sourceNodeName]]
			nodeB := nl.Net.NodeSlice[nl.Net.NodeIndeces[targetNodeName]]
			posRealA := Vec2Df32{nodeA.X, nodeA.Y}
			posAdjustedA := Vec2Df32{posRealA.X - com.X,
				posRealA.Y - com.Y}
			posRealB := Vec2Df32{nodeB.X, nodeB.Y}
			posAdjustedB := Vec2Df32{posRealB.X - com.X,
				posRealB.Y - com.Y}

			//rescale for image output
			posAdjustedA.X *= spaceScale
			posAdjustedA.Y *= spaceScale
			posAdjustedB.X *= spaceScale
			posAdjustedB.Y *= spaceScale

			posAdjustedA.X = cameraCenter.X + posAdjustedA.X
			posAdjustedA.Y = cameraCenter.Y + posAdjustedA.Y
			posAdjustedB.X = cameraCenter.X + posAdjustedB.X
			posAdjustedB.Y = cameraCenter.Y + posAdjustedB.Y

			edgeWidth := 10.0 //TODO: Make this an option

			rl.ImageDrawLineEx(img,
				rl.Vector2{X: posAdjustedA.X, Y: posAdjustedA.Y},
				rl.Vector2{X: posAdjustedB.X, Y: posAdjustedB.Y},
				int32(edgeWidth), rl.Black)
		}
	}
}

func (nl *NetworkLayer) DrawNodesImage(img *rl.Image, width, height uint, nodeScale, spaceScale float32) {
	frame := rl.Rectangle{0.0, 0.0, float32(width), float32(height)}
	cx, cy := nl.Net.GetCOM()
	com := Vec2Df32{cx, cy}
	for _, n := range nl.Net.NodeSlice {
		posReal := Vec2Df32{n.X, n.Y}
		posAdjusted := Vec2Df32{posReal.X - com.X,
			posReal.Y - com.Y}
		posAdjusted.X *= spaceScale
		posAdjusted.Y *= spaceScale
		cameraCenter := Vec2Df32{frame.X + frame.Width/2, frame.Y + frame.Height/2}
		posAdjusted.X = cameraCenter.X + posAdjusted.X
		posAdjusted.Y = cameraCenter.Y + posAdjusted.Y
		radius := n.Radius * nodeScale
		nodeColor := rl.NewColor(0, 0, 255, 255)
		rl.ImageDrawCircle(img, int32(posAdjusted.X), int32(posAdjusted.Y), int32(radius), nodeColor)
		rl.ImageDrawText(img, int32(posAdjusted.X), int32(posAdjusted.Y), n.Name, 8, rl.White)
	}
}
