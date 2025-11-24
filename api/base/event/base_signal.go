package event

import (
	"bytes"
	"encoding/binary"

	"github.com/bytedance/sonic/ast"
)

// SignalInterface 定义signal的接口
type SignalInterface interface {
	WithVersion(version int)
	WithSN(sn int64)
	WithHasSN(hasSN bool)
	WithPayload(payload []byte)
	WithJsonPayload(payload *ast.Node)
	WithJsonBytePayload(payload []byte)
	//ithIncludeLength(include bool)
	Encode() ([]byte, error)
	Decode(data []byte) error
	IsVersion0() bool
	IsVersion1() bool
	GetSN() int64
	HaveSN() bool
	IsMsgSignal() bool
}

// BaseSignal 是所有signal的基类
type BaseSignal struct {
	SignalType      int
	Version         int
	SN              int64
	HasSN           bool
	Payload         []byte
	JsonPayload     *ast.Node
	JsonBytePayload []byte
}

// NewBaseSignal 创建一个新的基础signal
func NewBaseSignal(signalType int, version int, hasSN bool) *BaseSignal {
	return &BaseSignal{
		SignalType: signalType,
		Version:    version,
		HasSN:      hasSN,
		Payload:    nil,
	}
}

// IsVersion0 检查是否是版本0
func (s *BaseSignal) IsVersion0() bool {
	return s.Version == 0
}

// IsVersion1 检查是否是版本1
func (s *BaseSignal) IsVersion1() bool {
	return s.Version == 1
}

// WithSN 设置序列号
func (s *BaseSignal) WithSN(sn int64) {
	s.SN = sn
	s.HasSN = true
}

func (s *BaseSignal) HaveSN() bool {
	return s.HasSN
}

// WithVersion 设置版本号
func (s *BaseSignal) WithVersion(version int) {
	s.Version = version
}

// WithHasSN 设置是否有序列号
func (s *BaseSignal) WithHasSN(hasSN bool) {
	s.HasSN = hasSN
}

// WithPayload 设置负载数据
func (s *BaseSignal) WithPayload(payload []byte) {
	s.Payload = payload
}

func (s *BaseSignal) WithJsonPayload(payload *ast.Node) {
	s.JsonPayload = payload
}

func (s *BaseSignal) WithJsonBytePayload(payload []byte) {
	s.JsonBytePayload = payload
}

//// WithIncludeLength 设置是否包含载荷长度
//func (s *BaseSignal) WithIncludeLength(include bool) {
//	// 暂不实现，后续可根据需要添加
//}

// GetSN 获取序列号
func (s *BaseSignal) GetSN() int64 {
	return s.SN
}

// IsMsgSignal 判断是否为消息类型signal
func (s *BaseSignal) IsMsgSignal() bool {
	return s.SignalType == 0 // SignalMsg
}

// Decode 实现SignalInterface接口的Decode方法
func (s *BaseSignal) Decode(data []byte) error {
	decoded, err := Decode(data, s.Version)
	if err != nil {
		return err
	}
	//s.SignalType = decoded.SignalType
	s.Version = decoded.Version
	s.SN = decoded.SN
	s.HasSN = decoded.HasSN
	s.Payload = decoded.Payload
	return nil
}

// Encode 编码signal为字节数组
func (s *BaseSignal) Encode() ([]byte, error) {
	if s.IsVersion0() {
		// 版本0没有header，直接返回payload
		if s.Payload != nil {
			return s.Payload, nil
		}

	}

	// 版本1需要封装header
	return s.encodeWithHeader()
}

// encodeWithHeader 编码带header的signal
func (s *BaseSignal) encodeWithHeader() ([]byte, error) {
	buf := &bytes.Buffer{}

	// 第一个字节：版本(8bit)
	// 版本占8位
	versionByte := uint8(s.Version)

	// 第二个字节：标志位
	flagByte := uint8(0)
	if s.HasSN {
		flagByte |= 0x80 // 最高位表示是否有序列号
	}
	// 保留的6位填0
	// 最低位表示下一个byte是否也是标志位，这里设置为0

	buf.WriteByte(versionByte)
	buf.WriteByte(flagByte)

	// 序列号长度标识：使用4位
	snLen := 2 // 默认使用4字节
	if s.SN <= 0xFF {
		snLen = 0 // 1字节
	} else if s.SN <= 0xFFFF {
		snLen = 1 // 2字节
	} else if s.SN <= 0xFFFFFFFF {
		snLen = 2 // 4字节
	} else {
		snLen = 3 // 8字节
	}

	snLenByte := uint8(snLen) << 4 // 高4位表示序列号长度
	// 保留的4位填0
	buf.WriteByte(snLenByte)

	// 根据长度写入序列号
	switch snLen {
	case 0:
		// 1字节
		buf.WriteByte(uint8(s.SN))
	case 1:
		// 2字节
		snBuf := make([]byte, 2)
		binary.BigEndian.PutUint16(snBuf, uint16(s.SN))
		buf.Write(snBuf)
	case 2:
		// 4字节
		snBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(snBuf, uint32(s.SN))
		buf.Write(snBuf)
	case 3:
		// 8字节
		snBuf := make([]byte, 8)
		binary.BigEndian.PutUint64(snBuf, uint64(s.SN))
		buf.Write(snBuf)
	}

	// 写入负载
	buf.Write(s.Payload)

	return buf.Bytes(), nil
}

// Decode 从字节数组解码signal
func Decode(data []byte, version int) (*BaseSignal, error) {
	if version == 0 {
		// 版本0直接返回，signal type需要根据payload解析
		return &BaseSignal{
			Version: version,
			Payload: data,
		}, nil
	}

	// 版本1需要解析header
	return decodeWithHeader(data)
}

// decodeWithHeader 解析带header的signal
func decodeWithHeader(data []byte) (*BaseSignal, error) {
	buf := bytes.NewReader(data)

	// 读取版本
	versionByte, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	version := int(versionByte)

	// 读取标志位
	flagByte, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	hasSN := (flagByte & 0x80) != 0

	signal := &BaseSignal{
		Version: version,
		HasSN:   hasSN,
	}

	if hasSN {
		// 读取序列号长度标识
		snLenByte, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}
		snLen := int((snLenByte >> 4) & 0x0F)

		// 根据长度读取序列号
		switch snLen {
		case 0:
			// 1字节
			snByte, err := buf.ReadByte()
			if err != nil {
				return nil, err
			}
			signal.SN = int64(snByte)
		case 1:
			// 2字节
			snBuf := make([]byte, 2)
			_, err := buf.Read(snBuf)
			if err != nil {
				return nil, err
			}
			signal.SN = int64(binary.BigEndian.Uint16(snBuf))
		case 2:
			// 4字节
			snBuf := make([]byte, 4)
			_, err := buf.Read(snBuf)
			if err != nil {
				return nil, err
			}
			signal.SN = int64(binary.BigEndian.Uint32(snBuf))
		case 3:
			// 8字节
			snBuf := make([]byte, 8)
			_, err := buf.Read(snBuf)
			if err != nil {
				return nil, err
			}
			signal.SN = int64(binary.BigEndian.Uint64(snBuf))
		}
	}

	// 读取payload
	remaining := buf.Len()
	actualLength := remaining

	payload := make([]byte, actualLength)
	_, err = buf.Read(payload)
	if err != nil {
		return nil, err
	}
	signal.Payload = payload

	return signal, nil
}
