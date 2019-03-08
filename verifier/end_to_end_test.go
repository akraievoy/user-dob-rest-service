package verifier

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/akraievoy/tsv_load/service_utils"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
	"time"
)

type verifierCtxKey string

var verifierCtxKeyDB = verifierCtxKey("db")
var verifierCtxKeyRetries = verifierCtxKey("retries")

func setupGlobal() (context.Context, func()) {
	retries := uint8(20)
	ctx, cancelFunc :=
		context.WithCancel(
			context.WithValue(
				context.Background(),
				verifierCtxKeyRetries,
				&retries,
			),
		)
	go func() { service_utils.CancelOnTermSignal(cancelFunc) }()

	connStr :=
		fmt.Sprintf(
			"postgres://upserter:bazinga@postgres:5432/upserter?sslmode=disable",
		)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("after sql.Open: err = %v", err)
	}

	dbCloseFunc := func() {
		err := db.Close()
		if err != nil {
			log.Print(err)
		}
	}

	for {
		timeout, cancelFunc := context.WithTimeout(ctx, time.Second*3)
		defer cancelFunc()

		err := db.PingContext(timeout)
		if err == nil {
			break
		} else {
			retries -= 1
			if retries == 0 {
				log.Fatalf("unable to connect to database: %v", err)
			}
			log.Printf("waiting for successful ping: %v", err)
			if service_utils.SleepCancellably(ctx, time.Second*3) {
				log.Fatalf("test cancelled, %d retries left", retries)
			}
		}
	}

	return context.WithValue(ctx, verifierCtxKeyDB, db), dbCloseFunc
}

var testCtx context.Context

func TestMain(m *testing.M) {
	log.Printf("started e2e test")
	var dbCloseFunc func()
	testCtx, dbCloseFunc = setupGlobal()
	defer dbCloseFunc()
	os.Exit(m.Run())
}

func TestAllRecordsImported(t *testing.T) {
	query := "select count(*) from users"
	expectedMinimum := 100

	awaitMinimum(t, query, expectedMinimum)
}

func TestDamienRecordsUpserted(t *testing.T) {
	query := "select count(*) from users where id=33 and name='Damien'"
	expectedMinimum := 1

	awaitMinimum(t, query, expectedMinimum)
}

func TestRecordsHaveCountryCode(t *testing.T) {
	query := "select count(*) from users where country_code='+44'"
	expectedMinimum := 100

	awaitMinimum(t, query, expectedMinimum)
}

func TestRecordsHaveNormalizedPhoneNumbers(t *testing.T) {
	query := "select count(*) from users where mobile_number not like '0%' and mobile_number not like '(%'"
	expectedMinimum := 100

	awaitMinimum(t, query, expectedMinimum)
}

func awaitMinimum(t *testing.T, query string, expectedMinimum int) {
	db := testCtx.Value(verifierCtxKeyDB).(*sql.DB)
	retries := testCtx.Value(verifierCtxKeyRetries).(*uint8)
	for *retries > 0 {
		success := func() bool {
			rows, err := db.Query(query)
			if err == nil {
				//	ouch, to force this defer before next retry I need an extra closure
				defer func() { _ = rows.Close() }()

				if !rows.Next() {
					t.Fatal("count query should return at least one row")
				}

				countRes := uint64(0)
				if err := rows.Scan(&countRes); err != nil {
					t.Fatalf("count query should return at least one numeric column: %v", err)
				}

				if countRes >= uint64(expectedMinimum) {
					return true
				} else {
					*retries -= 1
				}
			} else {
				*retries -= 1
			}
			return false
		}()
		if success {
			break
		}
		t.Logf("sleeping for 3 seconds for query '%s'", query)
		if service_utils.SleepCancellably(testCtx, time.Second*3) {
			t.Fatalf("test for query '%s' cancelled", query)
		}
	}
	if *retries == 0 {
		t.Errorf("haven't seen all records imported")
	}
}
