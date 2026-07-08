package main

import "fmt"

type Node struct {
	prev  *Node
	next  *Node
	key   int
	value int
}
type LRUCache struct {
	mp       map[int]*Node
	capacity int
	head     *Node
	tail     *Node
}

func Constructor(capacity int) LRUCache {
	lruchhe := LRUCache{
		mp:       make(map[int]*Node),
		capacity: capacity,
		head:     &Node{},
		tail:     &Node{},
	}
	lruchhe.head.prev = nil
	lruchhe.head.next = lruchhe.tail
	lruchhe.tail.next = nil
	lruchhe.tail.prev = lruchhe.head
	return lruchhe
}
func (l *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}
func (l *LRUCache) insert(node *Node) {
	node.next = l.head.next
	node.prev = l.head
	l.head.next.prev = node
	l.head.next = node

}

func (l *LRUCache) Get(key int) int {
	val, ok := l.mp[key]
	if !ok {
		return -1
	}
	l.remove(val)
	l.insert(val)
	return val.value
}
func (l *LRUCache) Put(key int, value int) {
	val, ok := l.mp[key]
	if !ok {
		node := Node{
			prev:  nil,
			next:  nil,
			key:   key,
			value: value,
		}
		l.insert(&node)
		l.mp[key] = &node

		if len(l.mp) > l.capacity {
			lru := l.tail.prev
			l.remove(lru)
			delete(l.mp, lru.key)
		}
		return
	}
	l.remove(val)
	val.value = value
	l.insert(val)

	return
}
func main() {
	LRUCache := Constructor(3)
	LRUCache.Put(1, 10)
	LRUCache.Put(2, 20)
	LRUCache.Put(3, 30)
	LRUCache.Put(4, 50)
	LRUCache.Put(1, 90)
	LRUCache.Put(7, 80)
	LRUCache.Put(8, 80)
	curr := LRUCache.head.next
	for curr != LRUCache.tail {
		fmt.Printf("%d:%d\n", curr.key, curr.value)
		curr = curr.next
	}

}
