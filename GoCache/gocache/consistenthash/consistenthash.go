package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 一致性哈希算法
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE

	}
	return m
}

func (m *Map) Add(keys ...string) {

	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)

}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	}) //顺时针找到匹配第一个节点对应的虚拟节点index，然后从m.keys中获取对应的一个hash值
	return m.hashMap[m.keys[index%len(m.keys)]]
	//获取一个真实的节点

}
