package information

import (
	"gowebfuzz/library/modules"
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
}

func (tree *Tree) FindPath(path string) bool {
	return true
}

func (tree *Tree) FindNode(name string) int {
	return -1
}

func (tree *Tree) Append(dir *TreeNode) {

}

func (tree *Tree) AppendPath(path string) *TreeNode {
	root := tree.Root
	GetDelimIndex := strings.Index(path, "?")
	GetValue := ""
	if GetDelimIndex == -1 {
		GetValue = ""
	} else {
		GetValue = path[GetDelimIndex+1:]
	}

	requestPath := path[0 : GetDelimIndex-1]
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

var sites map[string]Tree

func (siteMap SiteMap) MatchRequest(req *http.Request) error {
	return nil
}

func (siteMap SiteMap) MatchResponse(req *http.Request, res *http.Response, body *[]byte) error {
	host := res.Header.Get("Host")
	_, ok := sites[host]
	if ok != true {
		siteNode := Tree{
			Host: host,
		}
		sites[host] = siteNode
	}

	site, ok := sites[host]
	if res.StatusCode != 404 {
		site.AppendPath(req.RequestURI)
	}

	return nil
}