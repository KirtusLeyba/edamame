package networks

import (
	"errors"
)

type RGBA struct {
	R, G, B, A uint8
}
type Node struct {
	Name         string
	X, Y, Vx, Vy, Radius float32
	NodeColor    RGBA
	BorderColor  RGBA
	NameColor    RGBA
}
type Edge struct {
	NodeAName string
	NodeBName string
	Width float32
}

func SortedEdge(a, b string) Edge {
	if a < b {
		return Edge{NodeAName: a,
					NodeBName: b}
	}
	return Edge{NodeAName: b, NodeBName: a}
}

func (a Node) Equal(b Node) bool {
	return a.Name == b.Name
}
func (a Edge) Equal(b Edge) bool {
	return (a.NodeAName == b.NodeAName && a.NodeBName == b.NodeBName)
}

type Network struct {
	NodeCount, EdgeCount uint
	Nodes                []Node
	Edges                []Edge
	NodeSet              map[string]uint
	EdgeSet              map[string]uint
	Adjacencies          map[string][]string
}

func NewNetwork() *Network {
	return &Network{NodeCount: 0, EdgeCount: 0,
		NodeSet:     make(map[string]uint),
		EdgeSet:     make(map[string]uint),
		Adjacencies: make(map[string][]string)}
}

func (g Network) ContainsNode(aName string) bool {
	_, exists := g.NodeSet[aName]
	return exists
}

func (g Network) ContainsEdge(e Edge) bool {
	edgeName := e.NodeAName + "," + e.NodeBName
	_, exists := g.EdgeSet[edgeName]
	return exists
}

func (g *Network) AddNode(node Node) error {
	if g.ContainsNode(node.Name) {
		return errors.New("attempted to add node named " + node.Name + " twice.")
	}
	g.Nodes = append(g.Nodes, node)
	g.NodeSet[node.Name] = uint(len(g.Nodes) - 1)
	g.Adjacencies[node.Name] = make([]string, 0)
	g.NodeCount++
	return nil
}

func (g *Network) setAdjacent(a, b string){
	g.Adjacencies[a] = append(g.Adjacencies[a], b)
	g.Adjacencies[b] = append(g.Adjacencies[b], a)
}

func (g *Network) AddEdge(e Edge) error {
	edgeName := e.NodeAName + "," + e.NodeBName
	if g.ContainsEdge(e) {
		return errors.New("Attempted to add the edge (" + e.NodeAName + "," + e.NodeBName + ") twice!")
	}
	if(!g.ContainsNode(e.NodeAName) || !g.ContainsNode(e.NodeBName)){
		return errors.New("Cannot add edge between nodes that do not exist!")
	}
	g.Edges = append(g.Edges, e)
	g.EdgeSet[edgeName] = uint(len(g.Edges) - 1)
	g.setAdjacent(e.NodeAName, e.NodeBName)
	g.EdgeCount++
	return nil
}
