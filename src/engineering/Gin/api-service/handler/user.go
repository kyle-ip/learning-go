package handler

import (
	"api-service/model"
	"api-service/pkg/auth"
	"api-service/pkg/errno"
	"api-service/pkg/log"
	"api-service/pkg/log/lager"
	"api-service/pkg/token"
	"api-service/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

// VO 对象

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type ListResponse struct {
	TotalCount uint64            `json:"totalCount"`
	UserList   []*model.UserInfo `json:"userList"`
}

type SwaggerListResponse struct {
	TotalCount uint64           `json:"totalCount"`
	UserList   []model.UserInfo `json:"userList"`
}

// Create
// @Summary Add new user to the database
// @Description Add a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body CreateRequest true "Create a new user"
// @Success 200 {object} CreateResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user [post]
func Create(c *gin.Context) {
	log.Info("User Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})

	// 从 Context 中提取参数，并绑定到 CreateRequest struct 中（相当于 VO）。
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		util.SendResponse(c, errno.ErrBind, nil)
		return
	}

	u := model.UserModel{Username: r.Username, Password: r.Password}

	// Validate the data.
	if err := u.Validate(); err != nil {
		util.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		util.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}
	// Insert the user to the database.
	if err := u.Create(); err != nil {
		util.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	// Show the user information.
	util.SendResponse(c, nil, CreateResponse{Username: r.Username})
}

// Get
// @Summary Get a user by the user identifier
// @Description Get a user by username
// @Tags user
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Success 200 {object} model.UserModel "{"code":0,"message":"OK","data":{"username":"kong","password":"$2a$10$E0kwtmtLZbwW/bDQ8qI8e.eHPqhQOW9tvjwpyo/p05f/f4Qvr3OmS"}}"
// @Router /user/{username} [get]
func Get(c *gin.Context) {
	username := c.Param("username")
	// Get the user by the `username` from the database.
	user, err := model.GetUser(username)
	if err != nil {
		util.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	util.SendResponse(c, nil, user)
}

// List
// @Summary List the users in the database
// @Description List users
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body ListRequest true "List users"
// @Success 200 {object} SwaggerListResponse "{"code":0,"message":"OK","data":{"totalCount":1,"userList":[{"id":0,"username":"admin","random":"user 'admin' get random string 'EnqntiSig'","password":"$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG","createdAt":"2018-05-28 00:25:33","updatedAt":"2018-05-28 00:25:33"}]}}"
// @Router /user [get]
func List(c *gin.Context) {
	log.Info("List function called.")
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		util.SendResponse(c, errno.ErrBind, nil)
		return
	}

	infos, count, err := listUser(r.Username, r.Offset, r.Limit)
	if err != nil {
		util.SendResponse(c, err, nil)
		return
	}

	util.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		UserList:   infos,
	})
}

// Delete
// @Summary Delete a user by the user identifier
// @Description Delete user by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Success 200 {object} util.Response "{"code":0,"message":"OK","data":null}"
// @Router /user/{id} [delete]
func Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteUser(uint64(userId)); err != nil {
		util.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	util.SendResponse(c, nil, nil)
}

// Update
// @Summary Update a user info by the user identifier
// @Description Update a user by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The user's database id index num"
// @Param user body model.UserModel true "The user info"
// @Success 200 {object} util.Response "{"code":0,"message":"OK","data":null}"
// @Router /user/{id} [put]
func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	userId, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		util.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// We update the record based on the user id.
	u.Id = uint64(userId)

	// Validate the data.
	if err := u.Validate(); err != nil {
		util.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		util.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// Save changed fields.
	if err := u.Update(); err != nil {
		util.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	util.SendResponse(c, nil, nil)
}

// Login
// @Summary Login generates the authentication token
// @Tags common
// @Accept  json
// @Produce  json
// @Param user body LoginRequest true "Login with username and pasword"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MjgwMTY5MjIsImlkIjowLCJuYmYiOjE1MjgwMTY5MjIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LjxrK9DuAwAzUD8-9v43NzWBN7HXsSLfebw92DKd1JQ"}}"
// @Router /login [post]
func Login(c *gin.Context) {
	// Binding the data with the user struct.
	var u LoginRequest
	if err := c.Bind(&u); err != nil {
		util.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// Get the user information by the login username.
	d, err := model.GetUser(u.Username)
	if err != nil {
		util.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// Compare the login password with the user password.
	if err := auth.Compare(d.Password, u.Password); err != nil {
		util.SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: d.Id, Username: d.Username}, "")
	if err != nil {
		util.SendResponse(c, errno.ErrToken, nil)
		return
	}
	util.SendResponse(c, nil, model.Token{Token: t})
}

func listUser(username string, offset, limit int) ([]*model.UserInfo, uint64, error) {
	infos := make([]*model.UserInfo, 0)
	users, count, err := model.ListUser(username, offset, limit)
	if err != nil {
		return nil, count, err
	}

	var ids []uint64
	for _, user := range users {
		ids = append(ids, user.Id)
	}

	wg := sync.WaitGroup{}
	userList := model.UserList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.UserInfo, len(users)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// Improve query efficiency in parallel
	for _, u := range users {
		wg.Add(1)
		go func(u *model.UserModel) {
			defer wg.Done()

			shortId, err := util.GenShortId()
			if err != nil {
				errChan <- err
				return
			}

			userList.Lock.Lock()
			defer userList.Lock.Unlock()
			userList.IdMap[u.Id] = &model.UserInfo{
				Id:        u.Id,
				Username:  u.Username,
				SayHello:  fmt.Sprintf("Hello %s", shortId),
				Password:  u.Password,
				CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	for _, id := range ids {
		infos = append(infos, userList.IdMap[id])
	}

	return infos, count, nil
}
