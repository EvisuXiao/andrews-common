package exception

import "net/http"

const (
	SUCCESS_CODE = http.StatusOK
	FAILURE_CODE = http.StatusInternalServerError
	PARAM_CODE   = http.StatusBadRequest

	SUCCESS_MSG       = "操作成功"
	FAILURE_MSG       = "操作失败"
	INVALID_PARAM_MSG = "请求参数有误"
	DB_ERROR_MSG      = "数据库操作异常"
	SERVER_ERROR_MSG  = "系统内部异常"
)
