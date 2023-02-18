package health

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/kaepa3/healthplanet"
	"golang.org/x/oauth2"
)

const (
	tokenFileName = ".token"
)

func getClient(clientId string, clientSecret string, ctx context.Context) (*healthplanet.HealthPlanetClient, error) {

	conf := healthplanet.NewConfig(&healthplanet.HealthPlanetInit{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost",
		Scopes: []string{
			"innerscan",
		},
	})
	token, err := getToken(conf)
	if err != nil {
		return nil, err
	}

	client, err := conf.GetClient(ctx, token)
	return client, err
}

func GetHealthData(clientId string, clientSecret string, opt *healthplanet.HealthPlanetOption, ctx context.Context) (*healthplanet.JsonResponce, error) {
	client, err := getClient(clientId, clientSecret, ctx)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(healthplanet.Innerscan, opt)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("code error:" + strconv.Itoa(resp.StatusCode))
	}

	return healthplanet.ConvertToJson(resp.Body)
}

func getTokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(tokenFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {

		return nil, err
	}
	var p oauth2.Token
	err = json.Unmarshal(buf, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func getTokenFromWeb(conf healthplanet.HealthPlanetConfig) (*oauth2.Token, error) {
	url := conf.AuthCodeURL("state")
	fmt.Println(url)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	authCode := scanner.Text()

	token, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, err
	}

	wf, err := os.Create(tokenFileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	if err := encoder.Encode(token); err != nil {
		log.Fatal(err)
	}
	return token, nil
}

func getToken(conf healthplanet.HealthPlanetConfig) (*oauth2.Token, error) {

	if exists(tokenFileName) {
		token, err := getTokenFromFile()
		if err == nil {
			return token, nil
		}
	}

	return getTokenFromWeb(conf)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
