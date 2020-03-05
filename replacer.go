package replacer

import (
	"encoding/json"
	"fmt"
)

// Node structure to hold each vertex
// It can be either terminal or intermediate
// All character bytes are stored on edges
// Each key of map next is an edge leading to the next node
type Node struct {
	terminal bool
	next     map[byte]*Node
}

// AddString will add the given string into the Trie structure
// It marks the node of the last edge as terminal
func (n *Node) AddString(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("empty string not accepted")
	}

	ch := s[0]
	restOfString := s[1:]
	nextNode, exists := n.next[ch]
	if !exists {
		nextNode = &Node{next: make(map[byte]*Node)}
	}
	n.next[ch] = nextNode

	if len(s) == 1 {
		nextNode.terminal = true
		return nil
	}

	return nextNode.AddString(restOfString)
}

// Terminates returns true if the node is a terminal node
// and false otherwise.
func (n *Node) Terminates() bool {
	return n.terminal
}

// NewNode is a constructor to create a new node
func NewNode() *Node {
	return &Node{
		next: make(map[byte]*Node),
	}
}

// MarshalJSON is implemented to conform to Marshaler interface.
// It prints a JSON in a repeated manner till terminal node.
// Terminal node is marked with "terminal": true
//
// For example string: "wax"
//
// { "w": { "a": { "x": { "terminal": true } } } }
func (n *Node) MarshalJSON() ([]byte, error) {
	readableMap := make(map[string]interface{})
	for k, v := range n.next {
		readableMap[string(k)] = v
	}

	if n.terminal {
		readableMap["terminal"] = true
	}

	return json.Marshal(readableMap)
}

// Contains tests if the string k is present inside the
// trie structure.
func (n *Node) Contains(k string) bool {
	if len(k) == 0 {
		return false
	}

	ch := k[0]
	restOfString := k[1:]
	nextNode, exists := n.next[ch]
	if !exists {
		return false
	}

	if nextNode.Terminates() {
		return true
	}
	return nextNode.Contains(restOfString)
}
