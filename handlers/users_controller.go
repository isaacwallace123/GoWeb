package handlers

import (
	"github.com/isaacwallace123/GoWeb/handlers/Models"

	"github.com/isaacwallace123/GoWeb/core"
	"github.com/isaacwallace123/GoWeb/httpstatus"
	"github.com/isaacwallace123/GoWeb/response"
)

type UsersController struct{}

func (c *UsersController) BasePath() string {
	return "/api/v1/users"
}

func (c *UsersController) Routes() []core.RouteEntry {
	return []core.RouteEntry{
		{Method: "GET", Path: "/", Handler: "GetAll"},
		{Method: "GET", Path: "/{userid}", Handler: "Get"},
		{Method: "POST", Path: "/", Handler: "Post"},
		{Method: "PUT", Path: "/{userid}", Handler: "Put"},
		{Method: "DELETE", Path: "/{userid}", Handler: "Delete"},
	}
}

// key: user ID (int), value: name/email map
var fakeDB = map[int]map[string]string{}

func (c *UsersController) Get(userid int) *response.ResponseEntity {
	user, exists := fakeDB[userid]
	if !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	return response.Status(httpstatus.OK).Body(Models.UserResponse{
		Id:    userid,
		Name:  user["name"],
		Email: user["email"],
	})
}

func (c *UsersController) GetAll() *response.ResponseEntity {
	users := []Models.UserResponse{}
	for id, user := range fakeDB {
		users = append(users, Models.UserResponse{
			Id:    id,
			Name:  user["name"],
			Email: user["email"],
		})
	}
	return response.Status(httpstatus.OK).Body(users)
}

func (c *UsersController) Post(req Models.UserRequest) *response.ResponseEntity {
	if req.Name == "" || req.Email == "" {
		return response.Status(httpstatus.BAD_REQUEST).
			Body(map[string]string{"error": "Name and Email are required"})
	}

	newId := len(fakeDB) + 1

	fakeDB[newId] = map[string]string{
		"name":  req.Name,
		"email": req.Email,
	}

	return response.Status(httpstatus.CREATED).
		Body(map[string]any{"message": "User created", "id": newId})
}

func (c *UsersController) Put(userid int, req Models.UserRequest) *response.ResponseEntity {
	user, exists := fakeDB[userid]
	if !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	user["name"] = req.Name
	user["email"] = req.Email
	fakeDB[userid] = user

	return response.Status(httpstatus.OK).
		Body(map[string]string{"message": "User updated"})
}

func (c *UsersController) Delete(userid int) *response.ResponseEntity {
	if _, exists := fakeDB[userid]; !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	delete(fakeDB, userid)
	return response.Status(httpstatus.NO_CONTENT)
}
