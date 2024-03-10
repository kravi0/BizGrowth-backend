package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kravi0/BizGrowth-backend/controllers"
	"github.com/kravi0/BizGrowth-backend/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// Use the provided gin.Engine instead of creating a new router
	//incomingRoutes.Use(cors.Default())
	router := gin.Default()
	router.Use(gin.Logger())
	config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Token","token"} // Add "Token" header
    incomingRoutes.Use(cors.New(config))
	incomingRoutes.GET("/getcategory", controllers.GetCategory())
	incomingRoutes.POST("/addcategory", controllers.AddCategory())
	incomingRoutes.PUT("/updatecategory", controllers.EditCategory())
	incomingRoutes.GET("/getproduct", controllers.GetProduct())
	incomingRoutes.GET("/product", controllers.SearchProductByQuery())
	incomingRoutes.POST("/update-user", controllers.UpdateUserDetails())
	incomingRoutes.POST("/post-requirement",controllers.CreateRequirementMessage())
	incomingRoutes.GET("/get-productReference",controllers.FetchProductsAndReferencesHandler())
	incomingRoutes.GET("/product-reference", controllers.GetProductReferenceHandler())
	incomingRoutes.POST("/add-admin",controllers.RegisterAdmin())
	
	//incomingRoutes.GET("/getseller", controllers.GetSeller())
	incomingRoutes.POST("/sendOTP", controllers.SetOtpHandler())
	incomingRoutes.POST("/validate", controllers.ValidateOtpHandler())
	incomingRoutes.POST("/sellerOTPRegistration", controllers.SetSellerOtpHandler())
	incomingRoutes.POST("/validatesellerotpin", controllers.ValidateSellerOtpHandler())
	incomingRoutes.POST("/registerSellerDetaills", controllers.SellerRegistration())
	incomingRoutes.POST("/seller-login", controllers.Login())
	
	incomingRoutes.POST("/validatesellerotp", controllers.ValidatePasswordOTP())
	incomingRoutes.Use(middleware.UserAuthentication())
	incomingRoutes.POST("/product-enquiry", controllers.EnquiryHandler())

	
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.POST("/admin/approveSeller", controllers.ApproveSeller())
	incomingRoutes.PUT("/admin/updateProduct", controllers.UpdateProduct())
	incomingRoutes.POST("/admin/add-product", controllers.ProductViewerAdmin())
	incomingRoutes.GET("/admin/get-enquiry", controllers.GETEnquiryHandler())
	incomingRoutes.GET("/admin/getseller", controllers.GetSeller())
	incomingRoutes.GET("/admin/approve-product", controllers.ApproveProduct())
	incomingRoutes.GET("/seller/info", controllers.LoadSeller())
	incomingRoutes.DELETE("/admin/delete-product", controllers.DeleteProduct())
	incomingRoutes.POST("/seller/update-product", controllers.AddProductReferenceHandler())

	//incomingRoutes.GET("/getcategory", controllers.GetCategory())
}
