package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	value interface{}
	key   Key
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	cItem := &cacheItem{value: value, key: key}

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = cItem

		return true
	}

	if c.queue.Len() > c.capacity {
		itemBack := c.queue.Back()

		if itemBack != nil {
			c.queue.Remove(itemBack)

			delete(c.items, itemBack.Value.(*cacheItem).key)
		}
	}

	c.items[key] = c.queue.PushFront(cItem)

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.PushFront(item.Value)

		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
