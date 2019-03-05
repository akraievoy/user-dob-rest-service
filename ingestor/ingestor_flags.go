package main

import (
	"errors"
	"flag"
	"github.com/facebookgo/flagenv"
	"time"
)

type IngestorFlags struct {
	InFilePath        string
	BrokenLinesToFail uint32
	BatchSize         uint32
	UpserterHost      string
	UpserterPort      uint32
	FinalSleep        time.Duration
}

func ParseSenderFlags() (*IngestorFlags, error) {
	inFilePath := flag.String(
		"in-file-path",
		"-",
		"path to input file, dash denotes stdin",
	)

	brokenLinesToFail := flag.Uint(
		"broken-lines-to-fail",
		64,
		"after this many broken lines process fails with error, up to 8192, inclusive",
	)

	batchSize := flag.Uint(
		"batch-size",
		32,
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

	finalSleep := flag.Duration(
		"final-sleep",
		100*time.Millisecond,
		"duration of final sleep",
	)

	flagenv.Parse()
	flag.Parse()

	if *brokenLinesToFail > uint(8192) {
		return nil, errors.New("broken-lines-to-fail exceeds 8192")
	}
	if *batchSize > uint(1024) {
		return nil, errors.New("batch-size exceeds 1024")
	}
	if *upserterPort == uint(0) {
		return nil, errors.New("upserter-port is zero")
	}
	if *upserterPort > uint(65535) {
		return nil, errors.New("upserter-port exceeds 65535")
	}

	return &IngestorFlags{
		*inFilePath,
		uint32(*brokenLinesToFail),
		uint32(*batchSize),
		*upserterUrl,
		uint32(*upserterPort),
		*finalSleep,
	}, nil
}
