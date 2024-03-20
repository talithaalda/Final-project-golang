package service

import (
	"context"

	"final_project/internal/model"
	"final_project/internal/repository"
)

type PhotoService interface {
	GetPhotos(ctx context.Context) ([]model.Photo, error)
	GetPhotoByID(ctx context.Context, id uint64) (model.Photo, error)
	DeletePhotoByID(ctx context.Context, id uint64) (model.Photo, error)
	CreatePhoto(ctx context.Context, photo model.Photo, userID uint64) (model.Photo, error)
	EditPhoto(ctx context.Context, id uint64, photo model.Photo) (model.Photo, error)
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
        photos[i].User.Email = user.Email
        photos[i].User.Username = user.Username
    }
    
    return photos, nil
}


func (p *photoServiceImpl) GetPhotoByID(ctx context.Context, id uint64) (model.Photo, error) {
	photo, err := p.repoPhoto.GetPhotoByID(ctx, id)
	if err != nil {
		return model.Photo{}, err
	}
	return photo, err
}

func (p *photoServiceImpl) DeletePhotoByID(ctx context.Context, id uint64) (model.Photo, error) {
	photo, err := p.repoPhoto.GetPhotoByID(ctx, id)
	if err != nil {
		return model.Photo{}, err
	}
	// if photo doesn't exist, return
	if photo.ID == 0 {
		return model.Photo{}, nil
	}

	// delete photo by id
	err = p.repoPhoto.DeletePhotoByID(ctx, id)
	if err != nil {
		return model.Photo{}, err
	}

	return photo, err
}

func (p *photoServiceImpl) CreatePhoto(ctx context.Context, CreatePhoto model.Photo, userID uint64) (model.Photo, error) {
	// Set created_at and updated_at timestamp
	photo := model.Photo{
		Title: CreatePhoto.Title,
		Caption: CreatePhoto.Caption,
		PhotoURL: CreatePhoto.PhotoURL,
		UserID: userID,
	}

	// Call repoPhotository to create photo
	createdPhoto, err := p.repoPhoto.CreatePhoto(ctx, photo)
	if err != nil {
		return model.Photo{}, err
	}
	return createdPhoto, nil
}
func (p *photoServiceImpl) EditPhoto(ctx context.Context, id uint64, photo model.Photo) (model.Photo, error) {
    // Perform validation or additional checks here if necessary

    // Call repository to edit photo
    updatedPhoto, err := p.repoPhoto.EditPhoto(ctx, id, photo)
    if err != nil {
        return model.Photo{}, err
    }
    return updatedPhoto, nil
}
