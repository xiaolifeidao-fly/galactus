package controllers

import (
	"SignIos/common"
	pb "SignIos/proto"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

type MssdkController struct {
	BaseController
}

type MssdkReq struct {
	Hex string `json:"hex"`
}

// 抖音
func (c *MssdkController) DecodeTokenReqData() {
	var v MssdkReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Hex == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	bytesData, err := hex.DecodeString(v.Hex)
	if err != nil {
		c.responseFail("解码传入16进制字符串失败!", nil)
	}
	var token_req pb.TokenReq
	err = proto.Unmarshal(bytesData, &token_req)
	if err != nil {
		c.responseFail("解码protobuf失败!", nil)
		return
	}
	waitDecodeData := token_req.P_4
	resBytes, err := common.TokenDecrypt(waitDecodeData)
	if err != nil {
		c.responseFail("解密出错!", nil)
		return
	}
	var token_req_4 pb.TokenReq_4
	err = proto.Unmarshal(resBytes, &token_req_4)
	if err != nil {
		c.responseFail("解密出错!", nil)
		return
	}
	var res_token_req = &pb.TokenReq_Res{}
	res_token_req.P_1 = token_req.P_1
	res_token_req.P_2 = token_req.P_2
	res_token_req.P_3 = token_req.P_3
	res_token_req.P_4 = &token_req_4
	res_token_req.P_5 = token_req.P_5
	c.responseSuccess("解密成功!", res_token_req)
}
func (c *MssdkController) DecodeTokenResData() {
	var v MssdkReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Hex == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	bytesData, err := hex.DecodeString(v.Hex)
	if err != nil {
		c.responseFail("解码传入16进制字符串失败!", nil)
	}
	var token_res pb.TokenRes
	err = proto.Unmarshal(bytesData, &token_res)
	if err != nil {
		fmt.Println("解码protobuf失败!", err)
		return
	}
	waitDecodeData := token_res.P_6
	resBytes, err := common.TokenDecrypt(waitDecodeData)
	if err != nil {
		fmt.Println("解密出错!", err)
		return
	}
	var token_req_6 pb.TokenRes_6
	err = proto.Unmarshal(resBytes, &token_req_6)
	if err != nil {
		fmt.Println("解码protobuf失败!", err)
		return
	}
	var res_token_res = &pb.TokenRes_Res{}
	res_token_res.P_1 = token_res.P_1
	res_token_res.P_2 = token_res.P_2
	res_token_res.P_5 = token_res.P_5
	res_token_res.P_6 = &token_req_6
	c.responseSuccess("解密成功!", res_token_res)
}

// 极速版
func (c *MssdkController) JiSuDecodeTokenReqData() {
	var v MssdkReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Hex == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	bytesData, err := hex.DecodeString(v.Hex)
	if err != nil {
		c.responseFail("解码传入16进制字符串失败!", nil)
	}
	var token_req pb.TokenReq
	err = proto.Unmarshal(bytesData, &token_req)
	if err != nil {
		c.responseFail("解码protobuf失败!", nil)
		return
	}
	waitDecodeData := token_req.P_4
	resBytes, err := common.TokenDecryptByJiSu(waitDecodeData)
	if err != nil {
		c.responseFail("解密出错!", nil)
		return
	}
	var token_req_4 pb.TokenReq_4
	err = proto.Unmarshal(resBytes, &token_req_4)
	if err != nil {
		c.responseFail("解密出错!", nil)
		return
	}
	var res_token_req = &pb.TokenReq_Res{}
	res_token_req.P_1 = token_req.P_1
	res_token_req.P_2 = token_req.P_2
	res_token_req.P_3 = token_req.P_3
	res_token_req.P_4 = &token_req_4
	res_token_req.P_5 = token_req.P_5
	c.responseSuccess("解密成功!", res_token_req)
}
func (c *MssdkController) JiSuDecodeTokenResData() {
	var v MssdkReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Hex == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	bytesData, err := hex.DecodeString(v.Hex)
	if err != nil {
		c.responseFail("解码传入16进制字符串失败!", nil)
	}
	var token_res pb.TokenRes
	err = proto.Unmarshal(bytesData, &token_res)
	if err != nil {
		fmt.Println("解码protobuf失败!", err)
		return
	}
	waitDecodeData := token_res.P_6
	resBytes, err := common.TokenDecryptByJiSu(waitDecodeData)
	if err != nil {
		fmt.Println("解密出错!", err)
		return
	}
	var token_req_6 pb.TokenRes_6
	err = proto.Unmarshal(resBytes, &token_req_6)
	if err != nil {
		fmt.Println("解码protobuf失败!", err)
		return
	}
	var res_token_res = &pb.TokenRes_Res{}
	res_token_res.P_1 = token_res.P_1
	res_token_res.P_2 = token_res.P_2
	res_token_res.P_5 = token_res.P_5
	res_token_res.P_6 = &token_req_6
	c.responseSuccess("解密成功!", res_token_res)
}

// 多闪
func (c *MssdkController) DuoShanDecodeTokenReqData() {
	var v MssdkReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Hex == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	bytesData, err := hex.DecodeString(v.Hex)
	if err != nil {
		c.responseFail("解码传入16进制字符串失败!", nil)
	}
	var token_req pb.TokenReq
	err = proto.Unmarshal(bytesData, &token_req)
	if err != nil {
		c.responseFail("解码protobuf失败!", nil)
		return
	}
	waitDecodeData := token_req.P_4
	resBytes, err := common.TokenDecryptByDuoShan(waitDecodeData)
	if err != nil {
		c.responseFail("解密出错!", nil)
		return
	}
	var token_req_4 pb.TokenReq_4
	err = proto.Unmarshal(resBytes, &token_req_4)
	if err != nil {
		c.responseFail("解密出错!", nil)
		return
	}
	var res_token_req = &pb.TokenReq_Res{}
	res_token_req.P_1 = token_req.P_1
	res_token_req.P_2 = token_req.P_2
	res_token_req.P_3 = token_req.P_3
	res_token_req.P_4 = &token_req_4
	res_token_req.P_5 = token_req.P_5
	c.responseSuccess("解密成功!", res_token_req)
}
func (c *MssdkController) DuoShanDecodeTokenResData() {
	var v MssdkReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.responseFail("传入数据有误!", nil)
	}
	if v.Hex == "" {
		c.responseFail("请传入必要参数!", nil)
	}
	bytesData, err := hex.DecodeString(v.Hex)
	if err != nil {
		c.responseFail("解码传入16进制字符串失败!", nil)
	}
	var token_res pb.TokenRes
	err = proto.Unmarshal(bytesData, &token_res)
	if err != nil {
		fmt.Println("解码protobuf失败!", err)
		return
	}
	waitDecodeData := token_res.P_6
	resBytes, err := common.TokenDecryptByDuoShan(waitDecodeData)
	if err != nil {
		fmt.Println("解密出错!", err)
		return
	}
	var token_req_6 pb.TokenRes_6
	err = proto.Unmarshal(resBytes, &token_req_6)
	if err != nil {
		fmt.Println("解码protobuf失败!", err)
		return
	}
	var res_token_res = &pb.TokenRes_Res{}
	res_token_res.P_1 = token_res.P_1
	res_token_res.P_2 = token_res.P_2
	res_token_res.P_5 = token_res.P_5
	res_token_res.P_6 = &token_req_6
	c.responseSuccess("解密成功!", res_token_res)
}
