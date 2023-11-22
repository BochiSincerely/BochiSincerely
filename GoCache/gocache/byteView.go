package geecache

type ByteView struct {
	b []byte //b存储真实的缓存值
	//byte可以支持任意数据类型的存储，如字符串、图片等
}

func (v ByteView) Len() int {
	return len(v.b)
} //返回其所占的内存大小

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
} //因为b是只读，索引我们需要进行一个拷贝，防止缓存值被外部程序修改

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
