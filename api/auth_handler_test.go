package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func insertTestUser(t *testing.T, store db.UserStore) *types.User {
	encpw, err := bcrypt.GenerateFromPassword([]byte("supersecurepassword"), 12)
	if err != nil {
		t.Fatal(err)
	}
	user := &types.User{
		FirstName:         "james",
		LastName:          "foo",
		Email:             "james@foo.com",
		EncryptedPassword: string(encpw),
	}
	_, err = store.CreateUser(context.TODO(), user)
	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HanldeAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatalf("Expected the JWT token to be present in the auth response")
	}
}
