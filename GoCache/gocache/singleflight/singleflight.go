package singleflight

import (
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
} //表示正在进行中或者已经结束的请求。通过waitgroup锁来避免重入

type Group struct {
	mu sync.Mutex
	m  map[string]*call
} //管理不同key的请求

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         //有请求进行，等待
		return c.val, c.err //返回结果
	}
	c := new(call)
	c.wg.Add(1)  // 请求前加锁
	g.m[key] = c //添加到g.m中，表示有对应的请求
	g.mu.Unlock()

	c.val, c.err = fn() //调用fn函数
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key) //更新g.m
	g.mu.Unlock()

	return c.val, c.err
}

/*
用Do方法，接受2个参数，一个是key，另外一个是fn函数
作用：针对相同的key，无论Do被调用多少次，fn函数只会被调用一次，等待fn调用结束了
返回返回值或者错误
g.m是为了保护成员变量，不被并发读写而加上锁，
*/
