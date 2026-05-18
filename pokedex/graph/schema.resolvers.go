package graph

import (
	"context"
	"fmt"
	"strconv"

	"pokedex/database"
	"pokedex/graph/model"
)

func (r *mutationResolver) CreatePokemon(ctx context.Context, input model.CreatePokemonInput) (*model.Pokemon, error) {
	pokemon, err := r.PokemonRepo.CreatePokemon(ctx, database.CreatePokemonParams{
		Name:         input.Name,
		Description:  input.Description,
		Category:     input.Category,
		TypeNames:    input.Type,
		AbilityNames: input.Abilities,
	})
	if err != nil {
		return nil, err
	}

	return toGraphPokemon(pokemon), nil
}

func (r *mutationResolver) UpdatePokemon(ctx context.Context, id string, input model.UpdatePokemonInput) (*model.Pokemon, error) {
	pokemonID, err := parseGraphID(id)
	if err != nil {
		return nil, err
	}

	pokemon, err := r.PokemonRepo.UpdatePokemon(ctx, pokemonID, database.UpdatePokemonParams{
		Name:             input.Name,
		Description:      input.Description,
		Category:         input.Category,
		TypeNames:        input.Type,
		ReplaceTypes:     input.Type != nil,
		AbilityNames:     input.Abilities,
		ReplaceAbilities: input.Abilities != nil,
	})
	if err != nil {
		return nil, err
	}

	if pokemon == nil {
		return nil, fmt.Errorf("pokemon with id %s was not found", id)
	}

	return toGraphPokemon(pokemon), nil
}

func (r *mutationResolver) DeletePokemon(ctx context.Context, id string) (bool, error) {
	pokemonID, err := parseGraphID(id)
	if err != nil {
		return false, err
	}

	return r.PokemonRepo.DeletePokemon(ctx, pokemonID)
}

func (r *queryResolver) Pokemons(ctx context.Context) ([]*model.Pokemon, error) {
	pokemons, err := r.PokemonRepo.ListPokemons(ctx)
	if err != nil {
		return nil, err
	}

	return toGraphPokemons(pokemons), nil
}

func (r *queryResolver) Pokemon(ctx context.Context, id string) (*model.Pokemon, error) {
	pokemonID, err := parseGraphID(id)
	if err != nil {
		return nil, err
	}

	pokemon, err := r.PokemonRepo.FindPokemonByID(ctx, pokemonID)
	if err != nil {
		return nil, err
	}

	return toGraphPokemon(pokemon), nil
}

func (r *queryResolver) PokemonByName(ctx context.Context, name string) (*model.Pokemon, error) {
	pokemon, err := r.PokemonRepo.FindPokemonByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return toGraphPokemon(pokemon), nil
}

func parseGraphID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q", id)
	}

	return uint(parsed), nil
}

func toGraphPokemons(pokemons []database.Pokemon) []*model.Pokemon {
	result := make([]*model.Pokemon, 0, len(pokemons))
	for i := range pokemons {
		result = append(result, toGraphPokemon(&pokemons[i]))
	}

	return result
}

func toGraphPokemon(pokemon *database.Pokemon) *model.Pokemon {
	if pokemon == nil {
		return nil
	}

	types := make([]*model.PokemonType, 0, len(pokemon.Types))
	for _, pokemonType := range pokemon.Types {
		types = append(types, &model.PokemonType{
			ID:   formatGraphID(pokemonType.ID),
			Name: pokemonType.Name,
		})
	}

	abilities := make([]*model.Ability, 0, len(pokemon.Abilities))
	for _, ability := range pokemon.Abilities {
		abilities = append(abilities, &model.Ability{
			ID:   formatGraphID(ability.ID),
			Name: ability.Name,
		})
	}

	return &model.Pokemon{
		ID:          formatGraphID(pokemon.ID),
		Name:        pokemon.Name,
		Description: pokemon.Description,
		Category:    pokemon.Category,
		Type:        types,
		Abilities:   abilities,
	}
}

func formatGraphID(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
