package main

import "time"

func db() {
	defer wg.Done()
	for {
		time.Sleep(5 * time.Second)
	}
}
