package ios

import (
	pb "SignIos/proto"
	"SignIos/utils"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"

	"google.golang.org/protobuf/proto"
)

// ios 加密 解密相关

func EncryptArgus(info *pb.XArgus) (err error, res string) {
	SignKey := utils.Ios_DouYinSignKey
	key_4 := utils.GetRandomKey_4()
	key_bytes, err := base64.StdEncoding.DecodeString(SignKey)
	if err != nil {
		return err, ""
	}
	tmp_key8 := utils.Calc_key_8(key_4)
	sm3_buff_tmp := utils.BytesCombine(key_bytes, key_4, key_bytes)
	sm3_hash_tmp := utils.Sm3Hash(sm3_buff_tmp)
	aesKey := utils.BytesToMd5Byte(key_bytes[0:16])
	aesIv := utils.BytesToMd5Byte(key_bytes[16:])
	data, err := proto.Marshal(info)
	if err != nil {
		return err, ""
	}
	encryptData := utils.IosEncryptArgusData(data, sm3_hash_tmp)
	encryptBuf := utils.IosEncryptArgusBuf(encryptData)
	encryptBuf = utils.BytesCombine(info.Env, encryptBuf)
	xorData := utils.Xor(encryptBuf, tmp_key8)
	head_0 := utils.XArgusCalcHead_0(key_bytes)
	head_8 := utils.XArgusCalcHead_8(info.QueryHash[0], info.BodyHsh[0]) // ok
	beforeAesData := utils.BytesCombine(head_0, head_8, xorData, key_4[2:4])
	FirstAesData := utils.AES_CBC_Encrypt(beforeAesData, aesKey, aesIv)
	FirstAesData = utils.BytesCombine(key_4[0:2], FirstAesData)
	baseOver := base64.StdEncoding.EncodeToString(FirstAesData)
	return nil, baseOver
}
func EncryptLadon(Khronos int64, aid, lc_id string) (err error, res string) {
	int64Key_4 := utils.GetRandInt64()
	key_4 := utils.Int64ToBytes(int64Key_4)
	key_4 = key_4[4:]
	dataStr := strconv.FormatInt(Khronos, 10) + "-" + lc_id + "-" + aid
	dataByte := []byte(dataStr)
	dataPaddingNum := 16 - (len(dataByte) % 16)
	AddPadding := utils.GetPadding(dataPaddingNum)
	data := utils.BytesCombine(dataByte, AddPadding)
	byte_aid := []byte(aid)
	ByteKeyData := utils.BytesCombine(key_4, byte_aid)
	Md5Byte := utils.BytesToMd5Byte(ByteKeyData)
	Md5ByteHexStr := hex.EncodeToString(Md5Byte)
	NewByteMd5 := []byte(Md5ByteHexStr)
	aesBuf := NewByteMd5[16:]
	outData := utils.IosEncryptLadon(data, aesBuf, Md5Byte)
	outData = utils.BytesCombine(ByteKeyData[0:4], outData)
	base64DataStr := base64.StdEncoding.EncodeToString(outData)
	return nil, base64DataStr
}
func EncryptGorgon(Khronos int64, queryHash, bodyHash []byte, sdkVersionStr string) (err error, res string) {
	GorgonKey := utils.Ios_DouYinGorgonkey
	UintPtr := utils.GetRandPtr()
	Word := UintPtr & 0xFFF0
	key_2 := uint8(Word & 0xFF)        // ok
	key_3 := uint8((Word >> 8) & 0xFF) // ok
	tmpHead := []byte{0x84, 0x04, key_2, key_3, GorgonKey[1], GorgonKey[6]}
	GorgonKey[3] = tmpHead[3]
	GorgonKey[7] = tmpHead[2]
	if queryHash == nil {
		return errors.New("queryHash is nil"), ""
	}
	queryHash = queryHash[0:4]
	if bodyHash == nil {
		bodyHash = []byte{0x0, 0x0, 0x0, 0x0}
	} else {
		bodyHash = bodyHash[0:4]
	}
	unUsed := []byte{0x0, 0x0, 0x0, 0x0}
	sdkVersion, err := hex.DecodeString(sdkVersionStr)
	if err != nil {
		return err, ""
	}
	KhronosBytes := utils.Int64ToBytes(Khronos)
	inputData := utils.BytesCombine(queryHash, bodyHash, unUsed, sdkVersion, KhronosBytes[4:])
	fullKey := utils.XgprgonInitKey(GorgonKey)
	SecendData := utils.XGorgonDataEncryptWithKey(inputData, fullKey)
	LastData := utils.EncryptXGorgonData(SecendData)
	LastData = utils.BytesCombine(tmpHead, LastData)
	return nil, hex.EncodeToString(LastData)
}
func JiSuEncryptArgus(info *pb.XArgus) (err error, res string) {
	SignKey := "utils.Ios_DouYinJiSuSignKey"
	key_4 := utils.GetRandomKey_4()
	key_bytes, err := base64.StdEncoding.DecodeString(SignKey)
	if err != nil {
		return err, ""
	}
	tmp_key8 := utils.Calc_key_8(key_4)
	sm3_buff_tmp := utils.BytesCombine(key_bytes, key_4, key_bytes)
	sm3_hash_tmp := utils.Sm3Hash(sm3_buff_tmp)
	aesKey := utils.BytesToMd5Byte(key_bytes[0:16])
	aesIv := utils.BytesToMd5Byte(key_bytes[16:])
	data, err := proto.Marshal(info)
	if err != nil {
		return err, ""
	}
	encryptData := utils.IosEncryptArgusData(data, sm3_hash_tmp)
	encryptBuf := utils.IosEncryptArgusBuf(encryptData)
	encryptBuf = utils.BytesCombine(info.Env, encryptBuf)
	xorData := utils.Xor(encryptBuf, tmp_key8)
	head_0 := utils.XArgusCalcHead_0(key_bytes)
	head_8 := utils.XArgusCalcHead_8(info.QueryHash[0], info.BodyHsh[0]) // ok
	beforeAesData := utils.BytesCombine(head_0, head_8, xorData, key_4[2:4])
	FirstAesData := utils.AES_CBC_Encrypt(beforeAesData, aesKey, aesIv)
	baseOver := base64.StdEncoding.EncodeToString(FirstAesData)
	return nil, baseOver
}
func DuoShanEncryptArgus(info *pb.XArgus) (err error, res string) {
	SignKey := utils.Ios_DouYinDuoShanSignKey
	key_4 := utils.GetRandomKey_4()
	key_bytes, err := base64.StdEncoding.DecodeString(SignKey)
	if err != nil {
		return err, ""
	}
	tmp_key8 := utils.Calc_key_8(key_4)
	sm3_buff_tmp := utils.BytesCombine(key_bytes, key_4, key_bytes)
	sm3_hash_tmp := utils.Sm3Hash(sm3_buff_tmp)
	aesKey := utils.BytesToMd5Byte(key_bytes[0:16])
	aesIv := utils.BytesToMd5Byte(key_bytes[16:])
	data, err := proto.Marshal(info)
	if err != nil {
		return err, ""
	}
	encryptData := utils.IosEncryptArgusData(data, sm3_hash_tmp)
	encryptBuf := utils.IosEncryptArgusBuf(encryptData)
	encryptBuf = utils.BytesCombine(info.Env, encryptBuf)
	xorData := utils.Xor(encryptBuf, tmp_key8)
	head_0 := utils.XArgusCalcHead_0(key_bytes)
	head_8 := utils.XArgusCalcHead_8(info.QueryHash[0], info.BodyHsh[0]) // ok
	beforeAesData := utils.BytesCombine(head_0, head_8, xorData, key_4[2:4])
	FirstAesData := utils.AES_CBC_Encrypt(beforeAesData, aesKey, aesIv)
	FirstAesData = utils.BytesCombine(key_4[0:2], FirstAesData)
	baseOver := base64.StdEncoding.EncodeToString(FirstAesData)
	return nil, baseOver
}
