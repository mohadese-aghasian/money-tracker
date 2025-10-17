package usecase

import (
	"errors"
	"money-tracker/internal/dto"
	"money-tracker/internal/entity"
	"strings"
	"time"
)

type TagUseCase struct {
	Repo entity.TagRepository
}

func NewTagUseCase(repo entity.TagRepository) *TagUseCase {
	return &TagUseCase{Repo: repo}
}

func (uc *TagUseCase) Add(title string, status_id uint) (*entity.Tag, error) {

	if status_id < 0 || status_id == 0 {
		status_id = 1
	}

	existing, err_e := uc.Repo.FindByTitle(title)
	if err_e == nil || existing != nil {
		return nil, errors.New("tag duplicate")
	}

	tag, er := entity.NewTag(title, uint(status_id))
	if er != nil {
		return nil, er
	}

	err := uc.Repo.Insert(tag)
	if err != nil {
		return nil, err
	}

	return tag, err
}
func (uc *TagUseCase) Update(input dto.UpdateTagRequest) (*entity.Tag, error) {
	tag, err := uc.Repo.FindById(input.ID)
	if err != nil {
		return nil, errors.New("tag not found")
	}
	if input.Title != "" {
		tag.Title = input.Title
	}
	if input.StatusID != 0 {
		tag.StatusID = input.StatusID
	}

	tag.UpdatedAt = time.Now()

	return uc.Repo.Update(tag)
}
func (uc *TagUseCase) Remove(id uint) error {
	_, err := uc.Repo.FindById(id)
	if err != nil {
		return errors.New("tag not found")
	}

	return uc.Repo.Delete(id)
}
func (uc *TagUseCase) GetByID(id uint) (*entity.Tag, error) {
	return uc.Repo.FindById(id)
}
func (uc *TagUseCase) Get(input dto.ListTagsInput) ([]entity.Tag, int, error) {
	if input.Start < 0 {
		input.Start = 0
	}
	if input.Limit <= 0 {
		input.Limit = 30
	}
	if input.StatusID < 0 {
		input.StatusID = 0
	}

	if input.Sort == "" {
		input.Sort = "DESC"
	} else {
		if input.Sort != "ASC" && input.Sort != "DESC" {
			input.Sort = "DESC"
		}
	}

	allowedColumns := getModelColumns(entity.Tag{})
	if _, ok := allowedColumns[strings.ToLower(input.OrderBy)]; !ok || input.OrderBy == "" {
		input.OrderBy = "id"
	}

	return uc.Repo.FindAll(input.Start, input.Limit, input.OrderBy, input.Sort, uint(input.ID), input.StatusID, input.Title)
}