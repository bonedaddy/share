// Copyright 2019, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Package rbt provide a library for working with Red-Black tree.
//
package rbt

import (
	"bytes"
	"fmt"
	"strings"
)

//
// Tree contains the root.
//
type Tree struct {
	root *Node
	n    int
}

//
// Delete a node from tree by their key.
// If key found, the removed node will be returned; otherwise it will return
// nil.
//
func (tree *Tree) Delete(key []byte) *Node {
	if len(key) == 0 {
		return nil
	}

	x := tree.Find(key)
	if x == nil {
		return nil
	}

	return tree.remove(x)
}

//
// Find a node in tree that match with the key.
//
func (tree *Tree) Find(key []byte) (node *Node) {
	if len(key) == 0 {
		return nil
	}

	node = tree.root
	for node != nil {
		switch bytes.Compare(key, node.key) {
		case -1:
			node = node.left
		case 1:
			node = node.right
		default:
			return node
		}
	}

	return nil
}

//
// Insert a node to the tree.
//
func (tree *Tree) Insert(node *Node) {
	if tree.root == nil {
		tree.n++
		tree.root = node
		tree.root.color = colorBlack
		return
	}

	p := tree.root
	for {
		switch bytes.Compare(node.key, p.key) {
		case -1:
			if p.left == nil {
				p.left = node
				goto out
			}
			p = p.left
		case 1:
			if p.right == nil {
				p.right = node
				goto out
			}
			p = p.right
		default:
			// Replace the value if key matched.
			p.Value = node.Value
			return
		}
	}
out:
	tree.n++
	node.parent = p
	tree.insertFixup(node)
}

//
// Keys return all keys sorted in ascending order.
//
func (tree *Tree) Keys() (keys [][]byte) {
	if tree.root == nil {
		return
	}

	keys = make([][]byte, 0, tree.n)
	stack := make([]*Node, 0)
	node := tree.root

	for node != nil {
		for node.left != nil {
			stack = append(stack, node)
			node = node.left
		}

		keys = append(keys, node.key)

		if node.right != nil {
			node = node.right
			continue
		}

		for node.right == nil {
			if len(stack) == 0 {
				node = nil
				return
			}
			node = stack[len(stack)-1]
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
			keys = append(keys, node.key)
		}
		node = node.right
	}

	return
}

//
// Print each node in tree to standard output with identation based on node
// level.
//
func (tree *Tree) Print() {
	if tree.root == nil {
		return
	}

	node := tree.root
	stack := make([]*Node, 0)

	// We use for-loop to minimize call stack.
	for node != nil {
		for node.left != nil {
			stack = append(stack, node)
			node = node.left
		}

		printNode(node)

		if node.right != nil {
			node = node.right
			continue
		}

		for node.right == nil {
			if len(stack) == 0 {
				node = nil
				return
			}
			node = stack[len(stack)-1]
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
			printNode(node)
		}
		node = node.right
	}
}

//
// String return the text representation of tree like a slice, in ascending
// order.
//
func (tree *Tree) String() string {
	var sb strings.Builder

	node := tree.root
	stack := make([]*Node, 0)

	sb.WriteByte('[')

	// We use for-loop to minimize call stack.
	for node != nil {
		for node.left != nil {
			stack = append(stack, node)
			node = node.left
		}

		if sb.Len() > 1 {
			sb.WriteByte(' ')
		}
		sb.Write(node.key)

		if node.right != nil {
			node = node.right
			continue
		}

		for node.right == nil {
			if len(stack) == 0 {
				node = nil
				goto out
			}

			node = stack[len(stack)-1]
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]

			if sb.Len() > 1 {
				sb.WriteByte(' ')
			}
			sb.Write(node.key)
		}
		node = node.right
	}
out:
	sb.WriteByte(']')

	return sb.String()
}

func printNode(node *Node) {
	for p := node; p != nil; p = p.parent {
		fmt.Print(" ")
	}
	if node.parent != nil {
		if node.parent.left == node {
			fmt.Printf("/:")
		} else {
			fmt.Printf("\\:")
		}
	}
	fmt.Printf("%d:%s\n", node.color, node.key)
}

func (tree *Tree) insertFixup(node *Node) {
	var parent, gp, aunt *Node

	parent = node.parent

	for parent != nil && parent.color == colorRed {
		gp = parent.parent

		if gp == nil {
			break
		}

		if gp.left == parent {
			aunt = gp.right

			if aunt != nil && aunt.color == colorRed {
				parent.color = colorBlack
				aunt.color = colorBlack
				gp.color = colorRed

				node = gp
			} else {
				if parent.right == node {
					node = parent
					tree.rotateLeft(node)
				}

				node.parent.color = colorBlack
				gp.color = colorRed

				tree.rotateRight(gp)
			}
		} else {
			aunt = gp.left

			if aunt != nil && aunt.color == colorRed {
				parent.color = colorBlack
				aunt.color = colorBlack
				gp.color = colorRed

				node = gp
			} else {
				if parent.left == node {
					node = parent
					tree.rotateRight(node)
				}

				node.parent.color = colorBlack
				gp.color = colorRed

				tree.rotateLeft(gp)
			}
		}

		parent = node.parent
	}

	tree.root.color = colorBlack
}

func (tree *Tree) deleteFixup(x *Node) {
	var parent, sibling, siblingl, siblingr *Node

	for x != tree.root && x.color == colorBlack {
		parent = x.parent

		if x == parent.left {
			sibling = parent.right

			if sibling == nil {
				break
			}

			if sibling.color == colorRed {
				tree.rotateLeft(parent)

				sibling.color = colorBlack

				if parent.right != nil {
					parent.right.color = colorRed
				}
				return
			}

			siblingr = sibling.right
			if siblingr != nil && siblingr.color == colorRed {
				tree.rotateLeft(parent)

				if sibling.haveRedChilds() {
					sibling.color = colorRed
				} else {
					sibling.color = colorBlack
				}
				sibling.setChildsToBlack()

				return
			}

			siblingl = sibling.left
			if siblingl != nil && siblingl.color == colorRed {
				tree.rotateRight(sibling)
				sibling.swapColor(siblingl)
				continue
			}

			if sibling.haveNoChilds() || sibling.haveBlackChilds() {
				sibling.color = colorRed

				if parent.color == colorRed {
					parent.color = colorBlack
					return
				}

				if parent == tree.root {
					return
				}
			}
		} else {
			sibling = parent.left

			if sibling == nil {
				break
			}

			if sibling.color == colorRed {
				tree.rotateRight(parent)

				sibling.color = colorBlack

				if parent.left != nil {
					parent.left.color = colorRed
				}
				return
			}

			siblingl = sibling.left
			if siblingl != nil && siblingl.color == colorRed {
				tree.rotateRight(parent)

				if sibling.haveRedChilds() {
					sibling.color = colorRed
				} else {
					sibling.color = colorBlack
				}
				sibling.setChildsToBlack()
				return
			}

			siblingr = sibling.right
			if siblingr != nil && siblingr.color == colorRed {
				tree.rotateLeft(sibling)
				sibling.swapColor(siblingr)
				continue
			}

			if sibling.haveNoChilds() || sibling.haveBlackChilds() {
				sibling.color = colorRed

				if parent.color == colorRed {
					parent.color = colorBlack
					return
				}

				if parent == tree.root {
					return
				}
			}
		}

		x = parent
	}

	tree.root.setChildsToBlack()
}

func (tree *Tree) remove(x *Node) *Node {
	if x == nil {
		return nil
	}

	left := x.left
	right := x.right

	if left == nil && right == nil {
		tree.n--
		return tree.removeHaveNoChilds(x)
	}
	if left != nil && right == nil {
		x.swapContent(left)
		return tree.remove(left)
	}
	if left == nil && right != nil {
		x.swapContent(right)
		return tree.remove(right)
	}

	return tree.removeHaveBothChilds(x)
}

func (tree *Tree) removeHaveBothChilds(x *Node) *Node {
	heir := x.left

	for heir.right != nil {
		heir = heir.right
	}

	x.swapContent(heir)

	return tree.remove(heir)
}

func (tree *Tree) removeHaveNoChilds(x *Node) *Node {
	if x.parent == nil {
		tree.root = nil
		x.detach()
		return x
	}

	if x.color == colorBlack {
		tree.deleteFixup(x)
	}

	if x.parent.left == x {
		x.parent.left = nil
	} else {
		x.parent.right = nil
	}

	x.detach()

	return x
}

//
// rotateLeft will move the right-child of `x` as parent of `x`.
// if right-child have left-child, then it will be the right-child of `x`.
//
//       ...               ...
//        |                 |
//        x                right
//       / \               /
//    ...   right  ==>    x
//          /            / \
//       rleft         ..   rleft
//
func (tree *Tree) rotateLeft(x *Node) {
	// Make the right-child as parent of `x` first.
	x.right.parent = x.parent

	if x.parent == nil {
		tree.root = x.right
	} else {
		if x.parent.left == x {
			x.parent.left = x.right
		} else {
			x.parent.right = x.right
		}
	}
	x.parent = x.right

	// Replace the x's right with right-left.
	x.right = x.parent.left
	x.parent.left = x

	if x.right != nil {
		x.right.parent = x
	}
}

//
// rotateRight will move the left-child of `x` as parent of `x`.
// if left-child have right-child, then it will be the left-child of `x`.
//
//       ...               ...
//        |                 |
//        x                left
//       / \                 \
//    left  ...  ==>          x
//       \                   / \
//        lright        lright  ...
//
func (tree *Tree) rotateRight(x *Node) {
	// Move the left as parent of `x` first.
	x.left.parent = x.parent

	if x.parent == nil {
		tree.root = x.left
	} else {
		if x.parent.left == x {
			x.parent.left = x.left
		} else {
			x.parent.right = x.left
		}
	}
	x.parent = x.left

	// Replace the x's left with left-right.
	x.left = x.parent.right
	x.parent.right = x

	if x.left != nil {
		x.left.parent = x
	}
}
