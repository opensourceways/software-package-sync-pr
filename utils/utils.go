package utils

import "time"

func Retry(f func() error) (err error) {
	if err = f(); err == nil {
		return
	}

	m := 100 * time.Millisecond
	t := m

	for i := 1; i < 10; i++ {
		time.Sleep(t)
		t += m

		if err = f(); err == nil {
			return
		}
	}

	return
}
