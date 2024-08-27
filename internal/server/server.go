package server

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt"
	"github.com/ursuldaniel/go-market/internal/domain/models"
)

type Storage interface {
	RegisterUser(username, password, email string) error
	LoginUser(username, password string) (int, error)
	GetUserProfile(userId int) (models.User, error)

	AddProduct(name, description string, price, quantity int) error
	GetAllProducts() ([]models.Product, error)
	GetProductById(productId int) (models.Product, error)
	UpdateProduct(productId int, name, description string, price, quantity int) error
	DeleteProduct(productId int) error

	MakePurchase(userID, productID, quantity int) error
	GetUserPurchases(userID int) ([]models.Purchase, error)
	GetProductPurchases(productID int) ([]models.Purchase, error)
}

type Server struct {
	addr     string
	store    Storage
	validate *validator.Validate
}

func NewServer(addr string, store Storage) *Server {
	return &Server{
		addr:     addr,
		store:    store,
		validate: validator.New(),
	}
}

func (s *Server) Run() error {
	app := gin.Default()

	usersRoutes := app.Group("/users")
	usersRoutes.POST("/register", s.handleRegisterUser)
	usersRoutes.POST("/login", s.handleLoginUser)
	usersRoutes.GET("/:id", JWTAuthAdmin(s), s.handleGetUserProfile)
	usersRoutes.GET("/profile", JWTAuthUser(s), s.handleProfile)

	productsRoutes := app.Group("/products", JWTAuthUser(s))
	productsRoutes.POST("/", JWTAuthAdmin(s), s.handleAddProduct)
	productsRoutes.GET("/list", s.handleGetAllProducts)
	productsRoutes.GET("/:id", s.handleGetProductById)
	productsRoutes.PUT("/:id", JWTAuthAdmin(s), s.handleUpdateProduct)
	productsRoutes.DELETE(":id", JWTAuthAdmin(s), s.handleDeleteProduct)

	purchasesRoutes := app.Group("/purchases", JWTAuthUser(s))
	purchasesRoutes.POST("/:id", s.handleMakePurchase)
	purchasesRoutes.GET("/list", s.handleGetUserPurchases)
	purchasesRoutes.GET("/list/:id", JWTAuthAdmin(s), s.handleGetProductPurchases)

	return app.Run(s.addr)
}

func ParseId(idParam string) (int, error) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func CreateUserToken(id int) (string, error) {
	claims := &jwt.MapClaims{
		"id":        id,
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	}

	secret := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func JWTAuthUser(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header["Authorization"]
		if tokenString == nil {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Authorization token is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[0], func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Invalid token claims"})
			c.Abort()
			return
		}

		id, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusForbidden, models.Response{Message: "Unauthorized access to the account"})
			c.Abort()
			return
		}

		c.Set("id", int(id))
		c.Next()
	}
}

func CreateAdminToken(id int) (string, error) {
	claims := &jwt.MapClaims{
		"id":        id,
		"role":      "admin",
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	}

	secret := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func JWTAuthAdmin(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header["Authorization"]
		if tokenString == nil {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Authorization token is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[0], func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.Response{Message: "Invalid token claims"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, models.Response{Message: "Unauthorized access to the account"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, models.Response{Message: "Unauthorized access to the account"})
			c.Abort()
			return
		}

		c.Next()
	}
}
