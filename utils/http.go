package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"os"
	"time"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseData(status string, message string, data any) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func Logger() gin.HandlerFunc {
	// Create a new log file
	logFile, err := os.OpenFile(
		"./app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		panic(err)
	}

	// Set up log rotation every 5 days
	rotationTime := time.Now().AddDate(0, 0, 5)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./app.log",
		MaxSize:    10,   // Max size of each log file in megabytes
		MaxBackups: 30,   // Max number of log files to keep
		MaxAge:     5,    // Max number of days to keep old log files
		Compress:   true, // Compress the rotated log files
	}

	return func(c *gin.Context) {
		// Log the request time, method, and URL
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		logMessage := fmt.Sprintf(
			"[%s] %s %s %s %d %s\n",
			endTime.Format("2006/01/02 15:04:05"),
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.String(),
			c.Writer.Status(),
			endTime.Sub(startTime).String(),
		)

		// Write the log message to the file and to stdout
		_, err := logFile.WriteString(logMessage)
		if err != nil {
			panic(err)
		}

		// Rotate the log file if necessary
		if time.Now().After(rotationTime) {
			logFile.Close()
			lumberjackLogger.Rotate()
			logFile, err = os.OpenFile(
				"./app.log",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0666,
			)
			if err != nil {
				panic(err)
			}
			rotationTime = time.Now().AddDate(0, 0, 5)
		}
	}
}
