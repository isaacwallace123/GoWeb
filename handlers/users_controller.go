package handlers

import (
	"fmt"
	"github.com/isaacwallace123/GoWeb/httpstatus"
	"github.com/isaacwallace123/GoWeb/response"
)

type UsersController struct{}

// Replaced @RestController, @Controller anotations
func (controller *UsersController) Path() string {
	return "/users/{id}"
}

var fakeDB = map[string]map[string]string{} // key: id, value: {"name": ..., "email": ...}

func (controller *UsersController) Get(id string) *response.ResponseEntity {
	user, exists := fakeDB[id]

	fmt.Println(user, exists, id)

	if !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	return response.Status(httpstatus.OK).Body(user)
}

func (controller *UsersController) Post(id string, req CreateUserRequest) *response.ResponseEntity {
	if req.Name == "" || req.Email == "" {
		return response.Status(httpstatus.BAD_REQUEST).
			Body(map[string]string{"error": "Name and Email are required"})
	}

	fakeDB[id] = map[string]string{"name": req.Name, "email": req.Email}

	return response.Status(httpstatus.CREATED).Body(map[string]string{"message": "User created"})
}

func (controller *UsersController) Put(id string, req UpdateUserRequest) *response.ResponseEntity {
	user, exists := fakeDB[id]

	if !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	user["name"] = req.Name
	fakeDB[id] = user

	return response.Status(httpstatus.OK).Body(map[string]string{"message": "User updated"})
}

func (controller *UsersController) Delete(id string) *response.ResponseEntity {
	if _, exists := fakeDB[id]; !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	delete(fakeDB, id)

	return response.Status(httpstatus.NO_CONTENT)
}
