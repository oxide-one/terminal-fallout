package sleeper

import "time"

func Sleep(duration int) {
	timeDuration := time.Millisecond * time.Duration(duration)
	time.Sleep(timeDuration)
}
