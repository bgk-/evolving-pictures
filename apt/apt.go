package apt

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type Node interface {
	Eval(x, y float32) float32
	String() string
	SetParent(parent Node)
	GetParent() Node
	GetChildren() []Node
	AddRandom(node Node)
	AddLeaf(node Node) bool
	NodeCount() int
}

func Mutate(node Node) Node {
	r := rand.Intn(19)
	var mutatedNode Node
	if r <= 16 {
		mutatedNode = GetRandomNode()
	} else {
		mutatedNode = GetRandomLeaf()
	}

	// Change parents child pointer to point to new node
	if node.GetParent() != nil {
		for i, parentChild := range node.GetParent().GetChildren() {
			if parentChild == node {
				node.GetParent().GetChildren()[i] = mutatedNode
			}
		}
	}

	// Move children from old node to new
	for i, child := range node.GetChildren() {
		if i >= len(mutatedNode.GetChildren()) {
			break
		}
		mutatedNode.GetChildren()[i] = child
		child.SetParent(mutatedNode)
	}

	for i, child := range mutatedNode.GetChildren() {
		if child == nil {
			leaf := GetRandomLeaf()
			leaf.SetParent(mutatedNode)
			mutatedNode.GetChildren()[i] = leaf
		}
	}

	mutatedNode.SetParent(node.GetParent())
	return mutatedNode
}

type BaseNode struct {
	Parent   Node
	Children []Node
}

func (node *BaseNode) Eval(x, y float32) float32 {
	panic("tried to call eval on basenode")
}

func (node *BaseNode) String() string {
	panic("tried to call string on basenode")
}

func (node *BaseNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *BaseNode) GetParent() Node {
	return node.Parent
}

func (node *BaseNode) GetChildren() []Node {
	return node.Children
}

func (node *BaseNode) AddRandom(nodeToAdd Node) {
	addIndex := rand.Intn(len(node.Children))
	if node.Children[addIndex] == nil {
		nodeToAdd.SetParent(node)
		node.Children[addIndex] = nodeToAdd
	} else {
		node.Children[addIndex].AddRandom(nodeToAdd)
	}
}

func (node *BaseNode) AddLeaf(leaf Node) bool {
	for i, child := range node.Children {
		if child == nil {
			leaf.SetParent(node)
			node.Children[i] = leaf
			return true
		} else if node.Children[i].AddLeaf(leaf) {
			return true
		}
	}
	return false
}

func (node *BaseNode) NodeCount() int {
	count := 1
	for _, child := range node.Children {
		count += child.NodeCount()
	}
	return count
}

func GetNthNode(node Node, n, count int) (Node, int) {
	if n == count {
		return node, count
	}
	var result Node
	for _, child := range node.GetChildren() {
		count++
		result, count = GetNthNode(child, n, count)
		if result != nil {
			return result, count
		}
	}
	return nil, count
}

type OpLerp struct {
	BaseNode
}

func NewOpLerp() *OpLerp {
	return &OpLerp{BaseNode{nil, make([]Node, 3)}}
}

func (op *OpLerp) Eval(x, y float32) float32 {
	a := op.Children[0].Eval(x, y)
	b := op.Children[1].Eval(x, y)
	pct := op.Children[2].Eval(x, y)
	return a + pct*(b-a)
}

func (op *OpLerp) String() string {
	return fmt.Sprintf("( Lerp %s %s %s )", op.Children[0].String(), op.Children[1].String(), op.Children[2].String())
}

type OpClip struct {
	BaseNode
}

func NewOpClip() *OpClip {
	return &OpClip{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpClip) Eval(x, y float32) float32 {
	value := op.Children[0].Eval(x, y)
	max := float32(math.Abs(float64(op.Children[1].Eval(x, y))))
	if value > max {
		return max
	} else if value < -max {
		return -max
	}
	return value
}

func (op *OpClip) String() string {
	return fmt.Sprintf("( Clip %s %s )", op.Children[0].String(), op.Children[1].String())
}

type OpSin struct {
	BaseNode
}

func NewOpSin() *OpSin {
	return &OpSin{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Children[0].Eval(x, y))))
}

func (op *OpSin) String() string {
	return fmt.Sprintf("( Sin %s )", op.Children[0].String())
}

type OpCos struct {
	BaseNode
}

func NewOpCos() *OpCos {
	return &OpCos{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpCos) Eval(x, y float32) float32 {
	return float32(math.Cos(float64(op.Children[0].Eval(x, y))))
}

func (op *OpCos) String() string {
	return fmt.Sprintf("( Cos %s )", op.Children[0].String())
}

type OpAtan struct {
	BaseNode
}

func NewOpAtan() *OpAtan {
	return &OpAtan{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpAtan) Eval(x, y float32) float32 {
	return float32(math.Atan(float64(op.Children[0].Eval(x, y))))
}

func (op *OpAtan) String() string {
	return fmt.Sprintf("( Atan %s )", op.Children[0].String())
}

type OpWrap struct {
	BaseNode
}

func NewOpWrap() *OpWrap {
	return &OpWrap{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpWrap) Eval(x, y float32) float32 {
	f := op.Children[0].Eval(x, y)
	temp := (f - -1.0) / (2.0)
	return -1.0 + 2.0*(temp-float32(math.Floor(float64(temp))))
}
func (op *OpWrap) String() string {
	return fmt.Sprintf("( Wrap %s )", op.Children[0].String())
}

type OpLog2 struct {
	BaseNode
}

func NewOpLog2() *OpLog2 {
	return &OpLog2{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpLog2) Eval(x, y float32) float32 {
	return float32(math.Log2(float64(op.Children[0].Eval(x, y))))
}

func (op *OpLog2) String() string {
	return fmt.Sprintf("( Log2 %s )", op.Children[0].String())
}

type OpSquare struct {
	BaseNode
}

func NewOpSquare() *OpSquare {
	return &OpSquare{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpSquare) Eval(x, y float32) float32 {
	value := op.Children[0].Eval(x, y)
	return value * value
}

func (op *OpSquare) String() string {
	return fmt.Sprintf("( Square %s ) ", op.Children[0].String())
}

type OpNegate struct {
	BaseNode
}

func NewOpNegate() *OpNegate {
	return &OpNegate{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpNegate) Eval(x, y float32) float32 {
	return -op.Children[0].Eval(x, y)
}

func (op *OpNegate) String() string {
	return fmt.Sprintf("( Negate %s )", op.Children[0].String())
}

type OpCeil struct {
	BaseNode
}

func NewOpCeil() *OpCeil {
	return &OpCeil{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpCeil) Eval(x, y float32) float32 {
	return float32(math.Ceil(float64(op.Children[0].Eval(x, y))))
}

func (op *OpCeil) String() string {
	return fmt.Sprintf("( Ceil %s )", op.Children[0].String())
}

type OpFloor struct {
	BaseNode
}

func NewOpFloor() *OpFloor {
	return &OpFloor{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpFloor) Eval(x, y float32) float32 {
	return float32(math.Floor(float64(op.Children[0].Eval(x, y))))
}

func (op *OpFloor) String() string {
	return fmt.Sprintf("( Floor %s )", op.Children[0].String())
}

type OpAbs struct {
	BaseNode
}

func NewOpAbs() *OpAbs {
	return &OpAbs{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpAbs) Eval(x, y float32) float32 {
	return float32(math.Abs(float64(op.Children[0].Eval(x, y))))
}

func (op *OpAbs) String() string {
	return fmt.Sprintf("( Abs %s )", op.Children[0].String())
}

type OpPlus struct {
	BaseNode
}

func NewOpPlus() *OpPlus {
	return &OpPlus{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpPlus) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) + op.Children[1].Eval(x, y)
}

func (op *OpPlus) String() string {
	return fmt.Sprintf("( + %s %s )", op.Children[0].String(), op.Children[1].String())
}

type OpMinus struct {
	BaseNode
}

func NewOpMinus() *OpMinus {
	return &OpMinus{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpMinus) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) - op.Children[1].Eval(x, y)
}

func (op *OpMinus) String() string {
	return fmt.Sprintf("( - %s %s )", op.Children[0].String(), op.Children[1].String())
}

type OpMult struct {
	BaseNode
}

func NewOpMult() *OpMult {
	return &OpMult{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpMult) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) * op.Children[1].Eval(x, y)
}

func (op *OpMult) String() string {
	return fmt.Sprintf("( * %s %s )", op.Children[0].String(), op.Children[1].String())
}

type OpDiv struct {
	BaseNode
}

func NewOpDiv() *OpDiv {
	return &OpDiv{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpDiv) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) / op.Children[1].Eval(x, y)
}

func (op *OpDiv) String() string {
	return fmt.Sprintf("( / %s %s )", op.Children[0].String(), op.Children[1].String())
}

type OpAtan2 struct {
	BaseNode
}

func NewOpAtan2() *OpAtan2 {
	return &OpAtan2{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(op.Children[0].Eval(x, y)), float64(op.Children[1].Eval(x, y))))
}

func (op *OpAtan2) String() string {
	return fmt.Sprintf("( Atan2 %s %s )", op.Children[0].String(), op.Children[1].String())
}

type OpX struct {
	BaseNode
}

func NewOpX() *OpX {
	return &OpX{BaseNode{nil, make([]Node, 0)}}
}

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return fmt.Sprintf("x")
}

type OpY struct {
	BaseNode
}

func NewOpY() *OpY {
	return &OpY{BaseNode{nil, make([]Node, 0)}}
}

func (op *OpY) Eval(x, y float32) float32 {
	return y
}

func (op *OpY) String() string {
	return fmt.Sprintf("t")
}

type OpConstant struct {
	BaseNode
	value float32
}

func NewOpConstant() *OpConstant {
	return &OpConstant{BaseNode{nil, make([]Node, 0)}, rand.Float32()*2 - 1}
}

func (op *OpConstant) Eval(x, y float32) float32 {
	return op.value
}

func (op *OpConstant) String() string {
	return strconv.FormatFloat(float64(op.value), 'f', 9, 32)
}

func GetRandomNode() Node {
	r := rand.Intn(16)
	switch r {
	case 0:
		return NewOpPlus()
	case 1:
		return NewOpMinus()
	case 2:
		return NewOpMult()
	case 3:
		return NewOpDiv()
	case 4:
		return NewOpAtan()
	case 5:
		return NewOpAtan2()
	case 6:
		return NewOpCeil()
	case 7:
		return NewOpFloor()
	case 8:
		return NewOpClip()
	case 9:
		return NewOpCos()
	case 10:
		return NewOpSin()
	case 11:
		return NewOpLerp()
	case 12:
		return NewOpLog2()
	case 13:
		return NewOpAbs()
	case 14:
		return NewOpNegate()
	case 15:
		return NewOpSquare()
	case 16:
		return NewOpWrap()
	default:
		return NewOpPlus()
	}
}

func GetRandomLeaf() Node {
	r := rand.Intn(3)
	switch r {
	case 0:
		return NewOpX()
	case 1:
		return NewOpY()
	case 2:
		return NewOpConstant()
	default:
		return NewOpX()
	}
}
