package repository

import (
	"context"

	"final_project/internal/infrastructure"
	"final_project/internal/model"

	"gorm.io/gorm"
)

type SocialMediaQuery interface {
	GetSocialMediaByID(ctx context.Context, id uint64) (model.SocialMedia, error)
	DeleteSocialMediaByID(ctx context.Context, id uint64) error
	CreateSocialMedia(ctx context.Context, socialMedia model.CreateSocialMedia) (model.CreateSocialMedia, error)
	UpdateSocialMedia(ctx context.Context, id uint64, socialMedia model.UpdateSocialMedia) (model.UpdateSocialMedia, error)
	GetSocialMediaByID1(ctx context.Context, id uint64) (model.UpdateSocialMedia, error)
	GetSocialMediasByUserID(ctx context.Context, userID uint64) ([]model.SocialMedia, error)
	GetSocialMedias(ctx context.Context) ([]model.SocialMedia, error)
}	

type socialMediaQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewSocialMediaQuery(db infrastructure.GormPostgres) SocialMediaQuery {
	return &socialMediaQueryImpl{db: db}
}

func (s *socialMediaQueryImpl) GetSocialMediaByID(ctx context.Context, id uint64) (model.SocialMedia, error) {
	db := s.db.GetConnection()
	socialMedia := model.SocialMedia{}
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Where("id = ?", id).
		First(&socialMedia).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.SocialMedia{}, nil
		}
		return model.SocialMedia{}, err
	}
	return socialMedia, nil
}
func (c *socialMediaQueryImpl) GetSocialMedias(ctx context.Context) ([]model.SocialMedia, error) {
	db := c.db.GetConnection()
	socialMedias := []model.SocialMedia{}
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Find(&socialMedias).Error; err != nil {
		return nil, err
	}
	return socialMedias, nil
}
func (c *socialMediaQueryImpl) GetSocialMediasByUserID(ctx context.Context, userID uint64) ([]model.SocialMedia, error) {
	db := c.db.GetConnection()
	socialMedias := []model.SocialMedia{}
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Where("user_id = ?", userID).
		Find(&socialMedias).Error; err != nil {
		return nil, err
	}
	return socialMedias, nil
}
func (c *socialMediaQueryImpl) DeleteSocialMediaByID(ctx context.Context, id uint64) error {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Delete(&model.UpdateSocialMedia{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
func (c *socialMediaQueryImpl) GetSocialMediaByID1(ctx context.Context, id uint64) (model.UpdateSocialMedia, error) {
	db := c.db.GetConnection()
	socialMedia := model.UpdateSocialMedia{}
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Where("id = ?", id).
		First(&socialMedia).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UpdateSocialMedia{}, nil
		}
		return model.UpdateSocialMedia{}, err
	}
	return socialMedia, nil
}
func (c *socialMediaQueryImpl) CreateSocialMedia(ctx context.Context, socialMedia model.CreateSocialMedia) (model.CreateSocialMedia, error) {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Save(&socialMedia).Error; err != nil {
		return model.CreateSocialMedia{}, err
	}
	return socialMedia, nil
}

func (c *socialMediaQueryImpl) UpdateSocialMedia(ctx context.Context, id uint64, socialMedia model.UpdateSocialMedia) (model.UpdateSocialMedia, error) {
	db := c.db.GetConnection()
	updatedSocialMedia := model.UpdateSocialMedia{}
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Where("id = ?", id).
		Updates(&socialMedia).
		First(&updatedSocialMedia).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UpdateSocialMedia{}, nil
		}
		return model.UpdateSocialMedia{}, err
	}
	return updatedSocialMedia, nil
}
