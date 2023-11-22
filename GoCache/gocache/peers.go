package geecache

import pb "GoCache/gocache/gocachepb"

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
	//用于根据传入的key选择对应的节点
}
