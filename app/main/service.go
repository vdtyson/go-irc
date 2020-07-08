package main

type Service interface{}
type FirestoreService interface {
	Service
	CollectionPath() string
}

type AuthService struct{}

func (a *AuthService) CreateUser() {

}
func (a *AuthService) Login() {

}
