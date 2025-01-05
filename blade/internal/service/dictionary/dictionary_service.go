package dictionary

import (
	"galactus/blade/internal/service/dictionary/dto"
	"galactus/blade/internal/service/dictionary/repository"
	"galactus/common/middleware/db"
)

type DictionaryService struct {
	repo *repository.DictionaryRepository
}

func NewDictionaryService() *DictionaryService {
	return &DictionaryService{
		repo: db.GetRepository[repository.DictionaryRepository](),
	}
}

// GetByCode 根据编码获取字典
func (s *DictionaryService) GetByCode(code string) (*dto.DictionaryDTO, error) {
	dict, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	return &dto.DictionaryDTO{
		Code:        dict.Code,
		Value:       dict.Value,
		Description: dict.Description,
		Type:        dict.Type,
	}, nil
}

// GetByType 根据类型获取字典列表
func (s *DictionaryService) GetByType(typeStr string) ([]dto.DictionaryDTO, error) {
	dicts, err := s.repo.GetByType(typeStr)
	if err != nil {
		return nil, err
	}
	result := make([]dto.DictionaryDTO, len(dicts))
	for i, dict := range dicts {
		result[i] = dto.DictionaryDTO{
			Code:        dict.Code,
			Value:       dict.Value,
			Description: dict.Description,
			Type:        dict.Type,
		}
	}
	return result, nil
}

// Save 保存字典
func (s *DictionaryService) Save(dictDTO *dto.DictionaryDTO) error {
	dict := &repository.Dictionary{
		Code:        dictDTO.Code,
		Value:       dictDTO.Value,
		Description: dictDTO.Description,
		Type:        dictDTO.Type,
	}
	_, err := s.repo.SaveOrUpdate(dict)
	return err
}

// Delete 删除字典
func (s *DictionaryService) Delete(id int64) error {
	return s.repo.Delete(uint(id))
}
