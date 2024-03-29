package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kravi0/BizGrowth-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EnquiryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobileno, existse := c.Get("mobile")
		
		uid, exists := c.Get("uid")
	
		if !existse || !exists || uid == "" || mobileno == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var enquire models.Enquire
		if err := c.BindJSON(&enquire); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		

	

		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		enquire.Enquire_id = primitive.NewObjectID()
		enquire.User_id = uid.(string)

		fmt.Println(enquire.Enquiry_note)
		
		
		_, err := EnquireCollection.InsertOne(ctx, enquire)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Success": "enquiry registerd"})
	}
}

func GetUserEnquiries() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user ID from token
        uid, exists := c.Get("uid")
        if !exists || uid == "" {
            c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
            return
        }

        // Convert user ID to string
        userID := uid.(string)

        // Define filter to fetch enquiries for the specific user
        filter := bson.M{"user_id": userID}

        // Fetch enquiries from the database
        var enquiries []map[string]interface{}
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        cursor, err := EnquireCollection.Find(ctx, filter)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"Error": "something went wrong"})
            return
        }
        defer cursor.Close(ctx)
        for cursor.Next(ctx) {
            var enquire models.Enquire
            if err := cursor.Decode(&enquire); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"Error": "something went wrong"})
                return
            }

            // Fetch product details based on product_id
            var product models.Product

		prodID, err := primitive.ObjectIDFromHex(enquire.Product_id)
		if err != nil {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
			c.Abort()
			return 
		}
			
        errors := ProductCollection.FindOne(ctx, bson.M{"_id": prodID }).Decode(&product)
            if errors != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to fetch product details"})
                return
            }

            // Add product name and image to enquiry
            enquiryWithProduct := map[string]interface{}{
                "enquiry_id":   enquire.Enquire_id.Hex(),
                "user_id":      enquire.User_id,
                "product_name": product.Product_Name,
                "product_image": product.Image,
				"enquire_note": enquire.Enquiry_note,
				"enquire_quantity": enquire.Quantity,
                // Add other enquiry fields if needed
            }
            enquiries = append(enquiries, enquiryWithProduct)
        }
        if err := cursor.Err(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"Error": "something went wrong"})
            return
        }

        c.JSON(http.StatusOK, enquiries)
    }
}


func GETEnquiryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if !checkAdmin(ctx, c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin Token Not found"})
			return
		}

		var enquire []models.Enquire

		cursor, err := EnquireCollection.Find(ctx, primitive.M{})
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch"})
			return
		}
		if err := cursor.All(ctx, &enquire); err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		defer cursor.Close(ctx)

		// Enrich enquiry data with additional details
		enquiriesWithDetails := make([]map[string]interface{}, 0)

		for _, enquiry := range enquire {
			// Fetch product details based on product_id
			productDetails := getProductDetails(ctx, enquiry.Product_id)

			// Fetch user details based on user_id
			userDetails := getUserDetails(ctx, enquiry.User_id)

			// Construct enriched enquiry
			enquiryWithDetails := map[string]interface{}{
				"enquiry":   enquiry,
				"product":   productDetails,
				"user":      userDetails,
			}

			enquiriesWithDetails = append(enquiriesWithDetails, enquiryWithDetails)
		}

		c.JSON(http.StatusOK, enquiriesWithDetails)
	}
}

// Function to fetch product details based on product_id
func getProductDetails(ctx context.Context, productID string) map[string]interface{} {
	var productDetails map[string]interface{}

	id,err := primitive.ObjectIDFromHex(productID);

	if err != nil {
		log.Printf("Error parsing product ID %s: %s", productID, err.Error())
		return nil
	}

	errs := ProductCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&productDetails)
	if errs != nil {
		log.Printf("Error fetching product details for product ID %s: %s", productID, err.Error())
		return nil
	}

	return productDetails
}

// Function to fetch user details based on user_id
func getUserDetails(ctx context.Context, userID string) map[string]interface{} {
	var userDetails map[string]interface{}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Error parsing user ID %s: %s", userID, err.Error())
		return nil
	}

	errs := UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&userDetails)
	if errs != nil {
		log.Printf("Error fetching user details for user ID %s: %s", userID, err.Error())
		return nil
	}

	return userDetails
}
// GetAllRequirementMessages retrieves all RequirementMessages
func GetAllRequirementMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := RequirementMessageCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var messages []models.RequirementMessage
		if err := cursor.All(ctx, &messages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, messages)
	}
}

// CreateRequirementMessage creates a new RequirementMessage
func CreateRequirementMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var message models.RequirementMessage
		if err := c.BindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate new ObjectID
		message.Requirement_id = primitive.NewObjectID()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := RequirementMessageCollection.InsertOne(ctx, message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result.InsertedID)
	}
}
// UpdateRequirementMessage updates a RequirementMessage
func UpdateRequirementMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var message models.RequirementMessage
		if err := c.BindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": message.Requirement_id}
		update := bson.M{"$set": message}

		_, err := RequirementMessageCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "RequirementMessage updated successfully"})
	}
}

// DeleteRequirementMessage deletes a RequirementMessage
func DeleteRequirementMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var message models.RequirementMessage
		if err := c.BindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": message.Requirement_id}

		_, err := RequirementMessageCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "RequirementMessage deleted successfully"})
	}
}
