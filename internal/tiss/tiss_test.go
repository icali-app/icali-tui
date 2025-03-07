package tiss

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"testing"
)

func Test_RoomsParsedFromCsv_SerializeBackToOriginalContent(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)	
	writer.Comma = csvComma

	records := make([][]string, 0)
	for _, room := range rooms {
		record := make([]string, 0)
		record = append(record, room.name)
		record = append(record, strconv.Itoa(room.capacity))
		record = append(record, room.approvalInstitute)
		record = append(record, room.approvalMode)
		record = append(record, room.availability)
		record = append(record, room.building)
		record = append(record, room.address)
		record = append(record, room.roomNumber)
		record = append(record, room.tissLink)

		records = append(records, record)
	}

	if err := writer.WriteAll(records); err != nil {
		t.Fatalf("Couldn't write records: %v", err)
	}

	if buf.String() != roomsContent {
		t.Fatalf("Original content does not match parsed data.\nOriginal:\n\t%s\nParsed\n\t%s", roomsContent, buf.String())
	}
}
