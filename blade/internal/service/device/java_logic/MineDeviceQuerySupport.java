package com.blade.business.query.douyin.support.device;

import com.blade.business.query.douyin.entity.QueryUserEntity;
import com.blade.business.query.douyin.plugin.device.DyDeviceService;
import com.blade.business.query.tools.QueryDeviceSupport;
import com.blade.common.utils.DeviceUtil;
import com.blade.service.device.MineDeviceService;
import com.blade.service.device.dto.DeviceDTO;
import com.blade.service.dictionary.DictionaryService;
import com.blade.service.dictionary.dto.DictionaryConfig;
import com.blade.service.dictionary.dto.DictionaryDTO;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class MineDeviceQuerySupport extends QueryDeviceSupport<DeviceDTO, QueryUserEntity> {

    @Autowired
    private DictionaryService dictionaryService;

    @Value("${dy.mine.device.expire.times:30}")
    private int deviceExpireTimes;

    @Value("${dy.mine.device.page.size:5000}")
    private long devicePageSize;

    @Autowired
    private MineDeviceService mineDeviceService;

    @Autowired
    private DyDeviceService dyDeviceService;

    @Value("${dy.xiami.device.flag:false}")
    private boolean initFlag;

    @Value("${fill.version.flag:false}")
    private boolean fillVersionFlag;

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
        DictionaryDTO dictionaryDTO = dictionaryService.getByCode(DictionaryConfig.DEVICE_MINE_ID_RANGE.name());
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
        return DictionaryConfig.DEVICE_MINE_CURRENT_INDEX.name();
    }

    @Override
    protected Long getMinId(Long startIndex) {
        return mineDeviceService.minIdByStartIndex(startIndex);
    }

    @Override
    protected List<DeviceDTO> getDevices(Long startIndex, Long endIndex) {
        return mineDeviceService.getActiveByRangeId(startIndex, endIndex);
    }

    @Override
    protected Long getDevicePageSize() {
        return devicePageSize;
    }

    @Override
    public QueryUserEntity convert(DeviceDTO deviceDTO) {
        QueryUserEntity queryUserEntity = new QueryUserEntity();
        BeanUtils.copyProperties(deviceDTO, deviceDTO);
        deviceDTO.setChannel(DeviceUtil.getChannel());
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
        deviceDTO.setOsApi(DeviceUtil.getOsApi());
        deviceDTO.setOsVersion(DeviceUtil.getVersion(deviceDTO.getOsApi()));
        queryUserEntity.setDeviceDTO(deviceDTO);
        queryUserEntity.setId(deviceDTO.getId());
        return queryUserEntity;
    }

    @Override
    protected Long getFetchNum() {
        DictionaryDTO dictionaryDTO = dictionaryService.getByCode(DictionaryConfig.DEVICE_MINE_FETCH_NUM.name());
        return Long.valueOf(dictionaryDTO.getValue());
    }

    @Override
    public void save(DeviceDTO deviceDTO) {
        mineDeviceService.save(deviceDTO);
    }
}
