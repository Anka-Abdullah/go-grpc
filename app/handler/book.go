package handler

import (
	"context"
	"go-grpc-crud/app/model"
	"go-grpc-crud/app/repository"
	"go-grpc-crud/app/service"
	bookpb "go-grpc-crud/proto/book"
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type BookHandler struct {
	service service.BookService
	bookpb.UnimplementedBookServiceServer
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	return &BookHandler{service: bookService}
}

func (h *BookHandler) CreateBook(ctx context.Context, req *bookpb.CreateBookRequest) (*bookpb.CreateBookResponse, error) {

	if !isISBNDigitsOnly(req.Book.Isbn) {
		return nil, status.Errorf(codes.InvalidArgument, "ISBN must contain between 10 and 13 digits.")
	}

	bookModel := &model.Book{
		ISBN:              req.Book.Isbn,
		Title:             req.Book.Title,
		Author:            req.Book.Author,
		Publisher:         req.Book.Publisher,
		Year:              req.Book.Year,
		Quantity:          req.Book.Quantity,
		AvailableQuantity: req.Book.AvailableQuantity,
		CreatedBy:         1,
	}

	createdBook, err := h.service.CreateBook(bookModel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create book: %v", err)
	}

	return &bookpb.CreateBookResponse{Id: createdBook.ID}, nil
}

func (h *BookHandler) ReadBook(ctx context.Context, req *bookpb.ReadBookRequest) (*bookpb.ReadBookResponse, error) {
	book, err := h.service.ReadBook(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Book with ID '%d' not found", req.Id)
	}

	return &bookpb.ReadBookResponse{
		Book: &bookpb.Book{
			Id:     book.ID,
			Isbn:   book.ISBN,
			Title:  book.Title,
			Author: book.Author,
		},
	}, nil
}

func (h *BookHandler) UpdateBook(ctx context.Context, req *bookpb.UpdateBookRequest) (*bookpb.UpdateBookResponse, error) {
	bookModel := &model.Book{
		ID:     req.Book.Id,
		ISBN:   req.Book.Isbn,
		Title:  req.Book.Title,
		Author: req.Book.Author,
	}

	_, err := h.service.UpdateBook(bookModel)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Book with ID '%d' not found", req.Book.Id)

	}
	return &bookpb.UpdateBookResponse{Success: true}, nil
}

func (h *BookHandler) DeleteBook(ctx context.Context, req *bookpb.DeleteBookRequest) (*bookpb.DeleteBookResponse, error) {
	err := h.service.DeleteBook(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Book with ID '%d' not found", req.Id)

	}
	return &bookpb.DeleteBookResponse{Success: true}, nil
}

func (h *BookHandler) ListBooks(ctx context.Context, req *bookpb.ListBooksRequest) (*bookpb.ListBooksResponse, error) {
	// Convert protobuf filters to model filters
	var filters []model.Filter
	for _, f := range req.Filters {
		filters = append(filters, model.Filter{
			SearchBy: f.SearchBy,
			Value:    f.Value,
		})
	}

	books, totalCount, err := h.service.ListBooks(filters, int(req.Page), int(req.PageSize), req.SortBy, req.SortOrder)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve book list: %v", err)
	}

	// Map books
	var bookList []*bookpb.Book
	for _, b := range books {
		bookList = append(bookList, &bookpb.Book{
			Id:     b.ID,
			Isbn:   b.ISBN,
			Title:  b.Title,
			Author: b.Author,
		})
	}

	// Hitung total halaman
	totalPages := int32(0)
	if req.PageSize > 0 {
		totalPages = int32((totalCount + int64(req.PageSize) - 1) / int64(req.PageSize)) // ceiling division
	}

	return &bookpb.ListBooksResponse{
		Books:       bookList,
		TotalCount:  int32(totalCount),
		CurrentPage: req.Page,
		TotalPages:  totalPages,
		PageSize:    req.PageSize,
	}, nil
}

var onlyDigits = regexp.MustCompile(`^\d+$`)

func isISBNDigitsOnly(isbn string) bool {
	length := len(isbn)
	if length < 10 || length > 13 {
		return false
	}
	return onlyDigits.MatchString(isbn)
}
