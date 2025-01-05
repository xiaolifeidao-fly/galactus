package com.blade.service.device;

import com.blade.service.device.dto.DeviceDTO;

import java.util.List;

public interface MineDeviceService {

    DeviceDTO save(DeviceDTO deviceDTO);

    DeviceDTO getById(Long id);

    List<DeviceDTO> getActiveByRangeId(Long startIndex, Long endIndex);

    Long minIdByStartIndex(Long startIndex);
}
