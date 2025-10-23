package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "testing"

    "flux/database"
    "flux/models"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTaskDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil { t.Fatalf("open db: %v", err) }
    if err := db.AutoMigrate(&models.User{}, &models.Task{}); err != nil { t.Fatalf("migrate: %v", err) }
    database.DB = db
    return db
}

func TestCreateAndGetTask(t *testing.T) {
    setupTaskDB(t)

    // create user
    u := models.User{Name: "U", Email: "u@example.com", Password: "Password1!"}
    if err := database.DB.Create(&u).Error; err != nil { t.Fatal(err) }

    // create task
    body := models.Task{Title: "T", Description: "D", Status: "pending", UserID: u.ID}
    w, c := performJSONRequest(CreateTask, http.MethodPost, body)
    CreateTask(c)
    if w.Code != http.StatusCreated { t.Fatalf("expected 201, got %d", w.Code) }

    // list tasks
    w2, c2 := performJSONRequest(GetTasks, http.MethodGet, nil)
    GetTasks(c2)
    if w2.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w2.Code) }

    var tasks []models.Task
    if err := json.Unmarshal(w2.Body.Bytes(), &tasks); err != nil { t.Fatal(err) }
    if len(tasks) != 1 || tasks[0].Title != "T" { t.Fatalf("unexpected tasks: %+v", tasks) }
}

func TestGetTask_NotFoundAndFound(t *testing.T) {
    setupTaskDB(t)

    // not found
    w, c := performJSONRequest(GetTask, http.MethodGet, nil)
    c.Params = []gin.Param{{Key: "id", Value: "999"}}
    GetTask(c)
    if w.Code != http.StatusNotFound { t.Fatalf("expected 404, got %d", w.Code) }

    // create then found
    u := models.User{Name: "U", Email: "u2@example.com", Password: "Password1!"}
    if err := database.DB.Create(&u).Error; err != nil { t.Fatal(err) }
    task := models.Task{Title: "T2", UserID: u.ID}
    if err := database.DB.Create(&task).Error; err != nil { t.Fatal(err) }

    w2, c2 := performJSONRequest(GetTask, http.MethodGet, nil)
    c2.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(task.ID))}}
    GetTask(c2)
    if w2.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w2.Code) }
}

func TestUpdateAndDeleteTask(t *testing.T) {
    setupTaskDB(t)

    u := models.User{Name: "U", Email: "u3@example.com", Password: "Password1!"}
    if err := database.DB.Create(&u).Error; err != nil { t.Fatal(err) }
    task := models.Task{Title: "Old", UserID: u.ID}
    if err := database.DB.Create(&task).Error; err != nil { t.Fatal(err) }

    // update
    upd := models.Task{Title: "New"}
    w, c := performJSONRequest(UpdateTask, http.MethodPut, upd)
    c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(task.ID))}}
    UpdateTask(c)
    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }

    // delete
    w2, c2 := performJSONRequest(DeleteTask, http.MethodDelete, nil)
    c2.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(task.ID))}}
    DeleteTask(c2)
    if w2.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w2.Code) }

    // delete not found
    w3, c3 := performJSONRequest(DeleteTask, http.MethodDelete, nil)
    c3.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(task.ID))}}
    DeleteTask(c3)
    if w3.Code != http.StatusNotFound { t.Fatalf("expected 404, got %d", w3.Code) }
}
