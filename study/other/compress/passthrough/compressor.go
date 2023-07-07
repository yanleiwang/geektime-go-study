// Package passthrough 透传 不压缩
package passthrough

type Compressor struct {
}

func (c *Compressor) Compress(b []byte) ([]byte, error) {
	return b, nil
}

func (c *Compressor) DeCompress(b []byte) ([]byte, error) {
	return b, nil
}

func (c *Compressor) Code() byte {
	return 0
}
