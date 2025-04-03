package service

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
	"golang.org/x/crypto/bcrypt"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/domain"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/repository"
)

func userToGormModel(u domain.UserEntity) *repository.User {
	id, _ := uuid.Parse(u.ID)
	return &repository.User{
		ID:       id,
		Username: u.Username,
		Password: u.Password,
	}
}

func articleToGormModel(e domain.ArticleEntity) *repository.Article {
	return &repository.Article{
		ID:        e.ID,
		Title:     e.Title,
		Preview:   e.Preview,
		Author:    e.Author,
		Content:   e.Content,
		CreatedAt: e.CreatedAt,
	}
}

func gormToArticle(r repository.Article) domain.ArticleEntity {
	return domain.ArticleEntity{ID: r.ID,
		Title:     r.Title,
		Preview:   r.Preview,
		Author:    r.Author,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
	}
}

func extractInfo(mdText *string) (string, string, error) {

	header, preview := "", ""

	lines := strings.Split(*mdText, "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "# ") && header == "" {
			*mdText = strings.Join(append(lines[:i], lines[i+1:]...), "\n")
			header = line + "\n"
		}
		if strings.HasPrefix(line, "##") && preview == "" {
			*mdText = strings.Join(append(lines[:i], lines[i+1:]...), "\n")
			preview = line + "\n"
		}
	}

	if header == "" {
		return "", "", fmt.Errorf("no header")
	}

	return header, preview, nil
}

func prepareMd(article *domain.ArticleEntity) (domain.ArticleEntity, error) {
	title, preview, err := extractInfo(&article.Content)
	if err != nil {
		return *article, err
	}

	mdTitle, err := converMdToHTML(title)

	if err != nil {
		return *article, err
	}

	mdContent, err := converMdToHTML(article.Content)
	if err != nil {
		return *article, err
	}

	mdPreview, _ := converMdToHTML(preview)
	println(mdPreview)

	article.Title = mdTitle
	article.Preview = mdPreview
	article.Content = mdContent

	return *article, nil
}

func converMdToHTML(mdText string) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(mdText), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func compareHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateUUID() uuid.UUID {
	return uuid.New()
}
