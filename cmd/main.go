package main

import (
	"context"
	"flag"
	"hotel-reservation/api"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi    = "mongodb://localhost:27017"
	dbname   = "hotel-reservation"
	userColl = "users"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(map[string]string{"error": err.Error()})
		},
	})
	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication)

	userStore := db.NewMongoUserStore(client, db.DBNAME)
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(db.NewMongoHotelStore(client, db.DBNAME), db.NewMongoRoomStore(client, db.DBNAME))
	authHandler := api.NewAuthHandler(userStore)

	// auth
	auth.Post("/auth", authHandler.HanldeAuthenticate)

	// users handlers
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)
	log.Fatal(app.Listen(*listenAddr))
}
