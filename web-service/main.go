package main

//import required packages
import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// item represents data about a item.
type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

// receipt represents data about a receipt.
type Receipt struct {
	ID           string
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        float64 `json:"total"`
	Items        []Item  `json:"items"`
}

// main function with routes
func main() {
	router := gin.Default()
	router.GET("/receipts", getReceipts)
	router.POST("/receipts", postReceipts)
	router.GET("/receipts/:id", getReceiptByID)
	router.GET("/receipts/:id/points", getReceiptPoints)
	router.Run("localhost:8080")
}

// receipts example data.
var receipts = []Receipt{
	{
		ID:           "7fb1377b-b223-49d9-a31a-5a02701dd310",
		Retailer:     "Walgreens",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "08:13",
		Total:        2.65,
		Items: []Item{
			Item{
				ShortDescription: "Pepsi - 12-oz",
				Price:            1.25,
			},
			Item{
				ShortDescription: "Dasani",
				Price:            1.40,
			},
		},
	},

	{
		ID:           "3dc3154c-d423-57b9-c13d-2b07343bc504",
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        1.25,
		Items: []Item{
			Item{
				ShortDescription: "Pepsi - 12-oz",
				Price:            1.25,
			},
		},
	},
}

// getReceipts responds with the list of all receipts as JSON.
func getReceipts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, receipts)
}

// postReceipts adds a receipt from JSON received in the request body.
func postReceipts(c *gin.Context) {
	var newReceipt Receipt
	var id = uuid.New()

	if err := c.BindJSON(&newReceipt); err != nil {
		return
	}
	newReceipt.ID = id.String()

	// Add the new receipt to the slice.
	receipts = append(receipts, newReceipt)
	c.IndentedJSON(http.StatusCreated, newReceipt)
}

// extract the ID in the request path, then locate an receipt that matches.
func getReceiptByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of receipts for matching id
	for _, r := range receipts {
		if r.ID == id {
			c.IndentedJSON(http.StatusOK, r)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}

// extract the ID in the request path, then locate an receipt that matches.
// Than it calculates required points.
func getReceiptPoints(c *gin.Context) {
	id := c.Param("id")
	var points = 0
	var day = ""
	var time = ""

	// Loop over the list of receipts for matching id
	for _, r := range receipts {
		if r.ID == id {

			//converts date and time to int to use in calulations
			for i := 8; i < len(r.PurchaseDate); i++ {
				day = string(r.PurchaseDate[8]) + string(r.PurchaseDate[9])
			}
			dayAsInt, dayErr := strconv.Atoi(day)
			for i := 0; i < len(r.PurchaseTime); i++ {
				time = string(r.PurchaseTime[0]) + string(r.PurchaseTime[1])
			}
			timeAsInt, timeErr := strconv.Atoi(time)
			if dayErr != nil || timeErr != nil {
			}

			// calculations for various points

			//points for retailer name length
			points += len(r.Retailer)
			//points for evry pair items
			points += ((len(r.Items) / 2) * 5)
			//points for odd date purchase
			if (dayAsInt % 2) == 1 {
				points += 6
			}
			//points for time after 2 and before 4
			if timeAsInt >= 14 && timeAsInt < 16 {
				points += 10
			}
			//points for even dollar amount
			if r.Total-math.Floor(r.Total) == 0 {
				points += 50
			}
			//points for cents divisable by .25
			if math.Mod(r.Total, 0.25) == 0 && r.Total-math.Floor(r.Total) != 0 {
				points += 25
			}
			//points if description is divisable by 3
			for i := 0; i < len(r.Items); i++ {
				if (len(r.Items[i].ShortDescription) % 3) == 0 {
					points += int(math.Ceil(r.Items[i].Price * 0.2))
				}
			}

			c.IndentedJSON(http.StatusOK, gin.H{"points": points})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}
