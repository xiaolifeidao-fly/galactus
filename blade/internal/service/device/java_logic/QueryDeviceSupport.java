package com.blade.business.query.tools;

import com.blade.common.dto.BaseDTO;
import com.blade.lock.BladeLock;
import com.blade.redis.util.RedisUtil;
import com.blade.service.device.dto.DeviceDTO;
import com.blade.service.dictionary.DictionaryService;
import com.blade.service.dictionary.dto.DictionaryDTO;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.List;
import java.util.concurrent.ConcurrentLinkedQueue;

@Slf4j
public abstract class QueryDeviceSupport<I, O extends BaseDTO> implements InitializingBean {

    private ConcurrentLinkedQueue<O> userEntities = new ConcurrentLinkedQueue<>();

    @Autowired
    private RedisUtil redisUtil;

    @Autowired
    private DictionaryService dictionaryService;

    @Autowired
    private BladeLock bladeLock;

    private ConcurrentLinkedQueue<I> cacheDevices = new ConcurrentLinkedQueue<>();

    @Override
    public void afterPropertiesSet() throws Exception {
        try {
            if (!initFlag()) {
                return;
            }
            Long currentIndex = getCurrentIndex();
            String[] rangeIndex = getRangeIndex();
            List<I> list = getDevices(currentIndex, rangeIndex);
            if (list != null) {
                for (I i : list) {
                    userEntities.add(convert(i));
                }
            }
        } catch (Exception e) {
            log.error("{} init error : ", this.getClass().getSimpleName(), e);
            throw new RuntimeException(e);
        }
    }

    protected boolean initFlag() {
        return true;
    }

    public O get() {
        synchronized (userEntities) {
            O o = null;
            try {
                o = userEntities.poll();
                return o;
            } catch (Exception e) {
                log.error("{} get device error:", this.getClass().getName(), e);
                throw new RuntimeException(e);
            } finally {
                try {
                    if (o != null) {
                        String key = buildKey(o.getId());
                        Long currentFetchNum = redisUtil.incr(key);
                        boolean checkFlag = checkDevice(o, currentFetchNum);
                        if (checkFlag) {
                            redisUtil.expire(key, getExpireTimes());
                            userEntities.add(o);
                        } else {
                            fillDevice(key);
                        }
                    }
                } catch (Exception e) {
                    log.error("finally {} get device error:", this.getClass().getName(), e);
                    throw new RuntimeException(e);
                }
            }
        }
    }

    protected abstract int getExpireTimes();

    private void fillDevice(String key) {
        redisUtil.del(key);
        I input = cacheDevices.poll();
        if (input == null) {
            initCache();
        }
        O output = convert(input);
        userEntities.add(output);
    }

    private void initCache() {
        String deviceLock = buildBladeDeviceKey();
        String key = "DEVICE";
        try {
            bladeLock.lock(deviceLock, key);
            Long currentIndex = getCurrentIndex();
            String[] rangeIndex = getRangeIndex();
            List<I> list = getDevices(currentIndex, rangeIndex);
            if (list != null) {
                cacheDevices.addAll(list);
            }
        } finally {
            bladeLock.unLock(deviceLock, key);
        }
    }

    protected abstract String[] getRangeIndex();

    protected abstract Long getCurrentIndex();

    private List<I> getDevices(Long currentIndex, String[] rangeIndex) {
        Long minIndex = Long.valueOf(rangeIndex[0]);
        Long maxIndex = Long.valueOf(rangeIndex[1]);
        if (currentIndex == null || currentIndex < minIndex) {
            currentIndex = minIndex;
        }
        Long startIndex = currentIndex;
        Long endIndex = startIndex + getDevicePageSize();
        List<I> deviceDTOs = getDevices(startIndex, endIndex);
        if ((deviceDTOs == null || deviceDTOs.size() == 0) && endIndex < maxIndex) {
            startIndex = getMinId(startIndex);
            endIndex = startIndex + getDevicePageSize();
            deviceDTOs = getDevices(startIndex, endIndex);
        }
        if ((deviceDTOs == null || deviceDTOs.size() == 0) && endIndex > maxIndex) {
            startIndex = minIndex;
            endIndex = startIndex + getDevicePageSize();
            deviceDTOs = getDevices(startIndex, endIndex);
        }
        DictionaryDTO dictionaryDTO = dictionaryService.getByCode(getIndexName());
        dictionaryDTO.setValue(endIndex.toString());
        dictionaryService.save(dictionaryDTO);
        return deviceDTOs;
    }

    protected abstract String getIndexName();

    protected abstract Long getMinId(Long startIndex);

    protected abstract List<I> getDevices(Long startIndex, Long endIndex);

    protected abstract Long getDevicePageSize();

    public abstract O convert(I input);

    private String buildBladeDeviceKey() {
        return "BLADE_DEVICE_" + this.getClass().getSimpleName();
    }

    private String buildKey(Long id) {
        return "DEVICE_KEY_" + this.getClass().getSimpleName() + "_" + id;
    }

    protected boolean checkDevice(O o, Long currentFetchNum) {
        Long fetchNum = getFetchNum();
        return currentFetchNum <= fetchNum;
    }

    protected abstract Long getFetchNum();

    public abstract void save(DeviceDTO deviceDTO);
}
