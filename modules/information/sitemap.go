package information

import (
	"fmt"
	"gowebfuzz/modules"
	"net/http"
	"strings"
)

type SiteMap struct {
	modules.ModuleInfo
}

type SiteMapConfig struct {
}

type Tree struct {
	Root *TreeNode
	Host string
}

type TreeNode struct {
	Name     string
	NodeType string
	Get      string
	Nodes    map[string]*TreeNode
	GetList  []string
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

var sites map[string]*Tree

func (siteMap SiteMap) MatchRequest(req *http.Request) error {
	return nil
}

func (siteMap SiteMap) MatchResponse(req *http.Request, res *http.Response, body *[]byte) error {
	host := res.Header.Get("Host")
	_, ok := sites[host]
	if ok != true {
		siteNode := NewTree(host)
		sites[host] = siteNode
	}

	site, ok := sites[host]
	if res.StatusCode != 404 {
		site.AppendPath(req.RequestURI)
	}

	return nil
}
