# Glofox

Glofox is a simple class booking system built with Go (Golang) and the Gin web framework. It allows users to create classes, define schedules, and book slots for specific dates within the class schedule.

---

## 🚀 Features

- **Class Management**: Create classes with a name, start date, end date, and capacity.
- **Booking System**: Book slots for a specific class within its schedule.
- **Simple UI**: HTML forms for creating classes and booking slots.

---

## 📋 Prerequisites

Before start, the following must be installed:

- [Go](https://golang.org/dl/) (version 1.16 or later)
- A terminal or command-line interface

---

## 🛠️ Setup Instructions

### Clone the Repository

```bash
git clone https://github.com/your-username/glofox.git
cd glofox
```

---

## 📖 Usage

### Home Page

- Navigate to `/` to access the homepage with the **Create Class** form.

### API Endpoints

| **Method** | **Endpoint**                | **Description**                                                                                      |
|------------|-----------------------------|------------------------------------------------------------------------------------------------------|
| `GET`      | `/`                         | Serves the homepage with a form to create a class.                                                  |
| `POST`     | `/classes`                  | Creates a new class with details such as name, start date, end date, and capacity.                  |
| `GET`      | `/create-class-response`    | Displays the details of the created class and provides a form for creating bookings.                |
| `POST`     | `/bookings`                 | Creates a new booking for a specific class and date.                                                |
| `GET`      | `/create-booking-response`  | Displays booking confirmation details, including the class ID, user name, and booking date.         |

---

## 💻 Project Structure

```
glofox/
├── main.go        # Application entry point
├── go.mod         # Go module file
└── go.sum         # Go dependencies checksum file
```

---

## 🔴🟢🔵 API Endpoints

### POST API to create Class 
```
curl -X POST http://localhost:8080/classes \
-H "Content-Type: application/json" \
-d '{
  "name": "Yoga",
  "start_date": "2024-01-10",
  "end_date": "2024-01-20",
  "capacity": 15
}'
```

### POST API to create Booking
```
curl -X POST http://localhost:8080/bookings \
-H "Content-Type: application/json" \
-d '{
  "class_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "John Doe",
  "date": "2024-01-15"
}'
```
