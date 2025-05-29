package compress

import (
	"sync"
)

type CompressType int

const (
	CompressTypeNone           CompressType = 0
	CompressTypeZlibPerMessage CompressType = 1
	CompressTypeZstdPerMessage CompressType = 2
	CompressTypeZlibStream     CompressType = 3 // TODO: 暂未支持
	CompressTypeZstdStream     CompressType = 4
)

func ParseCompressType(compress bool, compressTypeStr string) CompressType {
	if compress == false {
		return CompressTypeNone
	}
	switch compressTypeStr {
	case "zlib":
		return CompressTypeZlibPerMessage
	case "zstd":
		return CompressTypeZstdPerMessage
	case "zlib_stream":
		return CompressTypeZlibStream
	case "zstd_stream":
		return CompressTypeZstdStream
	default:
		return CompressTypeZlibPerMessage
	}
}

func GetCompressTypeName(compressType CompressType) string {
	switch compressType {
	case CompressTypeZlibPerMessage:
		return "zlib"
	case CompressTypeZstdPerMessage:
		return "zstd"
	case CompressTypeZlibStream:
		return "zlib_stream"
	case CompressTypeZstdStream:
		return "zstd_stream"
	default:
		return "none"
	}
}

type CompressorInterface interface {
	Compress(data []byte) ([]byte, error)
	//Flush() ([]byte, error)
	Reset() error
	Recycle() error
	//Write(data []byte) (int, error)
}

type DecompressorInterface interface {
	Decompress(data []byte) ([]byte, error)
	Reset() error
	Recycle() error
}

var compressorZlibStreamPool = sync.Pool{
	New: func() interface{} {
		return newZlibStreamCompressor()
	},
}
var compressorZstdStreamPool = sync.Pool{
	New: func() interface{} {
		return newZstdStreamCompressor()
	},
}

func GetCompressor(compressType CompressType) CompressorInterface {
	if compressType == CompressTypeZstdStream {
		c := compressorZstdStreamPool.Get().(CompressorInterface)
		c.Reset()
		return c
	} else if compressType == CompressTypeZlibStream {
		c := compressorZlibStreamPool.Get().(CompressorInterface)
		c.Reset()
		return c
	}
	return nil
}
func RecycleCompressor(compressType CompressType, c CompressorInterface) {
	if c == nil || compressType == CompressTypeNone {
		return
	}
	// 确保所有数据已刷新
	// 放回池中
	if compressType == CompressTypeZstdStream {
		c.Recycle()
		compressorZstdStreamPool.Put(c)
		return
	} else if compressType == CompressTypeZlibStream {
		c.Recycle()
		compressorZlibStreamPool.Put(c)
		return

	}
}

var decompressorZstdStreamPool = sync.Pool{
	New: func() interface{} {
		return NewZstdStreamDecompressor()
	},
}
var decompressorZlibStreamPool = sync.Pool{
	New: func() interface{} {
		return NewZlibStreamDecompressor()
	},
}

func GetDecompressor(compressType CompressType) DecompressorInterface {
	if compressType == CompressTypeZstdStream {
		d := decompressorZstdStreamPool.Get().(DecompressorInterface)
		d.Reset()
		return d
	} else if compressType == CompressTypeZlibStream {
		d := decompressorZlibStreamPool.Get().(DecompressorInterface)
		d.Reset()
		return d
	} else if compressType == CompressTypeZlibPerMessage {
		return NewZlibPerMessageDecompressor()
	} else if compressType == CompressTypeZstdPerMessage {
		return NewZstdPerMessageDecompressor()
	}
	return nil
}

func RecycleDecompressor(compressType CompressType, decompressor DecompressorInterface) {
	if compressType == CompressTypeZstdStream {
		decompressor.Recycle()
		decompressorZstdStreamPool.Put(decompressor)
	} else if compressType == CompressTypeZlibStream {
		decompressor.Recycle()
		decompressorZlibStreamPool.Put(decompressor)
	}
}

var compressorZstdPerMessagePool = sync.Pool{
	New: func() interface{} {
		c := NewZstdPerMessageCompressor()
		c.Reset()
		return c
	},
}

var compressorZlibPerMessagePool = sync.Pool{
	New: func() interface{} {
		c := NewZlibPerMessageCompressor()
		c.Reset()
		return c
	},
}
