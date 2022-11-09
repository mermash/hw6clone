package main

import (
	"database/sql"
	"fmt"
)

type DictionaryRepo struct {
	DB *sql.DB
}

func NewDictionaryRepo(db *sql.DB) *DictionaryRepo {
	return &DictionaryRepo{
		DB: db,
	}
}

func (repo *DictionaryRepo) GetCategoryByName(name string) (*Category, error) {
	fmt.Println("Get category by name")
	category := &Category{}
	row := repo.DB.QueryRow(`SELECT category.* FROM category WHERE name = ?`, name)
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return category, nil
}
