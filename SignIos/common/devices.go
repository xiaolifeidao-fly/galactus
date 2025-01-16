package common

import (
	"SignIos/utils"
	"encoding/base64"
	"encoding/hex"
)

func Base64StrDecrypt(base64Str string) (err error, res string) {
	byteData, err := base64.URLEncoding.DecodeString(base64Str)
	if err != nil {
		return err, ""
	}
	err, resByte := Decrypt(byteData)
	if err != nil {
		return err, ""
	}
	return nil, string(resByte)
}
func HexStrDecrypt(hexStr string) (err error, res string) {
	byteInput, err := hex.DecodeString(hexStr)
	if err != nil {
		return err, ""
	}
	err, decryptData := Decrypt(byteInput)
	if err != nil {
		return err, ""
	}
	return nil, string(decryptData)
}
func Decrypt(data []byte) (err error, res []byte) {
	randomBuf := data[6:38]
	hashRandom := utils.Sha512Hash(randomBuf)
	KeyInfo := utils.BytesCombine(hashRandom, utils.Common_Devices_Key_One)
	hashKey := utils.Sha512Hash(KeyInfo)
	Key := hashKey[0:16]
	Iv := hashKey[16:32]
	aesDecryptData := utils.AES_CBC_Decrypt(data[38:], Key, Iv)
	unGzipData := UnGzipData(aesDecryptData[64:])
	return nil, unGzipData
}
func EncryptToBase64(input []byte) (err error, res string) {
	err, encryptData := Encrypt(input, 0)
	if err != nil {
		return err, ""
	}
	return nil, base64.URLEncoding.EncodeToString(encryptData)
}
func Encrypt(inputData []byte, os uint8) (err error, res []byte) {
	// android os=0 ios os=0x13
	err, GzipData := GzipData(inputData)
	if err != nil {
		return err, nil
	}
	GzipData[9] = os
	randomInt := utils.GetRandInt64()
	byteRandom := utils.Int64ToBytes(randomInt)
	randomBuf := utils.Sm3Hash(byteRandom)
	hashRandom := utils.Sha512Hash(randomBuf)
	KeyInfo := utils.BytesCombine(hashRandom, utils.Common_Devices_Key_One)
	hashKey := utils.Sha512Hash(KeyInfo)
	Key := hashKey[0:16]
	Iv := hashKey[16:32]
	GzipDataHash := utils.Sha512Hash(GzipData)
	waitAesData := utils.BytesCombine(GzipDataHash, GzipData)
	aesData := utils.AES_CBC_Encrypt(waitAesData, Key, Iv)
	LastData := utils.BytesCombine(utils.Common_Devices_Head, randomBuf, aesData)
	return nil, LastData
}
