package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type PokemonRepository struct {
	db *gorm.DB
}

type CreatePokemonParams struct {
	Name         string
	Description  string
	Category     string
	TypeNames    []string
	AbilityNames []string
}

type UpdatePokemonParams struct {
	Name             *string
	Description      *string
	Category         *string
	TypeNames        []string
	ReplaceTypes     bool
	AbilityNames     []string
	ReplaceAbilities bool
}

func NewPokemonRepository(db *gorm.DB) *PokemonRepository {
	return &PokemonRepository{db: db}
}

func (r *PokemonRepository) CreatePokemon(ctx context.Context, input CreatePokemonParams) (*Pokemon, error) {
	name, err := cleanRequired(input.Name, "name")
	if err != nil {
		return nil, err
	}

	description, err := cleanRequired(input.Description, "description")
	if err != nil {
		return nil, err
	}

	category, err := cleanRequired(input.Category, "category")
	if err != nil {
		return nil, err
	}

	typeNames, err := cleanNameList(input.TypeNames, "type")
	if err != nil {
		return nil, err
	}

	abilityNames, err := cleanNameList(input.AbilityNames, "ability")
	if err != nil {
		return nil, err
	}

	var pokemon Pokemon
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		types, err := findOrCreateTypes(tx, typeNames)
		if err != nil {
			return err
		}

		abilities, err := findOrCreateAbilities(tx, abilityNames)
		if err != nil {
			return err
		}

		pokemon = Pokemon{
			Name:        name,
			Description: description,
			Category:    category,
			Types:       types,
			Abilities:   abilities,
		}

		return tx.Create(&pokemon).Error
	})
	if err != nil {
		return nil, err
	}

	return r.FindPokemonByID(ctx, pokemon.ID)
}

func (r *PokemonRepository) UpdatePokemon(ctx context.Context, id uint, input UpdatePokemonParams) (*Pokemon, error) {
	var pokemon Pokemon

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&pokemon, id).Error; err != nil {
			return err
		}

		updates := map[string]any{}

		if input.Name != nil {
			name, err := cleanRequired(*input.Name, "name")
			if err != nil {
				return err
			}
			updates["name"] = name
		}

		if input.Description != nil {
			description, err := cleanRequired(*input.Description, "description")
			if err != nil {
				return err
			}
			updates["description"] = description
		}

		if input.Category != nil {
			category, err := cleanRequired(*input.Category, "category")
			if err != nil {
				return err
			}
			updates["category"] = category
		}

		if len(updates) > 0 {
			if err := tx.Model(&pokemon).Updates(updates).Error; err != nil {
				return err
			}
		}

		if input.ReplaceTypes {
			typeNames, err := cleanNameList(input.TypeNames, "type")
			if err != nil {
				return err
			}

			types, err := findOrCreateTypes(tx, typeNames)
			if err != nil {
				return err
			}

			if err := tx.Model(&pokemon).Association("Types").Replace(types); err != nil {
				return err
			}
		}

		if input.ReplaceAbilities {
			abilityNames, err := cleanNameList(input.AbilityNames, "ability")
			if err != nil {
				return err
			}

			abilities, err := findOrCreateAbilities(tx, abilityNames)
			if err != nil {
				return err
			}

			if err := tx.Model(&pokemon).Association("Abilities").Replace(abilities); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.FindPokemonByID(ctx, id)
}

func (r *PokemonRepository) DeletePokemon(ctx context.Context, id uint) (bool, error) {
	var deleted bool

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var pokemon Pokemon
		if err := tx.First(&pokemon, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&pokemon).Association("Types").Clear(); err != nil {
			return err
		}

		if err := tx.Model(&pokemon).Association("Abilities").Clear(); err != nil {
			return err
		}

		if err := tx.Delete(&pokemon).Error; err != nil {
			return err
		}

		deleted = true
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return deleted, nil
}

func (r *PokemonRepository) ListPokemons(ctx context.Context) ([]Pokemon, error) {
	var pokemons []Pokemon
	err := r.db.WithContext(ctx).
		Preload("Types").
		Preload("Abilities").
		Order("id asc").
		Find(&pokemons).Error

	return pokemons, err
}

func (r *PokemonRepository) FindPokemonByID(ctx context.Context, id uint) (*Pokemon, error) {
	var pokemon Pokemon
	err := r.db.WithContext(ctx).
		Preload("Types").
		Preload("Abilities").
		First(&pokemon, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &pokemon, nil
}

func (r *PokemonRepository) FindPokemonByName(ctx context.Context, name string) (*Pokemon, error) {
	name, err := cleanRequired(name, "name")
	if err != nil {
		return nil, err
	}

	var pokemon Pokemon
	err = r.db.WithContext(ctx).
		Preload("Types").
		Preload("Abilities").
		Where("lower(name) = lower(?)", name).
		First(&pokemon).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &pokemon, nil
}

func findOrCreateTypes(tx *gorm.DB, names []string) ([]PokemonType, error) {
	types := make([]PokemonType, 0, len(names))

	for _, name := range names {
		pokemonType := PokemonType{}
		if err := tx.Where("name = ?", name).FirstOrCreate(&pokemonType, PokemonType{Name: name}).Error; err != nil {
			return nil, err
		}
		types = append(types, pokemonType)
	}

	return types, nil
}

func findOrCreateAbilities(tx *gorm.DB, names []string) ([]Ability, error) {
	abilities := make([]Ability, 0, len(names))

	for _, name := range names {
		ability := Ability{}
		if err := tx.Where("name = ?", name).FirstOrCreate(&ability, Ability{Name: name}).Error; err != nil {
			return nil, err
		}
		abilities = append(abilities, ability)
	}

	return abilities, nil
}

func cleanRequired(value string, field string) (string, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return "", fmt.Errorf("%s is required", field)
	}

	return cleaned, nil
}

func cleanNameList(names []string, field string) ([]string, error) {
	cleaned := make([]string, 0, len(names))
	seen := map[string]struct{}{}

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			return nil, fmt.Errorf("%s cannot be empty", field)
		}

		key := strings.ToLower(name)
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		cleaned = append(cleaned, name)
	}

	if len(cleaned) == 0 {
		return nil, fmt.Errorf("at least one %s is required", field)
	}

	return cleaned, nil
}
