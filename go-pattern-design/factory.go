package main

import "fmt"

type Product interface {
	Use() string
}

type Creator interface {
	Create() Product
}

type Book struct {
}

func (b *Book) Use() string {
	return "Using a Book"
}

type NoteBook struct {
}

func (nb *NoteBook) Use() string {
	return "Using a NoteBook"
}

type BookFactory struct {
}

func (bf *BookFactory) Create() Product {
	return &Book{}
}

type NoteBookFactory struct {
}

func (nf *NoteBookFactory) Create() Product {
	return &NoteBook{}

}

func main() {
	var factory Creator
	// 使用书籍工厂生产书籍
	factory = &BookFactory{}
	book := factory.Create()
	fmt.Println(book.Use())

	// 使用笔记本工厂生产笔记本
	factory = &NoteBookFactory{}
	notebook := factory.Create()
	fmt.Println(notebook.Use())
}
