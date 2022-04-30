package service

import (
	"context"

	"github.com/ProjectUnion/project-backend.git/internal/domain"
	"github.com/ProjectUnion/project-backend.git/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userData struct {
	AccessToken  string
	Position     string
	RefreshToken string
	UserID       string
}

type Authorization interface {
	Register(ctx context.Context, inp domain.UserAuth) (userData, error)
	Login(ctx context.Context, email, password string) (userData, error)
	Refresh(ctx context.Context, refreshToken string) (userData, error)
	Logout(ctx context.Context, refreshToken string) error
	ParseToken(token string) (string, error)
}

type User interface {
	GetDataProfile(ctx context.Context, userID primitive.ObjectID) (domain.UserProfile, error)
	GetDataSettings(ctx context.Context, userID primitive.ObjectID) (domain.UserSettings, error)
	SaveData(ctx context.Context, userID primitive.ObjectID, inp domain.UserSettings) error
	ChangePassword(ctx context.Context, userID primitive.ObjectID, inp domain.ChangePassword) error
	DeleteAccount(ctx context.Context, userID primitive.ObjectID) error
}

type Project interface {
	CreateProject(ctx context.Context, inp domain.ProjectData) error
	GetProjects(ctx context.Context, userID primitive.ObjectID) ([]domain.ProjectData, error)
}

type Service struct {
	Authorization
	User
	Project
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthorizationService(repos.Authorization),
		User:          NewUserService(repos.User),
		Project:       NewProjectService(repos.Project),
	}
}
