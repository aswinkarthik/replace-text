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
	value    string
}

// ErrPrefixConflict represents an error where the prefix of the
// added string is already present in the trie.
// E.g if "hell" is already added, adding "hello" has a common prefix "hell"
// This causes ambiguity when finding for words on when to quit.
// Hence, this error is returned on AddString
var ErrPrefixConflict = fmt.Errorf("conflict: a prefix of the world already exists")

// ErrContainsConflict represents an error where the given added string
// is contained in a string in the trie.
// E.g if "hello" is already added, then adding "hell" becomes ambigous
// as it is contained inside the existing trie.
var ErrContainsConflict = fmt.Errorf("conflict: a prefix of the world already exists")

// AddString will add the given string into the Trie structure
// It marks the node of the last edge as terminal
func (n *Node) AddString(s string) error {
	return n.insertString(s, "")
}

// Terminates returns true if the node is a terminal node
// and false otherwise.
func (n *Node) Terminates() bool {
	return n.terminal
}

// NewNode is a constructor to create a new node
func NewNode() *Node {
	return &Node{
		next: make(map[byte]*Node, 1),
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

// AddReplacement provides a way to define what string
// needs to be found and what to be replaced with.
// The string to be found is the path.
// The string to be replaced is stored at the terminal node.
func (n *Node) AddReplacement(old, new string) error {
	return n.insertString(old, new)
}

func (n *Node) insertString(path string, leafValue string) error {
	if len(path) == 0 {
		return fmt.Errorf("empty string not accepted")
	}

	ch := path[0]
	restOfString := path[1:]
	nextNode, nextNodeExists := n.next[ch]
	if !nextNodeExists {
		nextNode = NewNode()
	}
	n.next[ch] = nextNode

	lastCharacter := len(path) == 1

	// Last character but is not a new node to be created
	if lastCharacter && nextNodeExists {
		return ErrContainsConflict
	}

	// Not the last character, but found a terminal node already
	if !lastCharacter && nextNode.Terminates() {
		return ErrPrefixConflict
	}

	// Last character and a new node
	if lastCharacter {
		nextNode.terminal = true
		nextNode.value = leafValue
		return nil
	}

	// Recurse the rest of the string otherwise
	return nextNode.AddString(restOfString)
}
