package CustomUtils

import (
	"fmt"
	"time"
)

var history map[string][]TimeTag = map[string][]TimeTag{}

func LogTime(tag, message string) {
	if _, ok := history[tag]; !ok {
		history[tag] = make([]TimeTag, 10)
	}
	timestamp := time.Now().UnixNano()
	fmt.Printf("%v %v %v\n", message, timestamp, tag)

	t := TimeTag{
		tag:  tag,
		time: timestamp,
	}

	history[tag] = append(history[tag], t)
}

type TimeTag struct {
	tag  string
	time int64
}
