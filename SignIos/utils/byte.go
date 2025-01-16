package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

// BytesCombine2 字节合并
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// BytesToUint64 字节数组转换成64位数组
func BytesToUint64(bytes []byte) []uint64 {
	var res_list []uint64
	for _, v := range bytes {
		res_list = append(res_list, uint64(v))
	}
	return res_list
}

// Bswap_4 单字节大小端转换
func Bswap_4(x uint8) uint8 {
	return ((x & 0xF) << 4) | ((x >> 4) & 0xF)
}

func Bswap_1(x uint8) uint8 {
	var tmp uint8 = 0
	for i := 0; i < 8; i++ {
		tmp |= (((x >> i) & 0x1) << (7 - i)) & 0xFF
	}
	return tmp
}

func Uint64toBytes(data uint64) []byte {
	hexStr := UInt64ListToHexStr([]uint64{data})
	byteData, _ := hex.DecodeString(hexStr)
	return byteData
}

func Int64ToBytes(n int64) []byte {
	data := int64(n)
	bytebuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytebuffer, binary.BigEndian, data)
	return bytebuffer.Bytes()
}
