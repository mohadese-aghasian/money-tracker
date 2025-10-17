package usecase

import (
	"errors"
	"fmt"
	"money-tracker/internal/constants"
	"money-tracker/internal/dto"
	"money-tracker/internal/entity"

	"os"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	Repo      entity.UserRepository
	TokenRepo entity.UserTokenRepository
}

func NewUserUseCase(repo entity.UserRepository, tokenRepo entity.UserTokenRepository) *UserUseCase {
	return &UserUseCase{
		Repo:      repo,
		TokenRepo: tokenRepo,
	}
}

func getModelColumns(model interface{}) map[string]bool {
	t := reflect.TypeOf(model)
	columns := make(map[string]bool)

	for i := 0; i < t.NumField(); i++ {
		// use struct field name in snake_case as DB column
		name := t.Field(i).Name
		columns[toSnakeCase(name)] = true
	}

	return columns
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_', r+32) // to lowercase
		} else {
			if r >= 'A' && r <= 'Z' {
				result = append(result, r+32)
			} else {
				result = append(result, r)
			}
		}
	}
	return string(result)
}

func generateToken(userID uint, user_name string) (string, error) {
	expirationTime := time.Now().Add(720 * time.Hour) // month

	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": user_name,
		"exp":       expirationTime.Unix(), // Expiration time in Unix format
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secretKey))
}

// /-------------------------- insert -------------------------------
func (uc *UserUseCase) Add(input dto.AddUserInput) (*entity.User, error) {

	if input.StatusID < 0 || input.StatusID == 0 {
		input.StatusID = 1
	}

	existedUser, er := uc.Repo.FindByUserName(input.UserName)
	fmt.Println("-*******--- %v ", existedUser, er)

	if existedUser != nil {
		return nil, errors.New("user duplicate!")
	}
	if input.LevelManage == 0 {
		input.LevelManage = 2 //user
	}

	user, err := entity.NewUser(input.UserName, input.Password, int8(input.LevelManage), input.StatusID)
	if err != nil {
		return nil, err
	}

	res := uc.Repo.Insert(user)
	if res != nil {
		return nil, res
	}

	return user, nil

}

// ---------------------------------------- Login ----------------------
func (uc *UserUseCase) Login(req dto.LoginRequest) (string, int8, uint, error) {
	user, err := uc.Repo.FindByUserName(req.UserName)
	if err != nil {
		return "", 0, 0, errors.New("user name not found!")
	}

	if !user.CheckPassword(req.Password) {
		return "", 0, 0, errors.New("invalid password")
	}

	// Generate JWT
	token, err := generateToken(user.ID, user.UserName)
	if err != nil {
		return "", 0, 0, err
	}

	userToken, e := entity.NewUserToken(token, user.ID)
	if e != nil {
		return "", 0, 0, e
	}

	er := uc.TokenRepo.Insert(userToken)
	if er != nil {
		return "", 0, 0, er
	}

	return token, user.LevelManage, user.ID, nil
}

// ------------------
func (uc *UserUseCase) Logout(input dto.LogoutInput) error {

	existedUser, _ := uc.TokenRepo.FindByToken(input.Token, []uint{constants.StatusActive})
	if existedUser == nil {
		return errors.New("token not found!")
	}

	res := uc.TokenRepo.Delete(existedUser.ID)
	if res != nil {
		return nil, res
	}

	return user, nil

}

// /////---------------------- delete ---------------------
func (uc *UserUseCase) Remove(id uint) error {
	_, err := uc.Repo.FindById(id)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.Repo.Delete(id)
}

// ////---------------------------update--------------------
func (uc *UserUseCase) Update(input dto.UpdateUserRequest) (*entity.User, error) {
	user, err := uc.Repo.FindById(input.ID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if input.UserName != "" {
		user.UserName = input.UserName
	}
	if input.Password != "" {
		// hash the new password
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}
	if input.LevelManage != 0 {
		user.LevelManage = input.LevelManage
	}
	if input.StatusID != 0 {
		user.StatusID = input.StatusID
	}

	user.UpdatedAt = time.Now()

	return uc.Repo.Update(user)
}

// /--------------------------------- GET -----------------------
func (uc *UserUseCase) Get(input dto.ListUsersInput) ([]entity.User, int, error) {
	if input.Start < 0 {
		input.Start = 0
	}
	if input.Limit <= 0 {
		input.Limit = 30
	}

	sort := strings.ToUpper(input.Sort)
	if sort != "ASC" && sort != "DESC" {
		sort = "DESC"
	}

	allowedColumns := getModelColumns(entity.User{})
	orderBy := strings.ToLower(input.OrderBy)
	if _, ok := allowedColumns[orderBy]; !ok || orderBy == "" {
		orderBy = "id"
	}

	return uc.Repo.FindAll(input)

}
