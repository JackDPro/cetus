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

type Controller struct {
}

func (ctr *Controller) ResponseUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, &map[string]string{
		"message": "unauthorized",
	})
}
func (ctr *Controller) ResponseForbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, &map[string]string{
		"message": "forbidden",
	})
}

func (ctr *Controller) ResponseAccepted(c *gin.Context) {
	c.JSON(http.StatusAccepted, &model.Success{Code: 0})
}
func (ctr *Controller) ResponseSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, &model.Success{Code: 0})
}
func (ctr *Controller) ResponseCreated(c *gin.Context, id uint64) {
	c.Writer.Header().Set("Location", strconv.FormatUint(id, 10))
	c.JSON(http.StatusCreated, map[string]uint64{
		"id": id,
	})
}

func (ctr *Controller) ResponseError(c *gin.Context, status int, code int, message string, err error) {
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

func (ctr *Controller) ResponseInternalError(c *gin.Context, code int, message string, err error) {
	_ = level.Error(provider.GetLogger()).Log(err)
	ctr.ResponseError(c, http.StatusInternalServerError, code, message, nil)
}

func (ctr *Controller) ResponseUnprocessable(c *gin.Context, code int, message string, err error) {
	ctr.ResponseError(c, http.StatusUnprocessableEntity, code, message, err)
}

func (ctr *Controller) ResponseBadRequest(c *gin.Context, code int, message string) {
	ctr.ResponseError(c, http.StatusBadRequest, code, message, nil)
}

func (ctr *Controller) ResponseNotFound(c *gin.Context, message string) {
	ctr.ResponseError(c, http.StatusNotFound, 1, message, nil)
}

func (ctr *Controller) ResponseItem(c *gin.Context, item model.IModel) {
	data, err := item.ToMap()
	if err != nil {
		ctr.ResponseInternalError(c, 9827, "convert data failed", err)
		return
	}
	c.JSON(http.StatusOK, &model.DataWrapper{
		Data: data,
		Meta: nil,
	})
}

func (ctr *Controller) ResponseCollection(c *gin.Context, items interface{}, meta *model.Meta) {
	if reflect.TypeOf(items).Kind() != reflect.Slice {
		return
	}
	itemArray := reflect.ValueOf(items)
	data := make([]interface{}, itemArray.Len())
	var err error

	for i := 0; i < itemArray.Len(); i++ {
		item := itemArray.Index(i)
		mod := reflect.ValueOf(&item).Interface().(model.IModel)
		data[i], err = mod.ToMap()
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
