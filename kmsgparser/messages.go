package kmsgparser

import (
	"fmt"
	"time"
)

type Messages []*Message

func (msgs *Messages) String() string {

	msg := ""
	for _, m := range *msgs {
		msg += fmt.Sprintf("(%d) - %s: %s", m.SequenceNumber, m.Timestamp.Format(time.RFC3339Nano), m.Message)
	}

	return msg
}
