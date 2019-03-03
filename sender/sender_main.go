package main

import (
	"log"
	"os"
	"fmt"
)

import (
	"github.com/akraievoy/tsv_load/parser"
	"google.golang.org/grpc"
	"bufio"
	"context"
	"github.com/akraievoy/tsv_load/proto"
	"strconv"
	"errors"
)

func main() {
	senderFlags, err := ParseSenderFlags()
	if err != nil {
		log.Fatal("failed to parse arguments", err)
	}

	var reader *bufio.Reader = nil
	if senderFlags.InFilePath != "" {
		opened, err := os.Open(senderFlags.InFilePath)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed on opening file %s", senderFlags.InFilePath), err)
		}
		defer opened.Close()

		reader = bufio.NewReader(opened)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	stream, err := parser.NewRecordStream(reader)
	if err != nil {
		log.Fatal("failed to parse CSV header", err)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", senderFlags.UpserterHost, senderFlags.UpserterPort),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("failed to dial", err)
	}
	defer conn.Close()

	upserterClient := proto.NewUpserterClient(conn)

	sendContext, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	stream.ForEachBatch(
		sendContext,
		senderFlags.BatchSize,
		func(ctx context.Context, rs []parser.Record) []error {
			userBatch := proto.UserBatch{}

			for _, r := range rs {
				idStr, idStrErr := r.GetValue("id")
				name, nameErr := r.GetValue("name")
				email, emailErr := r.GetValue("email")
				mobile, mobileErr := r.GetValue("mobile_number")
				if idStrErr != nil || nameErr != nil || emailErr != nil || mobileErr != nil {
					log.Printf("some fields missing: %v", []error{idStrErr, nameErr, emailErr, mobileErr})
				}

				id, err := strconv.ParseUint(idStr, 10, 32)
				if err != nil {
					log.Print("id field is not a valid uint32", err)
				}

				var user proto.User
				user.Id = int32(id)
				user.Name = name
				user.Email = email
				user.PhoneNumber = mobile
				userBatch.Batch = append(userBatch.Batch, &user)
			}

			result := make([]error, 0)
			feedback, err := upserterClient.Upsert(ctx, &userBatch)

			if err != nil {
				log.Printf("whole batch failed with", err)
				for range rs {
					result = append(result, err)
				}
			} else {
				for _, f := range feedback.Feedbacks {
					if f.Success {
						log.Printf("got one success")
						result = append(result, nil)
					} else {
						log.Printf("got one failure", f.ErrorMessage)
						result = append(result, errors.New(f.ErrorMessage))
					}
				}
			}

			return result
		},
	)
}
