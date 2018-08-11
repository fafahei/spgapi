package controllers

import (
	"spgapi/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"regexp"
	"fmt"
	"github.com/astaxie/beego/utils"
)

// Operations about Friends Management
type FriendController struct {
	beego.Controller
}

// @Title CreateFriend
// @Description create friend connection
// @Param	data		body 	models.Friends	true		"body for connection emails"
// @Success 200 {object} models.Status
// @Failure 400 bad request
// @Failure 403 empty argument in request body
// @router /create [post]
func (f *FriendController) CreateConnection() {
	var friends [2]*models.Info
	var emails models.Friends
	var status models.Status
	var err error
	err = json.Unmarshal(f.Ctx.Input.RequestBody, &emails)
	if err != nil {
		f.Ctx.Output.SetStatus(400)
	} else if len(emails.Friends) == 0 || utils.InSlice("", emails.Friends) {
		f.Ctx.Output.SetStatus(403)
	} else if len(emails.Friends) != 2 {
		status.Message = "Number of emails must be two!"
	} else if emails.Friends[0] == emails.Friends[1] {
		status.Message = "Two emails cannot be same!"
	} else {
		for i, email := range emails.Friends {
			friends[i] = &models.Info{0, email}
		}
		if status.Success, err = models.Connect(friends); err != nil {
			status.Message = err.Error()
		} else {
			status.Message = "New friend Connection has been created successfully!"
		}
	}
	f.Data["json"] = status
	f.ServeJSON()
}

// @Title ListFriends
// @Description list friends for email address
// @Param	data		body	models.Email	true		"body for email"
// @Success 200 {object} models.Response
// @Failure 400 bad request
// @Failure 403 empty argument in request body
// @router /list [post]
func (f *FriendController) ListFriends() {
	var email models.Email
	var response models.Response
	var err error
	err = json.Unmarshal(f.Ctx.Input.RequestBody, &email)
	if err != nil {
		f.Ctx.Output.SetStatus(400)
	} else if email.Email == "" {
		f.Ctx.Output.SetStatus(403)
	} else {
		friend := &models.Info{0, email.Email}
		if response.Status.Success, response.Friends, err = models.GetFriends(friend); err != nil {
			response.Status.Message = err.Error()
		} else {
			response.Count = len(response.Friends)
			response.Status.Message = "Get friend list from " + friend.Email + " successfully!"
		}
	}
	f.Data["json"] = response
	f.ServeJSON()
}

// @Title ListCommonFriends
// @Description list common friends for two email addresses
// @Param	data		body 	models.Friends	true		"body for two emails"
// @Success 200 {object} models.Response
// @Failure 400 bad request
// @Failure 403 empty argument in request body
// @router /listCommon [post]
func (f *FriendController) ListCommonFriends() {
	var friends [2]*models.Info
	var emails models.Friends
	var response models.Response
	var err error
	err = json.Unmarshal(f.Ctx.Input.RequestBody, &emails)
	if err != nil {
		f.Ctx.Output.SetStatus(400)
	} else if len(emails.Friends) == 0 || utils.InSlice("", emails.Friends) {
		f.Ctx.Output.SetStatus(403)
	} else if len(emails.Friends) != 2 {
		response.Status.Message = "Number of emails must be two!"
	} else if emails.Friends[0] == emails.Friends[1] {
		response.Status.Message = "Two emails cannot be same!"
	} else {
		for i, email := range emails.Friends {
			friends[i] = &models.Info{0, email}
		}
		if response.Status.Success, response.Friends, err = models.GetCommonFriends(friends); err != nil {
			response.Status.Message = err.Error()
		} else {
			response.Count = len(response.Friends)
			response.Status.Message = "Get common friend list from " + friends[0].Email + " and " + friends[1].Email + " successfully!"
		}
	}
	f.Data["json"] = response
	f.ServeJSON()
}

// @Title SubscribeUpdate
// @Description subscribe updates from requestor email
// @Param	data		body 	models.Request	true		"body for request content"
// @Success 200 {object} models.Status
// @Failure 400 bad request
// @Failure 403 empty argument in request body
// @router /subscribe [post]
func (f *FriendController) SubscribeUpdate() {
	var request models.Request
	var status models.Status
	var err error
	err = json.Unmarshal(f.Ctx.Input.RequestBody, &request)
	fmt.Println(request)
	if err != nil {
		f.Ctx.Output.SetStatus(400)
	} else if request.Requestor == "" || request.Target == "" {
		f.Ctx.Output.SetStatus(403)
	} else if request.Requestor == request.Target {
		status.Message = "Requestor and target emails cannot be same!"
	} else {
		requestor := &models.Info{0, request.Requestor}
		target := &models.Info{0, request.Target}
		r := &models.Relationship{0, requestor, target, 0}
		if status.Success, err = models.Subscribe(r); err != nil {
			status.Message = err.Error()
		} else {
			status.Message = r.Requestor.Email + " subscribes to updates from " + r.Target.Email + " successfully!"
		}
	}
	f.Data["json"] = status
	f.ServeJSON()
}

// @Title BlockUpdate
// @Description block updates from requestor email
// @Param	data		body 	models.Request	true		"body for request content"
// @Success 200 {object} models.Status
// @Failure 400 bad request
// @Failure 403 empty argument in request body
// @router /block [post]
func (f *FriendController) BlockUpdate() {
	var request models.Request
	var status models.Status
	var err error
	err = json.Unmarshal(f.Ctx.Input.RequestBody, &request)
	if err != nil {
		f.Ctx.Output.SetStatus(400)
	} else if request.Requestor == "" || request.Target == "" {
		f.Ctx.Output.SetStatus(403)
	} else if request.Requestor == request.Target {
		status.Message = "Requestor and target emails cannot be same!"
	} else {
		requestor := &models.Info{0, request.Requestor}
		target := &models.Info{0, request.Target}
		r := &models.Relationship{0, requestor, target, 0}
		if status.Success, err = models.Block(r); err != nil {
			status.Message = err.Error()
		} else {
			status.Message = r.Requestor.Email + " blocks updates from " + r.Target.Email + " successfully!"
		}
	}
	f.Data["json"] = status
	f.ServeJSON()
}

// @Title ListUpdate
// @Description list updates from requestor email
// @Param	data		body 	models.RequestUpdate	true		"body for request content"
// @Success 200 {object} models.ResponseUpdate
// @Failure 400 bad request
// @Failure 403 empty argument in request body
// @router /listUpdate [post]
func (f *FriendController) ListUpdate() {
	var request models.RequestUpdate
	var response models.ResponseUpdate
	var err error
	err = json.Unmarshal(f.Ctx.Input.RequestBody, &request)
	if err != nil {
		f.Ctx.Output.SetStatus(400)
	} else if request.Sender == "" || request.Text == "" {
		f.Ctx.Output.SetStatus(403)
	} else {
		target := &models.Info{0, request.Sender}
		m := regexp.MustCompile("[a-z0-9-]{1,30}@[a-z0-9-]{1,65}.[a-z]{1,}").FindAllString(request.Text, -1)
		requestor := make([]*models.Info, len(m))
		for i:=0;i<len(m);i++ {
			requestor[i] = &models.Info{0, m[i]}
		}
		if response.Status.Success, response.Recipients, err = models.GetUpdateFriends(target, requestor); err != nil {
			response.Status.Message = err.Error()
		} else {
			response.Status.Message = "Get update friend list from " + request.Sender + " successfully!"
		}
	}
	f.Data["json"] = response
	f.ServeJSON()
}
