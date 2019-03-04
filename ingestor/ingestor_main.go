package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/akraievoy/tsv_load/parser"
	"github.com/akraievoy/tsv_load/proto"
	"github.com/akraievoy/tsv_load/service_utils"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	sendContext, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go service_utils.CancelOnTermSignal(cancelFunc)

	iFlags, err := ParseSenderFlags()
	if err != nil {
		log.Fatalf("failed to parse arguments: %v", err)
	}

	closeFunc, upserterClient, err := newUpserterClient(sendContext, iFlags)
	if closeFunc != nil {
		defer closeFunc()
	}
	if err != nil {
		log.Fatalf("failed to establish gRPC connection: %v", err)
	}

	reader, readerCloseFunc, err := newReader(iFlags)
	if readerCloseFunc != nil {
		defer readerCloseFunc()
	}
	if err != nil {
		log.Fatalf("failed to open data reader: %v", err)
	}

	stream, err := parser.NewRecordStream(reader)
	if err != nil {
		log.Fatalf("failed to parse CSV header: %v", err)
	}

	stats, err := stream.ForEachBatch(
		sendContext,
		iFlags.BatchSize,
		iFlags.BrokenLinesToFail,
		ingestorFunc(upserterClient),
	)

	log.Printf(
		"Completed: %d records successfully processed, %d records failed",
		stats.Successes, len(stats.FailedLineNumbers),
	)
	if err != nil {
		log.Fatalf("failed to process data fully: %v", err)
	}
}

func ingestorFunc(upserterClient proto.UpserterClient) func(ctx context.Context, rs []parser.Record) ([]uint64, error) {
	return func(ctx context.Context, rs []parser.Record) ([]uint64, error) {
		userBatch := proto.UserBatch{}
		failedLineNumbers := make([]uint64, 0)
		rsSent := make([]parser.Record, 0)

		for _, r := range rs {
			user, err := tsvToDomain(r)
			if err == nil {
				userBatch.Batch = append(userBatch.Batch, user)
				rsSent = append(rsSent, r)
			} else {
				log.Printf("row #%d format error: %v", r.LineNumber(), err)
				failedLineNumbers = append(failedLineNumbers, r.LineNumber())
			}
		}

		feedback, err := upserterClient.Upsert(ctx, &userBatch)

		if err != nil {
			log.Printf(
				"whole batch of rows #%d..%d failed: %v",
				rs[0].LineNumber(), rs[len(rs)-1].LineNumber(), err,
			)
		}

		if feedback != nil {
			if len(rsSent) != len(feedback.Feedbacks) {
				message := fmt.Sprintf(
					"whole batch of rows #%d..%d failed: feedback length mismatch",
					rs[0].LineNumber(), rs[len(rs)-1].LineNumber(),
				)
				log.Printf(message)
				return failedLineNumbers, errors.New(message)
			}

			for i, f := range feedback.Feedbacks {
				if !f.Success {
					log.Printf("remote failure for row #%d: %v", rsSent[i].LineNumber(), f.ErrorMessage)
					failedLineNumbers = append(failedLineNumbers, rsSent[i].LineNumber())
				}
			}
		}

		return failedLineNumbers, err
	}
}

func tsvToDomain(r parser.Record) (*proto.User, error) {
	idStr, err := r.GetValue("id")
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("id field is not a valid uint32: %v", err))
	}
	name, err := r.GetValue("name")
	if err != nil {
		return nil, err
	}
	email, err := r.GetValue("email")
	if err != nil {
		return nil, err
	}
	mobile, err := r.GetValue("mobile_number")
	if err != nil {
		return nil, err
	}

	var user proto.User

	user.Id = int32(id)
	user.Name = name
	user.Email = email
	user.PhoneNumber = mobile

	return &user, nil
}

func newReader(senderFlags *IngestorFlags) (*bufio.Reader, func(), error) {
	var reader *bufio.Reader = nil
	var readerCloseFunc func() = nil
	if senderFlags.InFilePath != "-" {
		opened, err := os.Open(senderFlags.InFilePath)
		if err != nil {
			return nil, nil, err
		}
		readerCloseFunc = func() {
			err := opened.Close()
			if err != nil {
				log.Print("failed to close file being read:", err)
			}
		}
		reader = bufio.NewReader(opened)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	return reader, readerCloseFunc, nil
}

func newUpserterClient(sendContext context.Context, senderFlags *IngestorFlags) (func(), proto.UpserterClient, error) {
	var conn *grpc.ClientConn
	var closeFunc func()
	var err error
	for {
		dialAddr := fmt.Sprintf("%s:%d", senderFlags.UpserterHost, senderFlags.UpserterPort)
		log.Printf("dialling %s...", dialAddr)
		conn, err = grpc.Dial(dialAddr, grpc.WithInsecure())
		if err == nil {
			closeFunc = func() {
				err := conn.Close()
				if err != nil {
					log.Print("error while closing gRPC connection:", err)
				}
			}
			log.Printf("dialling %s successful", dialAddr)
			break
		} else {
			log.Printf("failed to dial %s (retrying in 5 sec): %v", dialAddr, err)
			if service_utils.SleepCancellably(sendContext, time.Second*5) {
				return nil, nil, errors.New("cancelled while dialling to upserter")
			}
		}
	}

	upserterClient := proto.NewUpserterClient(conn)

	for {
		versionResponse, err := upserterClient.Version(sendContext, &proto.VersionRequest{})
		if err == nil {
			log.Printf("connected to upserter version: %v", versionResponse.Version)
			break
		} else {
			log.Printf("failed to query upserter version (retrying in 5 sec): %v", err)
			if service_utils.SleepCancellably(sendContext, time.Second*5) {
				return closeFunc, nil, errors.New("cancelled while verifying upserter version")
			}
		}
	}

	return closeFunc, upserterClient, nil
}
