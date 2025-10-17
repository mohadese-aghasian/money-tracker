package usecase

import (
	"errors"

	"money-tracker/internal/constants"
	"money-tracker/internal/dto"
	"money-tracker/internal/entity"
	"strconv"
	"strings"
	"time"
)

type PurchaseUseCase struct {
	Repo    entity.PurchaseRepository
	TagRepo entity.TagRepository
	CatRepo entity.CategoryRepository
}

func NewPurchaseUseCase(repo entity.PurchaseRepository) *PurchaseUseCase {
	return &PurchaseUseCase{Repo: repo}
}

// /-----------------------add-----------------------------
func (uc *PurchaseUseCase) Add(input dto.AddPurchaseInput) (*entity.Purchase, error) {
	// default status
	if input.StatusID == 0 {
		input.StatusID = 1
	}

	category, err_c := uc.CatRepo.FindById(*input.CategoryId)
	if err_c != nil || category == nil {
		return nil, errors.New("category not found")
	}

	if input.SubCategoryId != nil {
		category, err_c := uc.CatRepo.FindById(*input.SubCategoryId)
		if err_c != nil || category == nil {
			return nil, errors.New("sub category not found")
		}
	}

	// create new category
	purchase, err := entity.NewPurchase(input.Amount, input.Date, input.CategoryId, input.StatusID)
	if err != nil {
		return nil, err
	}

	// category.Slug = input.Slug
	// category.CoverId = coverId

	// handle tag IDs if provided
	if input.TagIDs != "" {
		idsStr := strings.Split(input.TagIDs, ",")
		for _, idStr := range idsStr {
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, errors.New("invalid tag id: " + idStr)
			}

			// Check tag existence
			tag, err := uc.TagRepo.FindById(uint(id))
			if err != nil || tag == nil {
				return nil, errors.New("tag not found: " + strconv.Itoa(int(id)))
			}
		}
		purchase.TagIDs = input.TagIDs
	}

	purchase.SubCategoryId = input.SubCategoryId
	purchase.Note = input.Note
	purchase.Color = input.Color
	purchase.Method = input.Method

	// insert into repo
	res := uc.Repo.Insert(purchase)
	if res != nil {
		return nil, res
	}

	return purchase, nil
}

// ----------------------------------------------
func (uc *PurchaseUseCase) Get(input dto.PurchaseFindAll) ([]dto.PurchaseResponse, int, error) {

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

	allowedColumns := getModelColumns(entity.Purchase{})
	if _, ok := allowedColumns[strings.ToLower(input.OrderBy)]; !ok || input.OrderBy == "" {
		input.OrderBy = "id"
	}

	purchases, count, err := uc.Repo.FindAll(input)

	var responses []dto.PurchaseResponse
	for _, pur := range purchases {

		tagIDs := parseTagIDs(pur.TagIDs)
		var tags []dto.FetchedTag
		if len(tagIDs) > 0 {
			pureTags, _ := uc.TagRepo.FindByIDs(tagIDs)

			for _, tag := range pureTags {
				tags = append(tags, dto.FetchedTag{
					ID:        tag.ID,
					Title:     tag.Title,
					StatusID:  tag.StatusID,
					CreatedAt: tag.CreatedAt,
				})
			}
		}

		var category *dto.CategoryResponse
		if pur.Category != nil {
			category = &dto.CategoryResponse{
				ID:        pur.Category.ID,
				Title:     pur.Category.Title,
				StatusID:  pur.Category.StatusID,
				Slug:      pur.Category.Slug,
				CreatedAt: pur.Category.CreatedAt,
			}
		}
		var subcategory *dto.CategoryResponse
		if pur.SubCategory != nil {
			subcategory = &dto.CategoryResponse{
				ID:        pur.Category.ID,
				Title:     pur.Category.Title,
				StatusID:  pur.Category.StatusID,
				Slug:      pur.Category.Slug,
				CreatedAt: pur.Category.CreatedAt,
			}
		}

		responses = append(responses, dto.PurchaseResponse{
			ID:          pur.ID,
			Reason:      pur.Reason,
			StatusID:    pur.StatusID,
			Date:        pur.Date,
			Amount:      pur.Amount,
			Tags:        tags,
			Category:    category,
			SubCategory: subcategory,
			Note:        pur.Note,
			Color:       pur.Color,
			Method:      pur.Method,
			CreatedAt:   pur.CreatedAt,
		})
	}

	return responses, count, err

}

// /-----------------------------------------------
func (uc *PurchaseUseCase) Remove(id uint) error {
	_, err := uc.Repo.FindById(id, []uint{constants.StatusActive})
	if err != nil {
		return errors.New("purchase not found")
	}

	return uc.Repo.Delete(id)
}

// ---------------------------------------------------
func (uc *PurchaseUseCase) Update(input dto.UpdatePurchaseInput) (*entity.Purchase, error) {
	purchase, err := uc.Repo.FindById(input.ID, []uint{constants.StatusActive})
	if err != nil {
		return nil, errors.New("purchase not found")
	}

	if input.Reason != "" {
		purchase.Reason = input.Reason
	}
	if input.StatusID != 0 {
		purchase.StatusID = input.StatusID
	}
	if input.CategoryId != nil {
		_, err := uc.CatRepo.FindById(*input.CategoryId)
		if err != nil {
			return nil, errors.New("category not found")
		}
		purchase.CategoryId = input.CategoryId
	}

	if input.TagIDs != "" {
		idsStr := strings.Split(input.TagIDs, ",")
		for _, idStr := range idsStr {
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, errors.New("invalid tag id: " + idStr)
			}

			// Check tag existence
			tag, err := uc.TagRepo.FindById(uint(id))
			if err != nil || tag == nil {
				return nil, errors.New("tag not found: " + strconv.Itoa(int(id)))
			}

		}
		purchase.TagIDs = input.TagIDs
	}

	if input.SubCategoryId != nil && *input.SubCategoryId > 0 {
		_, err := uc.CatRepo.FindById(*input.SubCategoryId)
		if err != nil {
			return nil, errors.New("subcat not found")
		}

		purchase.SubCategoryId = input.SubCategoryId
	}

	purchase.Amount = input.Amount
	purchase.Color = input.Color
	purchase.Note = input.Note
	purchase.Date = input.Date
	purchase.Method = input.Method

	purchase.UpdatedAt = time.Now()

	// Save changes
	return uc.Repo.Update(purchase)
}

//----------------------------------------

func parseTagIDs(tagIDs string) []uint {
	var ids []uint
	for _, idStr := range strings.Split(tagIDs, ",") {
		if idStr = strings.TrimSpace(idStr); idStr != "" {
			if idVal, err := strconv.Atoi(idStr); err == nil {
				ids = append(ids, uint(idVal))
			}
		}
	}
	return ids
}
