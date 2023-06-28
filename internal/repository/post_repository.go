package repository

import (
	"context"

	"github.com/ekkinox/fx-template/internal/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) (*PostRepository, error) {

	err := db.AutoMigrate(&model.Post{})
	if err != nil {
		return nil, err
	}

	return &PostRepository{
		db: db,
	}, nil
}

func (r *PostRepository) Find(ctx context.Context, id int) (*model.Post, error) {

	var post model.Post

	res := r.db.WithContext(ctx).Take(&post, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &post, nil
}

func (r *PostRepository) FindAll(ctx context.Context) ([]model.Post, error) {

	var posts []model.Post

	res := r.db.WithContext(ctx).Find(&posts)
	if res.Error != nil {
		return nil, res.Error
	}

	return posts, nil
}

func (r *PostRepository) Create(ctx context.Context, post *model.Post) error {
	res := r.db.WithContext(ctx).Create(post)

	return res.Error
}

func (r *PostRepository) Update(ctx context.Context, post *model.Post, update *model.Post) error {
	res := r.db.WithContext(ctx).Model(post).Updates(update)

	return res.Error
}

func (r *PostRepository) Delete(ctx context.Context, post *model.Post) error {
	res := r.db.WithContext(ctx).Delete(post)

	return res.Error
}
