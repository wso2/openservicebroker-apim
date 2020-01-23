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
package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	ErrMsgTestCouldNotSetEnv  = "couldn't set the ENV: %v"
	ErrMsgTestIncorrectResult = "expected value: %v but then returned value: %v"
)

func TestValidateParam(t *testing.T) {
	valid := IsValidParams()
	if valid {
		t.Errorf(ErrMsgTestIncorrectResult, !valid, valid)
	}
	valid = IsValidParams("a", "b", "c")
	if !valid {
		t.Errorf(ErrMsgTestIncorrectResult, !valid, valid)
	}
	valid = IsValidParams("a", "b", "")
	if valid {
		t.Errorf(ErrMsgTestIncorrectResult, !valid, valid)
	}
}

func TestRawMSGToString(t *testing.T) {
	msg := `{"foo":"bar"}`
	raw := json.RawMessage(`{"foo":"bar"}`)
	result, err := RawMsgToString(&raw)
	if err != nil {
		t.Error(err)
	}
	if result != msg {
		t.Errorf(ErrMsgTestIncorrectResult, msg, result)
	}
}

func TestConstructURL(t *testing.T) {
	result, err := ConstructURL("https://localhost:9443", "carbon")
	if err != nil {
		t.Error(err)
	}
	expected := "https://localhost:9443/carbon"
	if result != expected {
		t.Errorf(ErrMsgTestIncorrectResult, expected, result)
	}

	result, err = ConstructURL("https://localhost:9443", "carbon", "publisher")
	if err != nil {
		t.Error(err)
	}
	expected = "https://localhost:9443/carbon/publisher"
	if result != expected {
		t.Errorf(ErrMsgTestIncorrectResult, expected, result)
	}

	result, err = ConstructURL("https://localhost:9443")
	if err != nil {
		t.Error(err)
	}
	expected = "https://localhost:9443"
	if result != expected {
		t.Errorf(ErrMsgTestIncorrectResult, expected, result)
	}

	result, err = ConstructURL()
	if err == nil {
		t.Error(fmt.Sprintf("Expecting an error: %v", ErrNoPaths))
	}
	if err != ErrNoPaths {
		t.Error(fmt.Sprintf("Expecting the error %v but got ", ErrNoPaths) + err.Error())
	}
}
