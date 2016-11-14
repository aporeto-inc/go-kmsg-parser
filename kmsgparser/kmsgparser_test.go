package kmsgparser

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockBootTime(bootTime time.Time) {
	sysInfoFunc = func(output *syscall.Sysinfo_t) error {
		output.Uptime = 0
		return nil
	}
	// It was the dawn of time
	timeNowFunc = func() time.Time {
		return bootTime
	}
}

func TestParseMessage(t *testing.T) {
	bootTime = time.Unix(0xb100, 0x5ea1).Round(time.Microsecond)
	msg, err := parseMessage("6,2565,102258085667,-;docker0: port 2(vethc1bb733) entered blocking state")
	if err != nil {
		t.Fatalf("error parsing: %v", err)
	}

	assert.Equal(t, msg.Message, "docker0: port 2(vethc1bb733) entered blocking state")

	assert.Equal(t, msg.Priority, 6)
	assert.Equal(t, msg.SequenceNumber, 2565)
	assert.Equal(t, msg.Timestamp, bootTime.Add(102258085667*time.Microsecond))
}

func TestParse(t *testing.T) {
	bootTime = time.Unix(0xb100, 0x5ea1).Round(time.Microsecond)
	mockBootTime(bootTime)
	f, err := os.Open(filepath.Join("test_data", "sample1.kmsg"))
	if err != nil {
		t.Fatalf("could not find sample data: %v", err)
	}
	defer f.Close()

	expectedMessages := []Message{
		{
			Priority:       6,
			SequenceNumber: 1804,
			Timestamp:      bootTime.Add(47700428483 * time.Microsecond),
			Message:        "wlp4s0: associated",
		},
		{
			Priority:       6,
			SequenceNumber: 1805,
			Timestamp:      bootTime.Add(51742248189 * time.Microsecond),
			Message:        "thinkpad_acpi: EC reports that Thermal Table has changed",
		},
		{
			Priority:       6,
			SequenceNumber: 2651,
			Timestamp:      bootTime.Add(106819644585 * time.Microsecond),
			Message:        "CPU1: Package temperature/speed normal",
		},
	}

	s := bufio.NewScanner(f)
	mockKmsg, mockKmsgInput := io.Pipe()
	go func() {
		for s.Scan() {
			_, err := mockKmsgInput.Write(s.Bytes())
			if err != nil {
				panic(err)
			}
		}
		mockKmsgInput.Close()
	}()

	lines, err := Parse(mockKmsg)
	if err != nil {
		t.Fatalf("err parsing: %v", err)
	}

	messages := []Message{}
	for line := range lines {
		messages = append(messages, line)
	}

	assert.Equal(t, expectedMessages, messages)
}
