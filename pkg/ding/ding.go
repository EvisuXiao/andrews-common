package ding

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/pkg/curl"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Ding struct {
	appKey           string
	appSecret        string
	token            string
	tokenExpiredTime time.Time
	dep              Departments
	depSub           map[int][]int
}

type IErrorResult interface {
	GetErrCode() int
	GetErrMsg() string
}

type ErrorResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type User struct {
	UserId    string `json:"userid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	OrgEmail  string `json:"orgEmail"`
	WorkPlace string `json:"workPlace"`
	Branch    string `json:"-"`
	Position  string `json:"position"`
	Active    bool   `json:"active"`
}
type Users []*User

type Department struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parentid"`
}
type Departments []*Department

const (
	host        = "https://oapi.dingtalk.com"
	successCode = 0

	funcToken   = "/gettoken"
	funcScope   = "/auth/scopes"
	funcDep     = "/department/list"
	funcUser    = "/user/listbypage"
	funcSendMsg = "/robot/send"
)

func Client(appKey, appSecret string) *Ding {
	d := &Ding{appKey: appKey, appSecret: appSecret}
	depList, err := d.GetDepList()
	if utils.HasErr(err) {
		d.dep = Departments{}
	} else {
		d.dep = depList
	}
	d.depSub = d.dep.GetSubList()
	return d
}

func (d *Ding) request(funcName, method string, data map[string]interface{}, result IErrorResult) error {
	if funcName != funcToken {
		token := d.getToken()
		if utils.IsEmpty(token) {
			return errors.New("token is required")
		}
		data["access_token"] = token
	}
	err := curl.Request(host+funcName, method, data, result)
	if utils.HasErr(err) {
		return err
	}
	if result.GetErrCode() != successCode {
		return errors.New(result.GetErrMsg())
	}
	return nil
}

func (d *Ding) getToken() string {
	now := time.Now()
	if !utils.IsEmpty(d.token) && now.Before(d.tokenExpiredTime) {
		return d.token
	}
	var result struct {
		ErrorResult
		AccessToken string `json:"access_token"`
	}
	err := d.request(funcToken, http.MethodGet, map[string]interface{}{"appkey": d.appKey, "appsecret": d.appSecret}, &result)
	if utils.HasErr(err) {
		logging.Error("Ding token request err: %+v", err)
	}
	d.token = result.AccessToken
	d.tokenExpiredTime = now.Add(30 * time.Minute)
	return d.token
}

func (d *Ding) GetDepScope() ([]int, error) {
	var result struct {
		ErrorResult
		AuthOrgScopes struct {
			Dept []int `json:"authed_dept"`
		} `json:"auth_org_scopes"`
	}
	err := d.request(funcScope, http.MethodGet, map[string]interface{}{}, &result)
	return result.AuthOrgScopes.Dept, err
}

func (d *Ding) GetDepList() (Departments, error) {
	var result struct {
		ErrorResult
		Department Departments `json:"department"`
	}
	err := d.request(funcDep, http.MethodGet, map[string]interface{}{}, &result)
	return result.Department, err
}

func (d *Ding) GetDirectUsersByDepId(depId int) (Users, error) {
	var result struct {
		ErrorResult
		HasMore  bool  `json:"hasMore"`
		UserList Users `json:"userlist"`
	}
	var err error
	offset := 0
	limit := 100
	users := Users{}
	for {
		err = d.request(funcUser, http.MethodGet, map[string]interface{}{"department_id": depId, "offset": offset, "size": limit}, &result)
		if utils.HasErr(err) {
			return nil, err
		}
		users = append(users, result.UserList...)
		if !result.HasMore {
			break
		}
		offset += limit
	}
	return users, nil
}

func (d *Ding) GetAllUsersByDepId(depId int) (Users, error) {
	users := Users{}
	err := d.getAllUsersByDepIds([]int{depId}, &users)
	if utils.HasErr(err) {
		return nil, err
	}
	return users, nil
}

func (d *Ding) getAllUsersByDepIds(depIds []int, users *Users) error {
	for _, depId := range depIds {
		depUsers, err := d.GetDirectUsersByDepId(depId)
		if utils.HasErr(err) {
			return err
		}
		branch := d.dep.GetFullName(depId)
		for _, user := range depUsers {
			user.Branch = branch
			*users = append(*users, user)
		}
		if subList, ok := d.depSub[depId]; ok && !utils.IsEmpty(subList) {
			err = d.getAllUsersByDepIds(subList, users)
			if utils.HasErr(err) {
				return err
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}

func (d *Ding) GetAllUsers() (Users, error) {
	depList, err := d.GetDepScope()
	if utils.HasErr(err) {
		return nil, err
	}
	users := Users{}
	err = d.getAllUsersByDepIds(depList, &users)
	if utils.HasErr(err) {
		return nil, err
	}
	return users, nil
}

func (e *ErrorResult) GetErrCode() int {
	return e.ErrCode
}

func (e *ErrorResult) GetErrMsg() string {
	return e.ErrMsg
}

func (d Departments) GetMapById() map[int]*Department {
	depMap := make(map[int]*Department)
	for _, v := range d {
		depMap[v.Id] = v
	}
	return depMap
}

func (d Departments) GetSubList() map[int][]int {
	subMap := make(map[int][]int)
	for _, v := range d {
		if _, ok := subMap[v.ParentId]; !ok {
			subMap[v.ParentId] = []int{}
		}
		subMap[v.ParentId] = append(subMap[v.ParentId], v.Id)
	}
	return subMap
}

func (d Departments) GetFullName(id int) string {
	var nameArr []string
	nameMap := d.GetMapById()
	for {
		if v, ok := nameMap[id]; ok {
			nameArr = append([]string{v.Name}, nameArr...)
			if v.ParentId > 1 {
				id = v.ParentId
			} else {
				break
			}
		} else {
			break
		}
	}
	return strings.Join(nameArr, "-")
}

func (u Users) GetUserMapByUsername() map[string]*User {
	uMap := make(map[string]*User)
	for _, user := range u {
		username := user.GetUsername()
		if !utils.IsEmpty(username) {
			uMap[username] = user
		}
	}
	return uMap
}

func (u *User) GetUsername() string {
	return u.getUsernameByEmail(utils.Or(u.OrgEmail, u.Email).(string))
}

func (u *User) getUsernameByEmail(email string) string {
	splitIndex := strings.Index(email, "@")
	if splitIndex > 0 {
		return email[0:splitIndex]
	}
	return ""
}
