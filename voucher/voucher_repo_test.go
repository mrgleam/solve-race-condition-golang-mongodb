package voucher_test

import (
	"context"
	"log"
	"solve-race-condition-golang-mongodb/voucher"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type result struct {
	int
	error
}

var db *mongo.Database

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")
	// Setup database
	ctx := context.Background()

	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:latest"))
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %s", err)
	}

	db = mongoClient.Database("voucher")

	return func(tb testing.TB) {
		log.Println("teardown test")
		mongodbContainer.Terminate(context.Background())
	}
}

func TestCreateVoucherRepository(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		t.Skip()
	}

	t.Run("ClaimVoucher", func(t *testing.T) {
		teardownTest := setupTest(t)
		defer teardownTest(t)

		userID := primitive.NewObjectID()
		ctx := context.Background()
		res, err := voucher.CreateVoucher(db)(ctx, voucher.Voucher{
			VoucherName: "ฟรีเครื่องดื่ม",
			Remaining:   10,
		})
		if err != nil {
			t.Error(err)
		}

		insertObjID := res.InsertedID.(primitive.ObjectID)

		results := make(map[int]error)
		resultChannel := make(chan result)

		for i := 0; i < 20; i++ {
			go func(i int) {
				resultChannel <- result{i, voucher.ClaimVoucher(db)(ctx, insertObjID, userID)}
			}(i)
		}

		for i := 0; i < 20; i++ {
			r := <-resultChannel
			results[r.int] = r.error
		}

		actual, err := voucher.GetVoucherByID(db)(ctx, insertObjID)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(0, actual.Remaining); diff != "" {
			t.Errorf("ClaimVoucher() mismatch (-want +got):\n%s", diff)
		}
	})
}
