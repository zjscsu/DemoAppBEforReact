package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/reiver/go-porterstemmer"
)

type Image struct {
	Thumbnail          *string   `json:"thumbnail"`
}
type Detailed struct {
	Rating             *float64  `json:"rating"`
	ReviewsCount       *int      `json:"reviewsCount"`
}
type Price struct {
	List               *float64  `json:"list"`
	PreviousPrice      *float64  `json:"previousPrice"`
	PriceUnitOfMeasure *string   `json:"priceUnitOfMeasure"`
	SalesUnitOfMeasure *string   `json:"salesUnitOfMeasure"`
	SalesQuantity      *float64  `json:"salesQuantity"`
	DisplayCondition   *string   `json:"displayCondition"`
	DisplayPrice       *float64  `json:"displayPrice"`
	DisplayUnitPrice   *string   `json:"displayUnitPrice"`
	IsClearance        *bool     `json:"isClearance"`
	IsRollback         *bool     `json:"isRollback"`
	Unit               *float64  `json:"unit"`
}
type Basic struct {
	SalesUnit          *string   `json:"salesUnit"`
	Name               *string   `json:"name"`
	Image              *Image    `json:"image"`
	WeightIncrement    *float64  `json:"weightIncrement"`
	AverageWeight      *float64  `jaon:"averageWeight"`
	MaxAllowed         *int      `json:"maxAllowed"`
	ProductUrl         *string   `json:"productUrl"`
	IsSnapEligible     *bool     `json:"isSnapEligible"`
	Type               *string   `json:"type"`
}
type Store struct {
	IsOutOfStock       *bool     `json:"isOutOfStock"`
	Price              *Price    `json:"price"`
}
type ObjItem struct {
	USItemId           string    `json:"USItemId"`
	OfferId            *string   `json:"offerId"`
	Sku                *string   `json:"sku"`
	Basic              *Basic    `json:"basic"`	
	Detailed           *Detailed `json:detailed`
	Store              *Store    `json:"store`
}
type Items struct {
	Items             []ObjItem `json:"products"`  
}
// Use pointers for all nullable fields (else defaults to 0, 0.0, "" in Go)
type Item struct {
	USItemId           string    `json:"USItemId"`
	OfferId            *string   `json:"offerId"`
	Sku                *string   `json:"sku"`
	SalesUnit          *string   `json:"salesUnit"`
	Name               *string   `json:"name"`
	NameLc             *string   `json:"nameLc"`
	Thumbnail          *string   `json:"thumbnail"`
	WeightIncrement    *float64  `json:"weightIncrement"`
	AverageWeight      *float64  `jaon:"averageWeight"`
	MaxAllowed         *int      `json:"maxAllowed"`
	ProductUrl         *string   `json:"productUrl"`
	IsSnapEligible     *bool     `json:"isSnapEligible"`
	Type               *string   `json:"type"`
	Rating             *float64  `json:"rating"`
	ReviewsCount       *int      `json:"reviewsCount"`
	IsOutOfStock       *bool     `json:"isOutOfStock"`
	List               *float64  `json:"list"`
	PreviousPrice      *float64  `json:"previousPrice"`
	PriceUnitOfMeasure *string   `json:"priceUnitOfMeasure"`
	SalesUnitOfMeasure *string   `json:"salesUnitOfMeasure"`
	SalesQuantity      *float64  `json:"salesQuantity"`
	DisplayCondition   *string   `json:"displayCondition"`
	DisplayPrice       *float64  `json:"displayPrice"`
	DisplayUnitPrice   *string   `json:"displayUnitPrice"`
	IsClearance        *bool     `json:"isClearance"`
	IsRollback         *bool     `json:"isRollback"`
	Unit               *float64  `json:"unit"`
}

/*
	{products: [{
		"USItemId": "342342432",
		"offerId": "GHJHFG645654FGD",
		"sky": "8767865674",
        "basic": {
			"salesUnit": "Each_Weight",
			"name": "Bananas, each",
			"image": {
				"thumbnail": "https://i5.walmartimages.com/asr/209bb8a0-30ab-46be-b38d-58c2feb93e4a_1.1a15fb5bcbecbadd4a45822a11bf6257.jpeg?odnHeight=150&odnWidth=150&odnBg=ffffff"
			},
			"weightIncrement": 1,
			"averageWeight": 0.4,
			"maxAllowed": 12,
			"productUrl": "/ip/Bananas-each/44390948",
			"isSnapEligible": true,
			"type": "REGULAR"
		},
		"detailed": {
			"rating": 4.5,
			"reviewsCount": 205
		},
		"store": {
			"isOutOfStock": false,
			"price": {
				"list": 0.52,
				"previousPrice": 0,
				"priceUnitOfMeasure": "lb",
				"salesUnitOfMeasure": "lb",
				"salesQuantity": 1,
				"displayCondition": "each",
				"displayPrice": 0.21,
				"displayUnitPrice": "(52.0 cents/LB)",
				"isClearance": false,
				"isRollback": false,
				"unit": 0.52
			}
		}
	}]
*/

func main() {
	db      := getDB()
	itemMap := make(map[string]Item)
	client  := http.Client{
		Timeout: time.Second * 5, // Maximum of 2 secs
	}
	// url := "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=fruit&count=20"
	urls    := []string{
		"https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=fruit&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=fruit&count=200&page=2",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=fruit&count=200&page=3",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=vegetables&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=vegetables&count=200&page=2",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=vegetables&count=200&page=3",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=vegetables&count=200&page=4",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=vegetables&count=200&page=5",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=household+essentials&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=household+essentials&count=200&page=2",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=household+essentials&count=200&page=3",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=meat&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=meat&count=200&page=2",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=meat&count=200&page=3",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=dairy&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=dairy&count=200&page=2",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=dairy&count=200&page=3",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=dairy&count=200&page=4",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=bread&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=pantry&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=pantry&count=200&page=2",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=snacks&count=200",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=snacks&count=200&page=2",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=snacks&count=200&page=3",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=snacks&count=200&page=4",
        "https://grocery.walmart.com/v4/api/products/search?storeId=2119&query=snacks&count=200&page=5",
	}

	for _, url := range urls {
		body := getUrltext(url, client)
		var items Items
		err := json.Unmarshal(body, &items)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Processing %d records ...\n", len(items.Items))
		for _, objItem := range items.Items {
			item := flattenItem(objItem)
			if _, ok := itemMap[item.USItemId]; !ok {
				// FlatItem does not exist, insert into itemMap
				itemMap[item.USItemId] = item
			}
		}
	}

	// save itemMap to db
	fmt.Printf("Saving to DB\n")
	saveToDB(db, itemMap)
}

func getDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@(localhost)/DemoApp?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
	  panic("Failed to connect: " + err.Error())
	}
	return db
}

func saveToDB(db *gorm.DB, itemMap map[string]Item) {
	for _, item := range itemMap {
		db.Create(&item)
	}
}

func flattenItem(item ObjItem) Item {
	nameLc := stemSentence(*item.Basic.Name)
	return Item{
		USItemId:           item.USItemId,
		OfferId:            item.OfferId,
		Sku:                item.Sku,
		SalesUnit:          item.Basic.SalesUnit,
		Name:               item.Basic.Name,
		NameLc:             &nameLc,
		Thumbnail:          item.Basic.Image.Thumbnail,
		WeightIncrement:    item.Basic.WeightIncrement,
		AverageWeight:      item.Basic.AverageWeight,
		MaxAllowed:         item.Basic.MaxAllowed,
		ProductUrl:         item.Basic.ProductUrl,
		IsSnapEligible:     item.Basic.IsSnapEligible,
		Type:               item.Basic.Type,
		Rating:             item.Detailed.Rating,
		ReviewsCount:       item.Detailed.ReviewsCount,
		IsOutOfStock:       item.Store.IsOutOfStock,
		List:               item.Store.Price.List,
		PreviousPrice:      item.Store.Price.PreviousPrice,
		PriceUnitOfMeasure: item.Store.Price.PriceUnitOfMeasure,
		SalesUnitOfMeasure: item.Store.Price.SalesUnitOfMeasure,
		SalesQuantity:      item.Store.Price.SalesQuantity,
		DisplayCondition:   item.Store.Price.DisplayCondition,
		DisplayPrice:       item.Store.Price.DisplayPrice,
		DisplayUnitPrice:   item.Store.Price.DisplayUnitPrice,
		IsClearance:        item.Store.Price.IsClearance,
		IsRollback:         item.Store.Price.IsRollback,
		Unit:               item.Store.Price.Unit,
	}
}

func stemSentence(str string) string {
	reg, err    := regexp.Compile("[^a-zA-Z]+")
	strs        := strings.Fields(str)
	if err != nil {
        log.Fatal(err)
    }
	stemmedStrs := []string{}
	for _, word := range strs {
		stemmedWord := reg.ReplaceAllString(word, "")
		stemmedWord  = porterstemmer.StemString(strings.ToLower(stemmedWord))
		stemmedStrs  = append(stemmedStrs, stemmedWord)
	}
	return strings.Join(stemmedStrs[:], " ")
}

func getUrltext(url string, client http.Client) []byte {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}