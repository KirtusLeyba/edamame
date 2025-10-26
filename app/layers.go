package app

import (
	"errors"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vec2Df32 struct {
	X float32
	Y float32
}
type Vec2Di struct {
	X int
	Y int
}

type Layer interface {
	OnCreate()
	OnRemove()
	OnEvent()
	OnUpdate()
	OnRender()
	SetLTNode(ltNode *LayerTreeNode)
	SetTransform(origin, size Vec2Df32)
	GetTransform() (Vec2Df32, Vec2Df32)
}

type LayerTreeNode struct {
	TreeSize uint
	UniqueID uint
	Data     Layer
	Children []*LayerTreeNode
	Parent   *LayerTreeNode
	Removed  bool
}

func NewRootLayerTreeNode(rootData Layer) *LayerTreeNode {
	var result = LayerTreeNode{TreeSize: 1, UniqueID: 1, Data: rootData, Parent: nil, Removed: false}
	result.Data.SetLTNode(&result)
	result.Data.OnCreate()
	return &result
}

func NewChildLayerTreeNode(parent *LayerTreeNode, data Layer, uniqueID uint) *LayerTreeNode {
	var result = LayerTreeNode{TreeSize: 1, UniqueID: uniqueID, Data: data, Parent: parent, Removed: false}
	result.Data.SetLTNode(&result)
	result.Data.OnCreate()
	return &result
}

func (ltNode *LayerTreeNode) IncrementTreeSize() {
	ltNode.TreeSize++
	if ltNode.Parent != nil {
		ltNode.Parent.IncrementTreeSize()
	}
}

func (ltNode *LayerTreeNode) DecrementTreeSize() {
	ltNode.TreeSize--
	if ltNode.Parent != nil {
		ltNode.Parent.DecrementTreeSize()
	}
}

func (ltNode *LayerTreeNode) AddChild(childData Layer) {
	ltNode.IncrementTreeSize()
	childNode := NewChildLayerTreeNode(ltNode, childData, ltNode.TreeSize)
	ltNode.Children = append(ltNode.Children, childNode)
}

/**
 * Remove a ltNode by recursively removing it's children.
 */
func (ltNode *LayerTreeNode) Remove() {
	if(ltNode.Children != nil){
		for _, child := range ltNode.Children {
			child.Remove()
		}
	}

	ltNode.DecrementTreeSize()
	ltNode.Data.OnRemove()
	ltNode.Children = nil

	if ltNode.Parent != nil {
		parent := ltNode.Parent
		var newChildren []*LayerTreeNode
		for _, child := range parent.Children {
			if child.UniqueID != ltNode.UniqueID {
				newChildren = append(newChildren, child)
			}
		}
		parent.Children = newChildren
	}

	ltNode.Removed = true
}

func (ltNode *LayerTreeNode) RemoveChildWithID(uniqueID uint) error {
	for _, child := range ltNode.Children {
		if child.UniqueID == uniqueID {
			child.Remove()
			return nil
		}
	}
	return errors.New("cannot remove non-existant child node")
}

func (ltRootNode *LayerTreeNode) UpdateTree() {
	for _, child := range ltRootNode.Children {
		child.UpdateTree()
	}
	ltRootNode.Data.OnUpdate()
}

func (ltRootNode *LayerTreeNode) RenderTree() {
	for _, child := range ltRootNode.Children {
		child.RenderTree()
	}
	ltRootNode.Data.OnRender()
}

func (ltNode *LayerTreeNode) GetRoot() {
	current := ltNode
	for current != nil {
		current = ltNode.Parent
	}
}

func (ltNode *LayerTreeNode) GetFrame() rl.Rectangle {
	origin, size := ltNode.Data.GetTransform()
	if ltNode.Parent == nil {
		screenWidth := float32(rl.GetScreenWidth())
		screenHeight := float32(rl.GetScreenHeight())

		return rl.Rectangle{(origin.X * screenWidth),
			(origin.Y * screenHeight),
			(size.X * screenWidth),
			(size.Y * screenHeight)}
	}
	parentFrame := ltNode.Parent.GetFrame()
	return rl.Rectangle{(origin.X * parentFrame.Width),
		(origin.Y * parentFrame.Height),
		(size.X * parentFrame.Width),
		(size.Y * parentFrame.Height)}
}
