package utils

import (
	"errors"
	"fmt"
	"strconv"
)

// HexStrToInt64List 将十六进制字符串转成 64位数字
func HexStrToInt64List(str string) (list []int64, err error) {
	if len(str)%8 != 0 {
		return nil, errors.New("十六进制字符串长度不是8的倍数!请检查!")
	}
	for i := 0; i < len(str); i += 16 {
		temp := "0x" + str[i+14:i+16] + str[i+12:i+14] + str[i+10:i+12] + str[i+8:i+10] + str[i+6:i+8] + str[i+4:i+6] + str[i+2:i+4] + str[i:i+2]
		temp_Int, err := strconv.ParseInt(temp, 0, 0)
		if err != nil {
			return nil, err
		}
		list = append(list, temp_Int)
	}
	return list, nil
}

// HexStrToUInt64List 将十六进制字符串转成 无符号64位数字
func HexStrToUInt64List(str string) (list []uint64, err error) {
	if len(str)%8 != 0 {
		return nil, errors.New("十六进制字符串长度不是8的倍数!请检查!")
	}
	for i := 0; i < len(str); i += 16 {
		temp := "0x" + str[i+14:i+16] + str[i+12:i+14] + str[i+10:i+12] + str[i+8:i+10] + str[i+6:i+8] + str[i+4:i+6] + str[i+2:i+4] + str[i:i+2]
		temp_Int, err := strconv.ParseUint(temp, 0, 0)
		if err != nil {
			return nil, err
		}
		list = append(list, uint64(temp_Int))
	}
	return list, nil
}
// HexToInt64 十六进制字符串转int64
func HexToInt64(str string) int64 {
	if len(str)>8{
		return 0
	}
	temp := "0x" + str
	temp_Int, _ := strconv.ParseUint(temp, 0, 0)
	return int64(temp_Int)
}

// Int64ToHex 将64位数字转换为16位长度的十六进制字符串
func Int64ToHex(int64_ int64) string {
	H := fmt.Sprintf("%016x", uint64(int64_))
	return H
}

// Int64ListToHexStr 64位数字数组转换成16进制字符串
func Int64ListToHexStr(list []int64) string {
	var str string
	for _, v := range list {
		H := fmt.Sprintf("%016x", uint64(v))
		temp := H[14:16] + H[12:14] + H[10:12] + H[8:10] + H[6:8] + H[4:6] + H[2:4] + H[0:2]
		str += temp
	}
	return str
}

// UInt64ListToHexStr 无符号64位数字数组转换成16进制字符串
func UInt64ListToHexStr(list []uint64) string {
	var str string
	for _, v := range list {
		H := fmt.Sprintf("%016x", v)
		temp := H[14:16] + H[12:14] + H[10:12] + H[8:10] + H[6:8] + H[4:6] + H[2:4] + H[0:2]
		str += temp
	}
	return str
}
func LogHex(v int64) {
	H := fmt.Sprintf("%016x", uint64(v))
	fmt.Println(H)
}
func UintLogHex(v uint64) {
	H := fmt.Sprintf("%016x", v)
	fmt.Println(H)
}
