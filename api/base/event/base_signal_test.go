package event

import (
	"testing"
)

func TestDecodeWithHeader(t *testing.T) {
	// 创建一个带header的signal数据用于测试
	// 版本1，带SN，SN为1234567890（8字节）
	payload := []byte("test payload")

	// 手动构建带header的数据
	// 版本1 (1字节)
	// 标志位：带SN (1字节)
	// SN长度：3 (8字节) (1字节)
	// SN：1234567890 (8字节，大端序)
	// payload：test payload
	data := []byte{
		1,                      // 版本1
		0x80,                   // 带SN标志
		0x30,                   // SN长度3（8字节）
		0x00, 0x00, 0x00, 0x00, // SN的前4字节
		0x49, 0x96, 0x02, 0xd2, // SN的后4字节（1234567890）
	}
	data = append(data, payload...)

	// 测试decodeWithHeader函数
	signal, err := decodeWithHeader(data)
	if err != nil {
		t.Fatalf("decodeWithHeader failed: %v", err)
	}

	// 验证结果
	if signal.Version != 1 {
		t.Errorf("expected version 1, got %d", signal.Version)
	}

	if !signal.HasSN {
		t.Error("expected HasSN to be true")
	}

	if signal.SN != 1234567890 {
		t.Errorf("expected SN 1234567890, got %d", signal.SN)
	}

	if string(signal.Payload) != string(payload) {
		t.Errorf("expected payload %q, got %q", payload, signal.Payload)
	}
}

func TestDecodeWithHeaderNoSN(t *testing.T) {
	// 创建一个不带SN的signal数据用于测试
	payload := []byte("test payload without SN")

	// 手动构建带header的数据
	// 版本1 (1字节)
	// 标志位：不带SN (1字节)
	// payload：test payload without SN
	data := []byte{
		1,    // 版本1
		0x00, // 不带SN标志
	}
	data = append(data, payload...)

	// 测试decodeWithHeader函数
	signal, err := decodeWithHeader(data)
	if err != nil {
		t.Fatalf("decodeWithHeader failed: %v", err)
	}

	// 验证结果
	if signal.Version != 1 {
		t.Errorf("expected version 1, got %d", signal.Version)
	}

	if signal.HasSN {
		t.Error("expected HasSN to be false")
	}

	if signal.SN != 0 {
		t.Errorf("expected SN 0, got %d", signal.SN)
	}

	if string(signal.Payload) != string(payload) {
		t.Errorf("expected payload %q, got %q", payload, signal.Payload)
	}
}
