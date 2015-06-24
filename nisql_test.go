package nisql

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type nullable struct {
	StringNVal NullString
	StringVal  string

	Int64NVal NullInt64
	Int64Val  int64

	Float64NVal NullFloat64
	Float64Val  float64

	BoolNVal NullBool
	BoolVal  bool

	TimeNVal NullTime
	TimeVal  time.Time
}

func TestInit(t *testing.T) {
	db, err := sql.Open(os.Getenv("NISQL_TEST_DIALECT"), os.Getenv("NISQL_TEST_DSN"))
	if err != nil {
		t.Fatalf("err while creating connection: %s", err.Error())
	}

	sql := `CREATE TABLE nullable (
    string_n_val VARCHAR (255) DEFAULT NULL,
    string_val VARCHAR (255) DEFAULT 'empty',
    int64_n_val BIGINT DEFAULT NULL,
    int64_val BIGINT DEFAULT 1,
    float64_n_val NUMERIC DEFAULT NULL,
    float64_val NUMERIC DEFAULT 1,
    bool_n_val BOOLEAN,
    bool_val BOOLEAN NOT NULL,
    time_n_val timestamp,
    time_val timestamp NOT NULL
)`

	if _, err = db.Exec(sql); err != nil {
		t.Fatalf("err while creating table: %s", err.Error())
	}

	sql = `INSERT INTO nullable
VALUES
    (
        NULL,
        'NULLABLE',
        NULL,
        42,
        NULL,
        12,
        NULL,
        true,
        NULL,
        NOW()
    )`

	if _, err := db.Exec(sql); err != nil {
		t.Fatalf("err while adding null item: %s", err.Error())
	}

	n := &nullable{}
	err = db.QueryRow("SELECT * FROM nullable").
		Scan(&n.StringNVal,
		&n.StringVal,
		&n.Int64NVal,
		&n.Int64Val,
		&n.Float64NVal,
		&n.Float64Val,
		&n.BoolNVal,
		&n.BoolVal,
		&n.TimeNVal,
		&n.TimeVal,
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if n.StringVal != "NULLABLE" {
		t.Fatalf("expected NULLABLE, got: ", n.StringVal)
	}

	if n.StringNVal.Valid {
		t.Fatalf("expected invalid, got valid for string_n_val")
	}

	if n.Int64Val != int64(42) {
		t.Fatalf("expected 42, got: %d", n.Int64Val)
	}

	if n.Int64NVal.Valid {
		t.Fatalf("expected invalid, got valid for int64_n_val")
	}

	if n.Float64Val != float64(12) {
		t.Fatalf("expected 12, got: %f", n.Float64Val)
	}

	if n.Float64NVal.Valid {
		t.Fatalf("expected invalid, got valid for float64_n_val")
	}

	if n.BoolVal != true {
		t.Fatalf("expected true, got: %t", n.BoolVal)
	}

	if n.BoolNVal.Valid {
		t.Fatalf("expected invalid, got valid for float64_n_val")
	}

	if n.TimeNVal.Valid {
		t.Fatalf("expected false, got: %t", n.TimeNVal)
	}

	if n.TimeVal.IsZero() {
		t.Fatalf("expected valid, got invalid for TimeVal: %+v", n.TimeVal)
	}
}
