package voucher

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Voucher struct {
	VoucherID   primitive.ObjectID `json:"voucherId" bson:"_id,omitempty"`
	VoucherName string             `json:"voucherName" bson:"voucherName"`
	Remaining   int                `json:"remaining" bson:"remaining"`
}
