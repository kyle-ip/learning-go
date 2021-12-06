package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

type node struct {
	isLast  bool              // 该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	segment string            // uri中的字符串
	handler ControllerHandler // 控制器
	childs  []*node           // 子节点
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root}
}

// isWildSegment 判断一个 segment 是否是通用 segment，即以 : 开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// filterChildNodes 过滤下一层满足 segment 规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	// 如果 segment 是通配符，则所有下一层子节点都满足需求。
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	// 过滤所有的下一层子节点。
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			// 如果下一层子节点有通配符，则满足需求。
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment {
			// 如果下一层子节点没有通配符，但是文本完全匹配，则满足需求。
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

// matchNode 判断路由是否已经在节点的所有子节点树中存在了
func (n *node) matchNode(uri string) *node {
	// 使用分隔符将 uri 切割为两个部分。
	segments := strings.SplitN(uri, "/", 2)
	// 第一个部分用于匹配下一层子节点。
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// 匹配符合的下一层子节点。
	cnodes := n.filterChildNodes(segment)

	// 如果当前子节点没有一个符合，说明这个 uri 之前不存在, 直接返回 nil
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	// 如果只有一个 segment，则是最后一个标记；否则递归每个子节点继续进行查找。
	if len(segments) == 1 {
		// 如果 segment 已经是最后一个节点，判断这些 cnode 是否有 isLast 标志
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}

		// 都不是最后一个节点。
		return nil
	} else {
		for _, tn := range cnodes {
			tnMatch := tn.matchNode(segments[1])
			if tnMatch != nil {
				return tnMatch
			}
		}
		return nil
	}
}

// AddRouter 增加路由节点, 路由节点有先后顺序
func (tree *Tree) AddRouter(uri string, handler ControllerHandler) error {
	// /book/list
	// /book/:id (冲突)
	// /book/:id/name
	// /book/:student/age
	// /:user/name
	// /:user/name/:age (冲突)
	n := tree.root
	if n.matchNode(uri) != nil {
		return errors.New("route exist: " + uri)
	}

	segments := strings.Split(uri, "/")
	// 对每个 segment
	for index, segment := range segments {

		// 最终进入 Node segment 的字段
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1

		var objNode *node // 标记是否有合适的子节点

		childNodes := n.filterChildNodes(segment)
		// 如果有匹配的子节点
		if len(childNodes) > 0 {
			// 如果有segment相同的子节点，则选择这个子节点
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			// 创建一个当前 node 的节点
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}

		n = objNode
	}

	return nil
}

// FindHandler 匹配 uri
func (tree *Tree) FindHandler(uri string) ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handler
}
