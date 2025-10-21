package app

import (
	ednet "edamame/core/networks"
	"math"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"sync"
)

func getCenterOfMass(g *ednet.Network) Vec2Df32 {
	var cx float32
	var cy float32
	for _, node := range g.Nodes {
		cx += node.X
		cy += node.Y
	}
	cx /= float32(g.NodeCount)
	cy /= float32(g.NodeCount)
	return Vec2Df32{cx, cy}
}

func computeForce(g *ednet.Network, i uint,
				  springConstant, ts, equilCon, equilDis, friction float32){
	var fx float32 = 0.0
	var fy float32 = 0.0
	for j := range g.NodeCount {
		if(i == j){
			continue
		}

		a := &g.Nodes[i]
		b := &g.Nodes[j]

		r := math.Sqrt(float64((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y)))
		xCon := float32(r) - equilCon
		// xDis := float32(r) - equilDis
		var f float32
		testEdge := ednet.SortedEdge(a.Name, b.Name)
		if g.ContainsEdge(testEdge) {
			f = float32(xCon)*springConstant
		} else {
			if(r < 1.0){
				r = 1.0
			}
			f = -1.0*equilDis/(float32(r)*float32(r))
		}
		theta := math.Atan2(float64(b.Y-a.Y), float64(b.X-a.X))
		fx += f * float32(math.Cos(theta))
		fy += f * float32(math.Sin(theta))
	}
	g.Nodes[i].Vx += fx*ts
	g.Nodes[i].Vy += fy*ts
}

func springUpdate(g *ednet.Network,
				  springConstant, ts, equilCon, equilDis, friction float32) {
	var wg sync.WaitGroup
	for i := range g.NodeCount {
		wg.Go(func() {
			computeForce(g, i, springConstant, ts, equilCon, equilDis, friction)
		})
	}
	wg.Wait()

	for i := range g.NodeCount {
		g.Nodes[i].X += g.Nodes[i].Vx * ts
		g.Nodes[i].Y += g.Nodes[i].Vy * ts
		g.Nodes[i].Vx -= friction*g.Nodes[i].Vx
		g.Nodes[i].Vy -= friction*g.Nodes[i].Vy
	}
}

// type Layer interface {
// 	OnCreate()
// 	OnRemove()
// 	OnEvent()
// 	OnUpdate()
// 	OnRender()
// 	SetLTNode(ltNode *LayerTreeNode)
// 	SetTransform(origin, size Vec2Df32)
// 	GetTransform() (Vec2Df32, Vec2Df32)
// }

type NetworkLayer struct {
	origin Vec2Df32
	size Vec2Df32
	ltNode *LayerTreeNode
	Net *ednet.Network
	SpringConstant, TimeStep, EquilCon, EquilDis, Friction float32
	RunningLayout bool
}

func (nl *NetworkLayer) OnCreate() {
	log.Printf("NetworkLayer Layer created with unique id: %v\n", nl.ltNode.UniqueID)
}
func (nl *NetworkLayer) OnRemove() {
	log.Printf("NetworkLayer Layer removed with unique id: %v\n", nl.ltNode.UniqueID)
}
func (nl *NetworkLayer) OnEvent() {}
func (nl *NetworkLayer) OnUpdate(){
	if(nl.RunningLayout){
		springUpdate(nl.Net, nl.SpringConstant, nl.TimeStep, nl.EquilCon, nl.EquilDis, nl.Friction)
	}
}
func (nl *NetworkLayer) OnRender(){
	nl.drawEdges()
	nl.drawNodes()
}
func (nl *NetworkLayer) SetLTNode(ltNode *LayerTreeNode){
	nl.ltNode = ltNode
}
func (nl *NetworkLayer) SetTransform(origin, size Vec2Df32){
	nl.origin = origin
	nl.size = size
}
func (nl *NetworkLayer) GetTransform() (Vec2Df32, Vec2Df32){
	return nl.origin, nl.size
}

func (nl *NetworkLayer) drawEdges(){
	frame := nl.ltNode.GetFrame()
	cameraCenter := Vec2Df32{frame.X + frame.Width/2, frame.Y + frame.Height/2}
	com := getCenterOfMass(nl.Net)
	for _, edge := range nl.Net.Edges{
		nodeA := nl.Net.Nodes[nl.Net.NodeSet[edge.NodeAName]]
		nodeB := nl.Net.Nodes[nl.Net.NodeSet[edge.NodeBName]]
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
		rl.DrawLineEx(rl.Vector2{X: posAdjustedA.X, Y: posAdjustedA.Y}, rl.Vector2{X: posAdjustedB.X, Y: posAdjustedB.Y}, edge.Width, rl.Black)
	}
}

func (nl *NetworkLayer) drawNodes(){
	frame := nl.ltNode.GetFrame()
	com := getCenterOfMass(nl.Net)
	for _, n := range nl.Net.Nodes{
		node := nl.Net.Nodes[nl.Net.NodeSet[n.Name]]
		posReal := Vec2Df32{node.X, node.Y}
		posAdjusted := Vec2Df32{posReal.X - com.X,
							    posReal.Y - com.Y}
		cameraCenter := Vec2Df32{frame.X + frame.Width/2, frame.Y + frame.Height/2}
		posAdjusted.X = cameraCenter.X + posAdjusted.X
		posAdjusted.Y = cameraCenter.Y + posAdjusted.Y
		nodeColor := rl.NewColor(node.NodeColor.R, node.NodeColor.G, node.NodeColor.B, node.NodeColor.A)
		rl.DrawCircle(int32(posAdjusted.X), int32(posAdjusted.Y), node.Radius, nodeColor)
		rl.DrawText(node.Name, int32(posAdjusted.X), int32(posAdjusted.Y), 8, rl.White)
	}
}
