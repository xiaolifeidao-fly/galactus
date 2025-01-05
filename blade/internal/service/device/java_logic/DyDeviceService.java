package com.blade.business.query.douyin.plugin.device;

import com.alibaba.fastjson.JSONObject;
import com.blade.business.query.douyin.plugin.device.config.DeviceConfig;
import com.blade.service.device.dto.DeviceDTO;

/**
 * @author xianglong
 * @date 2019/11/14
 */
public interface DyDeviceService {

    DeviceDTO registerDevice(DeviceDTO deviceDTO, DeviceConfig deviceConfig, String ip);

    void fillVersionCode(DeviceDTO deviceDTO);

    String registerXLog(DeviceDTO deviceDTO, String extParams, String ip);

    JSONObject getTokenId(DeviceDTO deviceDTO, String ip);

    JSONObject logSetting(DeviceDTO deviceDTO, String ip);

    JSONObject logSettingWithFingerprint(DeviceDTO deviceDTO, String ip);

    JSONObject getSignResult();

    void removeSignResult();

    void fill03VersionCode(DeviceDTO deviceDTO);

    void filDevice(DeviceDTO deviceDTO);
}
