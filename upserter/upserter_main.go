package main

import (
	"log"
	"github.com/akraievoy/tsv_load/proto"
	"context"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"database/sql"

	_ "github.com/lib/pq"
)

type mockDbRecord struct {
	name  string
	email string
	phone string
}

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
	uFlags, err := ParseUpserterFlags()
	if err != nil {
		log.Fatal("failed to parse arguments", err)
	}

	toBindTo := fmt.Sprintf("%s:%d", uFlags.BindHost, uFlags.BindPort)
	listen, err := net.Listen("tcp", toBindTo)
	if err != nil {
		log.Fatalf("failed to bind to %s", toBindTo)
	}

	connStr :=
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			uFlags.DBUsername, uFlags.DBPassword, uFlags.DBHost, uFlags.DBPort, uFlags.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("something is terribly fishy with db connection setup", err)
	}

	defer db.Close()

	result, err := db.Exec(
		"create table if not exists USERS (" +
			"ID int not null primary key, " +
			"NAME text not null, " +
			"EMAIL text not null, " +
			"COUNTRY_CODE text not null, " +
			"MOBILE_NUMBER text not null " +
			")",
	)
	if err != nil {
		log.Fatal("failed to connect/create users table", err)
	} else {
		log.Printf("create table result: %v", result)
	}

	server := UpserterServer{db}
	grpcServer := grpc.NewServer()

	proto.RegisterUpserterServer(grpcServer, &server)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("failed to serve %v", err)
	}

}
