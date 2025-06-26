package handlers

import (
	"context"
	"fmt"
	"github.com/isaacwallace123/GoWeb/httpstatus"
	"github.com/isaacwallace123/GoWeb/response"
)

type UsersController struct{}

// Replaced @RestController, @Controller anotations
func (c *UsersController) Path() string {
	return "/users/{id}"
}

var fakeDB = map[string]map[string]string{} // key: id, value: {"name": ..., "email": ...}

func (c *UsersController) Get(ctx context.Context, id string) *response.ResponseEntity {
	user, exists := fakeDB[id]

	fmt.Println(user, exists, id)

	if !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	return response.Status(httpstatus.OK).Body(user)
}

func (c *UsersController) Post(ctx context.Context, id string, req CreateUserRequest) *response.ResponseEntity {
	if req.Name == "" || req.Email == "" {
		return response.Status(httpstatus.BAD_REQUEST).
			Body(map[string]string{"error": "Name and Email are required"})
	}

	fakeDB[id] = map[string]string{"name": req.Name, "email": req.Email}

	return response.Status(httpstatus.CREATED).Body(map[string]string{"message": "User created"})
}

func (c *UsersController) Put(ctx context.Context, id string, req UpdateUserRequest) *response.ResponseEntity {
	user, exists := fakeDB[id]

	if !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	user["name"] = req.Name
	fakeDB[id] = user

	return response.Status(httpstatus.OK).Body(map[string]string{"message": "User updated"})
}

func (c *UsersController) Delete(ctx context.Context, id string) *response.ResponseEntity {
	if _, exists := fakeDB[id]; !exists {
		return response.Status(httpstatus.NOT_FOUND).
			Body(map[string]string{"error": "User not found"})
	}

	delete(fakeDB, id)

	return response.Status(httpstatus.NO_CONTENT)
}
