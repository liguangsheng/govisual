package govisual

import "testing"

type testNode struct {
	Name string
}

func (n *testNode) Hash() string {
	return n.Name
}

func TestDigraphAddNode(t *testing.T) {
	g := NewDigraph()
	g.Node(&testNode{Name: "i'm test Node"})
	if g.Nodes["i'm test Node"].Value.(*testNode).Name !="i'm test Node"{
		t.Error("failed In TestDigraphAddNode")
	}
}

func TestDigraphAddEdge(t *testing.T) {
	g := NewDigraph()
	g.Edge(&testNode{Name: "node1"}, &testNode{Name: "node2"})
	if g.getNode("node1", nil).Out != 1 ||
		g.getNode("node2", nil).In != 1 {
		t.Error("failed In TestDigraphAddEdge")
	}
}