package main

import (
	"encoding/json"
	"fmt"
	// "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/labstack/echo"
	"reflect"
	"github.com/robfig/cron/v3"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	// "encoding/json"
	"time"
	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var apikey = "eD10lo4rYvpuPOYtsccIPnNrMKUvuqLyRqPHPSes23g"
var address string
// for mongo
type State struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	state  string             `json:"state" bson:"state"`
	cases  string             `json:"cases" bson:"cases"`
}

// for getting cases
type Output struct{
	Statewise []Tot `json:"statewise"`
}

type Tot struct{
	Confirmed string `json:"confirmed"`
	State string `json:"state"`
}

type MapValues struct {
	API       string
	Address   string
}

type GPS struct{
	State Stat `json:"address"`
}

type Stat struct{
	State string `json:"state"`
}

// global state array
// var statwiseData[]State

func main() {
	cr := cron.New()
	e:= echo.New()
	fmt.Printf("  --end--")
	e.GET("/getCases", func(c echo.Context) error {
		lat := c.QueryParam("lat")
		lon := c.QueryParam("lon")
		url := "http://us1.locationiq.com/v1/reverse.php?key=pk.72726dfe49cec53af2e2c8ffce4233c8&lat=" + fmt.Sprint(lat) + "&lon=" + fmt.Sprint(lon) + "&format=json"
		fmt.Printf(url)
		res, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK{
			return c.String(res.StatusCode, "ERROR")
		}
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil{
			log.Fatal(err)
		}

		var dataState GPS
		json.Unmarshal(bodyBytes, &dataState)
  		
		fmt.Println(dataState)
	
		state := dataState.State.State


		clientOptions := options.Client().ApplyURI("mongodb+srv://abhnv:abhnv@cluster0.l4l4b.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, clientOptions)
		
		fmt.Println(1)
		if err != nil {
			fmt.Println(err)
		}
		
		fmt.Println(2)
		// var stateInfo bson.M
		collection := client.Database("Inshorts").Collection("covid19india")

		
		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var stateCases bson.M
		collection.FindOne(ctx, bson.M{"state" : state}).Decode(&stateCases)
		if len(stateCases)==0 {
			return c.String(200,"Sorry, your GPS co-ordinates do not lie in the Indian territory. This website only contains indian data.")
		}	
		fmt.Println(stateCases)
		// "There are " + stateCases["cases"].(string) + " active covid-19 cases in " + state
		return c.String(200,"There are " + stateCases["cases"].(string) + " active covid-19 cases in " + state	)


	})
	cr.AddFunc("@midnight", func() {
		
		urlCovidData := "http://api.covid19india.org/data.json"

		resCovid, errCovid := http.Get(urlCovidData)
		if errCovid != nil {
			println("HERE1")
			log.Fatalln(errCovid)
		}
		defer resCovid.Body.Close()

		if resCovid.StatusCode != http.StatusOK{
			println("Here2")
			return
		}
		bodyBytesCovid, errCovid := ioutil.ReadAll(resCovid.Body)
		if errCovid != nil{
			log.Fatal(errCovid)
		}
		bodyStringCovid := string(bodyBytesCovid)
		fmt.Println(reflect.TypeOf(bodyStringCovid))
		var data Output
		json.Unmarshal(bodyBytesCovid, &data)
		
		// converting to mongo format data
		var statewiseData[]interface{}

		for _, count := range data.Statewise {
			entry := bson.D {
				{"state" , count.State},
				{"cases", count.Confirmed},
			}
			statewiseData = append(statewiseData,entry)
			// entry.state = count.State
			// entry.cases = count.Confirmed
			// statwiseData = append(statwiseData,entry)
		}







		// Inserting statewiseData to mongodb
		clientOptions := options.Client().ApplyURI("mongodb+srv://abhnv:abhnv@cluster0.l4l4b.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, clientOptions)
		
		fmt.Println(1)
		if err != nil {
			fmt.Println(err)
		}
		
		fmt.Println(2)
		// var stateInfo bson.M
		collection := client.Database("Inshorts").Collection("covid19india")

		
		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		fmt.Println("before inserting into database")
		collection.Drop(ctx)
		collection.InsertMany(ctx,  statewiseData)
		fmt.Println("inserted into database")
		return
	}) 
	cr.Start()
	e.Logger.Fatal(e.Start(":8080"))

}
