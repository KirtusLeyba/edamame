package app

import (
	"encoding/csv"
	ednet "github.com/KirtusLeyba/edamame/core/networks"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

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

type HeadlessLayer struct {
	ltNode                                                     *LayerTreeNode
	opt                                                        *EdamameOptions
	Net                                                        *ednet.SpatialNet
	SpringConstant, StepSize, Equilibrium, Repulsion, Friction float32
	currentIteration, lastIteration                            int
	MaxIters                                                   int
	MaxWorkers                                                 uint
	finished                                                   bool
}

func logHeadless(msg string) {
	log.Printf("[HEADLESS]: %s\n", msg)
}

func (hl *HeadlessLayer) OnCreate() {
	logHeadless("Loading node data from: " +
		hl.opt.NodeFilePath +
		", and edge data from: " +
		hl.opt.EdgeFilePath)
	hl.loadNodeData(hl.opt.NodeFilePath)
	hl.loadEdgeData(hl.opt.EdgeFilePath)

	hl.currentIteration = 0
	hl.lastIteration = 0

	logHeadless("Computing layout")
	hl.finished = false
	go func() {
		for range hl.MaxIters {
			hl.Net.SpringUpdateParallel(
				hl.SpringConstant,
				hl.StepSize,
				hl.Equilibrium,
				hl.Repulsion,
				hl.Friction,
				hl.MaxWorkers)
			hl.currentIteration++
		}
		hl.finished = true
	}()

}
func (hl *HeadlessLayer) OnRemove() {
	//TODO: Make these options the user can select
	var imgSize uint = 8192
	var nodeScale float32 = 4.0
	var spaceScale float32 = 2.0
	var edgeScale float32 = 4.0

	img := rl.GenImageColor(int(imgSize), int(imgSize), rl.White)
	hl.DrawEdgesImage(img, imgSize, imgSize, edgeScale, spaceScale)
	hl.DrawNodesImage(img, imgSize, imgSize, nodeScale, spaceScale)
	rl.ExportImage(*img, hl.opt.OutputFilePath)
	logHeadless("go routine finished layout iterations")

	//write the node positions to a csv file for reuse
	//TODO: Add option to name file
	fp, err := os.Create("./node_positions.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	_, err = fp.WriteString("node,x,y\n")
	if err != nil {
		log.Fatal(err)
	}
	fp.Sync()

	for _, node := range hl.Net.NodeSlice {
		_, err = fp.WriteString(node.Name + "," +
			strconv.FormatFloat(float64(node.X), 'f', 4, 64) +
			strconv.FormatFloat(float64(node.Y), 'f', 4, 64) + "\n")
		if err != nil {
			log.Fatal(err)
		}
		fp.Sync()
	}
	logHeadless("Wrote ./node_positions.csv")
}

func (hl *HeadlessLayer) OnEvent() {}
func (hl *HeadlessLayer) OnUpdate() {

	if hl.currentIteration != hl.lastIteration {
		hl.lastIteration = hl.currentIteration
		if hl.lastIteration%50 == 0 {
			progress := float64(hl.lastIteration) / float64(hl.MaxIters)
			logHeadless("Progress: " + strconv.FormatFloat(progress, 'f', 4, 64))
		}
	}
	if hl.finished {
		hl.ltNode.Remove()
	}
}
func (hl *HeadlessLayer) OnRender() {}
func (hl *HeadlessLayer) SetLTNode(ltNode *LayerTreeNode) {
	hl.ltNode = ltNode
}
func (hl *HeadlessLayer) SetTransform(origin, size Vec2Df32) {}
func (hl *HeadlessLayer) GetTransform() (Vec2Df32, Vec2Df32) {
	return Vec2Df32{X: 0.0, Y: 0.0}, Vec2Df32{X: 0.0, Y: 0.0}
}

func (hl *HeadlessLayer) loadNodeData(fname string) {
	content, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(strings.NewReader(string(content)))
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	hl.Net = ednet.NewSpatialNet()
	for lineIDX, record := range records {
		//skip the header
		if lineIDX == 0 {
			continue
		}
		if len(record) != 6 {
			log.Fatal("Bad node data file!")
		}
		name := record[0]
		radius, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			log.Fatal(err)
		}
		//TODO: Store color data
		// r, err := strconv.ParseUint(record[2], 10, 8)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// g, err := strconv.ParseUint(record[3], 10, 8)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// b, err := strconv.ParseUint(record[4], 10, 8)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// a, err := strconv.ParseUint(record[5], 10, 8)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		hl.Net.AddNode(name)
		var node *ednet.SpatialNetNode = &hl.Net.NodeSlice[len(hl.Net.NodeSlice)-1]
		node.X = (100.0 * rand.Float32()) - 50.0
		node.Y = (100.0 * rand.Float32()) - 50.0
		node.Radius = float32(radius)
	}
}

func (hl *HeadlessLayer) loadEdgeData(fname string) {
	content, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(strings.NewReader(string(content)))
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//Reset edge data in the SpatialNet
	hl.Net.Adjacencies = make(ednet.EdgeSet)
	for _, node := range hl.Net.NodeSlice {
		hl.Net.Adjacencies[node.Name] = make(map[string]struct{})
	}
	for lineIDX, record := range records {
		//skip the header
		if lineIDX == 0 {
			continue
		}
		if len(record) != 3 {
			log.Fatal("Bad edge data file!")
		}
		nameA := record[0]
		nameB := record[1]
		//TODO: Store edge width
		// width, err := strconv.ParseFloat(record[2], 32)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		hl.Net.AddEdge(nameA, nameB)
	}
}

func (hl *HeadlessLayer) DrawEdgesImage(img *rl.Image, width, height uint, edgeScale, spaceScale float32) {
	frame := rl.Rectangle{0.0, 0.0, float32(width), float32(height)}
	cameraCenter := Vec2Df32{frame.X + frame.Width/2, frame.Y + frame.Height/2}
	cx, cy := hl.Net.GetCOM()
	com := Vec2Df32{cx, cy}
	for sourceNodeName, targetNodeSet := range hl.Net.Adjacencies {
		for targetNodeName, _ := range targetNodeSet {
			nodeA := hl.Net.NodeSlice[hl.Net.NodeIndeces[sourceNodeName]]
			nodeB := hl.Net.NodeSlice[hl.Net.NodeIndeces[targetNodeName]]
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

func (hl *HeadlessLayer) DrawNodesImage(img *rl.Image, width, height uint, nodeScale, spaceScale float32) {
	frame := rl.Rectangle{0.0, 0.0, float32(width), float32(height)}
	cx, cy := hl.Net.GetCOM()
	com := Vec2Df32{cx, cy}
	for _, n := range hl.Net.NodeSlice {
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
		// rl.ImageDrawText(img, int32(posAdjusted.X), int32(posAdjusted.Y), n.Name, 8, rl.White)
	}
}
