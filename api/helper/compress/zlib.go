package compress

import (
	"bytes"
	"github.com/klauspost/compress/zlib"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

type ZlibStreamCompressor struct {
	writer *zlib.Writer
	buffer bytes.Buffer
	mu     sync.Mutex
}

func newZlibStreamCompressor() *ZlibStreamCompressor {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	return &ZlibStreamCompressor{
		writer: writer,
		buffer: buf,
	}
}

func (z *ZlibStreamCompressor) Compress(data []byte) ([]byte, error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.buffer.Reset()
	if _, err := z.writer.Write(data); err != nil {
		return nil, err
	}
	if err := z.writer.Flush(); err != nil {
		return nil, err
	}
	return z.buffer.Bytes(), nil
}

func (z *ZlibStreamCompressor) Reset() error {
	z.buffer.Reset()
	z.writer.Reset(&z.buffer)
	return nil
}

func (z *ZlibStreamCompressor) Recycle() error {
	z.buffer.Reset()
	return z.writer.Close()
}

type ZlibStreamDecompressor struct {
	decoder io.ReadCloser
	buf     bytes.Buffer
	mu      sync.Mutex
}

func NewZlibStreamDecompressor() DecompressorInterface {
	buf := bytes.Buffer{}
	decoder, err := zlib.NewReader(&buf)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	z := &ZlibStreamDecompressor{decoder: decoder, buf: buf}
	return z
}

func (z *ZlibStreamDecompressor) Decompress(data []byte) ([]byte, error) {
	z.mu.Lock()
	defer z.mu.Unlock()
	// 使用缓冲区复制数据
	//// 创建带缓冲的写入器
	z.buf.Write(data)
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

func (z *ZlibStreamDecompressor) Reset() error {
	z.buf.Reset()
	return nil
}

func (z *ZlibStreamDecompressor) Recycle() error {
	z.buf.Reset()
	z.decoder = nil
	return nil
}

type ZlibPerMessageCompressor struct {
	writer *zlib.Writer
	buffer bytes.Buffer
}

func NewZlibPerMessageCompressor() *ZlibPerMessageCompressor {
	buffer := bytes.Buffer{}
	writer := zlib.NewWriter(&buffer)
	return &ZlibPerMessageCompressor{writer: writer, buffer: buffer}
}

func (z *ZlibPerMessageCompressor) Compress(data []byte) ([]byte, error) {
	_, err := z.writer.Write(data)
	if err != nil {
		return nil, err
	}
	z.writer.Flush()
	return z.buffer.Bytes(), nil
}

func (z *ZlibPerMessageCompressor) Reset() error {
	z.buffer.Reset()
	z.writer.Reset(&z.buffer)
	return nil
}

func (z *ZlibPerMessageCompressor) Recycle() error {
	z.Reset()
	return nil
}

type ZlibPerMessageDecompressor struct {
}

func NewZlibPerMessageDecompressor() *ZlibPerMessageDecompressor {
	z := &ZlibPerMessageDecompressor{}
	z.Reset()
	return z
}

func (z *ZlibPerMessageDecompressor) Decompress(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	decoder, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()
	res, err := io.ReadAll(decoder)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (z *ZlibPerMessageDecompressor) Reset() error {
	return nil
}

func (z *ZlibPerMessageDecompressor) Recycle() error {
	return nil
}
