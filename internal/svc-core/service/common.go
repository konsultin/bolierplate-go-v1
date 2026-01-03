package service

import (
	"context"

	"github.com/konsultin/project-goes-here/config"
	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/project-goes-here/internal/svc-core/repository"
	"github.com/konsultin/project-goes-here/libs/logk"
	logkOption "github.com/konsultin/project-goes-here/libs/logk/option"
)

type Service struct {
	Repo    *repository.Repository
	Log     logk.Logger
	Ctx     context.Context
	Config  *config.Config
	Subject *model.Subject
}

func (s *Service) WithSubject(subject *dto.Subject) *Service {
	newS := *s
	newS.Subject = model.NewSubject(subject)
	return &newS
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) Close() {
	// Returns connection to pool
	err := s.Repo.Close()
	if err != nil {
		s.Log.Error("Failed to close connection", logkOption.Error(err))
	} else {
		s.Log.Tracef("DB: Connection returned to pool")
	}
}
