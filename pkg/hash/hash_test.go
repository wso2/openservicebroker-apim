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

package hash

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

const ErrMsgTestIncorrectResult = "expected value: %v but then returned value: %v"

var m = &Manager{
	Hash: sha256.New(),
}

func TestAddUint32(t *testing.T) {
	u := uint32(100)
	m.ResetHash()
	err := m.AddUint32(u)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "40e736c02a102a050e1555781b4171020a4279adaa7ed9ca3cc9633a0ade9c37")
}

func TestAddUint64(t *testing.T) {
	u := uint64(100)
	m.ResetHash()
	err := m.AddUint64(u)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "26ab39150b6330152576e4c7fa7e0caa804b5e9db0476a3e48e6b53f1cda8279")
}

func TestAddArray(t *testing.T) {
	a1 := []string{"a", "b"}
	m.ResetHash()
	err := m.AddArray(a1)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "fb8e20fc2e4c3f248c60c39bd652f3c1347298bb977b8b4d5903b85055620603")

	a2 := []string{"b", "a"}
	m.ResetHash()
	err = m.AddArray(a2)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "fb8e20fc2e4c3f248c60c39bd652f3c1347298bb977b8b4d5903b85055620603")

	m.ResetHash()
	err = m.AddArray(nil)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}

func TestAddMap(t *testing.T) {
	m1 := map[string]string{"a": "a", "b": "b"}
	m.ResetHash()
	err := m.AddMap(m1)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "a667282675f4876021d392aa6592f39dabf718748c4b738563cb9d5dc8f21f24")

	m2 := map[string]string{"b": "b", "a": "a"}
	m.ResetHash()
	err = m.AddMap(m2)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "a667282675f4876021d392aa6592f39dabf718748c4b738563cb9d5dc8f21f24")

	m.ResetHash()
	err = m.AddMap(nil)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}

func TestAddBool(t *testing.T) {
	m.ResetHash()
	err := m.AddBool(true)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "4bf5122f344554c53bde2ebb8cd2b7e3d1600ad631c385a5d7cce23c7785459a")

	m.ResetHash()
	err = m.AddBool(false)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d")
}

func TestAddString(t *testing.T) {
	str := "Hash"
	m.ResetHash()
	err := m.AddString(str)
	if err != nil {
		t.Error(err)
	}
	checkHashResult(t, m, "a91069147f9bd9245cdacaef8ead4c3578ed44f179d7eb6bd4690e62ba4658f2")
}

func checkHashResult(t *testing.T, m *Manager, expected string) {
	result, err := m.Generate()
	if err != nil {
		t.Error(err)
	}
	if expected != result {
		t.Error(fmt.Sprintf(ErrMsgTestIncorrectResult, expected, result))
	}
}

func TestNilHash(t *testing.T) {
	mNil := &Manager{}
	err := mNil.AddBool(true)
	if err != errHashNil {
		t.Error(fmt.Sprintf("expecting the error: %v", errHashNil))
	}
}
