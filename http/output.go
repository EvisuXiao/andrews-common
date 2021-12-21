package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiOutput struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewOutput(code int) *ApiOutput {
	o := &ApiOutput{}
	o.Code = code
	return o
}

func (o *ApiOutput) ApiResponse(ctx *gin.Context) bool {
	ctx.AbortWithStatusJSON(http.StatusOK, o)
	return true
}

func (o *ApiOutput) SetCode(code int) {
	o.Code = code
}

func (o *ApiOutput) SetMessage(msg string) {
	o.Message = msg
}

func (o *ApiOutput) SetData(data interface{}) {
	o.Data = data
}
