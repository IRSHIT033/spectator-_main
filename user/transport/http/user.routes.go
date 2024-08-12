package http

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"spectator.main/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UsrUsecase domain.UserUsecase
}

func NewUserHandler(r *gin.RouterGroup, uu domain.UserUsecase, ) {
	handler := &UserHandler{
		UsrUsecase: uu,
	}
	r.POST("/user", handler.InsertOne)
	r.GET("/user", handler.FindOne)
	r.GET("/users", handler.GetAll)
	r.PUT("/user", handler.UpdateOne)
}

func isRequestValid(m *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (user *UserHandler) InsertOne(ctx *gin.Context) {
	var (
		usr domain.User
	)

	err := ctx.Bind(&usr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var ok bool
	if ok, err = isRequestValid(&usr); !ok {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := user.UsrUsecase.InsertOne(ctx, &usr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (user *UserHandler) FindOne(ctx *gin.Context) {

	id, _ := ctx.GetQuery("id")

	result, err := user.UsrUsecase.FindOne(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (user *UserHandler) GetAll(ctx *gin.Context) {

	type Response struct {
		Total       int64         `json:"total"`
		PerPage     int64         `json:"per_page"`
		CurrentPage int64         `json:"current_page"`
		LastPage    int64         `json:"last_page"`
		From        int64         `json:"from"`
		To          int64         `json:"to"`
		User        []domain.User `json:"users"`
	}

	var (
		res   []domain.User
		count int64
	)
	rp_ctx, _ := ctx.GetQuery("rp")
	rp, err := strconv.ParseInt(rp_ctx, 10, 64)
	if err != nil {
		rp = 25
	}

	p_ctx, _ := ctx.GetQuery("p")
	page, err := strconv.ParseInt(p_ctx, 10, 64)
	if err != nil {
		page = 1
	}

	name_ctx, _ := ctx.GetQuery("name")
	filters := bson.D{{Key: "name", Value: primitive.Regex{Pattern: ".*" + name_ctx + ".*", Options: "i"}}}

	res, count, err = user.UsrUsecase.GetAllWithPage(ctx, rp, page, filters, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := Response{
		Total:       count,
		PerPage:     rp,
		CurrentPage: page,
		LastPage:    int64(math.Ceil(float64(count) / float64(rp))),
		From:        page*rp - rp + 1,
		To:          page * rp,
		User:        res,
	}

	ctx.JSON(http.StatusOK, result)
}

func (user *UserHandler) UpdateOne(ctx *gin.Context) {

	id, _ := ctx.GetQuery("id")

	var (
		usr domain.User
		err error
	)

	err = ctx.Bind(&usr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	result, err := user.UsrUsecase.UpdateOne(ctx, &usr, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, result)
}
