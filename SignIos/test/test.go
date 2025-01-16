package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// 读取二进制文件到字节切片
	data, err := ioutil.ReadFile("Desktop")
	if err != nil {
		log.Fatal(err)
	}
	byteData := base64.RawStdEncoding.EncodeToString(data)
	// fmt.Println(hexString)
	// bytesData, err := hex.DecodeString(hexString)
	// if err != nil {
	// 	fmt.Println("解码传入16进制字符串失败!")
	// }
	// var token_req pb.TokenReq
	// resBytes, err := common.TokenDecrypt(data)
	// if err != nil {
	// 	fmt.Println("解密出错!")
	// }
	// var token_req_4 pb.TokenReq_4
	// err = proto.Unmarshal(resBytes, &token_req_4)
	// if err != nil {
	// 	fmt.Println("解密出错!")
	// }
	fmt.Println(byteData)

}
