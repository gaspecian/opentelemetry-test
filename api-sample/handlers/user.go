package handlers

import (
	"api-sample/config"
	"api-sample/database"
	"api-sample/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type UserHandler struct {
	db      *database.MongoDB
	metrics *config.Metrics
}

func NewUserHandler(db *database.MongoDB, metrics *config.Metrics) *UserHandler {
	return &UserHandler{db: db, metrics: metrics}
}

func (h *UserHandler) recordMetrics(ctx context.Context, method, endpoint string, statusCode int, duration float64) {
	attrs := attribute.NewSet(
		attribute.String("http.method", method),
		attribute.String("http.endpoint", endpoint),
		attribute.Int("http.status_code", statusCode),
	)
	
	h.metrics.RequestCounter.Add(ctx, 1, metric.WithAttributeSet(attrs))
	h.metrics.RequestDuration.Record(ctx, duration, metric.WithAttributeSet(attrs))
	
	if statusCode >= 400 {
		h.metrics.ErrorCounter.Add(ctx, 1, metric.WithAttributeSet(attrs))
	}
}

func (h *UserHandler) recordDBMetrics(ctx context.Context, operation string, duration float64, err error) {
	attrs := attribute.NewSet(
		attribute.String("db.operation", operation),
	)
	
	h.metrics.DBQueryCounter.Add(ctx, 1, metric.WithAttributeSet(attrs))
	h.metrics.DBQueryDuration.Record(ctx, duration, metric.WithAttributeSet(attrs))
	
	if err != nil {
		h.metrics.DBErrorCounter.Add(ctx, 1, metric.WithAttributeSet(attrs))
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	tracer := otel.Tracer("api-sample")
	
	h.metrics.IncrementActiveRequests()
	defer h.metrics.DecrementActiveRequests()
	
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()
	
	config.LogWithTrace(ctx, "Creating new user")
	
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.LogWithTrace(ctx, "Failed to decode user: "+err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.recordMetrics(ctx, "POST", "/users", http.StatusBadRequest, time.Since(start).Seconds()*1000)
		return
	}

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dbStart := time.Now()
	result, err := h.db.DB.Collection("users").InsertOne(dbCtx, user)
	h.recordDBMetrics(ctx, "insert", time.Since(dbStart).Seconds()*1000, err)
	
	if err != nil {
		config.LogWithTrace(ctx, "Failed to insert user: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.recordMetrics(ctx, "POST", "/users", http.StatusInternalServerError, time.Since(start).Seconds()*1000)
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	config.LogWithTrace(ctx, "User created successfully: "+user.ID.Hex())
	
	h.metrics.UserCreatedCounter.Add(ctx, 1)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	h.recordMetrics(ctx, "POST", "/users", http.StatusCreated, time.Since(start).Seconds()*1000)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	tracer := otel.Tracer("api-sample")
	
	h.metrics.IncrementActiveRequests()
	defer h.metrics.DecrementActiveRequests()
	
	ctx, span := tracer.Start(ctx, "GetUser")
	defer span.End()
	
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		config.LogWithTrace(ctx, "Invalid user ID: "+vars["id"])
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.recordMetrics(ctx, "GET", "/users/{id}", http.StatusBadRequest, time.Since(start).Seconds()*1000)
		return
	}
	
	span.SetAttributes(attribute.String("user.id", id.Hex()))
	config.LogWithTrace(ctx, "Fetching user: "+id.Hex())

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	dbStart := time.Now()
	err = h.db.DB.Collection("users").FindOne(dbCtx, bson.M{"_id": id}).Decode(&user)
	h.recordDBMetrics(ctx, "findOne", time.Since(dbStart).Seconds()*1000, err)
	
	if err != nil {
		config.LogWithTrace(ctx, "User not found: "+id.Hex())
		http.Error(w, "User not found", http.StatusNotFound)
		h.recordMetrics(ctx, "GET", "/users/{id}", http.StatusNotFound, time.Since(start).Seconds()*1000)
		return
	}

	config.LogWithTrace(ctx, "User found: "+id.Hex())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	h.recordMetrics(ctx, "GET", "/users/{id}", http.StatusOK, time.Since(start).Seconds()*1000)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	tracer := otel.Tracer("api-sample")
	
	h.metrics.IncrementActiveRequests()
	defer h.metrics.DecrementActiveRequests()
	
	ctx, span := tracer.Start(ctx, "ListUsers")
	defer span.End()
	
	config.LogWithTrace(ctx, "Listing all users")

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dbStart := time.Now()
	cursor, err := h.db.DB.Collection("users").Find(dbCtx, bson.M{})
	h.recordDBMetrics(ctx, "find", time.Since(dbStart).Seconds()*1000, err)
	
	if err != nil {
		config.LogWithTrace(ctx, "Failed to list users: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.recordMetrics(ctx, "GET", "/users", http.StatusInternalServerError, time.Since(start).Seconds()*1000)
		return
	}
	defer cursor.Close(dbCtx)

	var users []models.User
	if err := cursor.All(dbCtx, &users); err != nil {
		config.LogWithTrace(ctx, "Failed to decode users: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.recordMetrics(ctx, "GET", "/users", http.StatusInternalServerError, time.Since(start).Seconds()*1000)
		return
	}

	config.LogWithTrace(ctx, "Users listed successfully")
	span.SetAttributes(attribute.Int("user.count", len(users)))
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	h.recordMetrics(ctx, "GET", "/users", http.StatusOK, time.Since(start).Seconds()*1000)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	tracer := otel.Tracer("api-sample")
	
	h.metrics.IncrementActiveRequests()
	defer h.metrics.DecrementActiveRequests()
	
	ctx, span := tracer.Start(ctx, "UpdateUser")
	defer span.End()
	
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		config.LogWithTrace(ctx, "Invalid user ID: "+vars["id"])
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.recordMetrics(ctx, "PUT", "/users/{id}", http.StatusBadRequest, time.Since(start).Seconds()*1000)
		return
	}
	
	span.SetAttributes(attribute.String("user.id", id.Hex()))
	config.LogWithTrace(ctx, "Updating user: "+id.Hex())

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.LogWithTrace(ctx, "Failed to decode user: "+err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.recordMetrics(ctx, "PUT", "/users/{id}", http.StatusBadRequest, time.Since(start).Seconds()*1000)
		return
	}

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"name": user.Name, "email": user.Email}}
	dbStart := time.Now()
	result, err := h.db.DB.Collection("users").UpdateOne(dbCtx, bson.M{"_id": id}, update)
	h.recordDBMetrics(ctx, "updateOne", time.Since(dbStart).Seconds()*1000, err)
	
	if err != nil {
		config.LogWithTrace(ctx, "Failed to update user: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.recordMetrics(ctx, "PUT", "/users/{id}", http.StatusInternalServerError, time.Since(start).Seconds()*1000)
		return
	}

	if result.MatchedCount == 0 {
		config.LogWithTrace(ctx, "User not found: "+id.Hex())
		http.Error(w, "User not found", http.StatusNotFound)
		h.recordMetrics(ctx, "PUT", "/users/{id}", http.StatusNotFound, time.Since(start).Seconds()*1000)
		return
	}

	user.ID = id
	config.LogWithTrace(ctx, "User updated successfully: "+id.Hex())
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	h.recordMetrics(ctx, "PUT", "/users/{id}", http.StatusOK, time.Since(start).Seconds()*1000)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	tracer := otel.Tracer("api-sample")
	
	h.metrics.IncrementActiveRequests()
	defer h.metrics.DecrementActiveRequests()
	
	ctx, span := tracer.Start(ctx, "DeleteUser")
	defer span.End()
	
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		config.LogWithTrace(ctx, "Invalid user ID: "+vars["id"])
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		h.recordMetrics(ctx, "DELETE", "/users/{id}", http.StatusBadRequest, time.Since(start).Seconds()*1000)
		return
	}
	
	span.SetAttributes(attribute.String("user.id", id.Hex()))
	config.LogWithTrace(ctx, "Deleting user: "+id.Hex())

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dbStart := time.Now()
	result, err := h.db.DB.Collection("users").DeleteOne(dbCtx, bson.M{"_id": id})
	h.recordDBMetrics(ctx, "deleteOne", time.Since(dbStart).Seconds()*1000, err)
	
	if err != nil {
		config.LogWithTrace(ctx, "Failed to delete user: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.recordMetrics(ctx, "DELETE", "/users/{id}", http.StatusInternalServerError, time.Since(start).Seconds()*1000)
		return
	}

	if result.DeletedCount == 0 {
		config.LogWithTrace(ctx, "User not found: "+id.Hex())
		http.Error(w, "User not found", http.StatusNotFound)
		h.recordMetrics(ctx, "DELETE", "/users/{id}", http.StatusNotFound, time.Since(start).Seconds()*1000)
		return
	}

	h.metrics.UserDeletedCounter.Add(ctx, 1)
	
	config.LogWithTrace(ctx, "User deleted successfully: "+id.Hex())
	w.WriteHeader(http.StatusNoContent)
	h.recordMetrics(ctx, "DELETE", "/users/{id}", http.StatusNoContent, time.Since(start).Seconds()*1000)
}
