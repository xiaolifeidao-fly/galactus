package controllers

import (
	"SignIos/ios"
	pb "SignIos/proto"
	"SignIos/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

type DuoShanController struct {
	BaseController
}

func (c *DuoShanController) GetParams() {
	var v ReqStruct
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Url == "" || v.DeviceId == "" || v.Method == "" || v.Lanusk == "" {
		c.responseFail("请传入必要的参数!", nil)
	}
	if v.Method == "post" && v.Body == "" {
		c.responseFail("传入请求格式为get,未传入body!", nil)
	}
	time_10 := time.Now().Unix()
	url := v.Url
	queryParams := strings.Split(url, "?")[1]
	var bodyMd5 []byte
	var X_SS_Stub string
	if v.Method == "post" {
		bodyMd5 = utils.BytesToMd5Byte([]byte(v.Body))
		X_SS_Stub = strings.ToUpper(hex.EncodeToString(bodyMd5))
	} else {
		bodyMd5 = []byte{0x0, 0x0, 0x0, 0x0}
		X_SS_Stub = "00000000000000000000000000000000"
	}
	queryMd5 := utils.BytesToMd5Byte([]byte(queryParams))
	bodyMd5ByteDecode, _ := hex.DecodeString(X_SS_Stub)
	X_SS_Stub_Sm3 := utils.Sm3Hash(bodyMd5ByteDecode)
	query_hash_sm3 := utils.Sm3Hash([]byte(queryParams))
	var querySm3 []byte
	if v.Method == "post" {
		querySm3 = utils.GetQuerySm3([]byte(queryParams), bodyMd5)
	} else {
		querySm3 = utils.GetQuerySm3([]byte(queryParams), nil)
	}
	ClientKeyMd5 := utils.Md5(v.Lanusk)
	var actionInfo = &pb.XArgusV15{
		Unknown_1: 242,
		Unknown_2: 12,
		Unknown_3: 1388734,
	}
	if v.License == "" {
		v.License = "1341894232"
	}
	if v.GorgonSdkVersionStr == "" {
		v.GorgonSdkVersionStr = "01070304"
	}
	if v.Appid == "" {
		v.Appid = "1349"
	}
	if v.AppVersion == "" {
		v.AppVersion = "20.0.1"
	}
	if v.ArgusSdkVersionStr == "" {
		v.ArgusSdkVersionStr = "v04.03.07-ml-iOS"
	}
	if v.ArgusSdkVersion == 0 {
		v.ArgusSdkVersion = 0x4030701 << 1
	}
	err, HeaderInfo := DuoShanIosMakeHeader(uint64(time_10), queryMd5[0:4], bodyMd5[0:4], query_hash_sm3[0:6], X_SS_Stub_Sm3[0:6], querySm3, ClientKeyMd5, v.DeviceId, v.ArgusDevicesToken, actionInfo, v.License, v.GorgonSdkVersionStr, v.Appid, v.AppVersion, v.ArgusSdkVersionStr, v.ArgusSdkVersion)
	if err != nil {
		fmt.Println("获取请求头失败!error:", err)
		c.responseFail("出错了!", nil)
	}
	var resinfo = &DyResponse{}
	resinfo.XArgus = HeaderInfo.XArgus
	resinfo.XKhronos = HeaderInfo.XKhronos
	resinfo.XGorgon = HeaderInfo.XGorgon
	resinfo.XLadon = HeaderInfo.XLadon
	if v.Method == "post" {
		resinfo.XssStub = X_SS_Stub
	}
	c.responseSuccess("ok!", resinfo)
}

func DuoShanIosMakeHeader(ts uint64, gorgon_querHash, gorgon_bodyHash, argus_queryHash, argus_bodyHash, argus_querySm3 []byte, clientKey, devices_id, device_token string, actionInfo *pb.XArgusV15, license, sdkVersionStr, appid string, argus_app_version, argus_sdk_version_str string, argus_sdk_version uint64) (err error, res *MyHeader) {
	var header = &MyHeader{}
	time_10_str := strconv.FormatInt(int64(ts), 10)
	header.XKhronos = time_10_str
	err, header.XArgus = DuoShanIosGetXArgus(ts, argus_queryHash, argus_bodyHash, argus_querySm3, clientKey, devices_id, device_token, actionInfo, appid, license, argus_app_version, argus_sdk_version_str, argus_sdk_version)
	if err != nil {
		return err, nil
	}
	err, header.XGorgon = DuoShanIosGetXGorgon(int64(ts), gorgon_querHash, gorgon_bodyHash, sdkVersionStr)
	if err != nil {
		return err, nil
	}
	err, header.XLadon = DuoShanIosGetXLadon(int64(ts), appid, license)
	if err != nil {
		return err, nil
	}
	return nil, header
}
func DuoShanIosGetXArgus(ts uint64, argus_queryHash, argus_bodyHash, argus_querySm3 []byte, clientKey, devices_id, device_token string, actionInfo *pb.XArgusV15, Appid, License, argus_app_version, argus_sdk_version_str string, argus_sdk_version uint64) (err error, res string) {
	envCode, _ := hex.DecodeString("0000000000000000")
	Psk := "1"
	if device_token != "" {
		envCode, _ = hex.DecodeString("1001000000000000")
		//Psk = "1"
	}
	if actionInfo == nil {
		actionInfo = &pb.XArgusV15{
			Unknown_1: 4,
			Unknown_2: 1388734,
			Unknown_3: 1388734,
		}
	}
	p1 := pb.XArgus{
		Magic:         1077940818,
		Version:       2,
		Random:        utils.RandUint64(),
		Appid:         Appid,
		License:       License,
		AppVersion:    argus_app_version,
		SdkVersionStr: argus_sdk_version_str,
		SdkVersion:    argus_sdk_version,
		Env:           envCode,
		Platform:      1,
		CreateTime:    ts << 1,
		BodyHsh:       argus_bodyHash,
		QueryHash:     argus_queryHash,
		ActionInfo:    actionInfo,
		IsLicense:     ts << 1,
		Psk:           Psk,
		CallType:      738,
	}
	if devices_id != "" {
		p1.DevicesId = devices_id
	}
	if clientKey != "" {
		// 11c3cfc17d9c29800f9e8d057f293328
		ClientKeyMd5, _ := hex.DecodeString(clientKey)
		p1.ClientKeyMd5 = ClientKeyMd5
	}
	if argus_querySm3 != nil {
		p1.QuerySm3 = argus_querySm3
	}
	if device_token != "" {
		p1.DevicesToken = device_token
	}
	err, Xargus := ios.DuoShanEncryptArgus(&p1)
	if err != nil {
		return err, ""
	}
	return nil, Xargus
}
func DuoShanIosGetXGorgon(ts int64, gorgon_querHash, gorgon_bodyHash []byte, sdkVersionStr string) (err error, res string) {
	err, XGorgon := ios.EncryptGorgon(ts, gorgon_querHash, gorgon_bodyHash, sdkVersionStr)
	if err != nil {
		return err, ""
	}
	return nil, XGorgon
}
func DuoShanIosGetXLadon(ts int64, appid, license string) (err error, res string) {
	err, XLadon := ios.EncryptLadon(ts, appid, license)
	if err != nil {
		return err, ""
	}
	return nil, XLadon
}

func (c *DuoShanController) DecryptInfo() {
	var v DecryptInfoReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Gorgon == "" || v.Argus == "" || v.Ladon == "" {
		c.responseFail("请传入必要的参数!", nil)
	}
	err, ArgusInfo := DuoShanDecryptIosXargus(v.Argus)
	if err != nil {
		c.responseFail("解密[x-argus]失败!!", nil)
	}
	err, LadonInfo := DuoShanDecryptIosXladon(v.Ladon)
	if err != nil {
		c.responseFail("解密[x-ladon]失败!!", nil)
	}
	LadonInfoList := strings.Split(LadonInfo, "-")
	if len(LadonInfoList) < 2 {
		c.responseFail("解密[x-ladon]失败!!", nil)
	}
	err, GorgonInfo := DuoShanDecryptIosGorgon(v.Gorgon)
	if err != nil {
		c.responseFail("解密[x-gorgon]失败!!", nil)
	}
	if ArgusInfo.Appid != LadonInfoList[2] {
		c.responseFail("解密[x-ladon]失败!!", nil)
	}
	if int64(ArgusInfo.CreateTime>>1) != GorgonInfo.XKhronos {
		c.responseFail("解密失败!!", nil)
	}
	var resInfo = &DecryptInfoRes{}
	resInfo.DeviceId = ArgusInfo.DevicesId
	resInfo.License = ArgusInfo.License
	resInfo.GorgonSdkVersionStr = GorgonInfo.SdkVersion
	resInfo.Appid = ArgusInfo.Appid
	resInfo.AppVersion = ArgusInfo.AppVersion
	resInfo.ArgusSdkVersionStr = ArgusInfo.SdkVersionStr
	resInfo.ArgusSdkVersion = ArgusInfo.SdkVersion
	resInfo.ArgusDevicesToken = ArgusInfo.DevicesToken
	c.responseSuccess("解密成功!", resInfo)
}
func DuoShanDecryptIosXargus(argus string) (err error, resInfo *pb.XArgus) {
	SignKey := utils.Ios_DouYinDuoShanSignKey
	key_bytes, err := base64.StdEncoding.DecodeString(SignKey)
	if err != nil {
		return errors.New("解密key失败!"), nil
	}
	aesKey := utils.BytesToMd5Byte(key_bytes[0:16])
	aesIv := utils.BytesToMd5Byte(key_bytes[16:])
	input := argus
	inputData, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return errors.New("base解密data失败!"), nil
	}
	aesData := utils.AES_CBC_Decrypt(inputData[2:], aesKey, aesIv)
	key_4 := inputData[0:2]
	key_4 = utils.BytesCombine(key_4, aesData[len(aesData)-2:])
	key_8 := utils.Calc_key_8(key_4)
	sm3_buff := utils.BytesCombine(key_bytes, key_4, key_bytes)
	sm3_hash := utils.Sm3Hash(sm3_buff)
	deXorInfo := utils.DeXor(aesData[9:len(aesData)-2], key_8)
	decryptBufData := deXorInfo[8:]
	decryptBuf := utils.IosDecryptArgusBuf(decryptBufData)
	decryptDta := utils.IosDecryptArgusData(decryptBuf, sm3_hash)
	decryptDta = utils.ClearPadding(decryptDta)
	var xargus pb.XArgus
	err = proto.Unmarshal(decryptDta, &xargus)
	if err != nil {
		return errors.New("解密argus失败!!"), nil
	}
	return nil, &xargus

}
func DuoShanDecryptIosXladon(ladon string) (err error, resInfo string) {
	base64Str := ladon
	unbase, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return errors.New("解密base失败"), ""
	}
	key4 := unbase[0:4]
	undata := unbase[4:]
	aid := "1349"
	byte_aid := []byte(aid)
	md5_data_byte := utils.BytesCombine(key4, byte_aid)
	md5_data := utils.BytesToMd5Byte(md5_data_byte)
	md5_hex_str := hex.EncodeToString(md5_data)
	byte_hex_md5_ := []byte(md5_hex_str)
	aesBuf := byte_hex_md5_[16:]
	ladonOut := utils.IosDecryptLadon(undata, aesBuf, md5_data)
	return nil, string(ladonOut)
}
func DuoShanDecryptIosGorgon(gorgon string) (err error, resInfo *utils.XGorGonInfo) {
	GorgonKey := utils.Ios_DouYinGorgonkey
	hexInputStr := gorgon
	byteInputData, err := hex.DecodeString(hexInputStr)
	if err != nil {
		return errors.New("hex转码失败"), nil
	}
	key8 := byteInputData[0:6]
	waitDecryptInput := byteInputData[6:]
	GorgonKey[1] = key8[4]
	GorgonKey[3] = key8[3]
	GorgonKey[6] = key8[5]
	GorgonKey[7] = key8[2]
	withKeyData := utils.DecryptXGorgonData(waitDecryptInput)
	fullKey := utils.XgprgonInitKey(GorgonKey)
	lastData := utils.XGorgonDataEncryptWithKey(withKeyData, fullKey)
	var XGorgonInfo = utils.XGorGonInfo{}
	XGorgonInfo.QueryMd5 = hex.EncodeToString(lastData[0:4])
	XGorgonInfo.BodyMd5OrStub = hex.EncodeToString(lastData[4:8])
	XGorgonInfo.Unused = hex.EncodeToString(lastData[8:12])
	XGorgonInfo.SdkVersion = hex.EncodeToString(lastData[12:16])
	XGorgonInfo.XKhronos = utils.HexToInt64(hex.EncodeToString(lastData[16:]))
	return nil, &XGorgonInfo
}
