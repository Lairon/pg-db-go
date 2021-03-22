package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type User struct {
	ID     int64
	Name   string
	Emails []string
}

type Book struct {
	ID       int64
	Title    string
	ReaderID int64
	Reader   *User `pg:"rel:has-one"`
}

func main() {
	exampleDBModel()
}

func exampleDBModel() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "pg-db-go",
	})
	defer db.Close()

	err := createScheme(db)
	if err != nil {
		panic(err)
	}

	user1 := User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	_, err = db.Model(&user1).Insert()
	if err != nil {
		panic(err)
	}

	user2 := User{
		Name:   "root",
		Emails: []string{"root1@admin", "root2@admin"},
	}
	_, err = db.Model(&user2).Insert()
	if err != nil {
		panic(err)
	}

	book1 := Book{
		Title: "PabloC",
		ReaderID: user1.ID,
	}
	_, err = db.Model(&book1).Insert()
	if err != nil {
		panic(err)
	}

	//select user by pk
	user := User{
		ID:     user1.ID,
	}
	err = db.Model(&user).WherePK().Select()
	if err != nil {
		panic(err)
	}

	// Select all users.
	var users []User
	err = db.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	// Select book and associated reader in one query.
	book := Book{}
	err = db.Model(&book).
		Relation("Reader").
		Where("book.id = ?", &book1.ID).
		Select()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v \n", user)
	fmt.Printf("%+v \n", users)
	fmt.Printf("%+v \n", book)

}

func createScheme(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Book)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
