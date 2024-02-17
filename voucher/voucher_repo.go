package voucher

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateVoucher(db *mongo.Database) func(context.Context, Voucher) (*mongo.InsertOneResult, error) {
	return func(ctx context.Context, voucher Voucher) (*mongo.InsertOneResult, error) {
		voucherCollection := GetVoucherCollection(db)
		result, err := voucherCollection.InsertOne(ctx, voucher)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
}

func GetVoucherByID(db *mongo.Database) func(context.Context, primitive.ObjectID) (*Voucher, error) {
	return func(ctx context.Context, id primitive.ObjectID) (*Voucher, error) {
		voucherCollection := GetVoucherCollection(db)

		var voucher Voucher

		voucherFilter := bson.M{"_id": id}

		err := voucherCollection.FindOne(ctx, voucherFilter).Decode(&voucher)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("voucher not found")
			}
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		return &voucher, nil
	}
}

func ClaimVoucher(db *mongo.Database) func(context.Context, primitive.ObjectID, primitive.ObjectID) error {
	return func(ctx context.Context, voucherID, userID primitive.ObjectID) error {
		voucherCollection := GetVoucherCollection(db)
		var voucher Voucher

		voucherFilter := bson.M{"_id": voucherID}
		err := voucherCollection.FindOne(ctx, voucherFilter).Decode(&voucher)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fmt.Errorf("voucher not found")
			}
			return err
		}

		if voucher.Remaining <= 0 {
			return fmt.Errorf("it's fully claimed")
		}

		_, err = voucherCollection.UpdateOne(ctx, voucherFilter, bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "remaining", Value: voucher.Remaining - 1},
			}}})

		// voucherFindOneAndUpdateFilter := bson.M{"_id": voucherID, "remaining": bson.M{"$gt": 0}}
		// err := voucherCollection.FindOneAndUpdate(ctx, voucherFindOneAndUpdateFilter, bson.D{
		// 	{Key: "$inc", Value: bson.D{
		// 		{Key: "remaining", Value: -1},
		// 	}}}).Decode(&voucher)

		if err != nil {
			return err
		}

		return nil
	}
}
