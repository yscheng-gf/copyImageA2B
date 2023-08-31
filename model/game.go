package model

type Game struct {
	GameProviderCode string `bson:"game_provider_code"`
	GameBrandCode    string `bson:"game_brand_code"`
	GameCode         string `bson:"game_code"`
	GameType         string `bson:"game_type"`
	Name             string `bson:"name"`
	NameZhCN         string `bson:"name_zh-CN"`
	NameEn           string `bson:"name_en"`
	NameViVN         string `bson:"name_vi-VN"`
	NameZhHK         string `bson:"name_zh-HK"`
	NameTh           string `bson:"name_th"`
	ImageEn          string `bson:"image_en"`
	ImageTh          string `bson:"image_th"`
	ImageViVN        string `bson:"image_vi-VN"`
	ImageZhCN        string `bson:"image_zh-CN"`
	ImageZhHK        string `bson:"image_zh-HK"`
	NamePh           string `bson:"name_ph"`
	Image            string `bson:"image"`
}
