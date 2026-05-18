package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Pokemon{}, &PokemonType{}, &Ability{}); err != nil {
		return nil, err
	}

	return db, nil
}
