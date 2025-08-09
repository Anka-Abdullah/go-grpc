package repository

import (
	"errors"
	"fmt"
	"go-grpc-crud/app/model"

	"gorm.io/gorm"
)

type BookRepository interface {
	CreateBook(book *model.Book) (*model.Book, error)
	ReadBook(id uint32) (*model.Book, error)
	UpdateBook(book *model.Book) (*model.Book, error)
	DeleteBook(id uint32) error
	ListBooks(page, limit int, filters []model.Filter, sortBy, sortOrder string) ([]model.Book, error)
	CountBooks(filters []model.Filter) (int64, error)
	IsISBNUnique(isbn string) bool
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) CreateBook(book *model.Book) (*model.Book, error) {

	if err := r.db.Create(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *bookRepository) ReadBook(id uint32) (*model.Book, error) {
	var book model.Book
	if err := r.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) UpdateBook(book *model.Book) (*model.Book, error) {
	if err := r.db.Save(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *bookRepository) DeleteBook(id uint32) error {
	if err := r.db.Delete(&model.Book{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *bookRepository) CountBooks(filters []model.Filter) (int64, error) {
	var total int64
	query := r.db.Model(&model.Book{})

	for _, filter := range filters {
		query = query.Where(fmt.Sprintf("%s LIKE ?", filter.SearchBy), "%"+filter.Value+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *bookRepository) ListBooks(page, limit int, filters []model.Filter, sortBy, sortOrder string) ([]model.Book, error) {
	var books []model.Book
	query := r.db.Where("1 = 1")

	for _, filter := range filters {
		query = query.Where(fmt.Sprintf("%s LIKE ?", filter.SearchBy), "%"+filter.Value+"%")
	}

	if sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}

	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) IsISBNUnique(isbn string) bool {
	var book model.Book

	result := r.db.Where("isbn = ?", isbn).First(&book)

	return result.Error == gorm.ErrRecordNotFound
}
