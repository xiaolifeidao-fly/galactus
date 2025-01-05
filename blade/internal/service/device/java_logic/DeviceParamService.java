package com.blade.service.device;

import com.alibaba.fastjson.JSONObject;
import com.blade.service.device.dto.DeviceDTO;

import java.util.List;

/**
 * @author xiaofeidao
 * @date 2019/5/28
 */
public interface DeviceParamService {

    /**
     * 保存
     * @param deviceDTO
     */
    DeviceDTO save(DeviceDTO deviceDTO);

    /**
     * 查询设备信息
     * @param udId
     * @param openUdId
     * @param serial
     * @return
     */
    DeviceDTO getByUdIdAndOpenUdIdAndSerialNumber(String udId, String openUdId, String serial);

    /**
     * 查询设备信息
     * @param id
     * @return
     */
    DeviceDTO getById(Long id);

    /**
     * 获取今日未使用的设备信息
     * @param num 获取数量
     * @return
     */
    List<DeviceDTO> getByUnTodayUseData(int num);

    Long minIdByStartIndex(Long startIndex);

    List<DeviceDTO> getActiveByRangeId(Long startIndex, Long endIndex);

    JSONObject registerOldDevice(String registerUrl, JSONObject deviceParams, String ip);

    String registerXLog(DeviceDTO deviceDTO, String ip, String ext);

    List<DeviceDTO> findByVersionCode(String s);

    void deleteById(Long oriDeviceId);

    List<DeviceDTO> findByVersionCode(String s, int startId, int endId);

}
