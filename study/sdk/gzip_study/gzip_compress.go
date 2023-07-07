// Package gzip_study  gzip 压缩/解压缩写法
package gzip_study

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Compressor struct {
}

func (c *Compressor) Compress(b []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := gzip.NewWriter(buf)

	_, err := writer.Write(b)
	if err != nil {
		return nil, err
	}
	// 只能用Close, 不能用flush? 用flush 还是会得到unexpected EOF
	// 解释: https://stackoverflow.com/questions/60923654/closing-gzip-writer-in-defer-causes-data-loss
	// 因为flush 只会刷新缓冲数据，但是不写gzip footer 页脚? 可能是gzip的结束符? 而close会
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func (c *Compressor) DeCompress(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	defer r.Close()
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)

}
