// POST API to create Class 

curl -X POST http://localhost:8080/classes \
-H "Content-Type: application/json" \
-d '{
  "name": "Yoga",
  "start_date": "2024-01-10",
  "end_date": "2024-01-20",
  "capacity": 15
}'


// POST API to create Booking

curl -X POST http://localhost:8080/bookings \
-H "Content-Type: application/json" \
-d '{
  "class_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "John Doe",
  "date": "2024-01-15"
}'