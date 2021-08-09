package router

import (
	"path/filepath"
	"strings"
)

type Node struct {
	MNodes map[string][]*Node `json:"m_nodes"`
	PNodes map[string][]*Node `json:"p_nodes"`
}

func getParams(f string, p string) map[string]string {
	result := make(map[string]string)
	item1 := strings.Split(strings.Trim(f, "/"), "/")
	item2 := strings.Split(strings.Trim(p, "/"), "/")
	if len(item1) != len(item2) {
		return result
	}
	for k, v := range item1 {
		if strings.HasPrefix(v, ":") {
			result[v[1:]] = item2[k]
		}
	}
	return result

}

func (n *Node) Match(p string, px string) (string, bool) {
	if len(p) == 0 {
		return filepath.Join("/", px), len(n.MNodes) == 0 && len(n.PNodes) == 0
	}
	item := strings.SplitN(strings.Trim(p, "/"), "/", 2)

	for k, nodes := range n.PNodes {
		for _, mn := range nodes {
			if p, ok := mn.Match(strings.Join(item[1:], "/"), filepath.Join(px, k)); ok {
				return p, ok
			}
		}
	}

	if nodes, ok := n.MNodes[item[0]]; ok && len(nodes) > 0 {
		for _, mn := range nodes {
			if p, ok := mn.Match(strings.Join(item[1:], "/"), filepath.Join(px, item[0])); ok {
				return p, ok
			}
		}
	}
	return "", false
}

func NewTree(path ...string) *Node {
	tree := &Node{MNodes: make(map[string][]*Node), PNodes: make(map[string][]*Node)}
	for _, p := range path {
		n := strings.Split(strings.Trim(p, "/"), "/")
		ctree := tree
		current := ctree.MNodes
		for _, k := range n {
			nn := &Node{
				MNodes: make(map[string][]*Node),
				PNodes: make(map[string][]*Node),
			}

			if strings.HasPrefix(k, ":") {
				current = ctree.PNodes
			}

			_, ok := current[k]
			if ok {
				current[k] = append(current[k], nn)
			} else {
				current[k] = []*Node{nn}
			}
			ctree = nn
			current = nn.MNodes
		}
	}
	return tree
}
