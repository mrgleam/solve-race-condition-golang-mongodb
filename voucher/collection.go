package voucher

import "go.mongodb.org/mongo-driver/mongo"

func GetVoucherCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("voucher")
}
