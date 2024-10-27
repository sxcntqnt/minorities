package main

import (
	"fmt"
	"sync"
)

// TreeAlgorithm interface for tree algorithms
type TreeAlgorithm interface {
	Search(value interface{}) (bool, error) // Change to interface{} to allow flexibility
	Print() error
}

// Splay Tree implementation
type SplayTree struct {
	root *SplayNode
}

type SplayNode struct {
	value int
	left  *SplayNode
	right *SplayNode
}

// Insert method for Splay Tree
func (st *SplayTree) Insert(value int) {
	st.root = st.insert(st.root, value)
}

func (st *SplayTree) insert(node *SplayNode, value int) *SplayNode {
	if node == nil {
		return &SplayNode{value: value}
	}
	if value < node.value {
		node.left = st.insert(node.left, value)
	} else if value > node.value {
		node.right = st.insert(node.right, value)
	}
	return node
}

// Implementing Search method for Splay Tree
func (st *SplayTree) Search(value interface{}) (bool, error) {
	intValue, ok := value.(int) // Type assertion to int
	if !ok {
		return false, fmt.Errorf("value must be an int")
	}
	if st.root == nil {
		return false, fmt.Errorf("tree is empty")
	}
	st.root = st.splay(st.root, intValue) // Splay the tree before searching
	return st.root != nil && st.root.value == intValue, nil
}

// Implementing Splay operation
func (st *SplayTree) splay(node *SplayNode, value int) *SplayNode {
	if node == nil || node.value == value {
		return node
	}

	if value < node.value {
		if node.left == nil {
			return node
		}
		if value < node.left.value {
			node.left.left = st.splay(node.left.left, value)
			node = rightRotate(node)
		} else if value > node.left.value {
			node.left.right = st.splay(node.left.right, value)
			if node.left.right != nil {
				node.left = leftRotate(node.left)
			}
		}
		return rightRotate(node)
	} else {
		if node.right == nil {
			return node
		}
		if value > node.right.value {
			node.right.right = st.splay(node.right.right, value)
			node = leftRotate(node)
		} else if value < node.right.value {
			node.right.left = st.splay(node.right.left, value)
			if node.right.left != nil {
				node.right = rightRotate(node.right)
			}
		}
		return leftRotate(node)
	}
}

// Rotate functions for Splay Tree
func rightRotate(y *SplayNode) *SplayNode {
	x := y.left
	y.left = x.right
	x.right = y
	return x
}

func leftRotate(x *SplayNode) *SplayNode {
	y := x.right
	x.right = y.left
	y.left = x
	return y
}

// Implementing Print method for Splay Tree
func (st *SplayTree) Print() error {
	if st.root == nil {
		return fmt.Errorf("tree is empty")
	}
	st.printInOrder(st.root)
	fmt.Println()
	return nil
}

func (st *SplayTree) printInOrder(node *SplayNode) {
	if node != nil {
		st.printInOrder(node.left)
		fmt.Print(node.value, " ")
		st.printInOrder(node.right)
	}
}

// Minimum Spanning Tree implementation using Prim's algorithm
type MST struct {
	edges []Edge
}

// Edge structure for graph edges
type Edge struct {
	src, dest, weight int
}

// AddEdge method for MST
func (m *MST) AddEdge(src, dest, weight int) {
	m.edges = append(m.edges, Edge{src, dest, weight})
}

// Print method for MST
func (m *MST) Print() error {
	if len(m.edges) == 0 {
		return fmt.Errorf("no edges in the MST")
	}
	for _, edge := range m.edges {
		fmt.Printf("Edge from %d to %d with weight %d\n", edge.src, edge.dest, edge.weight)
	}
	return nil
}

// Search method for MST (just a placeholder, you may implement it as needed)
func (m *MST) Search(value interface{}) (bool, error) {
	// Placeholder implementation, adjust based on how you want to search
	return false, fmt.Errorf("search not implemented for MST")
}

// Suffix Tree implementation
type SuffixTree struct {
	root *SuffixTreeNode
}

type SuffixTreeNode struct {
	children map[rune]*SuffixTreeNode
	end      bool
}

// NewSuffixTree initializes a new Suffix Tree
func NewSuffixTree() *SuffixTree {
	return &SuffixTree{root: &SuffixTreeNode{children: make(map[rune]*SuffixTreeNode)}}
}

// Insert method for Suffix Tree
func (st *SuffixTree) Insert(s string) {
	for i := 0; i < len(s); i++ {
		st.insertSuffix(s[i:])
	}
}

func (st *SuffixTree) insertSuffix(suffix string) {
	node := st.root
	for _, char := range suffix {
		if _, exists := node.children[char]; !exists {
			node.children[char] = &SuffixTreeNode{children: make(map[rune]*SuffixTreeNode)}
		}
		node = node.children[char]
	}
	node.end = true
}

// Search method for Suffix Tree
func (st *SuffixTree) Search(value interface{}) (bool, error) {
	strValue, ok := value.(string) // Type assertion to string
	if !ok {
		return false, fmt.Errorf("value must be a string")
	}
	node := st.root
	for _, char := range strValue {
		if _, exists := node.children[char]; !exists {
			return false, nil
		}
		node = node.children[char]
	}
	return node.end, nil
}

// Print method for Suffix Tree
func (st *SuffixTree) Print() error {
	if st.root == nil {
		return fmt.Errorf("tree is empty")
	}
	st.printNode(st.root, "")
	return nil
}

func (st *SuffixTree) printNode(node *SuffixTreeNode, prefix string) {
	if node.end {
		fmt.Println(prefix)
	}
	for char, child := range node.children {
		st.printNode(child, prefix+string(char))
	}
}

// Main function
func main() {
	var wg sync.WaitGroup

	trees := []TreeAlgorithm{
		&SplayTree{},
		&MST{},
		NewSuffixTree(),
	}

	// Add data to trees (example)
	splayTree := trees[0].(*SplayTree)
	splayTree.Insert(10)
	splayTree.Insert(20)
	splayTree.Insert(5)

	mst := trees[1].(*MST)
	mst.AddEdge(1, 2, 3)
	mst.AddEdge(2, 3, 4)

	suffixTree := trees[2].(*SuffixTree)
	suffixTree.Insert("hello")
	suffixTree.Insert("world")

	for _, tree := range trees {
		wg.Add(1)
		go func(t TreeAlgorithm) {
			defer wg.Done()
			if err := t.Print(); err != nil {
				fmt.Println("Error printing tree:", err)
			}
			// Use different values for search based on tree type
			switch t := t.(type) {
			case *SplayTree:
				if found, err := t.Search(5); err != nil {
					fmt.Println("Error searching tree:", err)
				} else {
					fmt.Printf("Search result for value 5 in SplayTree: %v\n", found)
				}
			case *SuffixTree:
				if found, err := t.Search("hello"); err != nil {
					fmt.Println("Error searching tree:", err)
				} else {
					fmt.Printf("Search result for string 'hello' in SuffixTree: %v\n", found)
				}
			case *MST:
				if found, err := t.Search(1); err != nil {
					fmt.Println("Error searching tree:", err)
				} else {
					fmt.Printf("Search result for value 1 in MST: %v\n", found)
				}
			}
		}(tree)
	}

	wg.Wait()
	fmt.Println("All operations completed.")
}
