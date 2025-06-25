package helper

import (
	"encoding/base64"
	"github.com/bytedance/sonic"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestDecryptData(t *testing.T) {
	key := "sfjlfvx98sdlsdfsdfsdfD)9DJ&sdf"

	// 初始化向量（必须是 16 字节）
	iv := "9051dbfc59d6af1d"
	data := "cnUxOXlsakwvWWNzTG02bEZsOXJ3V2hQUURaNEhxKzV6aC9BaktiMHF4cCt3Mm9LUmpGcElCdkZmRURmZzh1NDZGNXpuN0VsME1DZG1QRktVV2FNc2Jacko4QkN0dTZ0MThvbDVUWk54RDA9"
	//data := "XY4RWRCSXRxOURLVGRNSXNURFZtZkFYQ3pMUjNCcHkvSEkxT1pyOTV2Nmd5VnVoVE1NV0R4eGFiemxHZkdqRTlPanVmVkdDTzZQNmd5MEZnQ1pNMHNkVVhkQm5hMGU1V2w1dmw2QXhoVjlGU3BvUDdjeHNaQitjdVdHRUVzMVk1bjRraWtWa1VTYmVYb0pCNHpTdVRBUWdsYXN0WFlaZ1FPcXVCdjF2NTFPaVdjQTVDZ2tjRnFMNGJMQk9MU3RJUlJObjU2elZXcE9rNmFscTN5a1FXcysrMFMyWWlIdG5KN1B4ZkFzb0k0VERMWHBTckFteHFBSVZ3U0hzdVZHS3RybTVTRjZnWmxqWWVUekpnU2xRa0JlVThDVTJNaEVJeHc0enJySnE2MlV4Vm5BM3pFRFBOWk15MmNlaXpvaWdCZHFpUFVKeE16QWdBSkcyS0dhaW9LOTVJOTI4WStyVzJVYzY2aHRmY1V3ZXE3NmtQZkxLd2d6M1FIRjAzQnc2aDllU3NOZGYwZFlRbHhacDJtSnR5M2hLbTFnT3FmWnlPc0NTeVRiVjF2cmY0R3pnZ3hpWlNyZC9rUzdkaytMRStXaXhqUTVpWFFzZ3p6a0tKTm84SXdjbUVBTHhFc01jQ045NWcvdEtqNkltUElDbjFJTW1tU0hBd3dnTGNrdEdPaVdXVUFmT1NuQUVjUmROWFRlaFRxS1dsdTJ1UDVSMVhaTDBZQmUwSFNZUkhoZUxibC9NM040N3dlOEx5MXZvNFJja1g4Ni9peXNrWXhPVkY2Zkpud0NZcFRzRVFCYXhZNWkxU0NHOGsvbFAvYjloOTlBaWFvOGRzOTFVQ1M3eGJteUVjZXdZVGoxT3hoTnRZYzNUbVR2azFTRlZyQTNsUlptbzVwbG1COCtod0tGU0dwbzg0MnNzcUFhLzd2bzN0Y0M0RU81RGp1R2NtcFArL1hhM0hWK1E3U01Sd2JIR2h4VCtiYXRuMHdwalhCM09JS3d6QUNZZTJKUEtkZVpjN1U1T1FxcXZwSUYzSWZSeVozSUd3emgwOXlaenQrb2kyNXFiNnloOUwvV3QydjZQeS9GVGRDWEkrSEtHbUxwS0ZWbWVmek1oQk5oelRuMDZBWlExRlJRamNQckJ5dTdtVUJFaFRmMjl0THZ0MjdiRGladE1LZjcwMnFNS1YvWXhsUzYzY1pmcnZkR3RPUk9HUGE5TU5sWndsZlFySU5jL2E4cmNUdFo2S1lwamptWnloa3pYdEl5T0JNeGxSZTFISm1SQ0MwcUJrbnVhYW5ubEVHajF4QVVzNDNHUzdMSEszSENCZWRlUEpxdGlXVHJwQzFhaUpYeFhnVkNpbzZFZ3oydnAzbVJDRjR4SWQzQklqUnRnNWw5Tm9BMnVVSHVrbE5KMGhyM2tvVzVvUEV6dnRyeWdtcVlEVjJtWmIxL3JLNDdDVVo4a0pGZTNsWWsrOFUxUmJLKytxL2ZEcHllYnkyM1lSRTMrNWdUa2FyWU9LdExtUDIyRWhYSzJHRlh0SnQrb0Q2NTJlWkdLbjBaMFgwbjJWcHczRGI0dWhGazUrZU81YmROelc3SCt3cW8yaHB0eXN1cERyN1NoNm84eVFyVW0rTlZ1TEFOeDA4ZDJ1SGE1d3FlbWdJcG0rdUl2b3FyY1kwMy9kVmtPbEtLWU1tRXZESmhRcE5DL1IyVG5UeDJkaThWNWJEMjVyc29QbFdDMVdYSEdVNEpGdHBUMGV0akl6VnhpSStrd0prdUFZNnNZUEY4cEFCYk16Q2YvV245SUFHdW1rWG54SkhyZkh0OW9QSVVTMldBWitxRE9FSTBDNTM5SExPVit2MUplL0xST3VTVWl3ZUF1eHFmS2VkZWlVajNib1FGQT0="
	//encryptKeyUsed := strings.Builder{}
	//encryptKeyUsed.WriteString(key)
	//if len(encryptKey) < 32 {
	//	encryptKeyUsed.Write(bytes.Repeat([]byte{byte(0)}, 32-len(encryptKey)))
	//}
	rawBase64Decoded, err := base64.StdEncoding.DecodeString(data)
	err, s := Ase256CBCDecode(string(rawBase64Decoded), []byte(key), []byte(iv))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(s))
	}
	//frame := event2.ParseFrameMapByData(decodedStr)
	//if frame == nil {
	//	t.Error("frame is empty")
	//} else {
	//	t.Log(fmt.Sprintf("%+v", frame))
	//}

}

func TestEncryptData2(t *testing.T) {
	//data := map[string]interface{}{
	//	"id":    1,
	//	"time":  1,
	//	"bot":   false,
	//	"extra": map[string]interface{}{"ip": "127.0.0.1"},
	//}
	//dataBytes, _ := sonic.Marshal(data)
	dataBytes := `{"id":1,"time":1,"bot":false,"extra":{"ip":"127.0.0.1"}}`
	t.Logf("json encode:%s", string(dataBytes))
	//encryptedData, err := Aes256Encode(dataBytes, "sfjlfvx98sdlsdfsdfsdfD)9DJ&sdf", "9051dbfc59d6af1d", 32)
	// 密钥（必须是 32 字节，256 位）
	key := "sfjlfvx98sdlsdfsdfsdfD)9DJ&sdf"

	// 初始化向量（必须是 16 字节）
	iv := "9051dbfc59d6af1d"

	encryptedData, err := encryptAES256CBC([]byte(dataBytes), []byte(key), []byte(iv))
	if err != nil {
		t.Error(err)
		return
	}
	cipherBytesBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(encryptedData)))
	base64.StdEncoding.Encode(cipherBytesBase64, encryptedData)
	t.Logf("%s", string(cipherBytesBase64))
}

func TestEncryptData3(t *testing.T) {
	data := map[string]interface{}{
		"id":    1,
		"time":  time.Now().Unix(),
		"bot":   false,
		"extra": map[string]interface{}{"ip": "127.0.0.1"},
	}
	dataBytes, _ := sonic.Marshal(data)
	err, encryptedData := Aes256CBCEncode(string(dataBytes), []byte("sfjlfvx98sdlsdfsdfsdfD)9DJ&sdf"), []byte("9051dbfc59d6af1d"))
	if err != nil {
		log.WithError(err).Error("Ase256Encode")
		return
	}
	cipherBytesBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(encryptedData)))
	base64.StdEncoding.Encode(cipherBytesBase64, encryptedData)
	t.Logf("%s", string(cipherBytesBase64))
}
