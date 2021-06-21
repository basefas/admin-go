package utils

import (
	"admin-go/internal/utils/log"
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func LogRequestBody(c *gin.Context) {
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody)))
	log.Debug("body: " + reqBody)
}

func LogRequestHeader(c *gin.Context) {
	log.Debug("header: " + fmt.Sprint(c.Request.Header))
}

func LogRequest(c *gin.Context) {
	LogRequestHeader(c)
	LogRequestBody(c)
}

func GetRequestBody(c *gin.Context) string {
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody)))
	return reqBody
}

func Intersect(a, b []uint64) []uint64 {
	m := make(map[uint64]int64)
	n := make([]uint64, 0)
	for _, v := range a {
		m[v]++
	}

	for _, v := range b {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}
	return n
}

func Difference(a, b []uint64) []uint64 {
	m := make(map[uint64]int)
	n := make([]uint64, 0)
	inter := Intersect(a, b)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range a {
		times, _ := m[value]
		if times == 0 {
			n = append(n, value)
		}
	}
	return n
}
