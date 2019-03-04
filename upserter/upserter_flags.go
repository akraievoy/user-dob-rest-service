package main

import (
	"errors"
	"flag"
	"github.com/facebookgo/flagenv"
)

type UpserterFlags struct {
	BindHost   string
	BindPort   uint32
	DBHost     string
	DBPort     uint32
	DBName     string
	DBUsername string
	DBPassword string
}

func ParseUpserterFlags() (*UpserterFlags, error) {
	bindHost := flag.String(
		"bind-host",
		"localhost",
		"host resolving to network interface to bind endpoint to",
	)

	bindPort := flag.Uint(
		"bind-port",
		8080,
		"port number to bind endpoint to",
	)

	dbHost := flag.String(
		"db-host",
		"localhost",
		"dial this host to talk to postgres",
	)

	dbPort := flag.Uint(
		"db-port",
		5432,
		"dial this port number to talk to postgres",
	)

	dbName := flag.String(
		"db-name",
		"upserter",
		"use given database after connecting to db service",
	)

	dbUsername := flag.String(
		"db-username",
		"upserter",
		"authienticate with database using given username",
	)

	dbPassword := flag.String(
		"db-password",
		"fillmein", // LATER ahaha it does not work as standard docker trusts local
		"authienticate with database using given password",
	)

	flagenv.Parse()
	flag.Parse()

	if *bindPort == uint(0) {
		return nil, errors.New("bind-port is zero")
	}
	if *bindPort > uint(65535) {
		return nil, errors.New("bind-port exceeds 65535")
	}

	if *dbPort == uint(0) {
		return nil, errors.New("db-port is zero")
	}
	if *dbPort > uint(65535) {
		return nil, errors.New("db-port exceeds 65535")
	}

	return &UpserterFlags{
		*bindHost,
		uint32(*bindPort),
		*dbHost,
		uint32(*dbPort),
		*dbName,
		*dbUsername,
		*dbPassword,
	}, nil
}
