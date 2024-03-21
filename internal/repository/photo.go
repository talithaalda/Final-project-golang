package repository

import (
	"context"

	"final_project/internal/infrastructure"
	"final_project/internal/model"

	"gorm.io/gorm"
)

type PhotoQuery interface {
	GetPhotos(ctx context.Context) ([]model.Photo, error)
	GetPhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error)
	DeletePhotoByID(ctx context.Context, id uint64) error
	CreatePhoto(ctx context.Context, photo model.CreatePhoto) (model.CreatePhoto, error)
	EditPhoto(ctx context.Context, id uint64, user model.UpdatePhoto) (model.UpdatePhoto, error)
	GetPhotoByUserID(ctx context.Context, photoID uint64) ([]model.GetPhoto, error)
}

type photoQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewPhotoQuery(db infrastructure.GormPostgres) PhotoQuery {
	return &photoQueryImpl{db: db}
}

func (p *photoQueryImpl) GetPhotos(ctx context.Context) ([]model.Photo, error) {
	db := p.db.GetConnection()
	photos := []model.Photo{}
	if err := db.
		WithContext(ctx).
		Table("photos").
		Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (p *photoQueryImpl) GetPhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error) {
	db := p.db.GetConnection()
	photo := model.UpdatePhoto{}
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", id).
		Find(&photo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UpdatePhoto{}, nil
		}
		return model.UpdatePhoto{}, err
	}
	return photo, nil
}

func (p *photoQueryImpl) DeletePhotoByID(ctx context.Context, id uint64) error {
	db := p.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("photos").
		Delete(&model.UpdatePhoto{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func (p *photoQueryImpl) CreatePhoto(ctx context.Context, photo model.CreatePhoto) (model.CreatePhoto, error) {
	db := p.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("photos").
		Save(&photo).Error; err != nil {
		return model.CreatePhoto{}, err
	}
	return photo, nil
}
func (u *photoQueryImpl) EditPhoto(ctx context.Context, id uint64, user model.UpdatePhoto) (model.UpdatePhoto, error) {
	db := u.db.GetConnection()
	updatedPhoto := model.UpdatePhoto{}
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", id).Updates(&user).First(&updatedPhoto).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return model.UpdatePhoto{}, nil
			}
		}
	return updatedPhoto, nil
}
func (p *photoQueryImpl) GetPhotoByUserID(ctx context.Context, userID uint64) ([]model.GetPhoto, error) {
	db := p.db.GetConnection()
	photos := []model.GetPhoto{}
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("user_id = ?", userID).
		Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}