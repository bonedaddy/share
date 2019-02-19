// Copyright 2019, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbt

//
// Node in tree.
//
type Node struct {
	Value interface{}

	parent *Node
	left   *Node
	right  *Node
	key    []byte
	color  color
}

//
// NewNode create a new node using key and value.
//
func NewNode(key []byte, value interface{}) *Node {
	if len(key) == 0 {
		return nil
	}
	return &Node{
		Value: value,
		key:   key,
		color: colorRed,
	}
}

func (node *Node) detach() {
	node.parent = nil
	node.left = nil
	node.right = nil
}

func (node *Node) haveBlackChilds() bool {
	if node.left == nil || node.right == nil {
		return false
	}
	if node.left.color != colorBlack || node.right.color != colorBlack {
		return false
	}
	return true
}

func (node *Node) haveNoChilds() bool {
	if node.left == nil && node.right == nil {
		return true
	}
	return false
}

func (node *Node) haveRedChilds() bool {
	if node.left == nil || node.right == nil {
		return false
	}
	if node.left.color != colorRed || node.right.color != colorRed {
		return false
	}
	return true
}

func (node *Node) setChildsToBlack() {
	if node.left != nil {
		node.left.color = colorBlack
	}
	if node.right != nil {
		node.right.color = colorBlack
	}
}

func (node *Node) swapColor(x *Node) {
	tmp := node.color
	node.color = x.color
	x.color = tmp
}

func (node *Node) swapContent(x *Node) {
	key := node.key
	val := node.Value
	node.key = x.key
	node.Value = x.Value
	x.key = key
	x.Value = val
}
