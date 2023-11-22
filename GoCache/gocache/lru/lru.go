package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes  int64      //允许使用的最大内存
	nbytes    int64      //当前已使用的内存
	ll        *list.List //双链表
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数，可以为nil
}

type entry struct { //双链表节点
	key   string
	value Value
}

type Value interface {
	Len() int
}

//返回值所占用的内存大小

// 为实例化Cache函数，
func New(maxBytes int64, onEvicted func(string2 string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找功能
/*
step1: 从字典中找到对应的双链表节点
step2：将该节点移动到队尾
*/
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) //将ele节点移动到队尾
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return

}

// 删除 缓存淘汰  移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() //取队首节点，从链表中删除
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                //从字典中删除该节点的一个映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //更新当前所用的内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) //若不为nil则调用回调函数
		}
	}

}

// 新增或修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
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

func (c *Cache) Len() int {
	return c.ll.Len()
}
