package repository

import (
	"github.com/ekkinox/fx-template/app/model"
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

func (r *PostRepository) Find(id int) (*model.Post, error) {

	var post model.Post

	res := r.db.Take(&post, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &post, nil
}

func (r *PostRepository) FindAll() ([]model.Post, error) {

	var posts []model.Post

	res := r.db.Find(&posts)
	if res.Error != nil {
		return nil, res.Error
	}

	return posts, nil
}

func (r *PostRepository) Create(post *model.Post) error {
	res := r.db.Create(post)

	return res.Error
}

func (r *PostRepository) Update(post *model.Post, update *model.Post) error {
	res := r.db.Model(post).Updates(update)

	return res.Error
}

func (r *PostRepository) Delete(post *model.Post) error {
	res := r.db.Delete(post)

	return res.Error
}
