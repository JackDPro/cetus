package controller

import (
	"github.com/JackDPro/cetus/model"
	"github.com/JackDPro/cetus/provider"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log/level"
	"net/http"
	"reflect"
	"strconv"
)

func ResponseUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, &map[string]string{
		"message": "unauthorized",
	})
}
func ResponseForbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, &map[string]string{
		"message": "forbidden",
	})
}

func ResponseAccepted(c *gin.Context) {
	c.JSON(http.StatusAccepted, &model.Success{Code: 0})
}
func ResponseSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, &model.Success{Code: 0})
}
func ResponseCreated(c *gin.Context, id uint64) {
	c.Writer.Header().Set("Location", strconv.FormatUint(id, 10))
	c.JSON(http.StatusCreated, map[string]uint64{
		"id": id,
	})
}

func ResponseError(c *gin.Context, status int, code int, message string, err error) {
	detail := ""
	if err != nil {
		detail = err.Error()
	}
	c.AbortWithStatusJSON(status, model.Error{
		Code:    code,
		Message: message,
		Detail:  detail,
	})
}

func ResponseInternalError(c *gin.Context, code int, message string, err error) {
	_ = level.Error(provider.GetLogger()).Log(err)
	ResponseError(c, http.StatusInternalServerError, code, message, nil)
}

func ResponseUnprocessable(c *gin.Context, code int, message string, err error) {
	ResponseError(c, http.StatusUnprocessableEntity, code, message, err)
}

func ResponseBadRequest(c *gin.Context, code int, message string) {
	ResponseError(c, http.StatusBadRequest, code, message, nil)
}

func ResponseNotFound(c *gin.Context, message string) {
	ResponseError(c, http.StatusNotFound, 1, message, nil)
}

func ResponseItem(c *gin.Context, item model.IModel) {
	data, err := item.ToMap()
	if err != nil {
		ResponseInternalError(c, 9827, "convert data failed", err)
		return
	}
	c.JSON(http.StatusOK, &model.DataWrapper{
		Data: data,
		Meta: nil,
	})
}
func ResponseCollection[T any](c *gin.Context, items []T, meta *model.Meta) {
	data := make([]interface{}, len(items))
	var err error
	for index, item := range items {
		mod := reflect.ValueOf(&item).Interface().(model.IModel)
		data[index], err = mod.ToMap()
		if err != nil {
			_ = level.Error(provider.GetLogger()).Log("message", "model to map failed", "error", err)
			continue
		}
	}
	jsonData := &model.DataWrapper{
		Data: data,
	}
	if meta != nil && !meta.IsNull() {
		jsonData.Meta = meta
	}
	c.JSON(http.StatusOK, jsonData)
}
