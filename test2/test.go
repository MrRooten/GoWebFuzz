package main

import (
	"fmt"
	"strings"
)

type Tree struct {
	Root *TreeNode
	Host string
}

type TreeNode struct {
	Name     string
	NodeType string
	Get      string
	Nodes    map[string]*TreeNode
}

func (tree *Tree) FindPath(path string) bool {
	return true
}

func (tree *Tree) FindNode(name string) int {
	return -1
}

func (tree *Tree) Append(dir *TreeNode) {

}

func _PrintNode(node *TreeNode, n int) {
	for i := 0; i < n; i++ {
		fmt.Print("  ")
	}
	fmt.Println(node.Name)
	for _, v := range node.Nodes {
		_PrintNode(v, n+1)
	}
}

func (tree *Tree) PrintTree() {
	_PrintNode(tree.Root, 0)
}
func (tree *Tree) AppendPath(path string) *TreeNode {
	root := tree.Root
	GetDelimIndex := strings.Index(path, "?")
	GetValue := ""
	var requestPath string
	if GetDelimIndex == -1 {
		GetValue = ""
		requestPath = path
	} else {
		GetValue = path[GetDelimIndex+1:]
		requestPath = path[0 : GetDelimIndex-1]
	}

	dirs := strings.Split(requestPath, "/")
	curNode := root
	dirsLength := len(dirs)
	var result *TreeNode
	for index, dir := range dirs {
		node, ok := curNode.Nodes[dir]
		if ok == false {
			node = &TreeNode{
				Name:     dir,
				NodeType: "node",
			}
			if curNode.Nodes == nil {
				curNode.Nodes = make(map[string]*TreeNode)
			}
			curNode.Nodes[dir] = node
		}
		if index == dirsLength-1 {
			curNode.Get = GetValue
			result = curNode
		}
		curNode = node
	}

	return result
}

func NewTree(name string) *Tree {
	tree := Tree{}
	tree.Root = &TreeNode{}
	tree.Root.Name = name
	tree.Root.Nodes = make(map[string]*TreeNode)
	return &tree
}
func main() {
	tree := NewTree("http://baidu.com")
	tree.AppendPath("abc/123/def/456")
	tree.AppendPath("abc/123/vds/456")
	tree.PrintTree()
}
