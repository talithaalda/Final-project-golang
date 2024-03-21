package repository

import (
	"context"

	"final_project/internal/infrastructure"
	"final_project/internal/model"

	"gorm.io/gorm"
)

type CommentQuery interface {
	GetCommentByID(ctx context.Context, id uint64) (model.GetCommentByID, error)
	DeleteCommentByID(ctx context.Context, id uint64) error
	CreateComment(ctx context.Context, comment model.CreateComment) (model.CreateComment, error)
	UpdateComment(ctx context.Context, id uint64, comment model.UpdateComment) (model.UpdateComment, error)
	GetCommentByID1(ctx context.Context, id uint64) (model.UpdateComment, error)
	GetCommentsByPhotoID(ctx context.Context, photoID uint64) ([]model.Comment, error)
	GetComments(ctx context.Context) ([]model.GetCommentByID, error)
}	

type commentQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewCommentQuery(db infrastructure.GormPostgres) CommentQuery {
	return &commentQueryImpl{db: db}
}

func (c *commentQueryImpl) GetCommentByID(ctx context.Context, id uint64) (model.GetCommentByID, error) {
	db := c.db.GetConnection()
	comment := model.GetCommentByID{}
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.GetCommentByID{}, nil
		}
		return model.GetCommentByID{}, err
	}
	return comment, nil
}
func (c *commentQueryImpl) GetComments(ctx context.Context) ([]model.GetCommentByID, error) {
	db := c.db.GetConnection()
	comments := []model.GetCommentByID{}
	if err := db.
		WithContext(ctx).
		Table("comments").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
func (c *commentQueryImpl) GetCommentsByPhotoID(ctx context.Context, photoID uint64) ([]model.Comment, error) {
	db := c.db.GetConnection()
	comments := []model.Comment{}
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("photo_id = ?", photoID).
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
func (c *commentQueryImpl) DeleteCommentByID(ctx context.Context, id uint64) error {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("comments").
		Delete(&model.UpdateComment{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
func (c *commentQueryImpl) GetCommentByID1(ctx context.Context, id uint64) (model.UpdateComment, error) {
	db := c.db.GetConnection()
	comment := model.UpdateComment{}
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UpdateComment{}, nil
		}
		return model.UpdateComment{}, err
	}
	return comment, nil
}
func (c *commentQueryImpl) CreateComment(ctx context.Context, comment model.CreateComment) (model.CreateComment, error) {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("comments").
		Save(&comment).Error; err != nil {
		return model.CreateComment{}, err
	}
	return comment, nil
}

func (c *commentQueryImpl) UpdateComment(ctx context.Context, id uint64, comment model.UpdateComment) (model.UpdateComment, error) {
	db := c.db.GetConnection()
	updatedComment := model.UpdateComment{}
	if err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", id).
		Updates(&comment).
		First(&updatedComment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UpdateComment{}, nil
		}
		return model.UpdateComment{}, err
	}
	return updatedComment, nil
}
