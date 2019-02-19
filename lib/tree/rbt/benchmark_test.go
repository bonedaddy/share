// Copyright 2019, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbt

import (
	"math/rand"
	"testing"
	"time"

	libbytes "github.com/shuLhan/share/lib/bytes"
)

var (
	_keys [][]byte
	m     = make(map[string]*Node)
	tree  = &Tree{}
)

func generateKeys(l, n int) (keys [][]byte) {
	rand.Seed(time.Now().Unix())

	keys = make([][]byte, 0, n)
	m = make(map[string]*Node, n)

	seeds := []byte(libbytes.Numbers)
	for x := 0; x < n; x++ {
		key := libbytes.Random(seeds, l)
		keys = append(keys, key)
	}

	return keys
}

func BenchmarkInsert(b *testing.B) {
	_keys = generateKeys(5, 10000000)
	b.ResetTimer()

	for x := 0; x < len(_keys); x++ {
		node := NewNode(_keys[x], struct{}{})
		tree.Insert(node)
	}

	b.Logf("len(keys):%d tree.n:%d\n", len(_keys), tree.n)
}

func BenchmarkMapInsert(b *testing.B) {
	for x := 0; x < len(_keys); x++ {
		key := string(_keys[x])
		m[key] = NewNode(_keys[x], struct{}{})
	}

	b.Logf("len(keys):%d len(map):%d\n", len(_keys), len(m))
}

func BenchmarkFind(b *testing.B) {
	for x := 0; x < len(_keys); x++ {
		tree.Find(_keys[x])
	}
}

func BenchmarkMapFind(b *testing.B) {
	for x := 0; x < len(_keys); x++ {
		key := string(_keys[x])
		_, _ = m[key]
	}
}

func BenchmarkDelete(b *testing.B) {
	for x := 0; x < len(_keys); x++ {
		tree.Delete(_keys[x])
	}
	b.Logf("len(keys):%d tree.n:%d\n", len(_keys), tree.n)
}

func BenchmarkMapDelete(b *testing.B) {
	for x := 0; x < len(_keys); x++ {
		key := string(_keys[x])
		delete(m, key)
	}
	b.Logf("len(keys):%d len(map):%d\n", len(_keys), len(m))
}

//
// Benchmark Result
//
// go version devel +96ffbecd95 Sat Feb 16 05:05:59 2019 +0700 linux/amd64
//
// goos: linux
// goarch: amd64
// pkg: github.com/shuLhan/share/lib/tree/rbt
// BenchmarkInsert-4              1        17850024798 ns/op       880001816 B/op  20000021 allocs/op
// --- BENCH: BenchmarkInsert-4
//     benchmark_test.go:41: x:10000000 tree.n:9941235
// BenchmarkMapInsert-4           1        7220891706 ns/op        1870890632 B/op 20307398 allocs/op
// --- BENCH: BenchmarkMapInsert-4
//     benchmark_test.go:45: keys: 10000000
//     benchmark_test.go:53: x:10000000 len(m):9941235
// BenchmarkFind-4                1        16726249706 ns/op              0 B/op          0 allocs/op
// BenchmarkMapFind-4             1        1036405704 ns/op               0 B/op          0 allocs/op
// BenchmarkDelete-4              1        18327255676 ns/op           1736 B/op         20 allocs/op
// --- BENCH: BenchmarkDelete-4
//     benchmark_test.go:73: x:10000000 tree.n:31191
// BenchmarkMapDelete-4           1        1476922644 ns/op            1640 B/op         18 allocs/op
// --- BENCH: BenchmarkMapDelete-4
//     benchmark_test.go:81: x:10000000 len(map):0
// PASS
// ok      github.com/shuLhan/share/lib/tree/rbt   76.532s
//
