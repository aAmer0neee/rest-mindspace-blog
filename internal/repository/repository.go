package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/aAmer0neee/rest-mindspace-blog/internal/config"
)

type Repository struct {
	Db *gorm.DB
}

func ConnectRepository(cfg config.Cfg) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Repository.Host,
		cfg.Repository.User,
		cfg.Repository.Password,
		cfg.Repository.Name,
		cfg.Repository.Port,
		cfg.Repository.Sslmode)

	r, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	fmt.Printf("[Repository] [INFO] Open Data Base %s\n", r.Name())

	if cfg.Repository.Migrate {
		if err = r.AutoMigrate(&Article{}, &User{}); err != nil {
			return nil, err
		}
		fmt.Printf("[Repository] [INFO] migrate %s\n", r.Name())
	}
	return &Repository{Db: r}, nil
}

func (r *Repository) Post(model interface{}) error {
	return r.Db.Create(model).Error
}

func (r *Repository) GetRecord(id int) (Article, error) {
	dst := Article{}
	if err := r.Db.Find(&dst, id).Error; err != nil {
		return Article{}, err
	}
	return dst, nil
}

func (r *Repository) GetRecords(page, limit int) ([]Article, error) {
	offset := (page - 1) * limit
	dst := []Article{}
	err := r.Db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&dst).Error
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *Repository) GetUserPassword(username string) (string, uuid.UUID, error) {
	user := User{}
	return user.Password, user.ID, r.Db.Model(&User{}).
		Select("Password", "id").Where("username = ?", username).
		Take(&user).Error
}

func (r *Repository) RecordsCount() (int, error) {
	var count int64
	return int(count), r.Db.Model(&Article{}).Count(&count).Error
}
