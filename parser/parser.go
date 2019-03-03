package parser

import (
	"context"
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type Record interface {
	GetValue(key string) (string, error)
	LineNumber() uint64 // LATER this feels redundant
}

type ProcessingStats struct {
	Successes         uint64
	FailedLineNumbers []uint64
}

type RecordStream interface {
	ForEachBatch(
		ctx context.Context,
		batchSize uint32,
		callback func(context.Context, []Record) []error,
	) (ProcessingStats, error)
}

func NewRecordStream(reader *bufio.Reader) (RecordStream, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = 0
	schemaList, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	schemaMap := make(map[string]uint8, len(schemaList))
	for colIndex, colName := range schemaList {
		schemaMap[colName] = uint8(colIndex)
	}
	return &csvRecordStream{schemaMap, csvReader, 2}, nil
}

type csvRecordStream struct {
	schemaMap  map[string]uint8
	reader     *csv.Reader
	lineNumber uint64
}

func (rs csvRecordStream) ForEachBatch(
	ctx context.Context,
	batchSize uint32,
	callback func(context.Context, []Record) []error,
) (ProcessingStats, error) {
	buffer := make([]Record, 0, batchSize)
	successes := uint64(0)
	failedLineNumbers := make([]uint64, 0)
	thereIsMore := true
	for thereIsMore {
		values, err := rs.reader.Read()

		select {
		case <-ctx.Done():
			thereIsMore = false
			return ProcessingStats{ successes, failedLineNumbers}, ctx.Err()
		default:
		}

		if err == io.EOF {
			thereIsMore = false
		} else if err != nil {
			//	LATER recover from some CSV-specific errors as well (should we send last pre-bail buffer or not)?
			return ProcessingStats{ successes, failedLineNumbers}, err
		}

		rs.lineNumber += 1
		if len(values) > 0 { // LATER test for this, KEK
			buffer = append(buffer, &csvRecord{values, &rs.schemaMap, rs.lineNumber})
		}

		if len(buffer) == int(batchSize) || !thereIsMore {
			callbackErrors := callback(ctx, buffer)
			for idx, err := range callbackErrors {
				if err == nil {
					successes += 1
					continue
				}
				failedLineNumbers = append(failedLineNumbers, buffer[idx].LineNumber())
			}
			buffer = buffer[:0]
		}
	}

	return ProcessingStats{successes, failedLineNumbers}, nil
}

type csvRecord struct {
	values     []string
	schemaMap  *map[string]uint8
	lineNumber uint64
}

func (rec *csvRecord) GetValue(key string) (string, error) {
	index, pres := (*rec.schemaMap)[key]
	if !pres {
		return "", errors.New(fmt.Sprintf("field '%s' not defined in schema", key))
	}
	if int(index) >= len(rec.values) {
		return "", errors.New(fmt.Sprintf("field '%s' not defined in line", key))
	}
	return rec.values[index], nil
}

func (rec *csvRecord) LineNumber() uint64 {
	return rec.lineNumber
}
