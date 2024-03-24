package service

import (
	"context"
	"errors"
	"time"

	"final_project/internal/model"
	"final_project/internal/repository"
)

type CommentService interface {
	GetCommentByID(ctx context.Context, id uint64) (model.GetCommentByID, error)
	DeleteCommentByID(ctx context.Context, id uint64) (model.UpdateComment, error)
	CreateComment(ctx context.Context, comment model.CreateCommentInput, user uint64) (model.CreateComment, error)
	UpdateComment(ctx context.Context, id uint64, comment model.UpdateCommentInput) (model.UpdateComment, error)
	GetCommentsByPhotoID(ctx context.Context, photoID uint64) ([]model.Comment, error)
	GetCommentByID1(ctx context.Context, id uint64) (model.UpdateComment, error)
	GetComments(ctx context.Context) ([]model.GetCommentByID, error)
}

type commentServiceImpl struct {
	repoComment repository.CommentQuery
	repoUser    repository.UserQuery
	repoPhoto   repository.PhotoQuery
}

func NewCommentService(repoComment repository.CommentQuery, repoUser repository.UserQuery, repoPhoto repository.PhotoQuery) CommentService {
	return &commentServiceImpl{
		repoComment: repoComment,
		repoUser:    repoUser,
		repoPhoto:   repoPhoto,

	}
}

func (c *commentServiceImpl) GetCommentByID(ctx context.Context, id uint64) (model.GetCommentByID, error) {
	comment, err := c.repoComment.GetCommentByID(ctx, id)
	if err != nil {
		return model.GetCommentByID{}, err
	}
	user, err := c.repoUser.GetUsersByID(ctx, comment.UserID)
	if err != nil {
		return model.GetCommentByID{}, err
	}
	
	// Menetapkan data pengguna ke dalam struktur data foto
	comment.User.ID = user.ID
	comment.User.Email = user.Email
	comment.User.Username = user.Username

	photo, err := c.repoPhoto.GetPhotoByID(ctx, comment.PhotoID)
	if err != nil {
		return model.GetCommentByID{}, err
	}
	comment.Photo.ID = photo.ID
	comment.Photo.Title = photo.Title
	comment.Photo.Caption = photo.Caption
	comment.Photo.PhotoURL = photo.PhotoURL
	comment.Photo.UserID = photo.UserID
    
	return comment, err
}
func (c *commentServiceImpl) GetCommentByID1(ctx context.Context, id uint64) (model.UpdateComment, error) {
	comment, err := c.repoComment.GetCommentByID1(ctx, id)
	if err != nil {
		return model.UpdateComment{}, err
	}
	return comment, err
}
func (c *commentServiceImpl) GetComments(ctx context.Context) ([]model.GetCommentByID, error) {
    // Mengambil semua data foto dari repocommentsitory
    comments, err := c.repoComment.GetComments(ctx)
    if err != nil {
        return nil, err
    }
    
    // Untuk setiap foto, ambil data pengguna yang sesuai
    for i, comment := range comments {
        user, err := c.repoUser.GetUsersByID(ctx, comment.UserID)
		
        if err != nil {
            return nil, err
        }
        
        // Menetapkan data pengguna ke dalam struktur data foto
        comments[i].User.Email = user.Email
        comments[i].User.Username = user.Username
		comments[i].User.ID = user.ID

		photo, err := c.repoPhoto.GetPhotoByID(ctx, comment.PhotoID)
		if err != nil {
			return nil, err
		}
		comments[i].Photo.ID = photo.ID
		comments[i].Photo.Title = photo.Title
		comments[i].Photo.Caption = photo.Caption
		comments[i].Photo.PhotoURL = photo.PhotoURL
		comments[i].Photo.UserID = photo.UserID
    }
	
    
    return comments, nil
}
func (c *commentServiceImpl) DeleteCommentByID(ctx context.Context, id uint64) (model.UpdateComment, error) {
	comment, err := c.repoComment.GetCommentByID1(ctx, id)
	if err != nil {
		return model.UpdateComment{}, err
	}
	// if comment doesn't exist, return
	if comment.ID == 0 {
		return model.UpdateComment{}, nil
	}

	// delete comment by id
	err = c.repoComment.DeleteCommentByID(ctx, id)
	if err != nil {
		return model.UpdateComment{}, err
	}

	return comment, err
}

func (c *commentServiceImpl) CreateComment(ctx context.Context, CreateComment model.CreateCommentInput, userID uint64) (model.CreateComment, error) {
	comment := model.CreateComment{
		Message: CreateComment.Message,
		PhotoID: CreateComment.PhotoID,
		UserID:  userID,
		CreatedAt: time.Now(),
	}
	commentPhotoID, _ := c.repoComment.GetCommentsByPhotoID(ctx, comment.PhotoID)
	if len(commentPhotoID) == 0 {
		return model.CreateComment{}, errors.New("photo not found")
	}
	createdComment, err := c.repoComment.CreateComment(ctx, comment)
	if err != nil {
		return model.CreateComment{}, err
	}
	return createdComment, nil
}

func (c *commentServiceImpl) UpdateComment(ctx context.Context, id uint64, comment model.UpdateCommentInput) (model.UpdateComment, error) {
	updateComment := model.UpdateComment{
		Message: comment.Message,
		UpdatedAt: time.Now(),

	}
	updatedComment, err := c.repoComment.UpdateComment(ctx, id, updateComment)
	if err != nil {
		return model.UpdateComment{}, err
	}
	return updatedComment, nil
}
func (c *commentServiceImpl) GetCommentsByPhotoID(ctx context.Context, photoID uint64) ([]model.Comment, error) {
    comments, err := c.repoComment.GetCommentsByPhotoID(ctx, photoID)
    if err != nil {
        return nil, err
    }
	for i, comment := range comments {
        user, err := c.repoUser.GetUsersByID(ctx, comment.UserID)
		
        if err != nil {
            return nil, err
        }
        
        // Menetapkan data pengguna ke dalam struktur data foto
        comments[i].User.Email = user.Email
        comments[i].User.Username = user.Username
		comments[i].User.ID = user.ID

		photo, err := c.repoPhoto.GetPhotoByID(ctx, comment.PhotoID)
		if err != nil {
			return nil, err
		}
		comments[i].Photo.ID = photo.ID
		comments[i].Photo.Title = photo.Title
		comments[i].Photo.Caption = photo.Caption
		comments[i].Photo.PhotoURL = photo.PhotoURL
		comments[i].Photo.UserID = photo.UserID
    }
    return comments, nil
}