package service

import (
	"log/slog"

	"github.com/aAmer0neee/rest-mindspace-blog/internal/domain"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/repository"
)

type Service struct {
	R      *repository.Repository
	Logger *slog.Logger
}

func ConfigureService(r *repository.Repository, l *slog.Logger) *Service {
	return &Service{R: r, Logger: l}
}

func (s *Service) NewPost(article domain.ArticleEntity) error {

	article, err := prepareMd(&article)
	if err != nil {
		s.Logger.Info("error parse md", "message", err.Error())
		return err
	}

	if err := s.R.Post(articleToGormModel(article)); err != nil {
		s.Logger.Warn("error posting", "error", err.Error())
		return err
	}
	return nil
}

func (s *Service) RegisterUser(user domain.UserEntity) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		s.Logger.Warn("error hashing password", "error", err.Error())
		return err
	}

	user.Password = hash

	if user.ID == "" {
		user.ID = GenerateUUID().String()
	}

	s.Logger.Info("user", slog.Any("user", user))
	if err := s.R.Post(userToGormModel(user)); err != nil {
		s.Logger.Warn("error posting", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) AuthUser(user domain.UserEntity) (string, error) {
	hash, id, err := s.R.GetUserPassword(user.Username)
	if err != nil {
		return "", err
	}
	err = compareHash(hash, user.Password)
	if err != nil {
		s.Logger.Warn("error compare hash password", "error", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (s *Service) GetFeed(page, limit int) ([]domain.ArticleEntity, int, error) {
	records, err := s.R.GetRecords(page, limit)
	if err != nil {
		s.Logger.Warn("error getting record", "error", err.Error())
		return nil, 0, err
	}

	e := []domain.ArticleEntity{}
	for _, record := range records {
		e = append(e, gormToArticle(record))
	}
	count, err := s.R.RecordsCount()

	if err != nil {
		s.Logger.Warn("error getting records count", "error", err.Error())
		return nil, 0, err
	}

	return e, count, nil
}

func (s *Service) GetArticle(id int) (domain.ArticleEntity, error) {
	record, err := s.R.GetRecord(id)
	if err != nil {
		s.Logger.Warn("error posting", "error", err.Error())
		return domain.ArticleEntity{}, err
	}
	return gormToArticle(record), nil
}
