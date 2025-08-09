package main

import (
	"context"
	"fmt"
	"log"
	"time"

	book "go-grpc-crud/proto/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Gagal terhubung ke server: %v", err)
	}
	defer conn.Close()

	client := book.NewBookServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// CREATE
	fmt.Println("--- Membuat buku baru ---")
	createRes, err := client.CreateBook(ctx, &book.CreateBookRequest{
		Book: &book.Book{
			Title:  "The Go Programming Language",
			Author: "Alan A. A. Donovan & Brian W. Kernighan",
		},
	})
	if err != nil {
		log.Fatalf("Gagal membuat buku: %v", err)
	}
	bookID := createRes.Id
	// fmt.Printf("Buku berhasil dibuat dengan ID: %s\n\n", bookID)

	// READ
	fmt.Println("--- Membaca buku ---")
	readRes, err := client.ReadBook(ctx, &book.ReadBookRequest{Id: bookID})
	if err != nil {
		log.Fatalf("Gagal membaca buku: %v", err)
	}
	fmt.Printf("Buku yang dibaca: %v\n\n", readRes.Book)

	// UPDATE
	fmt.Println("--- Memperbarui buku ---")
	updateRes, err := client.UpdateBook(ctx, &book.UpdateBookRequest{
		Book: &book.Book{
			Id:     bookID,
			Title:  "The Go Programming Language (Edisi Revisi)",
			Author: "Alan A. A. Donovan & Brian W. Kernighan",
		},
	})
	if err != nil {
		log.Fatalf("Gagal memperbarui buku: %v", err)
	}
	fmt.Printf("Buku berhasil diperbarui: %v\n\n", updateRes.Success)

	// LIST
	fmt.Println("--- Menampilkan semua buku ---")
	listRes, err := client.ListBooks(ctx, &book.ListBooksRequest{})
	if err != nil {
		log.Fatalf("Gagal menampilkan daftar buku: %v", err)
	}
	for _, b := range listRes.Books {
		fmt.Printf("- %v\n", b)
	}
	fmt.Println()

	// DELETE
	fmt.Println("--- Menghapus buku ---")
	deleteRes, err := client.DeleteBook(ctx, &book.DeleteBookRequest{Id: bookID})
	if err != nil {
		log.Fatalf("Gagal menghapus buku: %v", err)
	}
	fmt.Printf("Buku berhasil dihapus: %v\n\n", deleteRes.Success)
}
