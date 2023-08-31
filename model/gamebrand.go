package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameBrand struct {
	ID               primitive.ObjectID `bson:"_id"`
	Code             string             `bson:"code"`
	WalletCode       string             `bson:"wallet_code"`
	GameType         string             `bson:"game_type"`
	Brief            interface{}        `bson:"brief"`
	Status           string             `bson:"status"`
	Logo             string             `bson:"logo"`
	VendorImage      string             `bson:"vendor_image"`
	BrandImage       string             `bson:"brand_image"`
	ProductImg1      string             `bson:"product_img_1"`
	ProductImg2      string             `bson:"product_img_2"`
	GameProviderCode string             `bson:"game_provider_code"`
	NamePh           string             `bson:"name_ph"`
}
