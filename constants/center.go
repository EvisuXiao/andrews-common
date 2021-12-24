package constants

type Nacos struct {
	Hosts       []string `json:"hosts" binding:"required"`
	Namespace   string   `json:"namespace"`
	Username    string   `json:"username" binding:"required"`
	Password    string   `json:"password" binding:"required"`
	GroupName   string   `json:"group_name"`
	TempPath    string   `json:"temp_path"`
	ServiceName string
}
