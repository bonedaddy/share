// Copyright 2019, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbt

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/shuLhan/share/lib/test"
)

func TestInsertDelete(t *testing.T) {
	var (
		exps []string
		expb bytes.Buffer
		x    byte  = 'a'
		tree *Tree = &Tree{}
	)

	exps = append(exps, "[]")

	expb.WriteByte('[')
	for ; x <= 'z'; x++ {
		if x > 'a' {
			expb.WriteByte(' ')
		}

		node := NewNode([]byte{x}, nil)

		tree.Insert(node)

		expb.WriteByte(x)
		exp := append(expb.Bytes(), ']')
		exps = append(exps, string(exp))

		test.Assert(t, "exp", exps[x-'a'+1], tree.String(), true)
	}

	t.Logf("Tree.Keys: %s\n", tree.Keys())

	for x = 'z'; x >= 'a'; x-- {
		tree.Delete([]byte{x})

		test.Assert(t, "exp", exps[x-'a'], tree.String(), true)
	}
}

func TestInsert1024(t *testing.T) {
	_keys = generateKeys(4, 51)

	for x := 0; x < len(_keys); x++ {
		node := NewNode(_keys[x], struct{}{})
		tree.Insert(node)
	}

	t.Logf("len(_keys):%d tree.n:%d\n", len(_keys), tree.n)

	for x := 0; x < len(_keys); x++ {
		node := tree.Find(_keys[x])
		if node == nil {
			t.Fatalf("%s not found\n", _keys[x])
		}
	}

	sortedKeys := tree.Keys()

	for x := 0; x < len(_keys); x++ {
		fmt.Printf("---- BEFORE deleting %s:\n", _keys[x])
		tree.Print()
		tree.Delete(_keys[x])
		fmt.Printf("---- AFTER deleting %s:\n", _keys[x])
		tree.Print()
		got := tree.Keys()
		sortedKeys = sliceRemove(sortedKeys, _keys[x])
		if !sliceEqual(sortedKeys, got) {
			t.Fatal("tree delete mis-balanced")
		}
	}

	test.Assert(t, "Tree.String", "[]", tree.String(), true)
}

func sliceRemove(keys [][]byte, key []byte) [][]byte {
	for x := 0; x < len(keys); x++ {
		if bytes.Compare(keys[x], key) == 0 {
			keys = append(keys[:x], keys[x+1:]...)
			return keys
		}
	}

	return keys
}

func sliceEqual(a [][]byte, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for x := 0; x < len(a); x++ {
		if bytes.Compare(a[x], b[x]) != 0 {
			fmt.Printf("mis-balanced at : %s %s\n", a[x], b[x])
			return false
		}
	}
	return true
}
