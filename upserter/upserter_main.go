package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/akraievoy/tsv_load/proto"
	"github.com/akraievoy/tsv_load/service_utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type UpserterServer struct {
	db *sql.DB
}

func (us *UpserterServer) Upsert(ctx context.Context, users *proto.UserBatch) (*proto.BatchFeedback, error) {
	stmt, err := us.db.Prepare(
		"insert into USERS(ID, NAME, EMAIL, COUNTRY_CODE, MOBILE_NUMBER) " +
			"values ($1, $2, $3, $4, $5) " +
			"on conflict (id) do " +
			"update set " +
			"NAME=excluded.NAME, EMAIL=excluded.EMAIL, " +
			"COUNTRY_CODE=excluded.COUNTRY_CODE, MOBILE_NUMBER=excluded.MOBILE_NUMBER",
	)
	if err != nil {
		log.Printf("failed to prepare statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	result := make([]*proto.UpsertFeedback, 0)
	for _, u := range users.Batch {
		execRes, err := stmt.Exec(uint32(u.Id), u.Name, u.Email, "", u.PhoneNumber)
		feedback := proto.UpsertFeedback{}
		if err == nil {
			insertId, insertIdErr := execRes.LastInsertId()
			rowsAffected, rowsAffectedErr := execRes.RowsAffected()
			log.Printf("exec result %d (%v) %d (%v)", insertId, insertIdErr, rowsAffected, rowsAffectedErr)
			feedback.Success = true
		} else {
			log.Print("failed to upsert record", err)
			feedback.Success = false
			feedback.ErrorMessage = fmt.Sprintf("%v", err)
		}
		result = append(result, &feedback)
	}

	log.Printf("upserted %d records", len(users.Batch))

	batchFeedback := proto.BatchFeedback{}
	batchFeedback.Feedbacks = result
	return &batchFeedback, nil
}

func (us *UpserterServer) Version(ctx context.Context, verReq *proto.VersionRequest) (*proto.VersionResponse, error) {
	var versionResp proto.VersionResponse
	versionResp.Version = "0.0.1"
	return &versionResp, nil
}

func main() {
	mainContext, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go service_utils.CancelOnTermSignal(cancelFunc)

	uFlags, err := ParseUpserterFlags()
	if err != nil {
		log.Fatal("failed to parse arguments", err)
	}

	db, dbCloseFunc, err := initDatabaseConnPool(mainContext, uFlags)
	if dbCloseFunc != nil {
		defer dbCloseFunc()
	}
	if err != nil {
		log.Fatal("failed to init DB connection pool: %v", err)
	}

	toBindTo := fmt.Sprintf("%s:%d", uFlags.BindHost, uFlags.BindPort)
	listen, err := net.Listen("tcp", toBindTo)
	if err != nil {
		log.Fatalf("failed to bind to %s", toBindTo)
	}

	server := UpserterServer{db}
	grpcServer := grpc.NewServer()

	proto.RegisterUpserterServer(grpcServer, &server)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("failed to serve %v", err)
	}

}

func initDatabaseConnPool(mainContext context.Context, uFlags *UpserterFlags) (*sql.DB, func(), error) {
	connStr :=
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			uFlags.DBUsername, uFlags.DBPassword, uFlags.DBHost, uFlags.DBPort, uFlags.DBName,
		)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}
	dbCloseFunc := func() {
		err := db.Close()
		if err != nil {
			log.Printf("failure on closing DB: %v", err)
		}
	}

	for {
		err := db.PingContext(mainContext)
		if err != nil {
			log.Printf("database not available (will retry in 5 seconds): %v", err)
			if service_utils.SleepCancellably(mainContext, time.Second*5) {
				return db, dbCloseFunc, err
			}
		} else {
			break
		}
	}

	result, err := db.Exec(
		"create table if not exists USERS (" +
			"ID int not null primary key, " +
			"NAME text not null, " +
			"EMAIL text not null, " +
			"COUNTRY_CODE text not null, " +
			"MOBILE_NUMBER text not null " +
			")",
	)
	if err == nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("create table USER (if not exists), rows affected not available: %v", err)
		} else {
			log.Printf("create table USER (if not exists), %d rows affected", rowsAffected)
		}
	}
	return db, dbCloseFunc, err
}
