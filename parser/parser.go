package parser

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/akraievoy/tsv_load/service_utils"
	"io"
)

type Record interface {
	GetValue(key string) (string, error)
	LineNumber() uint64
}

type ProcessingStats struct {
	Successes         uint64
	FailedLineNumbers []uint64
}

type RecordStream interface {
	ForEachBatch(
		ctx context.Context,
		batchSize uint32,
		brokenLinesToFail uint32,
		callback func(context.Context, []Record) ([]uint64, error),
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
	brokenLinesToFail uint32,
	callback func(context.Context, []Record) ([]uint64, error),
) (ProcessingStats, error) {
	buffer := make([]Record, 0, batchSize)
	successes := uint64(0)
	brokenLineNumbers := make([]uint64, 0)
	thereIsMore := true
	for thereIsMore {
		values, err := rs.reader.Read()

		if service_utils.PeekDone(ctx) {
			thereIsMore = false
			return ProcessingStats{successes, brokenLineNumbers}, ctx.Err()
		}

		if err == io.EOF {
			thereIsMore = false
		} else if err != nil {
			//	LATER recover from some CSV-specific errors as well (should we send last pre-bail buffer or not)?
			return ProcessingStats{successes, brokenLineNumbers}, err
		}

		rs.lineNumber += 1
		if len(values) > 0 { // LATER test for this, KEK
			buffer = append(buffer, &csvRecord{values, &rs.schemaMap, rs.lineNumber})
		}

		if len(buffer) == int(batchSize) || !thereIsMore {
			batchBrokenLineNumbers, batchErr := callback(ctx, buffer)
			if batchErr == nil {
				successes += uint64(len(buffer) - len(batchBrokenLineNumbers))
			}
			brokenLineNumbers = append(brokenLineNumbers, batchBrokenLineNumbers...)
			buffer = buffer[:0]
		}

		if uint32(len(brokenLineNumbers)) > brokenLinesToFail {
			message := fmt.Sprintf(
				"seen %d broken lines, which is greated than limit of %d",
				len(brokenLineNumbers), brokenLinesToFail,
			)
			return ProcessingStats{successes, brokenLineNumbers}, errors.New(message)
		}
	}

	return ProcessingStats{successes, brokenLineNumbers}, nil
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
