package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation-test"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		t.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func TestGetUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	user := &types.User{
		FirstName: "Serge",
		LastName:  "Basa",
		Email:     "stikkas17@gmail.com",
	}

	createdUser, _ := db.UserStore.CreateUser(context.TODO(), user)

	app := fiber.New()
	userHandler := NewUserHandler(db.UserStore)
	app.Get("/:id", userHandler.HandleGetUser)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", createdUser.ID.Hex()), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		t.Fatal(err)
	}

	if m["ID"] != createdUser.ID.Hex() {
		t.Errorf("expecting a user id to be set")
	}
	if m["encryptedPassword"] != nil {
		t.Errorf("expecting the EncryptedPassword not to be included into the json response")
	}
	if m["lastName"] != createdUser.LastName {
		t.Errorf("expected lastname %s - actual lastname %s", createdUser.LastName, m["lastName"])
	}
	if m["firstName"] != createdUser.FirstName {
		t.Errorf("expected firstname %s - actual firstname %s", createdUser.FirstName, m["firstName"])
	}
	if m["email"] != createdUser.Email {
		t.Errorf("expected email %s - actual email %s", createdUser.Email, m["email"])
	}

}

func TestPostUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(db.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "skikkas17@gmail.com",
		FirstName: "Serge",
		LastName:  "Basa",
		Password:  "lafheoihafaiofhaoh323",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var m types.User
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		t.Fatal(err)
	}

	if len(m.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(m.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included into the json response")
	}
	if m.LastName != params.LastName {
		t.Errorf("expected lastname %s - actual lastname %s", params.LastName, m.LastName)
	}
	if m.FirstName != params.FirstName {
		t.Errorf("expected firstname %s - actual firstname %s", params.FirstName, m.FirstName)
	}
	if m.Email != params.Email {
		t.Errorf("expected email %s - actual email %s", params.Email, m.Email)
	}
}
