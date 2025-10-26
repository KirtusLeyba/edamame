package networks

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"sync"
)

type SpatialNetNode struct {
	Name                 string
	X, Y, Vx, Vy, Radius float32
}

func (a *SpatialNetNode) Equals(b *SpatialNetNode) bool {
	return a.Name == b.Name
}

type EdgeSet map[string]map[string]struct{}

type SpatialNet struct {
	//The slice pointing to the actual
	//allocated memory for the nodes.
	NodeSlice []SpatialNetNode

	//Maps the name of a node to its index in the NodeSlice
	NodeIndeces map[string]uint

	//Maps the name of a node to a slice containing the indeces
	//of nodes it is adjacent to in the NodeSlice
	Adjacencies EdgeSet

	//structures for spatial hashing
	//SpatialBins map the int coords of a bin
	//to a slice containing the indeces of that bin's
	//nodes in NodeSlice
	SpatialBins        map[[2]int][]uint
	SpatialAdjacencies map[[2]int]map[string]int
}

func NewSpatialNet() *SpatialNet {
	return &SpatialNet{NodeSlice: make([]SpatialNetNode, 0),
		NodeIndeces: make(map[string]uint),
		Adjacencies: make(map[string]map[string]struct{})}
}

func (n *SpatialNet) GetCOM() (float32, float32) {
	var cx, cy float32

	for _, node := range n.NodeSlice {
		cx += node.X
		cy += node.Y
	}
	cx /= float32(len(n.NodeSlice))
	cy /= float32(len(n.NodeSlice))
	return cx, cy
}

func (n *SpatialNet) ContainsNode(nodeName string) bool {
	_, exists := n.NodeIndeces[nodeName]
	return exists
}

func (n *SpatialNet) ContainsEdge(nodeAName string, nodeBName string) bool {
	//if either node is missing, cannot contain edge
	if !n.ContainsNode(nodeAName) || !n.ContainsNode(nodeBName) {
		return false
	}

	//grab the set of neighbors to node A, and check if node B exists in the set
	neighborSet := n.Adjacencies[nodeAName]
	_, exists := neighborSet[nodeBName]
	return exists
}

func (n *SpatialNet) AddNode(name string) error {
	if n.ContainsNode(name) {
		return errors.New("attempted to add node named " + name + " twice")
	}
	n.NodeSlice = append(n.NodeSlice, SpatialNetNode{Name: name})
	n.NodeIndeces[name] = uint(len(n.NodeSlice) - 1)
	n.Adjacencies[name] = make(map[string]struct{})
	return nil
}

func (n *SpatialNet) AddEdge(nameA, nameB string) error {
	if !n.ContainsNode(nameA) || !n.ContainsNode(nameB) {
		return errors.New("Cannot add edge between nodes that do not exist!")
	}
	n.Adjacencies[nameA][nameB] = struct{}{}
	n.Adjacencies[nameB][nameA] = struct{}{}
	return nil
}

func NewRandomSpatialNet(numNodes int, edgeProb float32) *SpatialNet {
	n := NewSpatialNet()

	//add the nodes
	for i := range numNodes {
		n.AddNode(strconv.Itoa(i))
	}

	//add the edges
	for i := range numNodes {
		for j := i + 1; j < numNodes; j++ {
			if rand.Float32() < edgeProb {
				n.AddEdge(strconv.Itoa(i), strconv.Itoa(j))
			}
		}
	}

	return n
}

func (n *SpatialNet) SpringUpdate(k,
	stepSize,
	equilibriumDist,
	repulsion,
	friction float32) {

	for i := range len(n.NodeSlice) {
		var fx, fy float32 = 0.0, 0.0
		for j := range len(n.NodeSlice) {
			if i == j {
				continue
			}
			var f float32 = 0.0
			nodeA := &n.NodeSlice[i]
			nodeB := &n.NodeSlice[j]
			dist := math.Hypot(float64(nodeB.X-nodeA.X), float64(nodeB.Y-nodeA.Y))
			x := float32(dist) - equilibriumDist
			if n.ContainsEdge(nodeA.Name, nodeB.Name) {
				f = float32(x) * k
			} else {
				if dist < 1.0 {
					dist = 1.0
				}
				f = -1.0 * repulsion / (float32(dist * dist))
			}
			theta := math.Atan2(float64(nodeB.Y-nodeA.Y), float64(nodeB.X-nodeA.X))
			fx += f * float32(math.Cos(theta))
			fy += f * float32(math.Sin(theta))
		}
		n.NodeSlice[i].Vx += stepSize * fx
		n.NodeSlice[i].Vy += stepSize * fy
	}

	for i := range len(n.NodeSlice) {
		n.NodeSlice[i].X += stepSize * n.NodeSlice[i].Vx
		n.NodeSlice[i].Y += stepSize * n.NodeSlice[i].Vy
	}
}

func (n *SpatialNet) SpringUpdateParallel(k,
	stepSize,
	equilibriumDist,
	repulsion,
	friction float32,
	maxWorkers uint) {

	actualWorkers := maxWorkers
	if(len(n.NodeSlice) < int(actualWorkers)){
		actualWorkers = uint(len(n.NodeSlice))
	}

	var wg = &sync.WaitGroup{}
	queue := make(chan int, actualWorkers)

	worker := func(wg *sync.WaitGroup, queue chan int) {
		defer wg.Done()
		for i := range queue {
			var fx, fy float32 = 0.0, 0.0
			for j := range len(n.NodeSlice) {
				if i == j {
					continue
				}
				var f float32 = 0.0
				nodeA := &n.NodeSlice[i]
				nodeB := &n.NodeSlice[j]
				dist := math.Hypot(float64(nodeB.X-nodeA.X), float64(nodeB.Y-nodeA.Y))
				x := float32(dist) - equilibriumDist
				if n.ContainsEdge(nodeA.Name, nodeB.Name) {
					f = float32(x) * k
				} else {
					if dist < 1.0 {
						dist = 1.0
					}
					f = -1.0 * repulsion / (float32(dist * dist))
				}
				theta := math.Atan2(float64(nodeB.Y-nodeA.Y), float64(nodeB.X-nodeA.X))
				fx += f * float32(math.Cos(theta))
				fy += f * float32(math.Sin(theta))
			}
			n.NodeSlice[i].Vx += stepSize * fx
			n.NodeSlice[i].Vy += stepSize * fy
		}
	}

	for range actualWorkers {
		wg.Add(1)
		go worker(wg, queue)
	}
	for i := range n.NodeSlice {
		queue <- i
	}
	close(queue)
	wg.Wait()

	for i := range len(n.NodeSlice) {
		wg.Go(func() {
			n.NodeSlice[i].X += stepSize * n.NodeSlice[i].Vx
			n.NodeSlice[i].Y += stepSize * n.NodeSlice[i].Vy
			n.NodeSlice[i].Vx -= stepSize*friction*n.NodeSlice[i].Vx
			n.NodeSlice[i].Vy -= stepSize*friction*n.NodeSlice[i].Vy
		})
	}
	wg.Wait()
}

func (snn *SpatialNetNode) GetBin(binSize float32) [2]int {
	var bin [2]int
	bin[0] = int(snn.X / binSize)
	bin[1] = int(snn.Y / binSize)
	return bin
}

func (n *SpatialNet) ResetSpatialHashing(binSize float32) float32 {

	n.SpatialBins = make(map[[2]int][]uint)
	n.SpatialAdjacencies = make(map[[2]int]map[string]int)
	for i := range len(n.NodeSlice) {
		bin := n.NodeSlice[i].GetBin(binSize)
		_, exists := n.SpatialBins[bin]
		if !exists {
			n.SpatialBins[bin] = make([]uint, 0)
		}
		_, exists = n.SpatialAdjacencies[bin]
		if !exists {
			n.SpatialAdjacencies[bin] = make(map[string]int)
			n.SpatialAdjacencies[bin][n.NodeSlice[i].Name] = 0
		}
		n.SpatialBins[bin] = append(n.SpatialBins[bin], uint(i))
		neighborSet := n.Adjacencies[n.NodeSlice[i].Name]
		for nbr, _ := range neighborSet {
			n.SpatialAdjacencies[bin][nbr]++
		}
	}
	return binSize
}

func (n *SpatialNet) SpringUpdateHashing(k,
	stepSize,
	equilibriumDist,
	repulsion,
	friction,
	binSize float32) {

	for i := range len(n.NodeSlice) {
		var fx, fy float32 = 0.0, 0.0

		//update with the nodes that are all within this bin
		localBin := n.NodeSlice[i].GetBin(binSize)
		localNodeIndeces := n.SpatialBins[localBin]
		for j := range len(localNodeIndeces) {
			var nbrIDX int = int(localNodeIndeces[j])
			if i == nbrIDX {
				continue
			}
			var f float32 = 0.0
			nodeA := &n.NodeSlice[i]
			nodeB := &n.NodeSlice[nbrIDX]
			dist := math.Hypot(float64(nodeB.X-nodeA.X), float64(nodeB.Y-nodeA.Y))
			x := float32(dist) - equilibriumDist
			if n.ContainsEdge(nodeA.Name, nodeB.Name) {
				f = float32(x) * k
			} else {
				if dist < 1.0 {
					dist = 1.0
				}
				f = -1.0 * repulsion / (float32(dist * dist))
			}
			theta := math.Atan2(float64(nodeB.Y-nodeA.Y), float64(nodeB.X-nodeA.X))
			fx += f * float32(math.Cos(theta))
			fy += f * float32(math.Sin(theta))
		}

		//update with other bins
		for bin, _ := range n.SpatialBins {
			if bin == localBin {
				continue
			}
			var fDis, fCon float32 = 0.0, 0.0
			nodeA := &n.NodeSlice[i]
			binX := float32(bin[0]) * binSize
			binY := float32(bin[1]) * binSize
			dist := math.Hypot(float64(binX-nodeA.X), float64(binY-nodeA.Y))
			x := float32(dist) - equilibriumDist

			totalInBin := len(n.SpatialBins[bin])
			totalCon := n.SpatialAdjacencies[bin][nodeA.Name]
			fCon = float32(totalCon) * float32(x) * k
			fDis = float32(totalInBin-totalCon) * (-1.0 * repulsion / (float32(dist * dist)))
			theta := math.Atan2(float64(binY-nodeA.Y), float64(binX-nodeA.X))
			f := fCon + fDis
			fx += f * float32(math.Cos(theta))
			fy += f * float32(math.Sin(theta))
		}

		n.NodeSlice[i].Vx += stepSize * fx
		n.NodeSlice[i].Vy += stepSize * fy
	}

	for i := range len(n.NodeSlice) {
		n.NodeSlice[i].X += stepSize * n.NodeSlice[i].Vx
		n.NodeSlice[i].Y += stepSize * n.NodeSlice[i].Vy
	}
}

func (n *SpatialNet) SpringUpdateHashingParallel(k,
	stepSize,
	equilibriumDist,
	repulsion,
	friction,
	binSize float32) {

	var wg sync.WaitGroup
	for i := range len(n.NodeSlice) {
		wg.Go(func() {
			var fx, fy float32 = 0.0, 0.0

			//update with the nodes that are all within this bin
			localBin := n.NodeSlice[i].GetBin(binSize)
			localNodeIndeces := n.SpatialBins[localBin]
			for j := range len(localNodeIndeces) {
				var nbrIDX int = int(localNodeIndeces[j])
				if i == nbrIDX {
					continue
				}
				var f float32 = 0.0
				nodeA := &n.NodeSlice[i]
				nodeB := &n.NodeSlice[nbrIDX]
				dist := math.Hypot(float64(nodeB.X-nodeA.X), float64(nodeB.Y-nodeA.Y))
				x := float32(dist) - equilibriumDist
				if n.ContainsEdge(nodeA.Name, nodeB.Name) {
					f = float32(x) * k
				} else {
					if dist < 1.0 {
						dist = 1.0
					}
					f = -1.0 * repulsion / (float32(dist * dist))
				}
				theta := math.Atan2(float64(nodeB.Y-nodeA.Y), float64(nodeB.X-nodeA.X))
				fx += f * float32(math.Cos(theta))
				fy += f * float32(math.Sin(theta))
			}

			//update with other bins
			for bin, _ := range n.SpatialBins {
				if bin == localBin {
					continue
				}
				var fDis, fCon float32 = 0.0, 0.0
				nodeA := &n.NodeSlice[i]
				binX := float32(bin[0]) * binSize
				binY := float32(bin[1]) * binSize
				dist := math.Hypot(float64(binX-nodeA.X), float64(binY-nodeA.Y))
				x := float32(dist) - equilibriumDist

				totalInBin := len(n.SpatialBins[bin])
				totalCon := n.SpatialAdjacencies[bin][nodeA.Name]
				fCon = float32(totalCon) * float32(x) * k
				fDis = float32(totalInBin-totalCon) * (-1.0 * repulsion / (float32(dist * dist)))
				theta := math.Atan2(float64(binY-nodeA.Y), float64(binX-nodeA.X))
				f := fCon + fDis
				fx += f * float32(math.Cos(theta))
				fy += f * float32(math.Sin(theta))
			}

			n.NodeSlice[i].Vx += stepSize * fx
			n.NodeSlice[i].Vy += stepSize * fy
		})
	}
	wg.Wait()

	//TODO: Nodes moving from bin to bin
	// for i := range len(n.NodeSlice) {
	// 	wg.Go(func() {
	// 		oldBin := n.NodeSlice[i].GetBin()
	// 		n.NodeSlice[i].X += stepSize * n.NodeSlice[i].Vx
	// 		n.NodeSlice[i].Y += stepSize * n.NodeSlice[i].Vy
	// 		newBin := n.NodeSlice[i].GetBin()
	// 		n.SpatialBins[oldBin]
	// 	})
	// }
	wg.Wait()
}
