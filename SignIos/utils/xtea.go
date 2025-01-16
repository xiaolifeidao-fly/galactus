package utils

import (
	"encoding/binary"
	"errors"
)

func encrypt(key []byte, block []byte, rounds uint32) []byte {
	var k [4]uint32
	var i uint32
	end := make([]byte, 8)
	v0 := binary.LittleEndian.Uint32(block[:4])
	v1 := binary.LittleEndian.Uint32(block[4:])

	k[0] = binary.LittleEndian.Uint32(key[:4])
	k[1] = binary.LittleEndian.Uint32(key[4:8])
	k[2] = binary.LittleEndian.Uint32(key[8:12])
	k[3] = binary.LittleEndian.Uint32(key[12:])

	delta := binary.LittleEndian.Uint32([]byte{0xb9, 0x79, 0x37, 0x9e})
	mask := binary.LittleEndian.Uint32([]byte{0xff, 0xff, 0xff, 0xff})

	var sum uint32 = 0

	for i = 0; i < rounds; i++ {
		v0 = (v0 + (((v1<<4 ^ v1>>5) + v1) ^ (sum + k[sum&3]))) & mask
		sum = (sum + delta) & mask
		v1 = (v1 + (((v0<<4 ^ v0>>5) + v0) ^ (sum + k[sum>>11&3]))) & mask
	}

	binary.LittleEndian.PutUint32(end[:4], v0)
	binary.LittleEndian.PutUint32(end[4:], v1)

	return end
}

func decrypt(key []byte, block []byte, rounds uint32) []byte {
	var k [4]uint32
	var i uint32
	end := make([]byte, 8)
	v0 := binary.BigEndian.Uint32(block[:4])
	v1 := binary.BigEndian.Uint32(block[4:])
	k[0] = binary.BigEndian.Uint32(key[:4])
	k[1] = binary.BigEndian.Uint32(key[4:8])
	k[2] = binary.BigEndian.Uint32(key[8:12])
	k[3] = binary.BigEndian.Uint32(key[12:])
	var delta uint32 = 0x9E3779B9
	sum := delta * rounds
	for i = 0; i < rounds; i++ {
		v1 = v1 - (((v0<<4 ^ v0>>5) + v0) ^ (sum + k[sum>>11&3]))
		sum = sum - delta
		v0 = v0 - (((v1<<4 ^ v1>>5) + v1) ^ (sum + k[sum&3]))
	}
	binary.BigEndian.PutUint32(end[:4], v0)
	binary.BigEndian.PutUint32(end[4:], v1)
	return end
}

func MyXteaEncrypt(data, key, iv []byte) (err error, res []byte) {
	tmp := binary.LittleEndian.Uint32(iv)
	N := 8 * (tmp % 5)
	N = 0x20 + N
	length := len(data)
	if length%8 != 0 {
		return errors.New("ciphertext is not a multiple of the block size"), nil
	}
	var index = 0
	for length > 0 {
		temp1 := data[index : index+8]
		temp_res := decrypt(key, temp1, N)
		for i := 0; i < 8; i++ {
			res = append(res, temp_res[i]^iv[i])
		}
		iv = temp1
		index += 8
		length -= 8
	}
	return nil, res
}
