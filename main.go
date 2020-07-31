package main

import "fmt"

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type IService interface {
	GetUser(string) *User
	PushAlertMessage(string, string) bool
}

type serviceImpl struct{}

func NewService() IService {
	return &serviceImpl{}
}

func (s *serviceImpl) GetUser(userID string) *User {
	return &User{ID: "uuid-001-abcd", Name: "Xuan Vo", Email: "xuan.vo@outlook.com"}
}

func (s *serviceImpl) PushAlertMessage(userID string, message string) bool {
	return true
}

func main() {
	s := NewService()
	usr := s.GetUser("uuid-001-002")
	fmt.Println(usr)
	m := "hi there!"
	fmt.Printf("Pushing message '%s' to %s - status: %t", m, usr.ID, s.PushAlertMessage(usr.ID, m))
}
