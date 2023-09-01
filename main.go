package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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

var Map sync.Map

type Job struct {
	ImageDest string
	Hostname  string
}

type Consumer struct {
	InputChan chan Job
	JobsChan  chan Job
}

func NewConsumer(inputSize, jobSize int) *Consumer {
	return &Consumer{
		InputChan: make(chan Job, inputSize),
		JobsChan:  make(chan Job, jobSize),
	}
}

func (c *Consumer) queue(job Job) {
	c.InputChan <- job
}

func (c *Consumer) process(job Job, num int) {
	checkExistNDownloadImage(job.ImageDest, job.Hostname, num)
}

func (c *Consumer) worker(ctx context.Context, num int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-c.JobsChan:
			if !ok || ctx.Err() != nil {
				return
			}
			c.process(job, num)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Consumer) startConsumer(ctx context.Context) {
	for {
		select {
		case job, ok := <-c.InputChan:
			if !ok || ctx.Err() != nil {
				close(c.JobsChan)
				return
			}
			c.JobsChan <- job
		case <-ctx.Done():
			close(c.JobsChan)
			return
		}
	}
}

func main() {
	c := new(config.Conf)
	mustLoadConfig(c)

	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer func() {
		log.Println("Done🥵💦")
		signal.Stop(cc)
		cancel()
	}()

	go func() {
		select {
		case <-ctx.Done():
		case <-cc:
			fmt.Println("Bye🥵")
			cancel()
		}
	}()

	var poolSize = 10
	wg := new(sync.WaitGroup)
	wg.Add(poolSize)
	consumer := NewConsumer(10000, poolSize)

	for i := 0; i < poolSize; i++ {
		go consumer.worker(ctx, i, wg)
	}

	go consumer.startConsumer(ctx)

	client := helper.InitMongo(ctx, c.BMongo.Host)
	downloadGameBrandImage(ctx, consumer, client, c.AHostname)
	downloadGameImage(ctx, consumer, client, c.AHostname)
	close(consumer.InputChan)
	wg.Wait()
}

func downloadGameBrandImage(ctx context.Context, consumer *Consumer, client *mongo.Client, aHostname string) {
	curB, err := client.
		Database(consts.MongoDatabase).
		Collection(consts.GameBrand).
		Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return
	}
	defer curB.Close(ctx)

	for curB.Next(ctx) {
		gameBrand := new(model.GameBrand)
		if err := curB.Decode(gameBrand); err != nil {
			log.Println(err)
			continue
		}

		if gameBrand.Logo != "" {
			imageDest := localDest + "/" + gameBrand.Logo
			consumer.queue(Job{imageDest, aHostname})
		}

		if gameBrand.VendorImage != "" {
			imageDest := localDest + "/" + gameBrand.VendorImage
			consumer.queue(Job{imageDest, aHostname})
		}

		if gameBrand.BrandImage != "" {
			imageDest := localDest + "/" + gameBrand.BrandImage
			consumer.queue(Job{imageDest, aHostname})
		}

		if gameBrand.ProductImg1 != "" {
			imageDest := localDest + "/" + gameBrand.ProductImg1
			consumer.queue(Job{imageDest, aHostname})
		}

		if gameBrand.ProductImg2 != "" {
			imageDest := localDest + "/" + gameBrand.ProductImg2
			consumer.queue(Job{imageDest, aHostname})
		}
	}
}

func downloadGameImage(ctx context.Context, consumer *Consumer, client *mongo.Client, aHostname string) {
	curB, err := client.
		Database(consts.MongoDatabase).
		Collection(consts.Game).
		Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return
	}
	defer curB.Close(ctx)

	for curB.Next(ctx) {
		game := new(model.Game)
		if err := curB.Decode(game); err != nil {
			log.Println(err)
			continue
		}

		if game.Image != "" {
			imageDest := localDest + "/" + game.Image
			consumer.queue(Job{imageDest, aHostname})
		}
	}
}

func checkExistNDownloadImage(imageDest string, hostname string, num int) {
	if _, ok := Map.Load(imageDest); ok {
		log.Println("Worker", num, "Duplicate: ", imageDest)
		return
	}

	Map.Store(imageDest, true)
	defer Map.Delete(imageDest)

	if _, err := os.Stat(imageDest); err == nil {
		log.Println("Worker", num, " Already Exist: ", imageDest)
		return
	}
	resp, err := http.Get(hostname + imageDest)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		log.Println("Worker", num, " Not Found: ", hostname+imageDest)
		return
	}
	defer resp.Body.Close()

	target, _ := os.Create(imageDest)
	defer target.Close()
	if _, err := io.Copy(target, resp.Body); err != nil {
		log.Println(err)
		return
	}

	log.Println("Worker", num, " Download Success: ", imageDest)
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
