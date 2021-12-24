package http

import (
	"fmt"
	"net/http"

	"github.com/EvisuXiao/andrews-common/exception"
	cValidator "github.com/EvisuXiao/andrews-common/pkg/validator"
	"github.com/EvisuXiao/andrews-common/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct{}

func (c *Controller) SuccessResponse(ctx *gin.Context, data ...interface{}) bool {
	output := NewOutput(exception.SUCCESS_CODE)
	if len(data) > 1 {
		output.SetMessage(fmt.Sprint(data[1]))
	} else {
		output.SetMessage(exception.SUCCESS_MSG)
	}
	if !utils.IsEmpty(data) {
		output.SetData(data[0])
	}
	return output.ApiResponse(ctx)
}

func (c *Controller) FailureResponseWithCode(ctx *gin.Context, code int, desc ...interface{}) bool {
	output := NewOutput(code)
	descLen := len(desc)
	if descLen > 0 {
		msg := desc[0]
		if utils.IsEmpty(msg) {
			output.SetMessage(exception.FAILURE_MSG)
		} else if e, ok := msg.(*exception.DbError); ok {
			output.SetMessage(exception.DB_ERROR_MSG)
			output.SetData(e.Error())
		} else if e, ok := msg.(*exception.CustomError); ok {
			output.SetMessage(e.Error())
		} else {
			output.SetMessage(exception.SERVER_ERROR_MSG)
			output.SetData(fmt.Sprint(msg))
		}
		if descLen > 1 && !utils.IsEmpty(desc[1]) {
			if err, ok := desc[1].(validator.ValidationErrors); ok {
				desc[1] = cValidator.Translate(err)
			}
			output.SetData(fmt.Sprint(desc[1]))
		}
	} else {
		output.SetMessage(exception.FAILURE_MSG)
	}
	return output.ApiResponse(ctx)
}

func (c *Controller) FailureResponse(ctx *gin.Context, desc ...interface{}) bool {
	return c.FailureResponseWithCode(ctx, exception.FAILURE_CODE, desc...)
}

func (c *Controller) InvalidParamResponse(ctx *gin.Context, err error) bool {
	return c.FailureResponseWithCode(ctx, exception.PARAM_CODE, exception.INVALID_PARAM_ERR, err)
}

func (c *Controller) StaticFile(ctx *gin.Context, filename string) bool {
	ctx.File(filename)
	ctx.Abort()
	return true
}

func (c *Controller) Next(ctx *gin.Context) bool {
	ctx.Next()
	return true
}

func IsGet(ctx *gin.Context) bool {
	return ctx.Request.Method == http.MethodGet
}

func IsPost(ctx *gin.Context) bool {
	return ctx.Request.Method == http.MethodPost
}

func IsPut(ctx *gin.Context) bool {
	return ctx.Request.Method == http.MethodPut
}

func IsDelete(ctx *gin.Context) bool {
	return ctx.Request.Method == http.MethodDelete
}
