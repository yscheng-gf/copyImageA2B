package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/yscheng-gf/copyImageA2B/consts"
	"github.com/yscheng-gf/copyImageA2B/helper"
	"github.com/yscheng-gf/copyImageA2B/internal/config"
	"github.com/yscheng-gf/copyImageA2B/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	configFile = "etc/config.yaml"
	envFile    = "etc/.env"
	localDest  = "uploads"
)

func main() {
	c := new(config.Conf)
	mustLoadConfig(c)
	client := helper.InitMongo(c.BMongo.Host)
	downloadGameBrandImage(client, c.AHostname)
	downloadGameImage(client, c.AHostname)
}

func downloadGameBrandImage(client *mongo.Client, aHostname string) {
	gameBrand := new(model.GameBrand)
	curB, err := client.
		Database(consts.MongoDatabase).
		Collection(consts.GameBrand).
		Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		return
	}
	defer curB.Close(context.TODO())

	for curB.Next(context.TODO()) {
		if err := curB.Decode(gameBrand); err != nil {
			log.Println(err)
			continue
		}

		if gameBrand.Logo != "" {
			imageDest := localDest + "/" + gameBrand.Logo
			checkExistNDownloadImage(imageDest, aHostname)
		}

		if gameBrand.VendorImage != "" {
			imageDest := localDest + "/" + gameBrand.VendorImage
			checkExistNDownloadImage(imageDest, aHostname)
		}

		if gameBrand.BrandImage != "" {
			imageDest := localDest + "/" + gameBrand.BrandImage
			checkExistNDownloadImage(imageDest, aHostname)
		}

		if gameBrand.ProductImg1 != "" {
			imageDest := localDest + "/" + gameBrand.ProductImg1
			checkExistNDownloadImage(imageDest, aHostname)
		}

		if gameBrand.ProductImg2 != "" {
			imageDest := localDest + "/" + gameBrand.ProductImg2
			checkExistNDownloadImage(imageDest, aHostname)
		}
	}
}

func downloadGameImage(client *mongo.Client, aHostname string) {
	game := new(model.Game)
	curB, err := client.
		Database(consts.MongoDatabase).
		Collection(consts.Game).
		Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		return
	}
	defer curB.Close(context.TODO())

	wg := sync.WaitGroup{}
	for curB.Next(context.TODO()) {
		if err := curB.Decode(game); err != nil {
			log.Println(err)
			continue
		}

		wg.Add(1)
		go func() {
			if game.Image != "" {
				imageDest := localDest + "/" + game.Image
				checkExistNDownloadImage(imageDest, aHostname)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func checkExistNDownloadImage(imageDest string, hostname string) {
	if _, err := os.Stat(imageDest); err == nil {
		log.Println("Already Exist")
		return
	}
	resp, err := http.Get(hostname + imageDest)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	target, _ := os.Create(imageDest)
	defer target.Close()
	if _, err := io.Copy(target, resp.Body); err != nil {
		log.Println(err)
		return
	}

	log.Println("Download Success: ", imageDest)
}

func mustLoadConfig(c *config.Conf) {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error load config file: %w", err))
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		panic(fmt.Errorf("error viper unmarshal: %w", err))
	}
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
