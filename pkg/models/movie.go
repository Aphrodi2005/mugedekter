// models/movie.go
package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

var (
	ErrNoMovie        = errors.New("models: no matching movie found")
	ErrDuplicate      = errors.New("models: duplicate movie title")
	ErrTicketNotFound = errors.New("models: ticket not found")
)

type Ticket struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	MovieID   primitive.ObjectID `bson:"movie_id"`
	Seat      string             `bson:"seat"`
	Status    string             `bson:"status"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Movie struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title"`
	Genre   string             `bson:"genre"`
	Rating  int                `bson:"rating"`
	Tickets []Ticket           `bson:"tickets"`
}

type MovieModel struct {
	Client *mongo.Client
}

func (m *MovieModel) BuyTicket(movieID primitive.ObjectID, seat string) error {
	collection := m.Client.Database("mugedekter").Collection("tickets")

	filter := bson.M{"_id": movieID, "tickets": bson.M{"$elemMatch": bson.M{"seat": seat, "status": "available"}}}
	update := bson.M{
		"$set": bson.M{
			"tickets.$.status":     "sold",
			"tickets.$.updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return ErrTicketNotFound
	}

	return nil
}

func (m *MovieModel) ReturnTicket(movieID primitive.ObjectID, seat string) error {
	collection := m.Client.Database("mugedekter").Collection("tickets")

	filter := bson.M{"_id": movieID, "tickets": bson.M{"$elemMatch": bson.M{"seat": seat, "status": "sold"}}}
	update := bson.M{
		"$set": bson.M{
			"tickets.$.status":     "available",
			"tickets.$.updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return ErrTicketNotFound
	}

	return nil
}

func (m *MovieModel) GetTickets(movieID primitive.ObjectID) ([]Ticket, error) {
	collection := m.Client.Database("mugedekter").Collection("tickets")

	var result Movie
	err := collection.FindOne(context.TODO(), bson.M{"_id": movieID}, options.FindOne().SetProjection(bson.M{"tickets": 1})).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Tickets, nil
}

func (m *MovieModel) GetUserTickets(userID primitive.ObjectID) ([]Ticket, error) {
	collection := m.Client.Database("mugedekter").Collection("tickets")

	filter := bson.M{"tickets": bson.M{"$elemMatch": bson.M{"status": "sold"}}}
	projection := bson.M{"tickets": 1}

	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var userTickets []Ticket
	for cursor.Next(context.TODO()) {
		var movie Movie
		if err := cursor.Decode(&movie); err != nil {
			return nil, err
		}
		userTickets = append(userTickets, movie.Tickets...)
	}

	return userTickets, nil
}

func (m *MovieModel) Create(title, genre string, rating int) error {
	collection := m.Client.Database("mugedekter").Collection("movies")

	movie := &Movie{
		Title:  title,
		Genre:  genre,
		Rating: rating,
	}

	_, err := collection.InsertOne(context.TODO(), movie)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (m *MovieModel) Update(title, genre string, rating int, id primitive.ObjectID) error {
	collection := m.Client.Database("mugedekter").Collection("movies")

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":  title,
			"genre":  genre,
			"rating": rating,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (m *MovieModel) Delete(id primitive.ObjectID) error {
	collection := m.Client.Database("mugedekter").Collection("movies")

	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (m *MovieModel) Get(id primitive.ObjectID) (*Movie, error) {
	collection := m.Client.Database("mugedekter").Collection("movies")

	filter := bson.M{"_id": id}
	result := collection.FindOne(context.TODO(), filter)

	movie := &Movie{}
	if err := result.Decode(movie); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoMovie
		}
		return nil, err
	}

	return movie, nil
}

func (m *MovieModel) Latest(limit int64) ([]*Movie, error) {
	collection := m.Client.Database("mugedekter").Collection("movies")

	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}}).SetLimit(limit)
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var movies []*Movie
	err = cursor.All(context.TODO(), &movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (m *MovieModel) GetMovieByGenre(genre string) ([]*Movie, error) {
	collection := m.Client.Database("mugedekter").Collection("movies")

	filter := bson.M{"genre": genre}
	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var movies []*Movie
	err = cursor.All(context.TODO(), &movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
