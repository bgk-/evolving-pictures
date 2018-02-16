package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bgk-/gameswithgo/evolving-pictures/apt"
	"github.com/gopherjs/gopherjs/js"
)

type rgba struct {
	r, g, b byte
}

// Picture with r,g,b apt Nodes
type Picture struct {
	r apt.Node
	g apt.Node
	b apt.Node
}

func (p *Picture) String() string {
	return fmt.Sprintf("R %s \nG %s \n B %s", p.r.String(), p.g.String(), p.b.String())
}

// NewPicture create pointer to new picture
func NewPicture() *Picture {
	p := &Picture{}

	p.r = apt.GetRandomNode()
	p.g = apt.GetRandomNode()
	p.b = apt.GetRandomNode()

	// Number of starting Nodes
	num := rand.Intn(4)
	for i := 0; i < num; i++ {
		p.r.AddRandom(apt.GetRandomNode())
	}

	num = rand.Intn(4)
	for i := 0; i < num; i++ {
		p.g.AddRandom(apt.GetRandomNode())
	}

	num = rand.Intn(4)
	for i := 0; i < num; i++ {
		p.b.AddRandom(apt.GetRandomNode())
	}

	for p.r.AddLeaf(apt.GetRandomLeaf()) {
	}
	for p.g.AddLeaf(apt.GetRandomLeaf()) {
	}
	for p.b.AddLeaf(apt.GetRandomLeaf()) {
	}

	return p
}

// Mutate changes a single r,g,b Node
func (p *Picture) Mutate() {
	r := rand.Intn(3)
	var nodeToMutate apt.Node
	switch r {
	case 0:
		nodeToMutate = p.r
	case 1:
		nodeToMutate = p.g
	case 2:
		nodeToMutate = p.b
	}

	count := nodeToMutate.NodeCount()
	r = rand.Intn(count)
	nodeToMutate, count = apt.GetNthNode(nodeToMutate, r, 0)
	mutation := apt.Mutate(nodeToMutate)

	if nodeToMutate == p.r {
		p.r = mutation
	} else if nodeToMutate == p.g {
		p.g = mutation
	} else if nodeToMutate == p.b {
		p.b = mutation
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func aptToTexture(p *Picture, w, h int) []byte {
	// -1.0 and 1.0
	scale := float32(255 / 2)
	offset := float32(-1.0 * scale)
	pixels := make([]byte, w*h*4)
	pixelIndex := 0
	for yi := 0; yi < h; yi++ {
		y := float32(yi)/float32(h)*2 - 1
		for xi := 0; xi < w; xi++ {
			x := float32(xi)/float32(w)*2 - 1

			r := p.r.Eval(x, y)
			g := p.g.Eval(x, y)
			b := p.b.Eval(x, y)

			pixels[pixelIndex] = byte(r*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(g*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(b*scale - offset)
			pixelIndex++
			pixelIndex++
		}
	}
	return pixels
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	width := js.Global.Get("innerWidth").Int()
	height := js.Global.Get("innerHeight").Int()
	fmt.Println(width, height)
	document := js.Global.Get("document")
	canvas := document.Call("getElementById", "picture")
	context := canvas.Call("getContext", "2d")
	canvas.Set("width", width)
	canvas.Set("height", height)

	imgData := context.Call("createImageData", width, height)
	pic := NewPicture()
	tex := aptToTexture(pic, width, height)
	for i := 0; i < width*height; i++ {
		data := imgData.Get("data")
		pD := i * 4
		data.SetIndex(pD+0, tex[pD+0])
		data.SetIndex(pD+1, tex[pD+1])
		data.SetIndex(pD+2, tex[pD+2])
		data.SetIndex(pD+3, 0xff)
	}
	context.Call("putImageData", imgData, 10, 10)

	js.Global.Get("picture").Call("addEventListener", "click", func() {
		go func() {
			pic.Mutate()
			tex = aptToTexture(pic, width, height)
			for i := 0; i < width*height; i++ {
				data := imgData.Get("data")
				pD := i * 4
				data.SetIndex(pD+0, tex[pD+0])
				data.SetIndex(pD+1, tex[pD+1])
				data.SetIndex(pD+2, tex[pD+2])
				data.SetIndex(pD+3, 0xff)
			}
			context.Call("putImageData", imgData, 10, 10)
		}()
	})

}
