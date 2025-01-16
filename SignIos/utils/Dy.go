package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"math/big"
)

// calc_key_8 计算第八位
func Calc_key_8(key4 []byte) []byte {
	k8 := (uint64(key4[3]) | (uint64(key4[2]) << 11)) ^ (uint64(key4[2]) >> 5) ^ uint64(key4[2]) ^ 0xFFFFFFFF
	U32Bytes := ToU32(k8)
	return U32Bytes
}

// ToU32 将数值转换为 byte32
func ToU32(d uint64) []byte {
	var res_bytes []byte
	res_bytes = append(res_bytes, uint8((d>>24)&0xFF))
	res_bytes = append(res_bytes, uint8((d>>16)&0xFF))
	res_bytes = append(res_bytes, uint8((d>>8)&0xFF))
	res_bytes = append(res_bytes, uint8((d>>0)&0xFF))
	return res_bytes
}

// Xor 加密Xor
func Xor(data, key []byte) []byte {
	var res_bytes []byte
	for i := 0; i < len(data); i++ {
		temp := data[len(data)-i-1] ^ key[i%4]
		res_bytes = append(res_bytes, temp)
	}
	return res_bytes
}

// DeXor 解密Xor
func DeXor(data, key []byte) []byte {
	var res_bytes []byte
	res_bytes = make([]byte, len(data))
	copy(res_bytes, data)
	for i := 0; i < len(data); i++ {
		temp := data[i] ^ key[i%4]
		res_bytes[len(res_bytes)-i-1] = temp
	}
	return res_bytes
}

// 加密 EncryptXGorgonData Xgorgon
func EncryptXGorgonData(data []byte) []byte {
	var resData = make([]byte, 20)
	copy(resData, data)
	for i := 0; i < len(resData); i++ {
		resData[i] = Bswap_4(resData[i]) ^ resData[(i+1)%20]
		resData[i] = Bswap_1(resData[i]) ^ 0xFF ^ 20
	}
	return resData
}

// 解密 DecryptXGorgonData Xgorgon
func DecryptXGorgonData(data []byte) []byte {
	var resData = make([]byte, 50)
	copy(resData, data)
	for i := 19; i >= 0; i-- {
		resData[i] = resData[i] ^ 0xFF ^ 20
		resData[i] = Bswap_1(resData[i])
		resData[i] = resData[i] ^ resData[(i+1)%20]
		resData[i] = Bswap_4(resData[i])
	}
	return resData
}

// XGorgonDataEncryptWithKey 加密数据
func XGorgonDataEncryptWithKey(data, key []byte) []byte {
	var tmp1 uint8
	var tmp2 uint8
	var resData = make([]byte, 20)
	var key256 = make([]byte, 20)
	copy(resData, data)
	copy(key256, key)
	for i := 1; i <= 20; i++ {
		tmp2 = (key[i] + tmp1) & 0xFF
		tmp1 = tmp2
		tmp2 = key[tmp2]
		key[i] = tmp2
		tmp2 = key[(tmp2+tmp2)&0xFF]
		resData[i-1] = resData[i-1] ^ tmp2
	}
	return resData
}

// XgorgonInitKey 初始化密钥
func XgprgonInitKey(key []byte) []byte {
	key256 := make([]byte, 256)
	for i := 0; i < 256; i++ {
		key256[i] = uint8(i)
	}
	var tmp uint8
	for i := 0; i < 256; i++ {
		tmp = (key[uint8(i)%8] + uint8(i) + tmp) & 0xFF
		key256[i] = key256[tmp]
	}
	return key256
}
func XArgusCalcHead_0(key []byte) []byte {
	var sum uint32
	for i := 0; i < 32; i++ {
		if (i & 1) == 0 {
			sum = uint32(key[i]) ^ (sum << 7) ^ (sum >> 3) ^ sum
		} else {
			sum = (uint32(key[i]) | (sum << 11)) ^ (sum >> 5) ^ sum ^ 0xFFFFFFFF
		}
	}
	return []byte{byte(sum & 0xFF)}
}

func XArgusCalcHead_8(queryHashHeader uint8, bodyHashHeader uint8) []byte {
	var temp uint64
	temp = ((uint64(queryHashHeader) & 0x3F) << 46) | ((uint64(bodyHashHeader) & 0x3F) << 40)
	randNum := RandUint64()
	temp = 0x1810000100000000 | temp | randNum
	return Uint64toBytes(temp)
}

func RandUint64() uint64 {
	return uint64(RangeRand(0, 999999999))
}
func GetRandPtr() uint64 {
	return uint64(RangeRand(412316860416, 549755813887))
}
func GetRandInt64() int64 {
	return RangeRand(0, 4294967295)
}

func RangeRand(min, max int64) int64 {
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

func GetRandomKey_4() []byte {
	var res_byte []byte
	for i := 0; i < 4; i++ {
		res_byte = append(res_byte, uint8(RandUint64()&0xFF))
	}
	return res_byte
}
func GetPadding(num int) []byte {
	var resByte []byte
	for i := 0; i < num; i++ {
		resByte = append(resByte, uint8(num))
	}
	return resByte
}

var QuerySm3AddPadding []byte = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31}

func GetQuerySm3(inputParams []byte, stub []byte) []byte {
	var newParams []byte
	if stub == nil {
		newParams = BytesCombine(inputParams, QuerySm3AddPadding)
	} else {
		newParams = BytesCombine(inputParams, stub, []byte{0x31})
	}
	return Sm3Hash(newParams)
}

// ios
func IosDecryptArgusBuf(data []byte) []byte {
	var res_data = make([]byte, len(data))
	copy(res_data, data)
	var sum = 0
	for i := 1; i < len(res_data); i++ {
		sum = (sum + int(res_data[i])) & 0xFF
	}
	res_data[0] = (res_data[0] - uint8(sum)) ^ res_data[1]
	res_data[len(res_data)-1] ^= res_data[len(res_data)-2]
	for i := len(res_data) - 2; i >= 2; i-- {
		res_data[i] -= ROL_8(res_data[i-1], 1) ^ res_data[i-2] ^ uint8(i) ^ 0xFF
	}
	res_data[1] -= res_data[0] ^ res_data[len(res_data)-1] ^ 0xFE
	res_data[0] -= res_data[len(res_data)-1] ^ res_data[len(res_data)-2] ^ 0xFF
	return res_data
}
func IosEncryptArgusBuf(data []byte) []byte {
	var res_data = make([]byte, len(data))
	copy(res_data, data)
	var sum uint8
	sum = 0
	res_data[0] += res_data[len(res_data)-1] ^ res_data[len(res_data)-2] ^ 0xFF
	res_data[1] += res_data[0] ^ res_data[len(res_data)-1] ^ 0xFE
	for i := 2; i < len(res_data)-1; i++ {
		res_data[i] += ROL_8(res_data[i-1], 1) ^ res_data[i-2] ^ uint8(i) ^ 0xFF
	}
	res_data[len(res_data)-1] ^= res_data[len(res_data)-2]
	for i := 1; i < len(res_data); i++ {
		sum = (sum + res_data[i]) & 0xFF
	}
	res_data[0] = (res_data[0] ^ res_data[1]) + sum
	return res_data
}

func IosDecryptArgusData(data, key []byte) []byte {
	var res_data = make([]byte, len(data))
	copy(res_data, data)
	var tmp uint8
	var k uint8
	for i := 0; i < len(res_data); i++ {
		tmp = data[len(data)-1-i]
		tmp = ROR_8((tmp^0xFF^key[k%32])-key[k%32+1], 5)
		tmp = ROR_8((tmp^0xFF^key[k%32+1])-key[k%32], 6)
		res_data[i] = tmp
		k += 4
	}
	return res_data
}
func IosEncryptArgusData(data, key []byte) []byte {
	var res_data = make([]byte, len(data))
	copy(res_data, data)
	var tmp uint8
	var k uint8
	for i := 0; i < len(data); i++ {
		tmp = data[i]
		tmp = (ROL_8(tmp, 6) + key[k%32]) ^ key[k%32+1] ^ 0xFF
		tmp = (ROL_8(tmp, 5) + key[k%32+1]) ^ key[k%32] ^ 0xFF
		res_data[len(res_data)-1-i] = tmp
		k += 4
	}
	return res_data
}
func IosDecryptLadon(data, key_1, key_2 []byte) []byte {
	var res_data []byte
	hexStrKey := hex.EncodeToString(key_2)
	bytesHexKey := []byte(string(hexStrKey))[0:16]
	for i := 0; i < len(data); i += 16 {
		aesData := EcbEncrypt(key_1, bytesHexKey)
		key_1 = aesData[:16]
		xorOut := IosLadonXor(data[i:i+16], aesData)
		res_data = BytesCombine(res_data, xorOut)
	}
	padding_num := res_data[len(res_data)-1]
	int_padding_num := int(padding_num)
	res_data = res_data[0 : len(res_data)-int_padding_num]
	return res_data
}
func IosEncryptLadon(data, key_1, key_2 []byte) []byte {
	var res_data []byte
	hexStrKey := hex.EncodeToString(key_2)
	bytesHexKey := []byte(string(hexStrKey))[0:16]
	for i := 0; i < len(data); i += 16 {
		aesData := EcbEncrypt(key_1, bytesHexKey)
		key_1 = aesData[:16]
		xorOut := IosLadonXor(data[i:i+16], aesData)
		res_data = BytesCombine(res_data, xorOut)
	}
	return res_data
}

func IosLadonXor(data, key []byte) []byte {
	var res_data = make([]byte, len(data))
	copy(res_data, data)
	for i := 0; i < len(data); i++ {
		res_data[i] = data[i] ^ key[i]
	}
	return res_data
}

func ClearPadding(data []byte) []byte {
	if len(data) == 0 || (data[len(data)-1] != data[len(data)-2]) {
		return data
	}
	lastByte := data[len(data)-1]
	paddNum := 1
	for i := len(data) - 2; i > -0; i-- {
		if data[i] == lastByte {
			paddNum++
		} else {
			break
		}
	}
	return data[0 : len(data)-paddNum]
}

// ROR_64 ROR 64位版本
func ROR_64(x, n uint64) uint64 {
	return x>>n | x<<(64-n)
}
func ROR_8(x, n uint8) uint8 {
	return (x >> n) | (x << (8 - n))
}

// ROL_64 ROL 64位版本
func ROL_64(x, n uint64) uint64 {
	return x<<n | x>>(64-n)
}
func ROL_8(x, n uint8) uint8 {
	return (x << n) | (x >> (8 - n))
}

// Host: api5-normal-lq.amemv.com
// Connection: keep-alive
// Content-Length: 99
// x-tt-multi-sids: 1115358952428345%3A2c2476fd01f18f6c9a5c7615fa857072
// X-SS-Cookie: store-region=cn-fj; store-region-src=uid; install_id=1871824318410835; ttreq=1$275454359d83880e4765589dfefebb1041cc3c84; multi_sids=1115358952428345%3A2c2476fd01f18f6c9a5c7615fa857072; n_mh=EJpGAl8UyrYuf8qp_T0CjxfWrnBWnkqZm_Fq-wJea5U; odin_tt=21ea15607196e79fa110c11b80466f643499f1c90d7f70cfc28c7c056e7845df02525a7a1e92fdbf54a38ed92a8ca0725ebfdbcf5c3097232b44727a3f79d5b3efb345189880db01fd8fcbe5a12b8496; passport_assist_user=CkELlcduQ6VcQe7ureljsdz4XrBrFT6Ggue3kqykxn_x5Jk9uRWr98Ntxl6aeHuzsYvkhIcfKtjgLmt6zIBiiPoWkxpKCjx5_Er76C7Pa4u29aURxtxTTbX1SApjMsmQy5OmD_mGkLCh-h5KZw0ohNJqK32SMYaRA7M2M0Qjmy-km1cQzMPXDRiJr9ZUIAEiAQN7Miu_; sessionid=2c2476fd01f18f6c9a5c7615fa857072; sessionid_ss=2c2476fd01f18f6c9a5c7615fa857072; sid_guard=2c2476fd01f18f6c9a5c7615fa857072%7C1721825582%7C5184000%7CSun%2C+22-Sep-2024+12%3A53%3A02+GMT; sid_tt=2c2476fd01f18f6c9a5c7615fa857072; uid_tt=7ee5995032326a3704ca3789e090f936; uid_tt_ss=7ee5995032326a3704ca3789e090f936; d_ticket=df117dbc38bcb4fbbd89cf21b62669670be58; ticket_guard_has_set_public_key=1; passport_csrf_token=3a218467d059ba2d72eb5c54a99638d4; passport_csrf_token_default=3a218467d059ba2d72eb5c54a99638d4
// sdk-version: 2
// Content-Type: application/x-www-form-urlencoded
// x-Tt-Token: 002c2476fd01f18f6c9a5c7615fa85707200aa4004bd57a564f12e66e9118999b282668eec4eb2264190e52a79fa6063974cd7b58a5a6a7958a03a01de9ab48b131a6d00421272452805d894ba18fcedcc519fc8db6872420818e0213971318301f94-1.0.1
// User-Agent: Aweme 19.0.0 rv: 190015 (iPhone; iOS 13.5.1; zh_CN) Cronet
// x-vc-bdturing-sdk-version: 2.2.3
// tt-request-time: 1721825902361
// Cookie: passport_csrf_token=3a218467d059ba2d72eb5c54a99638d4; passport_csrf_token_default=3a218467d059ba2d72eb5c54a99638d4; ticket_guard_has_set_public_key=1; d_ticket=df117dbc38bcb4fbbd89cf21b62669670be58; multi_sids=1115358952428345%3A2c2476fd01f18f6c9a5c7615fa857072; n_mh=EJpGAl8UyrYuf8qp_T0CjxfWrnBWnkqZm_Fq-wJea5U; odin_tt=21ea15607196e79fa110c11b80466f643499f1c90d7f70cfc28c7c056e7845df02525a7a1e92fdbf54a38ed92a8ca0725ebfdbcf5c3097232b44727a3f79d5b3efb345189880db01fd8fcbe5a12b8496; passport_assist_user=CkELlcduQ6VcQe7ureljsdz4XrBrFT6Ggue3kqykxn_x5Jk9uRWr98Ntxl6aeHuzsYvkhIcfKtjgLmt6zIBiiPoWkxpKCjx5_Er76C7Pa4u29aURxtxTTbX1SApjMsmQy5OmD_mGkLCh-h5KZw0ohNJqK32SMYaRA7M2M0Qjmy-km1cQzMPXDRiJr9ZUIAEiAQN7Miu_; sessionid=2c2476fd01f18f6c9a5c7615fa857072; sessionid_ss=2c2476fd01f18f6c9a5c7615fa857072; sid_guard=2c2476fd01f18f6c9a5c7615fa857072%7C1721825582%7C5184000%7CSun%2C+22-Sep-2024+12%3A53%3A02+GMT; sid_tt=2c2476fd01f18f6c9a5c7615fa857072; uid_tt=7ee5995032326a3704ca3789e090f936; uid_tt_ss=7ee5995032326a3704ca3789e090f936; install_id=1871824318410835; ttreq=1$275454359d83880e4765589dfefebb1041cc3c84; store-region=cn-fj; store-region-src=uid
// x-tt-passport-csrf-token: 3a218467d059ba2d72eb5c54a99638d4
// x-tt-dt: AAAQ4WJS73ZVVSNCWESBA54VFQB4RCHYFLHJ3DD72RDGXDZKKQ43B4ZFGNFKUMNPZSBW42CPGD52JKJX4SYK346NIPPZOHCGAYEJBN325ULZY3XUXQUSXUAKGDSVS
// passport-sdk-version: 5.12.1
// X-SS-STUB: 28DF655EB957A67C17694DB9C36EFD79
// x-bd-kmsv: 1
// x-tt-request-tag: s=1
// X-SS-DP: 1128
// x-tt-trace-id: 00-e4d240b30dde665daec2b8b01fb30468-e4d240b30dde665d-01
// Accept-Encoding: gzip, deflate, br
// X-Argus: Mc+UNgplTZwpTTQ6zx/VWmzMcs483CHkHalpP8JkZiB+x8erMjTBxgjUFvKYWxJw9C3DCsJ1QUfs+/IXnU9SfQxAtgXpRRfNw2BFyFIx6K0klpn+IdGMZLDAWgwawhvnHlTdDcmD+gY/Ma7eIhBLeu2Mv0P/43vvLECu19PLC79HSRnfjEujAhDBI2fYPmBbleLVgVCaWkFOFWisTR+MbtnQOlE/DQeCbaxe/lF8thq8JQh9kBGWusO5SGh6eZ+WM2TNEizyb0LAtBOKobPYuSwqQvttMooZSGIYvME8wNzJSxwXiDrBZmIeseDHyvwNBO3ipowtCswesmS7qzk3RQdwt5y66FP6QnGVWmTIf06KlTWxNQyBnHvx7Os1fScSr4wsDIUtlzWhE5VPfqxa3ShusJXNgQv4kb4Qv/b+Z7/wAipiJ0SoWQO81I1zhyuspkmYttMqiBBC/4hEHqE8znvDerGn4UbTG6ozXNQBSH5uYtUcSQ37hRU9hiAFFUYYCDq5bZ05D6qh9SIwB7UMo/PrdWsw9LDfQXDK9grV0qCJiA==
// X-Gorgon: 840460622001afdbe42ea1a0701ffaccbd5216c2c7fd39675cc1
// X-Khronos: 1721825902
// X-Ladon: EPMeIFJKyE0rJ5xHpMAuMBkkFggrzlgLVxyjzj0Pnxvxz5qA
