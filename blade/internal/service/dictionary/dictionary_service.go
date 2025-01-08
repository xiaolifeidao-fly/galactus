package dictionary

import (
	"errors"
	"galactus/blade/internal/service/dictionary/dto"
	"galactus/blade/internal/service/dictionary/repository"
	"galactus/common/middleware/db"
)

type DictionaryService interface {
	GetByCode(code string) (*dto.DictionaryDTO, error)
	GetByType(typeStr string) ([]dto.DictionaryDTO, error)
	Save(dictDTO *dto.DictionaryDTO) error
	Delete(id int64) error
}

type dictionaryService struct {
	repo *repository.DictionaryRepository
}

func NewDictionaryService() DictionaryService {
	return &dictionaryService{
		repo: db.GetRepository[repository.DictionaryRepository](),
	}
}

// GetByCode 根据编码获取字典
func (s *dictionaryService) GetByCode(code string) (*dto.DictionaryDTO, error) {
	dict, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	if dict == nil {
		return nil, errors.New("not found")
	}
	return db.ToDTO[dto.DictionaryDTO](dict), nil
}

// GetByType 根据类型获取字典列表
func (s *dictionaryService) GetByType(typeStr string) ([]dto.DictionaryDTO, error) {
	dicts, err := s.repo.GetByType(typeStr)
	if err != nil {
		return nil, err
	}
	dtoList := db.ToDTOs[dto.DictionaryDTO](dicts)
	result := make([]dto.DictionaryDTO, len(dtoList))
	for i, dto := range dtoList {
		result[i] = *dto
	}
	return result, nil
}

// Save 保存字典
func (s *dictionaryService) Save(dictDTO *dto.DictionaryDTO) error {
	dict := db.ToPO[repository.Dictionary](dictDTO)
	_, err := s.repo.SaveOrUpdate(dict)
	return err
}

// Delete 删除字典
func (s *dictionaryService) Delete(id int64) error {
	return s.repo.Delete(uint(id))
}
