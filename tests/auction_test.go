package auction_test

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAutomaticAuctionClosure(t *testing.T) {
	// Set up MongoDB in-memory
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	assert.NoError(t, err)
	defer client.Disconnect(context.Background())

	db := client.Database("auction_test")
	auctionRepo := auction.NewAuctionRepository(db)

	// Set environment variable for auction duration to 1 minute
	os.Setenv("AUCTION_DURATION_MINUTES", "1")

	// Create auction
	auctionEntity := &auction_entity.Auction{
		Id:          "auction1",
		ProductName: "Test Product",
		Category:    "TestCategory",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	err = auctionRepo.CreateAuction(context.Background(), auctionEntity)
	assert.NoError(t, err)

	// Wait for 2 minutes to ensure the auction has time to close
	time.Sleep(2 * time.Minute)

	// Verify the auction is closed
	var result auction.AuctionEntityMongo
	err = db.Collection("auctions").FindOne(context.Background(), bson.M{"_id": "auction1"}).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, auction_entity.Completed, result.Status)
}
