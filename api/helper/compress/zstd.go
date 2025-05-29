package compress

import (
	"bytes"
	"github.com/valyala/gozstd"
	"io"
	"sync"
)

type ZstdStreamCompressor struct {
	encoder *gozstd.Writer // 直接作为Writer使用
	buffer  bytes.Buffer   // 输出缓冲区
	mu      sync.Mutex     // 并发控制
}

func newZstdStreamCompressor() *ZstdStreamCompressor {
	var buf bytes.Buffer
	encoder := gozstd.NewWriter(&buf)
	//encoder, err := zstd.NewWriter(&buf, zstd.WithEncoderLevel(zstd.SpeedDefault))

	return &ZstdStreamCompressor{
		encoder: encoder, // encoder本身即是Writer
		buffer:  buf,
	}
}

func (z *ZstdStreamCompressor) Reset() error {
	z.buffer.Reset()                   // 重置缓冲区
	z.encoder.Reset(&z.buffer, nil, 6) // 重置编码器状态
	return nil
}

func (z *ZstdStreamCompressor) Recycle() error {
	z.buffer.Reset()
	z.encoder.Reset(nil, nil, 6) // 释放资源
	return nil
}
func (z *ZstdStreamCompressor) Compress(data []byte) ([]byte, error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.buffer.Reset()
	if _, err := z.encoder.Write(data); err != nil {
		return nil, err
	}
	if err := z.encoder.Flush(); err != nil {
		return nil, err
	}
	res := z.buffer.Bytes()
	return res, nil
}

func (z *ZstdStreamCompressor) Close() error {
	return z.encoder.Close()
}

type ZstdStreamDecompressor struct {
	decoder *gozstd.Reader
	src     *bytes.Buffer
	mu      sync.Mutex
}

func NewZstdStreamDecompressor() DecompressorInterface {
	buf := new(bytes.Buffer)
	decoder := gozstd.NewReader(buf)
	return &ZstdStreamDecompressor{
		decoder: decoder,
		src:     buf,
	}
}

func (z *ZstdStreamDecompressor) Decompress(data []byte) ([]byte, error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	// 使用缓冲区复制数据
	//// 创建带缓冲的写入器
	z.src.Write(data)
	out := bytes.Buffer{}
	for {
		buf := make([]byte, 1024)
		n, err := z.decoder.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n > 0 {
			out.Write(buf[:n])
		} else if n == 0 || err == io.EOF {
			break
		}

	}
	return out.Bytes(), nil

}

func (z *ZstdStreamDecompressor) Reset() error {
	z.src.Reset()
	z.decoder.Reset(z.src, nil)
	return nil
}

func (z *ZstdStreamDecompressor) Recycle() error {
	z.src.Reset()
	z.decoder.Reset(nil, nil)
	return nil
}

type ZstdPerMessageCompressor struct {
}

func NewZstdPerMessageCompressor() *ZstdPerMessageCompressor {
	return &ZstdPerMessageCompressor{}
}
func (z *ZstdPerMessageCompressor) Compress(data []byte) ([]byte, error) {
	out := gozstd.Compress(nil, data)
	return out, nil
}

func (z *ZstdPerMessageCompressor) Reset() error {
	return nil
}

func (z *ZstdPerMessageCompressor) Recycle() error {
	return nil
}

type ZstdPerMessageDecompressor struct {
}

func NewZstdPerMessageDecompressor() DecompressorInterface {
	z := &ZstdPerMessageDecompressor{}
	return z
}

func (z *ZstdPerMessageDecompressor) Decompress(data []byte) ([]byte, error) {
	return gozstd.Decompress(nil, data)
}

func (z *ZstdPerMessageDecompressor) Reset() error {
	return nil
}

func (z *ZstdPerMessageDecompressor) Recycle() error {
	return nil
}
