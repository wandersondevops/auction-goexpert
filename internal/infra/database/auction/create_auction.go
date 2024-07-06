package auction

import (
	"context"
	"os"
	"strconv"
	"time"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// Inicia a goroutine para verificar o fechamento automático
	go ar.startAuctionExpirationChecker(auctionEntity.Id, auctionEntity.Timestamp)

	return nil
}

func (ar *AuctionRepository) startAuctionExpirationChecker(auctionId string, startTime time.Time) {
	durationStr := os.Getenv("AUCTION_DURATION_MINUTES")
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		logger.Error("Error parsing AUCTION_DURATION_MINUTES", err)
		return
	}

	auctionDuration := time.Duration(duration) * time.Minute
	expirationTime := startTime.Add(auctionDuration)

	time.Sleep(time.Until(expirationTime))

	ctx := context.Background()
	err = ar.closeExpiredAuction(ctx, auctionId)
	if err != nil {
		logger.Error("Error closing expired auction", err)
	}
}

func (ar *AuctionRepository) closeExpiredAuction(ctx context.Context, auctionId string) error {
	filter := bson.M{"_id": auctionId, "status": auction_entity.Active}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	return err
}
