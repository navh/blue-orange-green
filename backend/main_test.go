// main_test.go
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	pb "buoyboy/proto"

	"google.golang.org/protobuf/proto"
)

func TestMain(m *testing.M) {
	// Set up
	var err error
	db, err = initDB(":memory:") // Use in-memory SQLite for testing
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Run tests
	code := m.Run()

	// Clean up
	db.Close()
	os.Exit(code)
}

func TestHandleBuoyReading(t *testing.T) {
	// Create test message
	testBuoy := &pb.BuoyStatus{
		BuoyId:              42,
		ReportId:            12345,
		Timestamp:           time.Now().Unix(),
		Latitude:            37.805444,
		Longitude:           -122.890564,
		DepthMeters:         0,
		TempCelsius:         8.5,
		AccelX:              0.03,
		AccelY:              -0.01,
		AccelZ:              0.01,
		BatteryLevelPercent: 85,
	}

	// Serialize it
	data, err := proto.Marshal(testBuoy)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf: %v", err)
	}

	// Send it to test server
	req := httptest.NewRequest(http.MethodPost, "/buoy", bytes.NewReader(data))
	w := httptest.NewRecorder()
	handleBuoyReading(w, req)

	// Check we got OK response
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", w.Result().Status)
	}

	// Verify the data was actually saved
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM buoy_readings WHERE buoy_id = ? AND report_id = ?",
		testBuoy.BuoyId, testBuoy.ReportId).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 record in database, got %d", count)
	}
}
