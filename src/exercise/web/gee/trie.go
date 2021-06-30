package gee

import (
    "fmt"
    "strings"
)

// node: node of trie
type node struct {
    // request path to be matched, like: /p/:lang
    pattern  string

    // a part of request path, like: :lang
    part     string

    // child node, like: doc, tutorial, intro
    children []*node

    // wild matching or not. while part contains ':' or '*', it equals to true.
    isWild   bool
}

func (n *node) String() string {
    return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// insert node to trie
func (n *node) insert(pattern string, parts []string, height int) {
    // only on the last level, pattern will be set to full path, like '/p/:lang/doc'.
    if len(parts) == height {
        n.pattern = pattern
        return
    }
    // on another level (like 'p' and ':lang') pattern is "".

    // get child node's part.
    part := parts[height]

    // try to find child node under n. if no matched child, create new node and add to 'children'.
    child := n.matchChild(part)
    if child == nil {
        child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
        n.children = append(n.children, child)
    }

    // insert next level part of the patterns to the child node.
    child.insert(pattern, parts, height+1)
}

// search
func (n *node) search(parts []string, height int) *node {
    if len(parts) == height || strings.HasPrefix(n.part, "*") {
        // if n.pattern equals to "", it means matching failed and return nil.
        if n.pattern == "" {
            return nil
        }
        return n
    }

    // recursive search in children, which are matched the current part.
    children := n.matchChildren(parts[height])
    for _, child := range children {
        result := child.search(parts, height+1)
        if result != nil {
            return result
        }
    }

    return nil
}

// travel returns a list which contains all nodes of n.
func (n *node) travel(list *[]*node) {
    if n.pattern != "" {
        *list = append(*list, n)
    }
    for _, child := range n.children {
        child.travel(list)
    }
}

// matchChild returns the first node which is successfully matched in part.
func (n *node) matchChild(part string) *node {
    for _, child := range n.children {
        if child.part == part || child.isWild {
            return child
        }
    }
    return nil
}

// matchChildren returns all nodes which are successfully matched, so as to support searching.
func (n *node) matchChildren(part string) []*node {
    nodes := make([]*node, 0)
    for _, child := range n.children {
        if child.part == part || child.isWild {
            nodes = append(nodes, child)
        }
    }
    return nodes
}
