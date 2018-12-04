package controller

import (
	"log"
	"testing"
)

func TestSendMobile(t *testing.T) {
	SendMobile("18612660953", "123456")
}

func TestGetRandomSalt(t *testing.T) {
	log.Println(GetRandomSalt())
}
