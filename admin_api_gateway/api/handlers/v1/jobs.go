package v1

import (
	"context"
	"exam_5/admin_api_gateway/api/models"
	pb "exam_5/admin_api_gateway/genproto/jobProto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"google.golang.org/protobuf/encoding/protojson"
)

// GetCategories
// @Router /v1/jobs/create [post]
// @Summary Get categories
// @Description Get categories
// @Tags job
// @Accept json
// @Produce json
// @Param User body models.JobCreateReq true "createJobModel"
// @Success 200 {object} models.Job
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
func (h HandlerV1) Create(c *gin.Context) {

	var (
		body        models.JobCreateReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	if body.StartDate[2] != '-' {
		_, err = time.Parse("2006-01-02", body.StartDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Noto'g'ri sana malumoti",
			})
			return
		}

		_, err = time.Parse("2006-01-02", body.EndDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Noto'g'ri sana malumoti",
			})
			return
		}
	} else {
		_, err = time.Parse("02-01-2006", body.StartDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Noto'g'ri sana malumoti",
			})
			return
		}

		_, err = time.Parse("02-01-2006", body.EndDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Noto'g'ri sana malumoti",
			})
			return
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration((10)))
	defer cancel()
	response, err := h.Service.JobService().Create(ctx, &pb.Job{
		ClientId:  body.Client_id,
		Name:      body.Name,
		CompName:  body.CompName,
		Status:    body.Status,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		Location:  body.Location,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	respJob := pb.Job{
		Id:        response.Id,
		ClientId:  response.ClientId,
		Name:      response.Name,
		CompName:  response.CompName,
		Status:    response.Status,
		StartDate: response.StartDate,
		EndDate:   response.EndDate,
		Location:  response.Location,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	}

	c.JSON(http.StatusCreated, respJob)

}

// GetUser get job list
// @Summary ListJob
// @Description Api for getting job list
// @Tags job
// @Accept json
// @Produce json
// @Param page query integer true "page"
// @Param limit query integer true "limit"
// @Success 200 {object} models.JobList
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/jobs/list [get]
func (h *HandlerV1) List(c *gin.Context) {
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

	var filter map[string]string
	response, err := h.Service.JobService().GetAll(ctx, &pb.JobGetAllRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Filter: filter,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

	}

	c.JSON(http.StatusOK, response)
}

// Get gets job by id
// @Summary GetJob
// @Description Api for getting job by id
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "id or email"
// @Success 200 {object} models.Job
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/jobs/get/{id} [get]
func (h *HandlerV1) Get(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	response, err := h.Service.JobService().Get(
		ctx, &pb.JobGetRequest{
			Id: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respProduct := pb.Job{
		Id:        response.Id,
		ClientId:  response.ClientId,
		Name:      response.Name,
		CompName:  response.CompName,
		Status:    response.Status,
		StartDate: response.StartDate,
		EndDate:   response.EndDate,
		Location:  response.Location,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	}

	c.JSON(http.StatusOK, respProduct)
}

// GetUser update job
// @Summary UpdateJob
// @Description Api for update job
// @Tags job
// @Accept json
// @Produce json
// @Param User body models.JobUpdateRequest true "createUserModel"
// @Success 200 {object} models.Job
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/jobs/update/:id [put]
func (h *HandlerV1) Update(c *gin.Context) {
	var (
		body        models.JobUpdateRequest
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

	response, err := h.Service.JobService().Update(ctx, &pb.Job{
		ClientId:  body.Client_id,
		Name:      body.Name,
		CompName:  body.CompName,
		Status:    body.Status,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		Location:  body.Location,
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

// GetUser delete job
// @Summary DeleteJob
// @Description Api for delete Job
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Job
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/jobs/delete/{id} [delete]
func (h *HandlerV1) Delete(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	response, err := h.Service.JobService().Delete(
		ctx, &pb.JobDeleteRequest{
			Id: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUser get job list
// @Summary ListJob
// @Description Api for getting job list
// @Tags job
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param page query integer true "page"
// @Param limit query integer true "limit"
// @Success 200 {object} models.JobList
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/jobs/listbyclientid/{id} [get]
func (h *HandlerV1) ListByClientId(c *gin.Context) {

	id := c.Param("id")
	pp.Println(id)

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

	filter := map[string]string{
		"client_id": id,
	}
	response, err := h.Service.JobService().GetAll(ctx, &pb.JobGetAllRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Filter: filter,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

	}

	c.JSON(http.StatusOK, response)
}

// GetUser get job list
// @Summary ListJob
// @Description Api for getting job list
// @Tags job
// @Accept json
// @Produce json
// @Param page query integer true "page"
// @Param limit query integer true "limit"
// @Success 200 {object} models.GetAllJobByClientIdResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/jobs/listwithowner [get]
func (h *HandlerV1) ListWithOwner(c *gin.Context) {
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

	var filter map[string]string
	response, err := h.Service.JobService().GetAllJobWithOwner(ctx, &pb.JobGetAllRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Filter: filter,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

	}

	c.JSON(http.StatusOK, response)
}
