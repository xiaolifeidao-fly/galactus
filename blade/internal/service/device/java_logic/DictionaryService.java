package com.blade.service.dictionary;

import com.blade.service.dictionary.dto.DictionaryDTO;

import java.util.Dictionary;
import java.util.List;

/**
 * @author xiaofeidao
 * @date 2019/5/29
 */
public interface DictionaryService {

    /**
     * 保存
     * @param dictionaryDTO
     */
    void save(DictionaryDTO dictionaryDTO);

    /**
     * @param code
     * @return
     */
    DictionaryDTO getByCode(String code);


    List<DictionaryDTO> getByType(String name);
}
