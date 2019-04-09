package main

import (
	"math/rand"
	"strconv"
	"time"
)

var mc map[string]string

func genCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(999999))
}
