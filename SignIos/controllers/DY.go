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

	"github.com/astaxie/beego/context"
	"google.golang.org/protobuf/proto"
)

type DyController struct {
	BaseController
}

type ReqStruct struct {
	Url                 string `json:"url"`
	Method              string `json:"method"`
	Body                string `json:"body"`
	DeviceId            string `json:"device_id"`
	Lanusk              string `json:"lanusk"`
	License             string `json:"license"`
	GorgonSdkVersionStr string `json:"gorgon_sdk_version_str"`
	Appid               string `json:"appid"`
	AppVersion          string `json:"app_version"`
	ArgusSdkVersionStr  string `json:"argus_sdk_version_str"`
	ArgusSdkVersion     uint64 `json:"argus_sdk_version"`
	ArgusDevicesToken   string `json:"argus_devices_token"`
}
type DyResponse struct {
	XArgus   string `json:"x_argus"`
	XKhronos string `json:"x_khronos"`
	XGorgon  string `json:"x_gorgon"`
	XLadon   string `json:"x_ladon"`
	XssStub  string `json:"x_ss_stub"`
}

var RequestNum = 0

func (c *DyController) GetParams() {
	RequestNum++
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
	ip := getClientIP(c.Ctx)
	fmt.Println("接受到ip:", ip, "的请求,链接为:", v.Url, "请求方式为:", v.Method, ",当前请求次数:", RequestNum)
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
	fmt.Println("x-ss-stub:", X_SS_Stub)
	fmt.Println("x-gorgon->query_hash:", hex.EncodeToString(queryMd5[0:4]))
	fmt.Println("x-gorgon->body_hash:", hex.EncodeToString(bodyMd5[0:4]))
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
	fmt.Println("x-argus->query_hash:", hex.EncodeToString(query_hash_sm3[0:6]))
	fmt.Println("x-argus->body_hash:", hex.EncodeToString(X_SS_Stub_Sm3[0:6]))
	fmt.Println("x-argus->query_sm3:", hex.EncodeToString(querySm3))
	var actionInfo = &pb.XArgusV15{
		Unknown_1: 242,
		Unknown_2: 12,
		Unknown_3: 1388734,
	}
	if v.License == "" {
		v.License = "712198790"
	}
	if v.GorgonSdkVersionStr == "" {
		v.GorgonSdkVersionStr = "01040304"
	}
	if v.Appid == "" {
		v.Appid = "1128"
	}
	if v.AppVersion == "" {
		v.AppVersion = "21.3.0"
	}
	if v.ArgusSdkVersionStr == "" {
		v.ArgusSdkVersionStr = "v04.03.07-ml-iOS"
	}
	if v.ArgusSdkVersion == 0 {
		v.ArgusSdkVersion = 0x4030701 << 1
	}
	err, HeaderInfo := IosMakeHeader(uint64(time_10), queryMd5[0:4], bodyMd5[0:4], query_hash_sm3[0:6], X_SS_Stub_Sm3[0:6], querySm3, ClientKeyMd5, v.DeviceId, v.ArgusDevicesToken, actionInfo, v.License, v.GorgonSdkVersionStr, v.Appid, v.AppVersion, v.ArgusSdkVersionStr, v.ArgusSdkVersion)
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
	fmt.Println("算法成功!x-argus:", resinfo.XArgus, ",x-khronos:", resinfo.XKhronos, ",x-gorgon:", resinfo.XGorgon, "x-ladon:", resinfo.XLadon)
	c.responseSuccess("ok!", resinfo)
}
func getClientIP(ctx *context.Context) string {
	ip := ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = ctx.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		return "127.0.0.1"
	}
	return ip
}

type MyHeader struct {
	XArgus   string `json:"x_argus"`
	XKhronos string `json:"x_khronos"`
	XGorgon  string `json:"x_gorgon"`
	XLadon   string `json:"x_ladon"`
}

func IosMakeHeader(ts uint64, gorgon_querHash, gorgon_bodyHash, argus_queryHash, argus_bodyHash, argus_querySm3 []byte, clientKey, devices_id, device_token string, actionInfo *pb.XArgusV15, license, sdkVersionStr, appid string, argus_app_version, argus_sdk_version_str string, argus_sdk_version uint64) (err error, res *MyHeader) {
	var header = &MyHeader{}
	time_10_str := strconv.FormatInt(int64(ts), 10)
	header.XKhronos = time_10_str
	err, header.XArgus = IosGetXArgus(ts, argus_queryHash, argus_bodyHash, argus_querySm3, clientKey, devices_id, device_token, actionInfo, appid, license, argus_app_version, argus_sdk_version_str, argus_sdk_version)
	if err != nil {
		return err, nil
	}
	err, header.XGorgon = IosGetXGorgon(int64(ts), gorgon_querHash, gorgon_bodyHash, sdkVersionStr)
	if err != nil {
		return err, nil
	}
	err, header.XLadon = IosGetXLadon(int64(ts), appid, license)
	if err != nil {
		return err, nil
	}
	return nil, header
}
func IosGetXArgus(ts uint64, argus_queryHash, argus_bodyHash, argus_querySm3 []byte, clientKey, devices_id, device_token string, actionInfo *pb.XArgusV15, Appid, License, argus_app_version, argus_sdk_version_str string, argus_sdk_version uint64) (err error, res string) {
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
	err, Xargus := ios.EncryptArgus(&p1)
	if err != nil {
		return err, ""
	}
	return nil, Xargus
}
func IosGetXGorgon(ts int64, gorgon_querHash, gorgon_bodyHash []byte, sdkVersionStr string) (err error, res string) {
	err, XGorgon := ios.EncryptGorgon(ts, gorgon_querHash, gorgon_bodyHash, sdkVersionStr)
	if err != nil {
		return err, ""
	}
	return nil, XGorgon
}
func IosGetXLadon(ts int64, appid, license string) (err error, res string) {
	err, XLadon := ios.EncryptLadon(ts, appid, license)
	if err != nil {
		return err, ""
	}
	return nil, XLadon
}

type DecryptInfoReq struct {
	Gorgon string `json:"gorgon"`
	Argus  string `json:"argus"`
	Ladon  string `json:"ladon"`
}
type DecryptInfoRes struct {
	DeviceId            string `json:"device_id"`
	License             string `json:"license"`
	GorgonSdkVersionStr string `json:"gorgon_sdk_version_str"`
	Appid               string `json:"appid"`
	AppVersion          string `json:"app_version"`
	ArgusSdkVersionStr  string `json:"argus_sdk_version_str"`
	ArgusSdkVersion     uint64 `json:"argus_sdk_version"`
	ArgusDevicesToken   string `json:"argus_devices_token"`
}

func (c *DyController) DecryptInfo() {
	var v DecryptInfoReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Gorgon == "" || v.Argus == "" || v.Ladon == "" {
		c.responseFail("请传入必要的参数!", nil)
	}
	err, ArgusInfo := DecryptIosXargus(v.Argus)
	if err != nil {
		c.responseFail("解密[x-argus]失败!!", nil)
	}
	err, LadonInfo := DecryptIosXladon(v.Ladon)
	if err != nil {
		c.responseFail("解密[x-ladon]失败!!", nil)
	}
	LadonInfoList := strings.Split(LadonInfo, "-")
	if len(LadonInfoList) < 2 {
		c.responseFail("解密[x-ladon]失败!!", nil)
	}
	err, GorgonInfo := DecryptIosGorgon(v.Gorgon)
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
func DecryptIosXargus(argus string) (err error, resInfo *pb.XArgus) {
	SignKey := utils.Ios_DouYinSignKey
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
func DecryptIosXladon(ladon string) (err error, resInfo string) {
	base64Str := ladon
	unbase, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return errors.New("解密base失败"), ""
	}
	key4 := unbase[0:4]
	undata := unbase[4:]
	aid := "1128"
	byte_aid := []byte(aid)
	md5_data_byte := utils.BytesCombine(key4, byte_aid)
	md5_data := utils.BytesToMd5Byte(md5_data_byte)
	md5_hex_str := hex.EncodeToString(md5_data)
	byte_hex_md5_ := []byte(md5_hex_str)
	aesBuf := byte_hex_md5_[16:]
	ladonOut := utils.IosDecryptLadon(undata, aesBuf, md5_data)
	return nil, string(ladonOut)
}
func DecryptIosGorgon(gorgon string) (err error, resInfo *utils.XGorGonInfo) {
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
