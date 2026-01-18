package domain

import "time"

type Alert struct {
}

type Warnings struct {
	Timestamp time.Time
	Warnings  []string
	Faults    []string
}
