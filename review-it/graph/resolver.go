package graph

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/piranatw/review-it/graph/model"
)

type Database struct {
	MoviesTable  map[string]model.Movie
	ReviewsTable map[string]model.Review
}

func (db *Database) FindMovieByID(id string) (*model.Movie, error) {
	if ret, ok := db.MoviesTable[id]; ok {
		return &ret, nil
	}
	return nil, fmt.Errorf("movie id: %s was not found", id)
}

func (db *Database) FindAllMovies() []*model.Movie {
	allMovies := []*model.Movie{}
	for _, m := range db.MoviesTable {
		newMovie := m
		allMovies = append(allMovies, &newMovie)
	}
	return allMovies
}

func (db *Database) AddMovie(input *model.Movie) error {
	newID := uuid.New().String()
	input.ID = newID
	db.MoviesTable[newID] = *input
	return nil
}

func (db *Database) UpdateMovie(input model.Movie) error {
	if _, ok := db.MoviesTable[input.ID]; !ok {
		return fmt.Errorf("movie id: %s was not found", input.ID)
	}
	db.MoviesTable[input.ID] = input
	return nil
}

func (db *Database) DeleteMovie(ID string) error {
	if _, ok := db.MoviesTable[ID]; !ok {
		return fmt.Errorf("movie id: %s was not found", ID)
	}
	delete(db.MoviesTable, ID)
	return nil
}

func (db *Database) FindReviewByID(ID string) (*model.Review, error) {
	if ret, ok := db.ReviewsTable[ID]; ok {
		return &ret, nil
	}
	return nil, fmt.Errorf("review id: %s was not found", ID)
}

func (db *Database) FindAllReviews() []*model.Review {
	allReviews := []*model.Review{}
	for _, m := range db.ReviewsTable {
		newReview := m
		allReviews = append(allReviews, &newReview)
	}
	return allReviews
}

func (db *Database) AddReview(input *model.Review) error {
	newID := uuid.New().String()
	input.ID = newID
	db.ReviewsTable[newID] = *input
	return nil
}

func (db *Database) UpdateReview(input model.Review) error {
	if _, ok := db.ReviewsTable[input.ID]; !ok {
		return fmt.Errorf("review id: %s was not found", input.ID)
	}
	db.ReviewsTable[input.ID] = input
	return nil
}

func (db *Database) DeleteReview(ID string) error {
	if _, ok := db.ReviewsTable[ID]; !ok {
		return fmt.Errorf("review id: %s was not found", ID)
	}
	delete(db.ReviewsTable, ID)
	return nil
}

type Resolver struct {
	DB Database
}
