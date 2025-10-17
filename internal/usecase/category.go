package usecase

import (
	"errors"
	"fmt"
	"money-tracker/internal/constants"
	"money-tracker/internal/dto"
	"money-tracker/internal/entity"

	"strings"
	"time"
)

type CategoryUseCase struct {
	Repo entity.CategoryRepository
}

func NewCategoryUseCase(repo entity.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{Repo: repo}
}

// /----------------------------------------------------
func (uc *CategoryUseCase) Add(input dto.AddCategoryInput) (*entity.Category, error) {
	// check slug duplication
	existing_category, e_err := uc.Repo.FindBySlug(input.Slug, []uint{constants.ArticleActive})
	if e_err == nil && existing_category != nil {
		return nil, errors.New("slug(title) duplicate")
	}

	// default status
	if input.StatusID == 0 {
		input.StatusID = 1
	}

	// create new category
	category, err := entity.NewCategory(input.Title, input.Slug, input.StatusID, input.Color)
	if err != nil {
		return nil, err
	}
	fmt.Printf("$$$$#-------------%v ", category)

	// category.Slug = input.Slug
	// category.CoverId = coverId

	// handle tag IDs if provided
	// if input.TagIDs != "" {
	// 	idsStr := strings.Split(input.TagIDs, ",")
	// 	for _, idStr := range idsStr {
	// 		idStr = strings.TrimSpace(idStr)
	// 		id, err := strconv.Atoi(idStr)
	// 		if err != nil {
	// 			return nil, errors.New("invalid tag id: " + idStr)
	// 		}

	// 		// Check tag existence
	// 		tag, err := uc.TagRepo.FindById(uint(id))
	// 		if err != nil || tag == nil {
	// 			return nil, errors.New("tag not found: " + strconv.Itoa(int(id)))
	// 		}
	// 	}
	// 	category.TagIDs = input.TagIDs
	// }

	// insert into repo
	res := uc.Repo.Insert(category)
	if res != nil {
		return nil, res
	}

	return category, nil
}

// ----------------------------------------------
func (uc *CategoryUseCase) Get(input dto.ListCategoriesInput) ([]dto.CategoryResponse, int, error) {

	if input.Start < 0 {
		input.Start = 0
	}
	if input.Limit <= 0 {
		input.Limit = 30
	}

	if input.Sort == "" {
		input.Sort = "DESC"
	} else {
		if input.Sort != "ASC" && input.Sort != "DESC" {
			input.Sort = "DESC"
		}
	}

	allowedColumns := getModelColumns(entity.Category{})
	if _, ok := allowedColumns[strings.ToLower(input.OrderBy)]; !ok || input.OrderBy == "" {
		input.OrderBy = "id"
	}

	categories, count, err := uc.Repo.FindAll(dto.CategoryFindAll{
		Start:    input.Start,
		Limit:    input.Limit,
		OrderBy:  input.OrderBy,
		Sort:     input.Sort,
		ID:       uint(input.ID),
		Title:    input.Title,
		StatusID: input.StatusID,
		// TagIDs:   input.TagIds,
		Slug:  input.Slug,
		Color: input.Color,
	})

	var responses []dto.CategoryResponse
	for _, cat := range categories {

		// tagIDs := parseTagIDs(cat.TagIDs)
		// var tags []dto.FetchedTag
		// if len(tagIDs) > 0 {
		// 	pureTags, _ := uc.TagRepo.FindByIDs(tagIDs)

		// 	for _, tag := range pureTags {
		// 		tags = append(tags, dto.FetchedTag{
		// 			ID:        tag.ID,
		// 			Title:     tag.Title,
		// 			StatusID:  tag.StatusID,
		// 			CreatedAt: tag.CreatedAt,
		// 		})
		// 	}
		// }

		// var cover *dto.FetchedCover
		// if cat.Cover != nil && *cat.CoverId > 0 {
		// 	cover = &dto.FetchedCover{
		// 		ID:   cat.Cover.ID,
		// 		Name: cat.Cover.Name,
		// 		Link: cat.Cover.Link,
		// 	}
		// }

		responses = append(responses, dto.CategoryResponse{
			ID:       cat.ID,
			Title:    cat.Title,
			StatusID: cat.StatusID,
			Slug:     cat.Slug,
			// Tags:         tags,
			CreatedAt: cat.CreatedAt,
			// CoverId:      cat.CoverId,
			// Cover:        cover,
		})
	}

	return responses, count, err

}

// /-----------------------------------------------
func (uc *CategoryUseCase) Remove(id uint) error {
	_, err := uc.Repo.FindById(id)
	if err != nil {
		return errors.New("category not found")
	}

	return uc.Repo.Delete(id)
}

// ---------------------------------------------------
func (uc *CategoryUseCase) Update(input dto.UpdateCategoryInput) (*entity.Category, error) {
	category, err := uc.Repo.FindById(input.ID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if input.Title != "" {
		category.Title = input.Title
	}
	if input.StatusID != 0 {
		category.StatusID = input.StatusID
	}

	if input.Slug != "" {
		fmt.Println("____________", input.Slug)
		existing_cat, e_err := uc.Repo.FindBySlug(input.Slug, []uint{constants.StatusActive})
		if e_err == nil && existing_cat != nil {
			if existing_cat.ID != input.ID {
				return nil, errors.New("slug(title) duplicate")
			}
		}
		category.Slug = input.Slug
	}

	// if input.TagIDs != "" {
	// 	idsStr := strings.Split(input.TagIDs, ",")
	// 	for _, idStr := range idsStr {
	// 		idStr = strings.TrimSpace(idStr)
	// 		id, err := strconv.Atoi(idStr)
	// 		if err != nil {
	// 			return nil, errors.New("invalid tag id: " + idStr)
	// 		}

	// 		// Check tag existence
	// 		tag, err := uc.TagRepo.FindById(uint(id))
	// 		if err != nil || tag == nil {
	// 			return nil, errors.New("tag not found: " + strconv.Itoa(int(id)))
	// 		}

	// 	}
	// }

	// if input.CoverId != nil && *input.CoverId > 0 {
	// 	_, err := uc.FileRepo.FindById(*input.CoverId, []uint{constants.StatusActive})
	// 	if err != nil {
	// 		return nil, errors.New("file not found")
	// 	}

	// 	category.CoverId = input.CoverId
	// }

	// category.TagIDs = input.TagIDs
	category.UpdatedAt = time.Now()

	// Save changes
	return uc.Repo.Update(category)
}

//----------------------------------------

// func parseTagIDs(tagIDs string) []uint {
// 	var ids []uint
// 	for _, idStr := range strings.Split(tagIDs, ",") {
// 		if idStr = strings.TrimSpace(idStr); idStr != "" {
// 			if idVal, err := strconv.Atoi(idStr); err == nil {
// 				ids = append(ids, uint(idVal))
// 			}
// 		}
// 	}
// 	return ids
// }
