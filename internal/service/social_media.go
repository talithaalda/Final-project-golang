package service

import (
	"context"

	"final_project/internal/model"
	"final_project/internal/repository"
)

type SocialMediaService interface {
	GetSocialMediaByID(ctx context.Context, id uint64) (model.SocialMedia, error)
	DeleteSocialMediaByID(ctx context.Context, id uint64) (model.UpdateSocialMedia, error)
	CreateSocialMedia(ctx context.Context, socialMedia model.CreateSocialMedia, user uint64) (model.CreateSocialMedia, error)
	UpdateSocialMedia(ctx context.Context, id uint64, socialMedia model.UpdateSocialMedia) (model.UpdateSocialMedia, error)
	GetSocialMediasByUserID(ctx context.Context, userID uint64) ([]model.SocialMedia, error)
	GetSocialMediaByID1(ctx context.Context, id uint64) (model.UpdateSocialMedia, error)
	GetSocialMedias(ctx context.Context) ([]model.SocialMedia, error)
}

type socialMediaServiceImpl struct {
	repoSocialMedia repository.SocialMediaQuery
	repoUser    repository.UserQuery
}

func NewSocialMediaService(repoSocialMedia repository.SocialMediaQuery, repoUser repository.UserQuery) SocialMediaService {
	return &socialMediaServiceImpl{
		repoSocialMedia: repoSocialMedia,
		repoUser:    repoUser,
	}
}

func (c *socialMediaServiceImpl) GetSocialMediaByID(ctx context.Context, id uint64) (model.SocialMedia, error) {
	socialMedia, err := c.repoSocialMedia.GetSocialMediaByID(ctx, id)
	if err != nil {
		return model.SocialMedia{}, err
	}
	user, err := c.repoUser.GetUsersByID(ctx, socialMedia.UserID)
	if err != nil {
		return model.SocialMedia{}, err
	}
	
	// Menetapkan data pengguna ke dalam struktur data foto
	socialMedia.User.ID = user.ID
	socialMedia.User.Email = user.Email
	socialMedia.User.Username = user.Username
    
	return socialMedia, err
}
func (c *socialMediaServiceImpl) GetSocialMediaByID1(ctx context.Context, id uint64) (model.UpdateSocialMedia, error) {
	socialMedia, err := c.repoSocialMedia.GetSocialMediaByID1(ctx, id)
	if err != nil {
		return model.UpdateSocialMedia{}, err
	}
	return socialMedia, err
}
func (c *socialMediaServiceImpl) GetSocialMedias(ctx context.Context) ([]model.SocialMedia, error) {
    // Mengambil semua data foto dari reposocialMediasitory
    socialMedias, err := c.repoSocialMedia.GetSocialMedias(ctx)
    if err != nil {
        return nil, err
    }
    
    // Untuk setiap foto, ambil data pengguna yang sesuai
    for i, socialMedia := range socialMedias {
        user, err := c.repoUser.GetUsersByID(ctx, socialMedia.UserID)
		
        if err != nil {
            return nil, err
        }
        
        // Menetapkan data pengguna ke dalam struktur data foto
        socialMedias[i].User.Email = user.Email
        socialMedias[i].User.Username = user.Username
		socialMedias[i].User.ID = user.ID

    }
    return socialMedias, nil
}
func (c *socialMediaServiceImpl) DeleteSocialMediaByID(ctx context.Context, id uint64) (model.UpdateSocialMedia, error) {
	socialMedia, err := c.repoSocialMedia.GetSocialMediaByID1(ctx, id)
	if err != nil {
		return model.UpdateSocialMedia{}, err
	}
	// if socialMedia doesn't exist, return
	if socialMedia.ID == 0 {
		return model.UpdateSocialMedia{}, nil
	}

	// delete socialMedia by id
	err = c.repoSocialMedia.DeleteSocialMediaByID(ctx, id)
	if err != nil {
		return model.UpdateSocialMedia{}, err
	}

	return socialMedia, err
}

func (c *socialMediaServiceImpl) CreateSocialMedia(ctx context.Context, CreateSocialMedia model.CreateSocialMedia, userID uint64) (model.CreateSocialMedia, error) {
	socialMedia := model.CreateSocialMedia{
		Name: CreateSocialMedia.Name,
		SocialMediaURL : CreateSocialMedia.SocialMediaURL ,
		UserID:  userID,
	}
	createdSocialMedia, err := c.repoSocialMedia.CreateSocialMedia(ctx, socialMedia)
	if err != nil {
		return model.CreateSocialMedia{}, err
	}
	return createdSocialMedia, nil
}

func (c *socialMediaServiceImpl) UpdateSocialMedia(ctx context.Context, id uint64, socialMedia model.UpdateSocialMedia) (model.UpdateSocialMedia, error) {
	updatedSocialMedia, err := c.repoSocialMedia.UpdateSocialMedia(ctx, id, socialMedia)
	if err != nil {
		return model.UpdateSocialMedia{}, err
	}
	return updatedSocialMedia, nil
}
func (c *socialMediaServiceImpl) GetSocialMediasByUserID(ctx context.Context, userID uint64) ([]model.SocialMedia, error) {
    socialMedias, err := c.repoSocialMedia.GetSocialMediasByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
	for i, socialMedia := range socialMedias {
        user, err := c.repoUser.GetUsersByID(ctx, socialMedia.UserID)
		
        if err != nil {
            return nil, err
        }
        
        // Menetapkan data pengguna ke dalam struktur data foto
        socialMedias[i].User.Email = user.Email
        socialMedias[i].User.Username = user.Username
		socialMedias[i].User.ID = user.ID

    }
    return socialMedias, nil
}