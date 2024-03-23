package service

import (
	"context"

	"final_project/internal/model"
	"final_project/internal/repository"
)

type PhotoService interface {
	GetPhotos(ctx context.Context) ([]model.Photo, error)
	GetPhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error)
	DeletePhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error)
	CreatePhoto(ctx context.Context, photo model.CreatePhoto, userID uint64) (model.CreatePhoto, error)
	EditPhoto(ctx context.Context, id uint64, photo model.UpdatePhoto) (model.UpdatePhoto, error)
	GetPhotoByUserID(ctx context.Context, photoID uint64) ([]model.GetPhoto, error)
}

type photoServiceImpl struct {
	repoPhoto repository.PhotoQuery
	repoUser repository.UserQuery
}

func NewPhotoService(repoPhoto repository.PhotoQuery, repoUser repository.UserQuery) PhotoService {
	return &photoServiceImpl{
		repoPhoto: repoPhoto,
		repoUser:  repoUser,
	}
}

func (p *photoServiceImpl) GetPhotos(ctx context.Context) ([]model.Photo, error) {
    // Mengambil semua data foto dari repoPhotository
    photos, err := p.repoPhoto.GetPhotos(ctx)
    if err != nil {
        return nil, err
    }
    
    // Untuk setiap foto, ambil data pengguna yang sesuai
    for i, photo := range photos {
        user, err := p.repoUser.GetUsersByID(ctx, photo.UserID)
		
        if err != nil {
            return nil, err
        }
        
        // Menetapkan data pengguna ke dalam struktur data foto
		photos[i].User.ID = user.ID
        photos[i].User.Email = user.Email
        photos[i].User.Username = user.Username
    }
    
    return photos, nil
}


func (p *photoServiceImpl) GetPhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error) {
	photo, err := p.repoPhoto.GetPhotoByID(ctx, id)
	if err != nil {
		return model.UpdatePhoto{}, err
	}
	return photo, err
}

func (p *photoServiceImpl) DeletePhotoByID(ctx context.Context, id uint64) (model.UpdatePhoto, error) {
	photo, err := p.repoPhoto.GetPhotoByID(ctx, id)
	if err != nil {
		return model.UpdatePhoto{}, err
	}
	// if photo doesn't exist, return
	if photo.ID == 0 {
		return model.UpdatePhoto{}, nil
	}

	// delete photo by id
	err = p.repoPhoto.DeletePhotoByID(ctx, id)
	if err != nil {
		return model.UpdatePhoto{}, err
	}

	return photo, err
}

func (p *photoServiceImpl) CreatePhoto(ctx context.Context, CreatePhoto model.CreatePhoto, userID uint64) (model.CreatePhoto, error) {
	// Set created_at and updated_at timestamp
	photo := model.CreatePhoto{
		Title: CreatePhoto.Title,
		Caption: CreatePhoto.Caption,
		PhotoURL: CreatePhoto.PhotoURL,
		UserID: userID,
	}

	// Call repoPhotository to create photo
	createdPhoto, err := p.repoPhoto.CreatePhoto(ctx, photo)
	if err != nil {
		return model.CreatePhoto{}, err
	}
	return createdPhoto, nil
}
func (p *photoServiceImpl) EditPhoto(ctx context.Context, id uint64, photo model.UpdatePhoto) (model.UpdatePhoto, error) {
    // Perform validation or additional checks here if necessary

    // Call repository to edit photo
    updatedPhoto, err := p.repoPhoto.EditPhoto(ctx, id, photo)
    if err != nil {
        return model.UpdatePhoto{}, err
    }
    return updatedPhoto, nil
}
func (p *photoServiceImpl) GetPhotoByUserID(ctx context.Context, photoID uint64) ([]model.GetPhoto, error) {
    photo, err := p.repoPhoto.GetPhotoByUserID(ctx, photoID)
    if err != nil {
        return nil, err
    }
    return photo, nil
}