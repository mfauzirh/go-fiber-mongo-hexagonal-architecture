package repository

import (
	"context"
	"log"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database, collectionName string) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		log.Println("error when try to inserting new product collection", err)
		return nil, domain.ErrInternal
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return product, nil
}

func (r *ProductRepository) GetProductById(ctx context.Context, id string) (*domain.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error when try to convert object id", err)
		return nil, err
	}

	var product domain.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("error when try to retrieve product from collection, product not found", err)
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context, page int64, limit int64) ([]domain.Product, error) {
	skip := (page - 1) * limit

	cursor, err := r.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		log.Println("error when try to retrieve products from collection", err)
		return nil, domain.ErrInternal
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Println("error when try to retrieve all product cursor", err)
		return nil, domain.ErrInternal
	}

	return products, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": bson.M{
		"name":  product.Name,
		"stock": product.Stock,
		"price": product.Price,
	}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("error when try to update product in product collection", err)
		return nil, domain.ErrInternal
	}

	if result.MatchedCount == 0 {
		log.Println("error there is no match product to update", err)
		return nil, domain.ErrProductNotFound
	}

	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInternal
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Println("error when try to delete product in product collection", err)
		return domain.ErrInternal
	}

	if result.DeletedCount == 0 {
		log.Println("error there is no match product to delete", err)
		return domain.ErrProductNotFound
	}

	return nil
}
