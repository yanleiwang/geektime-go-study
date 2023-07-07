package compress

type Compressor interface {
	// Compress 返回压缩后的数据
	Compress(b []byte) ([]byte, error)
	// DeCompress 返回解压缩后的数据
	DeCompress(b []byte) ([]byte, error)
	// Code 压缩算法的 唯一标识符
	Code() byte
}
