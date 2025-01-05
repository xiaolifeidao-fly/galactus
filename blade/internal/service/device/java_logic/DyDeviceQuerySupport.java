package com.blade.business.query.douyin.support.device;

import com.blade.business.query.douyin.entity.QueryUserEntity;
import com.blade.business.query.douyin.plugin.device.DyDeviceService;
import com.blade.business.query.tools.QueryDeviceSupport;
import com.blade.service.device.DeviceParamService;
import com.blade.service.device.dto.DeviceDTO;
import com.blade.service.dictionary.DictionaryService;
import com.blade.service.dictionary.dto.DictionaryConfig;
import com.blade.service.dictionary.dto.DictionaryDTO;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.concurrent.atomic.AtomicLong;

@Component
public class DyDeviceQuerySupport extends QueryDeviceSupport<DeviceDTO, QueryUserEntity> {

    @Autowired
    private DictionaryService dictionaryService;

    @Value("${dy.device.expire.times:30}")
    private int deviceExpireTimes;

    @Value("${dy.device.page.size:5000}")
    private long devicePageSize;

    @Autowired
    private DeviceParamService deviceParamService;

    @Value("${fill.version.flag:false}")
    private boolean fillVersionFlag;

    @Autowired
    private DyDeviceService dyDeviceService;

    @Value("${dy.device.flag:false}")
    private boolean initFlag;

    @Value("${fill.high.version.flag:false}")
    private boolean fillHighVersionFlag;

    @Value("${fill.device.flag:false}")
    private boolean fillDeviceFlag;

    @Override
    protected int getExpireTimes() {
        return deviceExpireTimes * 60;
    }

    @Override
    protected boolean initFlag() {
        return initFlag;
    }

    @Override
    protected String[] getRangeIndex() {
        DictionaryDTO dictionaryDTO = dictionaryService.getByCode(DictionaryConfig.DEVICE_PARAM_ID_RANGE.name());
        return dictionaryDTO.getValue().split(",");
    }

    @Override
    protected Long getCurrentIndex() {
        DictionaryDTO dictionaryDTO = dictionaryService.getByCode(getIndexName());
        String value = dictionaryDTO.getValue();
        if (StringUtils.isBlank(value)) {
            return null;
        }
        return Long.valueOf(value);
    }

    @Override
    protected String getIndexName() {
        return DictionaryConfig.DEVICE_PARAM_CURRENT_INDEX.name();
    }

    @Override
    protected Long getMinId(Long startIndex) {
        return deviceParamService.minIdByStartIndex(startIndex);
    }

    @Override
    protected List<DeviceDTO> getDevices(Long startIndex, Long endIndex) {
        return deviceParamService.getActiveByRangeId(startIndex, endIndex);
    }

    @Override
    protected Long getDevicePageSize() {
        return devicePageSize;
    }

    @Override
    public QueryUserEntity convert(DeviceDTO deviceDTO) {
        QueryUserEntity queryUserEntity = new QueryUserEntity();
        if (fillVersionFlag) {
            dyDeviceService.fill03VersionCode(deviceDTO);
        }
        deviceDTO.setOs("android");
        deviceDTO.setAppName("aweme");
        deviceDTO.setAccess("wifi");
        deviceDTO.setCookie("");
        if (fillVersionFlag) {
            dyDeviceService.fill03VersionCode(deviceDTO);
        }
        if (fillHighVersionFlag) {
            dyDeviceService.fillVersionCode(deviceDTO);
        }
        if (fillDeviceFlag) {
            dyDeviceService.filDevice(deviceDTO);
        }
        queryUserEntity.setDeviceDTO(deviceDTO);
        queryUserEntity.setPreQueryResult(true);
        queryUserEntity.setLianXuQueryNum(new AtomicLong());
        queryUserEntity.setQueryErrorNum(new AtomicLong());
        queryUserEntity.setId(deviceDTO.getId());
        return queryUserEntity;
    }

    @Override
    protected Long getFetchNum() {
        DictionaryDTO dictionaryDTO = dictionaryService.getByCode(DictionaryConfig.DEVICE_PARAM_FETCH_NUM.name());
        return Long.valueOf(dictionaryDTO.getValue());
    }

    @Override
    public void save(DeviceDTO deviceDTO) {
        deviceParamService.save(deviceDTO);
    }
}
