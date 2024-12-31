package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Class represents the details of a class
type Class struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Capacity  int    `json:"capacity"`
}

// Booking represents a booking made by a member
type Booking struct {
	ClassID string `json:"class_id"`
	Name    string `json:"name"`
	Date    string `json:"date"`
}

// In-memory storage for classes and bookings
var (
	classes  = []Class{}
	bookings = []Booking{}
	mu       sync.Mutex
)

func main() {
	router := gin.Default()

	// Root route with HTML content for Create Class form
	router.GET("/", handleHomePage)

	// API routes
	router.POST("/classes", handleCreateClass)
	router.POST("/bookings", handleCreateBooking)

	// Create Class Response Page
	router.GET("/create-class-response", handleCreateClassResponse)

	// Create Booking Response Page
	router.GET("/create-booking-response", handleCreateBookingResponse)

	// Start the server
	router.Run(":8080")
}

// handleHomePage serves the homepage with the "Create Class" form
func handleHomePage(c *gin.Context) {
	// HTML content for creating a class
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Class Booking System</title>
		</head>
		<body>
			<h1>Create a Class</h1>
			<form action="/classes" method="POST">
				<label for="name">Class Name:</label>
				<input type="text" id="name" name="name" required><br><br>

				<label for="start_date">Start Date (YYYY-MM-DD):</label>
				<input type="text" id="start_date" name="start_date" required><br><br>

				<label for="end_date">End Date (YYYY-MM-DD):</label>
				<input type="text" id="end_date" name="end_date" required><br><br>

				<label for="capacity">Capacity:</label>
				<input type="number" id="capacity" name="capacity" required><br><br>

				<button type="submit">Create Class</button>
			</form>
		</body>
		</html>
	`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

// handleCreateClass handles the creation of classes
func handleCreateClass(c *gin.Context) {
	// Parse form data
	name := c.DefaultPostForm("name", "")
	startDate := c.DefaultPostForm("start_date", "")
	endDate := c.DefaultPostForm("end_date", "")
	capacityStr := c.DefaultPostForm("capacity", "")

	// Validate form fields
	if name == "" || startDate == "" || endDate == "" || capacityStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields"})
		return
	}

	// Convert capacity to integer
	capacity, err := strconv.Atoi(capacityStr)
	if err != nil || capacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid capacity"})
		return
	}

	// Validate date format
	start, errStart := time.Parse("2006-01-02", startDate)
	end, errEnd := time.Parse("2006-01-02", endDate)
	if errStart != nil || errEnd != nil || start.After(end) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range. Use YYYY-MM-DD format"})
		return
	}

	// Assign a unique ID and save the class
	class := Class{
		ID:        uuid.New().String(),
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  capacity,
	}

	// Store the class in memory
	mu.Lock()
	classes = append(classes, class)
	mu.Unlock()

	// Redirect to Create Class Response page with class ID
	c.Redirect(http.StatusFound, "/create-class-response?id="+class.ID)
}

// handleCreateClassResponse shows the response of the Create Class API and the booking form
func handleCreateClassResponse(c *gin.Context) {
	classID := c.DefaultQuery("id", "")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class ID is missing"})
		return
	}

	// Find the created class
	var createdClass Class
	mu.Lock()
	for _, class := range classes {
		if class.ID == classID {
			createdClass = class
			break
		}
	}
	mu.Unlock()

	// HTML content for the Create Class Response page with class details
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Class Created</title>
		</head>
		<body>
			<h1>Class Created Successfully</h1>
			<p><strong>Class ID:</strong> ` + createdClass.ID + `</p>
			<p><strong>Class Name:</strong> ` + createdClass.Name + `</p>
			<p><strong>Start Date:</strong> ` + createdClass.StartDate + `</p>
			<p><strong>End Date:</strong> ` + createdClass.EndDate + `</p>
			<p><strong>Capacity:</strong> ` + strconv.Itoa(createdClass.Capacity) + `</p>

			<h2>Create a Booking</h2>
			<form action="/bookings" method="POST">
				<input type="hidden" name="class_id" value="` + createdClass.ID + `">

				<label for="booking_name">Your Name:</label>
				<input type="text" id="booking_name" name="name" required><br><br>

				<label for="date">Booking Date (YYYY-MM-DD):</label>
				<input type="text" id="date" name="date" required><br><br>

				<button type="submit">Create Booking</button>
			</form>

			<br>
			<a href="/">Go Back</a>
		</body>
		</html>
	`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

// handleCreateBooking handles the creation of bookings
// handleCreateBooking handles the creation of bookings
func handleCreateBooking(c *gin.Context) {
	// Parse form data
	classID := c.DefaultPostForm("class_id", "")
	name := c.DefaultPostForm("name", "")
	date := c.DefaultPostForm("date", "")

	// Validate form fields
	if classID == "" || name == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields"})
		return
	}

	// Validate date format
	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Check if the class exists and validate the booking date
	var classExists bool
	mu.Lock()
	for _, class := range classes {
		if class.ID == classID {
			classStart, _ := time.Parse("2006-01-02", class.StartDate)
			classEnd, _ := time.Parse("2006-01-02", class.EndDate)
			if bookingDate.Before(classStart) || bookingDate.After(classEnd) {
				mu.Unlock()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Booking date is outside the class schedule"})
				return
			}
			classExists = true
			break
		}
	}
	mu.Unlock()

	if !classExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// Save booking
	booking := Booking{
		ClassID: classID,
		Name:    name,
		Date:    date,
	}

	mu.Lock()
	bookings = append(bookings, booking)
	mu.Unlock()

	// Redirect to Create Booking Response page
	c.Redirect(http.StatusFound, "/create-booking-response?class_id="+classID+"&name="+name+"&date="+date)
}

// handleCreateBookingResponse shows the response of the Create Booking API
func handleCreateBookingResponse(c *gin.Context) {
	classID := c.DefaultQuery("class_id", "")
	name := c.DefaultQuery("name", "")
	date := c.DefaultQuery("date", "")

	// HTML content for the Create Booking Response page
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Booking Created</title>
		</head>
		<body>
			<h1>Booking Created Successfully</h1>
			<p><strong>Class ID:</strong> ` + classID + `</p>
			<p><strong>Name:</strong> ` + name + `</p>
			<p><strong>Booking Date:</strong> ` + date + `</p>

			<br>
			<a href="/">Go Back</a>
		</body>
		</html>
	`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}
