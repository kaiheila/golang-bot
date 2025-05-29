package compress

import (
	"bytes"
	"github.com/valyala/gozstd"
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

// StreamReader 实现ZSTD流式解压
type StreamReader struct {
	reader    *gozstd.Reader
	src       io.Reader
	pool      *sync.Pool
	bytesRead int64
}

// NewStreamReader 创建新的流式解压器
func NewStreamReader(src io.Reader) *StreamReader {
	return &StreamReader{
		reader: gozstd.NewReader(src),
		src:    src,
	}
}

// // Read 实现io.Reader接口
func (sr *StreamReader) Read(p []byte) (int, error) {
	// 首次使用或需要重置时初始化
	if sr.reader == nil {
		sr.reader = sr.getReader()
		sr.reader.Reset(sr.src, nil)
	}

	n, err := sr.reader.Read(p)
	sr.bytesRead += int64(n)
	return n, err
}

// // Close 释放资源
func (sr *StreamReader) Close() error {
	if sr.reader != nil {
		sr.reader.Release()
		if sr.pool != nil {
			sr.pool.Put(sr.reader)
		}
		sr.reader = nil
	}
	return nil
}

//
//// Reset 重置解压器以处理新流
//func (sr *StreamReader) Reset(r io.Reader) error {
//	if sr.reader == nil {
//		sr.reader = sr.getReader()
//	}
//	sr.src = r
//	sr.bytesRead = 0
//	sr.reader.Reset(r)
//	return nil
//}

// BytesRead 返回已解压的字节数
func (sr *StreamReader) BytesRead() int64 {
	return sr.bytesRead
}

func (sr *StreamReader) getReader() *gozstd.Reader {
	if sr.pool != nil {
		return sr.pool.Get().(*gozstd.Reader)
	}
	return gozstd.NewReader(nil)
}

// SetReaderPool 设置Reader池以重用解压器
func (sr *StreamReader) SetReaderPool(pool *sync.Pool) {
	sr.pool = pool
}

func TestCompress(t *testing.T) {
	buf0 := bytes.Buffer{}
	w := gozstd.NewWriter(&buf0)
	w.Write([]byte("hello"))
	w.Flush()
	buf := bytes.Buffer{}
	reader := gozstd.NewReader(&buf)
	reader.Reset(nil, nil)
	reader.Reset(&buf, nil)
	buf.Write(buf0.Bytes())
	out := make([]byte, 32*1024)
	n, err := reader.Read(out)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", string(out[:n]))
}
func TestStream(t *testing.T) {
	// 1. 准备压缩数据
	compressedData := compressTestData(t)
	buf := new(bytes.Buffer)
	// 2. 创建流式解压器
	r := NewStreamReader(buf)
	defer r.Close()

	for i := 0; i < len(compressedData); i++ {
		t.Logf("compressed size:%d\n", len(compressedData[i]))
		t.Logf("%x\n", compressedData[i])
		buf.Write(compressedData[i])
		// 3. 流式读取解压数据
		buf := make([]byte, 32*1024) // 32KB缓冲区
		for {
			n, err := r.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal("Decompress error:", err)
			}

			// 处理解压后的数据
			processChunk(buf[:n])
		}

		t.Logf("Decompressed %d bytes\n", r.BytesRead())
	}
}

func compressTestData(t *testing.T) [][]byte {
	buf := bytes.Buffer{}
	var res [][]byte = make([][]byte, 0)
	w := gozstd.NewWriter(&buf)
	for i := 0; i < 2; i++ {
		buf.Reset()
		_, err := w.Write([]byte("test data to compress"))
		if err != nil {
			panic(err)
		}
		w.Flush()
		out := make([]byte, buf.Len())
		copy(out, buf.Bytes())
		res = append(res, out)
		t.Logf("0: %x", out)
	}
	//for i := 0; i < len(res); i++ {
	//	t.Logf("%d: %x", i, res[i])
	//}
	w.Close()
	return res

}

func TestCompressTestData(t *testing.T) {
	buf := bytes.Buffer{}
	var res [][]byte = make([][]byte, 0)
	w := gozstd.NewWriter(&buf)
	for i := 0; i < 2; i++ {
		buf.Reset()
		_, err := w.Write([]byte("test data to compress"))
		if err != nil {
			panic(err)
		}
		w.Flush()
		out := make([]byte, buf.Len())
		copy(out, buf.Bytes())
		res = append(res, out)
		t.Logf("0: %x", out)
	}
	//for i := 0; i < len(res); i++ {
	//	t.Logf("%d: %x", i, res[i])
	//}
	w.Close()

}

func processChunk(data []byte) {
	os.Stdout.Write(data)
}
