package service

import (
	"errors"
	"go-grpc-crud/app/model"
	"go-grpc-crud/app/repository"
	"go-grpc-crud/utils"
	"time"
)

type BookService interface {
	CreateBook(book *model.Book) (*model.Book, error)
	ReadBook(id uint32) (*model.Book, error)
	UpdateBook(book *model.Book) (*model.Book, error)
	DeleteBook(id uint32) error
	ListBooks(filters []model.Filter, page, pageSize int, sortBy, sortOrder string) ([]model.Book, int64, error)
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo}
}

func (s *bookService) CreateBook(book *model.Book) (*model.Book, error) {
	task := func() (any, error) {
		if !s.repo.IsISBNUnique(book.ISBN) {
			return nil, errors.New("ISBN sudah ada")
		}
		createdBook, err := s.repo.CreateBook(book)
		return createdBook, err
	}

	res, err := utils.ExecuteAsync(task, 10*time.Second)
	if err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return res.Result.(*model.Book), nil
}

func (s *bookService) ReadBook(id uint32) (*model.Book, error) {
	return s.repo.ReadBook(id)
}

func (s *bookService) UpdateBook(book *model.Book) (*model.Book, error) {
	return s.repo.UpdateBook(book)
}

func (s *bookService) DeleteBook(id uint32) error {
	return s.repo.DeleteBook(id)
}

func (s *bookService) ListBooks(filters []model.Filter, page, pageSize int, sortBy, sortOrder string) ([]model.Book, int64, error) {
	countTask := func() (any, error) {
		total, err := s.repo.CountBooks(filters)
		return total, err
	}

	listTask := func() (any, error) {
		books, err := s.repo.ListBooks(page, pageSize, filters, sortBy, sortOrder)
		return books, err
	}

	countRes, err := utils.ExecuteAsync(countTask, 10*time.Second)
	if err != nil {
		return nil, 0, err
	}

	listRes, err := utils.ExecuteAsync(listTask, 10*time.Second)
	if err != nil {
		return nil, 0, err
	}

	total := countRes.Result.(int64)
	books := listRes.Result.([]model.Book)

	return books, total, nil
}
