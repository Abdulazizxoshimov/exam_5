package v1

import (
	"context"
	_ "exam_5/api_gateway/api/docs"
	"exam_5/api_gateway/internal/usecase/refresh_token"
	"net/http"
	"strconv"
	"time"

	// "github.com/casbin/casbin/v2"
	// "exam_5/api_gateway/api/middleware"
	"exam_5/api_gateway/api/models"
	pbu "exam_5/api_gateway/genproto/clientProto"

	// grpcClient "exam_5/api_gateway/internal/infrastructure/grpc_service_client"
	// "exam_5/api_gateway/internal/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"google.golang.org/protobuf/encoding/protojson"
)

// GetCategories
// @Security ApiKeyAuth
// @Router /v1/users/create [post]
// @Summary Get categories
// @Description Get categories
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.CreateUser true "createUserModel"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
func (h HandlerV1) CreateUser(c *gin.Context) {
	
	
	var (
		body        models.CreateUser
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true


	err := c.ShouldBindJSON(&body)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((10)))
	defer cancel()

	id := uuid.NewString()

	h.RefreshToken = refresh_token.JWTHandler{
		Sub:        id,
		Role:       "user",
		SigningKey: h.Config.Token.SignInKey,
		Log:        h.Logger,
		Email:      body.Email,
		Name:       body.Name,
	}

	_, refresh, err := h.RefreshToken.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating jwt",
		})
		return
	}

	

	response, err := h.Service.ClientService().Create(ctx, &pbu.User{
		Id: id,
		LastName: body.LastName,
		Name: body.Name,
		Email: body.Email,
		Password: body.Password,
		Role: "user",
		RefreshToken: refresh,

	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// err = h.BrokerProducer.ProduceUserToCreate(ctx, h.Config.Kafka.Topic.UserCreateTopic, &models.User{
	// 	LastName: body.LastName,
	// 	Name: body.Name,
	// 	Email: body.Email,
	// 	Password: body.Password,
	// 	// Role: "user",
	// 	// RefreshToken: "dasdfasdfdfdfasdfasdfda;sldkfja;sdflkajjsdf",
	// })
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }


	respProduct := pbu.User{
		Id: response.Id,
		LastName: response.LastName,
		Name: response.Name,
		Email: response.Email,
		Password: response.Password,
		Role: "user",
		RefreshToken: response.RefreshToken,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	}

	c.JSON(http.StatusCreated, respProduct)

}

// GetUser get users list
// @Summary ListUser
// @Security ApiKeyAuth
// @Description Api for getting user list
// @Tags user
// @Accept json
// @Produce json
// @Param page query integer true "page"
// @Param limit query integer true "limit"
// @Success 200 {object} models.UserList
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/list [get]
func (h *HandlerV1) ListUsers(c *gin.Context) {
	p := c.Query("page")

	page, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	lm := c.Query("limit")

	limit, err := strconv.Atoi(lm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()


	
	response, err := h.Service.ClientService().List(ctx, &pbu.GetAllUsersRequest{
		Page: int64(page),
			Limit: int64(limit),
		})
	
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		
	}
	pp.Println(response)

	c.JSON(http.StatusOK, response)
}


// GetUser gets user by id
// @Summary GetUser
// @Security ApiKeyAuth
// @Description Api for getting user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "id or email"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/get/{id} [get]
func (h *HandlerV1) GetUser(c *gin.Context)  {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	response, err := h.Service.ClientService().Get(
		ctx, &pbu.GetRequest{
			Field: "id",
			Useroremail: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respProduct := pbu.User{
		Id: response.Id,
		LastName: response.LastName,
		Name: response.Name,
		Email: response.Email,
		Password: response.Password,
		Role: "user",
		RefreshToken: response.RefreshToken,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	}

	c.JSON(http.StatusOK, respProduct)
}

// GetUser update user
// @Summary UpdateUser
// @Security ApiKeyAuth
// @Description Api for update user
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.User true "createUserModel"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/update/:id [put]
func (h *HandlerV1) UpdateUser(c *gin.Context)  {
	var (
		body      models.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	response, err := h.Service.ClientService().Update(ctx, &pbu.User{
		Id: body.Id,
		LastName: body.LastName,
		Name: body.Name,
		Email: body.Email,
		Password: body.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
    // err = h.BrokerProducer.ProduceUserToCreate(ctx, h.Config.Kafka.Topic.UserCreateTopic, &body)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, response)
}

// GetUser delete user
// @Summary DeleteUser
// @Security ApiKeyAuth
// @Description Api for delete user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/delete/{id} [delete]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

		
		response, err := h.Service.ClientService().Delete(
			ctx, &pbu.DeleteRequest{
				Field: "id",
				Value: id,
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response)
}
