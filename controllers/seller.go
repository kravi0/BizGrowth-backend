package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kravi0/BizGrowth-backend/database"
	"github.com/kravi0/BizGrowth-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var SellerCollection *mongo.Collection = database.ProductData(database.Client, "seller")
var ProductReference *mongo.Collection = database.ProductData(database.Client, "ProductReference")


//var  *mongo.Collection = database.ProductData(database.Client, "seller")

// get all seller if no id is passesed all details if id id passed it will return sepcific seller
func GetSeller() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		// Check if the user is a seller
		if checkSeller(ctx, c) {
			var sellerDetail models.Seller
			uid, _ := c.Get("uid")
			uids := fmt.Sprintf("%v", uid)
			sellerID, _ := primitive.ObjectIDFromHex(string(uids))
			filter := primitive.M{"_id": sellerID}
			if err := SellerCollection.FindOne(ctx, filter).Decode(&sellerDetail); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seller details"})
				return
			}
			c.JSON(http.StatusOK, sellerDetail)
			return
		}
		
		// Check if the user is an admin
		if !checkAdmin(ctx, c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		// Parse seller ID from query
		sellerID := c.Query("sellerId")
		if sellerID != "" {
			// Fetch details of the specific seller
			var sellerDetail models.Seller
			sellerObjID, err := primitive.ObjectIDFromHex(sellerID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
				return
			}
			filter := primitive.M{"_id": sellerObjID}
			if err := SellerCollection.FindOne(ctx, filter).Decode(&sellerDetail); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
				return
			}
			c.JSON(http.StatusOK, sellerDetail)
			return
		}

		// Fetch details of all sellers
		var sellerDetails []models.Seller
		cur, err := SellerCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find sellers"})
			return
		}
		if err := cur.All(ctx, &sellerDetails); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seller details"})
			return
		}
		c.JSON(http.StatusOK, sellerDetails)
	}
}


// func AddSeller() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		var sellerDetails models.Seller
// 		defer cancel()
// 	}
// }

// delete specific seller
func DeleteSeller() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		seller_ID := c.Query("sellerID")
		if seller_ID == "" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid seller id"})
			c.Abort()
			return
		}
		sellerID, err := primitive.ObjectIDFromHex(seller_ID)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}
		filter := bson.D{primitive.E{Key: "_id", Value: sellerID}}
		_, err = SellerCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "fail to delete"})
			c.Abort()
			return
		}
		c.Header("content-type", "application/json")
		c.JSON(http.StatusOK, gin.H{"success": "deleted successfully"})

	}
}


func AddProductReferenceHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            SellerID    string `json:"seller_id" binding:"required"`
            ProductID   string `json:"product_id" binding:"required"`
            Price       string `json:"price" binding:"required"`
            MinQuantity int    `json:"min_quantity" binding:"required"`
            MaxQuantity int    `json:"max_quantity" binding:"required"`
        }
        ctx := context.Background()

        if err := c.ShouldBindJSON(&input); err != nil {
            fmt.Println(err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        sellerID, err := primitive.ObjectIDFromHex(input.SellerID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller ID"})
            return
        }

        productID, err := primitive.ObjectIDFromHex(input.ProductID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
            return
        }

        productReference := models.ProductReference{
            ID:          primitive.NewObjectID(),
            ProductID:   productID,
            SellerID:    sellerID,
            Price:       input.Price,
            MinQuantity: input.MinQuantity,
            MaxQuantity: input.MaxQuantity,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
            Approved:    false, // You may set the default value as needed
            Archived:    false, // You may set the default value as needed
        }

        // Insert product reference into the ProductReferenceCollection
        _, err = ProductReference.InsertOne(ctx, productReference)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Update seller model with the product reference ID
        update := bson.M{"$push": bson.M{"product_references": productReference.ID}}
        _, err = SellerCollection.UpdateOne(ctx, bson.M{"_id": sellerID}, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

		 // Update product model with the product reference ID
        updateProduct := bson.M{"$push": bson.M{"product_references": productReference.ID}}
        _, err = ProductCollection.UpdateOne(ctx, bson.M{"_id": productID}, updateProduct)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Product reference added successfully"})
    }
}

