package app

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	ednet "edamame/core/networks"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type UIState int

const (
	UIMain UIState = iota
	UILoad
)

type UILayer struct {
	currentState UIState
	currentFPS   int
	origin       Vec2Df32
	size         Vec2Df32
	ltNode       *LayerTreeNode
}

func (u *UILayer) SetLTNode(ltNode *LayerTreeNode) {
	u.ltNode = ltNode
}
func (u *UILayer) OnCreate() {
	log.Printf("UI Layer created with unique id: %v\n", u.ltNode.UniqueID)
}
func (u *UILayer) OnRemove() {
	log.Printf("UI Layer removed with unique id: %v\n", u.ltNode.UniqueID)
}
func (u *UILayer) OnEvent() {}
func (u *UILayer) OnUpdate() {
	u.currentFPS = int(rl.GetFPS())
}
func (u *UILayer) OnRender() {

	u.drawStats()
	loadNodeFile := u.drawNodeButton()
	if loadNodeFile && u.currentState == UIMain {
		var fileLoadLayer FileLoadLayer
		fileLoadLayer.SetTransform(Vec2Df32{0.2, 0.2}, Vec2Df32{0.6, 0.6})
		loadCallback := func(fname string) {
			u.currentState = UIMain
			u.loadNodeData(fname)
		}
		fileLoadLayer.SetCallback(loadCallback)
		u.currentState = UILoad

		u.ltNode.AddChild(&fileLoadLayer)
	}

	loadEdgeFile := u.drawEdgeButton()
	if loadEdgeFile && u.currentState == UIMain {
		var fileLoadLayer FileLoadLayer
		fileLoadLayer.SetTransform(Vec2Df32{0.2, 0.2}, Vec2Df32{0.6, 0.6})
		loadCallback := func(fname string) {
			u.currentState = UIMain
			u.loadEdgeData(fname)
		}
		fileLoadLayer.SetCallback(loadCallback)
		u.currentState = UILoad

		u.ltNode.AddChild(&fileLoadLayer)
	}

	runLayout := u.drawRunLayoutButton()
	if runLayout && u.currentState == UIMain {
		for _, child := range u.ltNode.Children {
			value, isType := child.Data.(*NetworkLayer)
			if isType {
				value.StartLayout = true
			}
		}
	}

	export := u.drawExportButton()
	if export && u.currentState == UIMain {
		//TODO: Make these options the user can select
		var imgSize uint = 8192
		var nodeScale float32 = 4.0
		var spaceScale float32 = 2.0
		var edgeScale float32 = 4.0

		img := rl.GenImageColor(int(imgSize), int(imgSize), rl.White)

		for _, child := range u.ltNode.Children {
			value, isType := child.Data.(*NetworkLayer)
			if isType {
				value.DrawEdgesImage(img, imgSize, imgSize, edgeScale, spaceScale)
				value.DrawNodesImage(img, imgSize, imgSize, nodeScale, spaceScale)
			}
		}

		rl.ExportImage(*img, "result.png")
	}
}

func (u *UILayer) drawStats() {

	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	pixelOrigin := Vec2Di{int(u.origin.X * screenWidth), int(u.origin.Y * screenHeight)}
	pixelSize := Vec2Di{int(u.size.X * screenWidth), int(u.size.Y * screenHeight)}

	infoBoxOrigin := Vec2Df32{float32(pixelOrigin.X) + 0.025*float32(pixelSize.X),
		float32(pixelOrigin.Y) + 0.05*float32(pixelSize.Y)}
	infoBoxSize := Vec2Df32{0.15 * float32(pixelSize.X),
		0.9 * float32(pixelSize.Y)}

	gui.GroupBox(rl.Rectangle{infoBoxOrigin.X, infoBoxOrigin.Y, infoBoxSize.X, infoBoxSize.Y}, "Info")
	fpsStr := strconv.Itoa(u.currentFPS)
	rl.DrawText("FPS: "+fpsStr, int32(infoBoxOrigin.X+8), int32(infoBoxOrigin.Y+8), 16, rl.White)
}

func (u *UILayer) drawNodeButton() bool {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	pixelOrigin := Vec2Di{int(u.origin.X * screenWidth), int(u.origin.Y * screenHeight)}
	pixelSize := Vec2Di{int(u.size.X * screenWidth), int(u.size.Y * screenHeight)}
	infoBoxOrigin := Vec2Df32{float32(pixelOrigin.X) + 0.025*float32(pixelSize.X),
		float32(pixelOrigin.Y) + 0.05*float32(pixelSize.Y)}
	infoBoxSize := Vec2Df32{0.15 * float32(pixelSize.X),
		0.9 * float32(pixelSize.Y)}

	buttonOrigin := Vec2Df32{X: infoBoxOrigin.X + 0.1*infoBoxSize.X,
		Y: infoBoxOrigin.Y + 0.05*infoBoxSize.Y}
	buttonSize := Vec2Df32{X: 0.8 * infoBoxSize.X,
		Y: 0.05 * infoBoxSize.Y}

	loadFilePressed := gui.Button(rl.Rectangle{buttonOrigin.X, buttonOrigin.Y, buttonSize.X, buttonSize.Y}, "Load Node Data")
	return loadFilePressed
}

func (u *UILayer) drawEdgeButton() bool {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	pixelOrigin := Vec2Di{int(u.origin.X * screenWidth), int(u.origin.Y * screenHeight)}
	pixelSize := Vec2Di{int(u.size.X * screenWidth), int(u.size.Y * screenHeight)}
	infoBoxOrigin := Vec2Df32{float32(pixelOrigin.X) + 0.025*float32(pixelSize.X),
		float32(pixelOrigin.Y) + 0.05*float32(pixelSize.Y)}
	infoBoxSize := Vec2Df32{0.15 * float32(pixelSize.X),
		0.9 * float32(pixelSize.Y)}

	buttonOrigin := Vec2Df32{X: infoBoxOrigin.X + 0.1*infoBoxSize.X,
		Y: infoBoxOrigin.Y + 0.15*infoBoxSize.Y}
	buttonSize := Vec2Df32{X: 0.8 * infoBoxSize.X,
		Y: 0.05 * infoBoxSize.Y}

	loadFilePressed := gui.Button(rl.Rectangle{buttonOrigin.X, buttonOrigin.Y, buttonSize.X, buttonSize.Y}, "Load Edge Data")
	return loadFilePressed
}

func (u *UILayer) drawRunLayoutButton() bool {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	pixelOrigin := Vec2Di{int(u.origin.X * screenWidth), int(u.origin.Y * screenHeight)}
	pixelSize := Vec2Di{int(u.size.X * screenWidth), int(u.size.Y * screenHeight)}
	infoBoxOrigin := Vec2Df32{float32(pixelOrigin.X) + 0.025*float32(pixelSize.X),
		float32(pixelOrigin.Y) + 0.05*float32(pixelSize.Y)}
	infoBoxSize := Vec2Df32{0.15 * float32(pixelSize.X),
		0.9 * float32(pixelSize.Y)}

	buttonOrigin := Vec2Df32{X: infoBoxOrigin.X + 0.1*infoBoxSize.X,
		Y: infoBoxOrigin.Y + 0.25*infoBoxSize.Y}
	buttonSize := Vec2Df32{X: 0.8 * infoBoxSize.X,
		Y: 0.05 * infoBoxSize.Y}

	runLayoutPressed := gui.Button(rl.Rectangle{buttonOrigin.X, buttonOrigin.Y, buttonSize.X, buttonSize.Y}, "Toggle Layout")
	return runLayoutPressed
}

func (u *UILayer) drawExportButton() bool {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	pixelOrigin := Vec2Di{int(u.origin.X * screenWidth), int(u.origin.Y * screenHeight)}
	pixelSize := Vec2Di{int(u.size.X * screenWidth), int(u.size.Y * screenHeight)}
	infoBoxOrigin := Vec2Df32{float32(pixelOrigin.X) + 0.025*float32(pixelSize.X),
		float32(pixelOrigin.Y) + 0.05*float32(pixelSize.Y)}
	infoBoxSize := Vec2Df32{0.15 * float32(pixelSize.X),
		0.9 * float32(pixelSize.Y)}

	buttonOrigin := Vec2Df32{X: infoBoxOrigin.X + 0.1*infoBoxSize.X,
		Y: infoBoxOrigin.Y + 0.35*infoBoxSize.Y}
	buttonSize := Vec2Df32{X: 0.8 * infoBoxSize.X,
		Y: 0.05 * infoBoxSize.Y}

	exportPressed := gui.Button(rl.Rectangle{buttonOrigin.X, buttonOrigin.Y, buttonSize.X, buttonSize.Y}, "Export Image")
	return exportPressed
}

func (u *UILayer) SetTransform(origin, size Vec2Df32) {
	u.origin = origin
	u.size = size
}

func (u *UILayer) GetTransform() (Vec2Df32, Vec2Df32) {
	return u.origin, u.size
}

func (u *UILayer) loadNodeData(fname string) {
	content, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(strings.NewReader(string(content)))
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//Apply the new data to all children layers that are of the NetworkLayer type
	var netLayers []*NetworkLayer
	for _, child := range u.ltNode.Children {
		value, isType := child.Data.(*NetworkLayer)
		if isType {
			netLayers = append(netLayers, value)
		}
	}

	// name,radius,r,g,b,a
	// A,1.0,40,94,150,255

	for _, netLayer := range netLayers {
		netLayer.Net = ednet.NewNetwork()
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
			r, err := strconv.ParseUint(record[2], 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			g, err := strconv.ParseUint(record[3], 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			b, err := strconv.ParseUint(record[4], 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			a, err := strconv.ParseUint(record[5], 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			rx := (100.0 * rand.Float32()) - 50.0
			ry := (100.0 * rand.Float32()) - 50.0
			nodeColor := ednet.RGBA{R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a)}
			node := ednet.Node{Name: name,
				Radius:    float32(radius),
				NodeColor: nodeColor,
				X:         rx,
				Y:         ry}
			netLayer.Net.AddNode(node)
		}
	}
}

func (u *UILayer) loadEdgeData(fname string) {
	content, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(strings.NewReader(string(content)))
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//Apply the new data to all children layers that are of the NetworkLayer type
	var netLayers []*NetworkLayer
	for _, child := range u.ltNode.Children {
		value, isType := child.Data.(*NetworkLayer)
		if isType {
			netLayers = append(netLayers, value)
		}
	}

	for _, netLayer := range netLayers {

		netLayer.Net.Adjacencies = make(map[string][]string)
		netLayer.Net.EdgeCount = 0
		netLayer.Net.EdgeSet = make(map[string]uint)
		netLayer.Net.Edges = make([]ednet.Edge, 0)

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
			width, err := strconv.ParseFloat(record[2], 32)
			if err != nil {
				log.Fatal(err)
			}
			edge := ednet.SortedEdge(nameA, nameB)
			edge.Width = float32(width)
			netLayer.Net.AddEdge(edge)
		}
	}
}
