package common

import (
	"SignIos/utils"
	"encoding/base64"
)

// 抖音
func TokenDecrypt(input []byte) (res []byte, err error) {
	step_one, err := tokenDecryptStep(input)
	if err != nil {
		return nil, err
	}
	overData := DoZlibUnCompress(step_one[4:])
	return overData, nil
}
func tokenDecryptStep(input []byte) (res []byte, err error) {
	XteaKey := utils.XteaKey
	CommonKey := utils.Ios_DouYinCommonKey
	key_bytes, err := base64.StdEncoding.DecodeString(CommonKey)
	if err != nil {
		return nil, err
	}
	aesKey := key_bytes[0:16]
	aesIv := key_bytes[16:]
	aesOver := utils.AES_CBC_Decrypt(input, aesKey, aesIv)
	var xtea_iv []byte
	xtea_iv = utils.BytesCombine(xtea_iv, aesOver[len(aesOver)-4:])
	xtea_iv = utils.BytesCombine(xtea_iv, XteaKey[16:])
	err, resData := utils.MyXteaEncrypt(aesOver[1:len(aesOver)-4], XteaKey[:16], xtea_iv)
	if err != nil {
		return nil, err
	}
	pad := resData[0] & 0x7
	randz := (pad + 7) & 3
	tmp_len := len(resData) - (int(pad) + 8)
	return resData[int(randz)+1 : int(randz)+tmp_len+1], nil
}

// 抖音极速版
func TokenDecryptByJiSu(input []byte) (res []byte, err error) {
	step_one, err := tokenDecryptStepByJisu(input)
	if err != nil {
		return nil, err
	}
	overData := DoZlibUnCompress(step_one[4:])
	return overData, nil
}
func tokenDecryptStepByJisu(input []byte) (res []byte, err error) {
	XteaKey := utils.XteaKey
	CommonKey := utils.Ios_DouYinJiSuCommonKey
	key_bytes, err := base64.StdEncoding.DecodeString(CommonKey)
	if err != nil {
		return nil, err
	}
	aesKey := key_bytes[0:16]
	aesIv := key_bytes[16:]
	aesOver := utils.AES_CBC_Decrypt(input, aesKey, aesIv)
	var xtea_iv []byte
	xtea_iv = utils.BytesCombine(xtea_iv, aesOver[len(aesOver)-4:])
	xtea_iv = utils.BytesCombine(xtea_iv, XteaKey[16:])
	err, resData := utils.MyXteaEncrypt(aesOver[1:len(aesOver)-4], XteaKey[:16], xtea_iv)
	if err != nil {
		return nil, err
	}
	pad := resData[0] & 0x7
	randz := (pad + 7) & 3
	tmp_len := len(resData) - (int(pad) + 8)
	return resData[int(randz)+1 : int(randz)+tmp_len+1], nil
}

// 抖音多闪
func TokenDecryptByDuoShan(input []byte) (res []byte, err error) {
	step_one, err := tokenDecryptStepByDuoShan(input)
	if err != nil {
		return nil, err
	}
	overData := DoZlibUnCompress(step_one[4:])
	return overData, nil
}
func tokenDecryptStepByDuoShan(input []byte) (res []byte, err error) {
	XteaKey := utils.XteaKey
	CommonKey := utils.Ios_DouYinDuoShanCommonKey
	key_bytes, err := base64.StdEncoding.DecodeString(CommonKey)
	if err != nil {
		return nil, err
	}
	aesKey := key_bytes[0:16]
	aesIv := key_bytes[16:]
	aesOver := utils.AES_CBC_Decrypt(input, aesKey, aesIv)
	var xtea_iv []byte
	xtea_iv = utils.BytesCombine(xtea_iv, aesOver[len(aesOver)-4:])
	xtea_iv = utils.BytesCombine(xtea_iv, XteaKey[16:])
	err, resData := utils.MyXteaEncrypt(aesOver[1:len(aesOver)-4], XteaKey[:16], xtea_iv)
	if err != nil {
		return nil, err
	}
	pad := resData[0] & 0x7
	randz := (pad + 7) & 3
	tmp_len := len(resData) - (int(pad) + 8)
	return resData[int(randz)+1 : int(randz)+tmp_len+1], nil
}
