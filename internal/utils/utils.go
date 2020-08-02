package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func LogRequestBody(c *gin.Context) {
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody))) // Write body back
	fmt.Println("body: " + reqBody)
}

func LogRequestHeader(c *gin.Context) {
	fmt.Println("header: " + fmt.Sprint(c.Request.Header))
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

func Intersect(a, b []uint) []uint {
	m := make(map[uint]int)
	n := make([]uint, 0)
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

func Difference(a, b []uint) []uint {
	m := make(map[uint]int)
	n := make([]uint, 0)
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
