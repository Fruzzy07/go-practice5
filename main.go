package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Job struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Company   string    `db:"company" json:"company"`
	Salary    int       `db:"salary" json:"salary"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type JobsResponse struct {
	Items       []Job `json:"items"`
	NextAfterID *int  `json:"next_after_id,omitempty"`
}

var db *sqlx.DB

func main() {
	var err error
	dsn := "host=localhost port=5432 user=postgres password=05110202 dbname=jobs sslmode=disable"
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	http.HandleFunc("/jobs", getJobsHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")

	company := strings.TrimSpace(r.URL.Query().Get("company"))
	afterIDStr := strings.TrimSpace(r.URL.Query().Get("after_id"))
	limitStr := strings.TrimSpace(r.URL.Query().Get("limit"))

	limit := 10
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}

	query := "SELECT id, title, company, salary, created_at FROM jobs WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if company != "" {
		query += " AND company = $" + strconv.Itoa(argIdx)
		args = append(args, company)
		argIdx++
	}

	if afterIDStr != "" {
		afterID, err := strconv.Atoi(afterIDStr)
		if err == nil {
			query += " AND (created_at, id) < (SELECT created_at, id FROM jobs WHERE id = $" + strconv.Itoa(argIdx) + ")"
			args = append(args, afterID)
			argIdx++
		}
	}

	query += " ORDER BY created_at DESC, id DESC"

	query += " LIMIT $" + strconv.Itoa(argIdx)
	args = append(args, limit)
	argIdx++

	var jobs []Job
	if err := db.Select(&jobs, query, args...); err != nil {
		http.Error(w, fmt.Sprintf("query error: %v", err), http.StatusInternalServerError)
		return
	}

	var nextAfterID *int
	if len(jobs) > 0 {
		last := jobs[len(jobs)-1].ID
		nextAfterID = &last
	}

	elapsed := time.Since(start)
	w.Header().Set("X-Query-Time", fmt.Sprintf("%dms", elapsed.Milliseconds()))

	resp := JobsResponse{
		Items:       jobs,
		NextAfterID: nextAfterID,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
		return
	}
}
