package demo

import (
	"time"

	"gorm.io/datatypes"

	"github.com/EvisuXiao/andrews-common/database"
)

type User struct {
	ModelUpdatable
	Username      string         `json:"username"`
	Nickname      string         `json:"nickname"`
	Realname      string         `json:"realname"`
	Password      string         `json:"password"`
	Email         string         `json:"email"`
	Phone         string         `json:"phone"`
	Oauth         datatypes.JSON `json:"oauth" gorm:"default:'{}'"`
	Branch        string         `json:"branch"`
	Position      string         `json:"position"`
	Source        string         `json:"source"`
	Salt          string         `json:"salt"`
	LastLoginIp   string         `json:"last_login_ip" gorm:"default:0.0.0.0"`
	LastLoginTime time.Time      `json:"last_login_time"`
	Enabled       bool           `json:"enabled"`
}

type Users []*User

var userModel = &User{}

func init() {
	database.RegisterModel(userModel)
}

func NewUserModel() *User {
	return userModel
}

func (m *User) TableName() string {
	return "rbac_user"
}

func (m *User) GetRows(options *database.Options) (Users, error) {
	if !options.HasOrder() {
		options.AddAscOrder("username")
	}
	var rows Users
	err := m.GetAnyRows(options, &rows)
	return rows, err
}
