package lru

import "container/list"

// Cache is a LRU cache (non-thread-safe)
type Cache struct {

	// maxBytes: max memory
	maxBytes int64

	// nbytes: used memory
	nbytes int64

	// ll: a builtin Doubly Linked List
	ll *list.List

	// cache: key -> list
	cache map[string]*list.Element

	// OnEvicted: a callback function, which is optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// entry is type of list value
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	// If key in hashmap, take out the element and move it to head (next) of the list.
	// It means that value has been used recently.
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		if kv, ok := ele.Value.(*entry); ok {
			return kv.value, ok
		}
	}
	return
}

// Add adds a value to the cache
func (c *Cache) Add(key string, value Value) {

	// the key exists in cache:
	if ele, ok := c.cache[key]; ok {
		// take it out and move it to head of list
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		// recalculate memory usage
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())

		// update value
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {

	// the last element in list, take out.
	ele := c.ll.Back()
	if ele == nil {
		return
	}
	// remove from list
	c.ll.Remove(ele)

	// remove from hashmap
	kv := ele.Value.(*entry)
	delete(c.cache, kv.key)

	// recalculate memory usage
	c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

	// callback
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Len returns the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
