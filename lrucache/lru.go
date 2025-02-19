//go:build !solution

package lrucache

import "container/list"

type LruCache struct {
	data     map[int]*list.Element
	list     *list.List
	capacity int
}

type Node struct {
	key   int
	value int
}

func New(cap int) Cache {
	return &LruCache{
		data:     make(map[int]*list.Element, cap),
		list:     list.New(),
		capacity: cap,
	}
}

func (lru *LruCache) Get(key int) (int, bool) {
	if nodeToGet, ok := lru.data[key]; ok {
		lru.list.MoveToFront(nodeToGet)
		return nodeToGet.Value.(Node).value, true
	}

	return 0, false
}

func (lru *LruCache) Set(key, value int) {
	if nodeToUpd, ok := lru.data[key]; ok {
		lru.list.MoveToFront(nodeToUpd)
		nodeToUpd.Value = Node{key: key, value: value}
		return
	}

	if lru.list.Len() == lru.capacity && lru.capacity != 0 {
		nodeToDelete := lru.list.Back()
		lru.list.Remove(nodeToDelete)
		delete(lru.data, nodeToDelete.Value.(Node).key)
	}

	if lru.list.Len() < lru.capacity {
		nodeToAdd := lru.list.PushFront(Node{key: key, value: value})
		lru.data[key] = nodeToAdd
	}
}

func (lru *LruCache) Range(f func(key, value int) bool) {
	node := lru.list.Back()
	for node != nil {
		if !f(node.Value.(Node).key, node.Value.(Node).value) {
			return
		}
		node = node.Prev()
	}
}

func (lru *LruCache) Clear() {
	lru.data = make(map[int]*list.Element, lru.capacity)
	lru.list = lru.list.Init()
}
