package main

import (
	"flag"
	"errors"
)

type SenderFlags struct {
	InFilePath  string
	BatchSize   uint32
	UpserterHost string
	UpserterPort uint32
}

func ParseSenderFlags() (*SenderFlags, error) {
	inFilePath := flag.String(
		"in-file-path",
		"",
		"path to input file, emptry string denotes stdin",
	)

	batchSize := flag.Uint(
		"batch-size",
		64,
		"upsert batch size, up to 1024, inclusive",
	)

	upserterUrl := flag.String(
		"upserter-host",
		"upserter",
		"upserter service host to send your data to",
	)

	upserterPort := flag.Uint(
		"upserter-port",
		8080,
		"upserter service port to send your data to",
	)

	flag.Parse()

	if *batchSize > uint(1024) {
		return nil, errors.New("batch-size exceeds 1024")
	}
	if *upserterPort == uint(0) {
		return nil, errors.New("upserter-port is zero")
	}
	if *upserterPort > uint(65535) {
		return nil, errors.New("upserter-port exceeds 65535")
	}

	return &SenderFlags{
		*inFilePath,
		uint32(*batchSize),
		*upserterUrl,
		uint32(*upserterPort),
	}, nil
}
