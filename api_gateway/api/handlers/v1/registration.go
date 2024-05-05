package v1

import (
	"context"
	"encoding/json"
	"exam_5/api_gateway/api/models"
	pb "exam_5/api_gateway/genproto/clientProto"
	regtool "exam_5/api_gateway/internal/pkg/regtool"
	"exam_5/api_gateway/internal/usecase/refresh_token"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	govalidator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
)

// register register users
// @Summary RegisterUser
// @Description Api for register user
// @Tags registration
// @Accept json
// @Produce json
// @Param User body models.CreateUser true "RegisterUser"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/register [post]
func (h *HandlerV1) Register(c *gin.Context) {

	var (
		body        models.CreateUser
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	emailErr := validation.Validate(body.Email, validation.Required)

	if emailErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})
	}

	passwordErr := validation.Validate(body.Password, validation.Required, validation.Length(8, 20), validation.Match(regexp.MustCompile("[a-z]|[A-Z][0-9]")))

	if passwordErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password should be 8-20 characters long and contain at least one lowercase letter, one uppercase letter, and one digit",
		})
	}

	exists, err := h.Service.ClientService().IsUnique(ctx, &pb.IsUniqueRequest{
		Field: "email",
		Value: body.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	if exists.IsUnique {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email already in use:",
		})
	}
	valid := govalidator.IsEmail(body.Email)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad email",
		})
	}

	err = h.redisStorage.Set(ctx, body.Email, body, time.Second*300)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	radomNumber, err := regtool.SendCodeGmail(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	randNum, err := strconv.Atoi(radomNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Println(err)
		return
	}

	err = h.redisStorage.Set(ctx, radomNumber, randNum, time.Second*300)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, body)
}

// register verify users
// @Summary RegisterUser
// @Description Api for verify user
// @Tags registration
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Param code query string true "code"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/verify [post]
func (h *HandlerV1) Verify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	code := c.Query("code")
	email := c.Query("email")

	number, err := strconv.Atoi(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while parsing string to int",
		})
	}

	userData, err := h.redisStorage.Get(ctx, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while getting user from redis",
		})
	}
	rand, err := h.redisStorage.Get(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while getting code from redis",
		})
	}
	var user models.User

	err = json.Unmarshal(userData, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error parsing user data from Redis",
		})
		log.Println(err)
		return
	}
	// n := string(rand)
	// randNum := n[1 : len(n)-1]

	if user.Email != email && strconv.Itoa(number) != string(rand) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email yoki kod xato",
		})
		log.Println(err)
		return
	}

	id := uuid.NewString()

	h.RefreshToken = refresh_token.JWTHandler{
		Sub:        id,
		Role:       "user",
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      user.Email,
		Name:       user.Name,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating jwt",
		})
		return
	}

	// hashPassword, err := regtool.HashPassword(user.Password)
	// if err != nil {
	// 	c.JSON(http.StatusBadGateway, gin.H{
	// 		"error": "Ooops something went wrong",
	// 	})
	// }

	claims, err := refresh_token.ExtractClaim(access, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Ooops something went wrong",
		})
	}

	response, err := h.Service.ClientService().Create(ctx, &pb.User{
		Id:           id,
		Name:         user.Name,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: refresh,
		Role:         cast.ToString(claims["role"]),
	})

	respUser := &models.UserResponse{
		Id:           response.Id,
		Name:         response.Name,
		LastName:     response.LastName,
		Email:        response.Email,
		Role:         response.Role,
		Password:     response.Password,
		RefreshToken: response.RefreshToken,
		AccessToken:  access,
		CreatedAt:    response.CreatedAt,
		UpdatedAt:    response.UpdatedAt,
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, respUser)

}

// Login  users
// @Summary LoginUser
// @Description Api for user user
// @Tags registration
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Param password query string true "password"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/login [post]
func (h *HandlerV1) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	email := c.Query("email")
	password := c.Query("password")
	response, err := h.Service.ClientService().Get(
		ctx, &pb.GetRequest{
			Field:       "email",
			Useroremail: email,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Println(err)
		return
	}
	// if !(regtool.CheckHashPassword(password, response.Password)) {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Noto'gri parol!!!!",
	// 	})
	// 	log.Println(err)
	// 	return
	// }

	if password != response.Password {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Noto'gri parol!!!!",
		})
		log.Println(err)
		return
	}

	h.RefreshToken = refresh_token.JWTHandler{
		Sub:        response.Id,
		Role:       response.Role,
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      response.Email,
		Name:       response.Name,
	}

	access, refresh, err := h.RefreshToken.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating jwt",
		})
		log.Println()
		return
	}
	respUser := &models.UserResponse{
		Id:           response.Id,
		Name:         response.Name,
		LastName:     response.LastName,
		Email:        response.Email,
		Password:     response.Password,
		RefreshToken: refresh,
		AccessToken:  access,
		Role:         response.Role,
		CreatedAt:    response.CreatedAt,
		UpdatedAt:    response.UpdatedAt,
	}
	_, err = h.Service.ClientService().UpdateUserRefreshToken(ctx, &pb.UpdateRefreshToken{
		Id:           response.Id,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, respUser)

}

// Login  users
// @Summary LoginUser
// @Description Api for user user
// @Tags registration
// @Accept json
// @Produce json
// @Param refreshToken query string true "Refresh Token"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/updatetoken [post]
func (h *HandlerV1) UpdateAccesTokenWithRefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	RToken := c.Query("refreshToken")

	respUser, err := h.Service.ClientService().Get(ctx, &pb.GetRequest{
		Field:       "refresh_token",
		Useroremail: RToken,
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": "error while getting user",
		})
		log.Println(err)
		return
	}

	resclaim, err := refresh_token.ExtractClaim(respUser.RefreshToken, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "oops something went wrong",
		})
		log.Println(err)
		return
	}
	Now_time := time.Now().Unix()
	exp := (resclaim["exp"])
	if exp.(float64)-float64(Now_time) > 0 {

		h.RefreshToken = refresh_token.JWTHandler{
			Sub:        respUser.Id,
			Role:       respUser.Role,
			SigningKey: h.Config.Token.SignInKey,
			Log:        h.Logger,
			Email:      respUser.Email,
			Name:       respUser.Name,
		}

		access, _, err := h.RefreshToken.GenerateAuthJWT()

		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "error while generating jwt",
			})
			log.Println(err)
			return
		}
		respUser := &models.UserResponse{
			Id:           respUser.Id,
			Name:         respUser.Name,
			LastName:     respUser.LastName,
			Email:        respUser.Email,
			Password:     respUser.Password,
			RefreshToken: respUser.RefreshToken,
			AccessToken:  access,
			Role:         respUser.Role,
			CreatedAt:    respUser.CreatedAt,
			UpdatedAt:    respUser.UpdatedAt,
		}

		c.JSON(http.StatusCreated, respUser)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid login",
		})
		println(err)
		return
	}
}
