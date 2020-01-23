/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package log handles logging.
package log

import (
	"code.cloudfoundry.org/lager"
	"github.com/pkg/errors"
	"io"
	"os"
)

const (
	// LoggerName is used to specify the source of the logger.
	LoggerName = "wso2-apim-broker"

	// FilePerm is the permission for the server log file.
	FilePerm = 0644

	ErrMsgUnableToOpenLogFile = "unable to open the Log file: %s"
)

var logger = lager.NewLogger(LoggerName)
var ioWriter io.Writer = os.Stdout

// Data represents the information in the logs.
type Data struct {
	lData lager.Data
}

// Configure initializes lager logging object,
// 1. Setup log level
// 2. Setup log file
// Returns configured logger and any error encountered.
func Configure(logFile, logLevelS string) (lager.Logger, error) {
	logL, err := lager.LogLevelFromString(logLevelS)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, FilePerm)
	if err != nil {
		return nil, errors.Wrapf(err, ErrMsgUnableToOpenLogFile, logFile)
	}
	ioWriter = io.MultiWriter(os.Stdout, f)
	logger.RegisterSink(lager.NewWriterSink(ioWriter, logL))
	return logger, nil
}

// IoWriterLog returns the IO writer object for logging. By default it is pointed to STDOUT.
func IoWriterLog() io.Writer {
	return ioWriter
}

// Info logs Info level messages using configured lager.Logger.
func Info(msg string, data *Data) {
	if data == nil {
		data = NewData()
	}
	logger.Info(msg, data.lData)
}

// Error logs Error level messages using configured lager.Logger.
func Error(msg string, err error, data *Data) {
	if data == nil {
		data = NewData()
	}
	logger.Error(msg, err, data.lData)
}

// Debug logs Debug level messages using configured lager.Logger.
func Debug(msg string, data *Data) {
	if data == nil {
		data = NewData()
	}
	logger.Debug(msg, data.lData)
}

// HandleErrorAndExit prints an error to STDOUT and invoke a panic.
func HandleErrorAndExit(errMsg string, err error) {
	logger.Fatal(errMsg, err, lager.Data{})
}

// NewData returns a pointer a Data struct.
func NewData() *Data {
	return &Data{}
}

// Add adds data to current data obj.
// Returns a reference to the current Log data obj.
func (l *Data) Add(key string, val interface{}) *Data {
	if l.lData == nil {
		l.lData = lager.Data{}
	}
	if key != "" && val != nil {
		l.lData[key] = val
	}
	return l
}
