package database

import "time"

type Pokemon struct {
	ID          uint          `gorm:"primaryKey"`
	Name        string        `gorm:"uniqueIndex;not null"`
	Description string        `gorm:"not null"`
	Category    string        `gorm:"not null"`
	Types       []PokemonType `gorm:"many2many:pokemon_type_links;constraint:OnDelete:CASCADE;"`
	Abilities   []Ability     `gorm:"many2many:pokemon_ability_links;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PokemonType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

type Ability struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}
