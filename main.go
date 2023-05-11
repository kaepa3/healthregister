package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kaepa3/healthplanet"
	"github.com/kaepa3/healthregister/config"
	"github.com/kaepa3/healthregister/db"
	"github.com/kaepa3/healthregister/health"
)

const (
	mondodbUri = "mongodb://mongo:pass@127.0.0.1:27017/"
)

var conf *config.HealthConfig

func main() {
	var err error
	conf, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	data, err := GetData()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = RegisterDB(data)
	if err != nil {
		fmt.Println(err)
	}
}

func RegisterDB(data *healthplanet.JsonResponce) error {
	ctx := context.Background()
	dbOpt := db.RegisterOption{
		Url: mondodbUri,
	}
	return db.RegisterDB(data, &dbOpt, conf.ClientID, ctx)
}

func GetData() (*healthplanet.JsonResponce, error) {
	from := time.Now().AddDate(0, -1, 0)
	ctx := context.Background()
	opt := healthplanet.HealthPlanetOption{
		Format: healthplanet.Json,
		From:   from,
		To:     time.Time{},
	}
	return health.GetHealthData(conf.ClientID, conf.ClientSecret, &opt, ctx)
}
