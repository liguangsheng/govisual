package govisual

import "strings"

type value interface {
	Hash() string
}

type Node struct {
	Value   value
	In, Out int
}

type Edge struct {
	From, To *Node
	Weight   int
}

func NewDigraph() *Digraph {
	return &Digraph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string]*Edge),
		froms: make(map[string]*Node),
	}
}

type Digraph struct {
	Nodes map[string]*Node
	Edges map[string]*Edge
	froms map[string]*Node
}

// HasFrom returns true if Digraph contains From node
func (g *Digraph) HasFrom(hash string) bool {
	_, has := g.froms[hash]
	return has
}

// Node returns node From Digraph, if not exist, add it.
func (g *Digraph) Node(v value) *Node {
	return g.getNode(v)
}

// Edge add a edge To Digraph
func (g *Digraph) Edge(from, to value) {
	fromNode := g.getNode(from)
	fromNode.Out += 1
	g.froms[from.Hash()] = fromNode

	toNode := g.getNode(to)
	toNode.In += 1

	g.getEdge(fromNode, toNode).Weight += 1
}

func (g *Digraph) getNode(v value) *Node {
	node, exist := g.Nodes[v.Hash()]
	if !exist {
		node = &Node{Value: v}
		g.Nodes[v.Hash()] = node
	}
	return node
}

func (g *Digraph) getEdge(from, to *Node) *Edge {
	hash := g.hashEdge(from.Value.Hash(), to.Value.Hash())
	edge, exist := g.Edges[hash]
	if !exist {
		edge = &Edge{From: from, To: to, Weight: 1}
		g.Edges[hash] = edge
	}
	return edge
}

func (g *Digraph) hashEdge(from, to string) string {
	return strings.Join([]string{from, to}, "->")
}
