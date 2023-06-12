package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int // 12345, 23456,  89654
	nodeHashMap map[int]string
}

func NewNodeMap(fn HashFunc) *NodeMap {
	m := &NodeMap{
		hashFunc:    fn,
		nodeHashs:   nil,
		nodeHashMap: make(map[int]string),
	}

	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}

	return m
}

func (n *NodeMap) IsEmpty() bool {
	return len(n.nodeHashs) == 0
}

func (n *NodeMap) AddNode(keys ...string) {
	for _, k := range keys {
		if k == "" {
			continue
		}

		hash := int(n.hashFunc([]byte(k)))
		n.nodeHashs = append(n.nodeHashs, hash)
		n.nodeHashMap[hash] = k
	}
	sort.Ints(n.nodeHashs)
}

func (n *NodeMap) PickNode(key string) string {
	if n.IsEmpty() {
		return ""
	}

	hashKey := int(n.hashFunc([]byte(key)))
	idx := sort.Search(len(n.nodeHashs), func(i int) bool {
		return n.nodeHashs[i] >= hashKey
	})

	// 是否落在最后
	if idx == len(n.nodeHashs) {
		idx = 0
	}

	return n.nodeHashMap[n.nodeHashs[idx]]
}
