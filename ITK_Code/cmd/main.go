package main

import "fmt"

// Базовый интерфейс
type User interface {
	GetUsername() string
	HasPermission(permission string) bool
	GetRole() string
}

type BasicUser struct {
	username   string
	role       string
	permission map[string]struct{}
}

type Moderator struct {
	BasicUser
}

type Admin struct {
	Moderator
}

func NewBasicUser(username string) BasicUser {
	return BasicUser{
		username: username,
		role:     "User",
		permission: map[string]struct{}{
			"read": {},
		},
	}
}

func (b BasicUser) GetUsername() string {
	return b.username
}

func (b BasicUser) HasPermission(permission string) bool {
	_, ok := b.permission[permission]
	return ok
}

func (b BasicUser) GetRole() string {
	return b.role
}

func NewModerator(username string) Moderator {
	moderator := NewBasicUser(username)

	moderator.role = "Moderator"
	moderator.permission["edit"] = struct{}{}
	moderator.permission["ban_user"] = struct{}{}

	return Moderator{
		BasicUser: moderator,
	}
}

func NewAdmin(username string) Admin {
	admin := NewModerator(username)

	admin.BasicUser.role = "admin"
	admin.BasicUser.permission["delete"] = struct{}{}
	admin.BasicUser.permission["manage_roles"] = struct{}{}

	return Admin{
		Moderator: admin,
	}
}

func main() {
	users := []User{
		NewBasicUser("user1"),
		NewModerator("mod1"),
		NewAdmin("admin1"),
	}

	permissionsToCheck := []string{
		"read",
		"edit",
		"ban_user",
		"delete",
		"manage_roles",
	}

	for _, u := range users {
		fmt.Println("------")
		fmt.Println("Username:", u.GetUsername())
		fmt.Println("Role:", u.GetRole())

		for _, p := range permissionsToCheck {
			fmt.Printf("Has %-15s: %t\n", p, u.HasPermission(p))
		}
	}
}
