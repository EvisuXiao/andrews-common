package demo

import (
	"net/http"

	"github.com/gin-gonic/gin"

	cHttp "github.com/EvisuXiao/andrews-common/http"
)

func InitRouter() *cHttp.MainRouterGroup {
	groups := []*cHttp.RouterGroup{
		{
			"",
			[]cHttp.RouterHandler{},
			[]*cHttp.RouterItem{
				{http.MethodGet, "test", (&routerTest{}).test1},
			},
		},
	}
	return &cHttp.MainRouterGroup{Groups: groups}
}

type routerTest struct {
	cHttp.Controller
}

func (r *routerTest) test1(ctx *gin.Context) bool {
	return r.SuccessResponse(ctx, "test")
}
