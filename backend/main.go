// main.go
package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	pb "buoyboy/proto"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/proto"
)

var db *sql.DB

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS buoy_readings (
        buoy_id INTEGER,
        report_id INTEGER,
        timestamp INTEGER,
        latitude REAL,
        longitude REAL,
        depth_meters REAL,
        temp_celsius REAL,
        accel_x REAL,
        accel_y REAL,
        accel_z REAL,
        battery_level_percent INTEGER,
        PRIMARY KEY (buoy_id, report_id)
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var err error
	db, err = initDB("./bouy_readings.db")
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	http.HandleFunc("/buoy", handleBuoyReading)
	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBuoyReading(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	message := &pb.BuoyStatus{}
	if err := proto.Unmarshal(body, message); err != nil {
		http.Error(w, "Error parsing protobuf message", http.StatusBadRequest)
		return
	}

	// Insert data into database
	insertSQL := `
    INSERT INTO buoy_readings (
        buoy_id, report_id, timestamp, latitude, longitude, 
        depth_meters, temp_celsius, accel_x, accel_y, accel_z, 
        battery_level_percent
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(insertSQL,
		message.BuoyId,
		message.ReportId,
		message.Timestamp,
		message.Latitude,
		message.Longitude,
		message.DepthMeters,
		message.TempCelsius,
		message.AccelX,
		message.AccelY,
		message.AccelZ,
		message.BatteryLevelPercent,
	)

	if err != nil {
		log.Printf("Error inserting data: %v", err)
		http.Error(w, "Error saving data", http.StatusInternalServerError)
		return
	}

	log.Printf("Saved buoy reading: %+v", message)
	w.WriteHeader(http.StatusOK)
}
