package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/pkg/errors"
)

const beersTableName = "beers"

// BeerRepo represents a boltdb repository for managing beers
type BeerRepo struct {
	db *dynamodb.DynamoDB
}

// NewBeerRepo returns a new Beer repository using the given bolt database. It also creates the Beers
// bucket if it is not yet created on disk.
func NewBeerRepo(db *dynamodb.DynamoDB) *BeerRepo {
	return &BeerRepo{db}
}

// Get returns a Beer by its ID
func (r *BeerRepo) Get(id string) (*internal.Beer, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(beersTableName),
		Key:       mapID(id),
	}

	result, err := r.db.GetItem(input)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get Beer %s from database", id)
	}

	var beer *internal.Beer

	err = dynamodbattribute.UnmarshalMap(result.Item, beer)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal Beer %s", id)
	}

	return beer, nil
}

// GetAll returns all Beers
func (r *BeerRepo) GetAll() ([]internal.Beer, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(beersTableName),
	}

	result, err := r.db.Scan(input)
	if err != nil {
		return nil, errors.Wrap(err, "Could not Beers from database")
	}

	beers := []internal.Beer{}

	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &beers)
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal Beers")
	}

	return beers, nil
}

// Save creates or updates a Beer
func (r *BeerRepo) Save(beer *internal.Beer) error {

	b, err := dynamodbattribute.MarshalMap(beer)
	if err != nil {
		return errors.Wrap(err, "Could not marshal Beer")
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(beersTableName),
		Item:      b,
	}

	if _, err := r.db.PutItem(input); err != nil {
		return errors.Wrap(err, "Could not put Beer")
	}

	return nil
}

// Delete permanently removes a Beer
func (r *BeerRepo) Delete(id string) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(beersTableName),
		Key:       mapID(id),
	}

	if _, err := r.db.DeleteItem(input); err != nil {
		return errors.Wrapf(err, "Could not delete Beer %s", id)
	}

	return nil
}
