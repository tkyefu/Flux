package handlers

import (
    "net/http"
    "strconv"
    "testing"

    "flux/database"
    "flux/models"
    "github.com/gin-gonic/gin"
)

func setupUserDB(t *testing.T) {
    t.Helper()
    db := newTestDB(t)
    database.DB = db
    // Ensure tasks table exists for Preload("Tasks") in GetUser
    if err := db.AutoMigrate(&models.Task{}); err != nil {
        t.Fatalf("migrate task: %v", err)
    }
}

func TestCreateUser_And_GetUsers(t *testing.T) {
    setupUserDB(t)

    u := models.User{Name: "U", Email: "x@example.com", Password: "Password1!"}
    w, c := performJSONRequest(CreateUser, http.MethodPost, u)
    CreateUser(c)
    if w.Code != http.StatusCreated { t.Fatalf("expected 201, got %d", w.Code) }

    w2, c2 := performJSONRequest(GetUsers, http.MethodGet, nil)
    GetUsers(c2)
    if w2.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w2.Code) }
}

func TestGetUser_NotFound_And_DeleteUser(t *testing.T) {
    setupUserDB(t)

    // not found
    w, c := performJSONRequest(GetUser, http.MethodGet, nil)
    c.Params = []gin.Param{{Key: "id", Value: "999"}}
    GetUser(c)
    if w.Code != http.StatusNotFound { t.Fatalf("expected 404, got %d", w.Code) }

    // create
    u := models.User{Name: "U2", Email: "y@example.com", Password: "Password1!"}
    if err := database.DB.Create(&u).Error; err != nil { t.Fatal(err) }

    // get
    w2, c2 := performJSONRequest(GetUser, http.MethodGet, nil)
    c2.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(u.ID))}}
    GetUser(c2)
    if w2.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w2.Code) }

    // delete
    w3, c3 := performJSONRequest(DeleteUser, http.MethodDelete, nil)
    c3.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(u.ID))}}
    DeleteUser(c3)
    if w3.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w3.Code) }
}
