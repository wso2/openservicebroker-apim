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

// Package Hash can be used to construct Hash values for data types.
package hash

// TODO use previous lib
import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash"
	"sort"
)

var errHashNil = errors.New("hash cannot be nil")

// Manager struct represents the Hash Manager.
type Manager struct {
	Hash hash.Hash
}

// AddMap method adds the given Map content for hashing irrespective of the element order.
// Returns an error if the Hash is not initialized.
func (m *Manager) AddMap(mVal map[string]string) error {
	if m.Hash == nil {
		return errHashNil
	}
	if m != nil {
		keys := make([]string, len(mVal))
		var i = 0
		for key, _ := range mVal {
			keys[i] = key
			i++
		}
		sort.Strings(keys)

		values := make([]string, len(mVal))
		i = 0
		for _, val := range keys {
			values[i] = mVal[val]
			i++
		}
		m.addSortedArray(keys)
		m.addSortedArray(values)
		return nil
	}
	return nil
}

// addSortedArray method adds the given sorted array content for hashing. Must do a nil check for m.Hash before calling.
func (m *Manager) addSortedArray(sortedArr []string) {
	if sortedArr != nil {
		for _, val := range sortedArr {
			m.Hash.Write([]byte(val))
		}

	}
}

// AddArray method adds the given unsorted array content for hashing irrespective of the element order.
// Returns an error if the Hash is not initialized.
func (m *Manager) AddArray(arr []string) error {
	if m.Hash == nil {
		return errHashNil
	}
	arrCopy := make([]string, len(arr))
	copy(arrCopy, arr)
	if arr != nil {
		sort.Strings(arrCopy)
		m.addSortedArray(arrCopy)
	}
	return nil
}

// AddString method adds the given string for hashing.
// Returns an error if encountered.
func (m *Manager) AddString(str string) error {
	if m.Hash == nil {
		return errHashNil
	}
	_, err := m.Hash.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}

// AddUint32 method adds the given uint32 for hashing.
// Returns an error if encountered.
func (m *Manager) AddUint32(u uint32) error {
	if m.Hash == nil {
		return errHashNil
	}
	return binary.Write(m.Hash, binary.LittleEndian, u)
}

// AddInt32 method adds the given int32 for hashing.
// Returns an error if encountered.
func (m *Manager) AddInt32(i int32) error {
	if m.Hash == nil {
		return errHashNil
	}
	return binary.Write(m.Hash, binary.LittleEndian, i)
}

// AddUint64 method adds the given uint64 for hashing.
// Returns an error if encountered.
func (m *Manager) AddUint64(u uint64) error {
	if m.Hash == nil {
		return errHashNil
	}
	return binary.Write(m.Hash, binary.LittleEndian, u)
}

// AddInt64 method adds the given int64 for hashing.
// Returns an error if encountered.
func (m *Manager) AddInt64(i int64) error {
	if m.Hash == nil {
		return errHashNil
	}
	return binary.Write(m.Hash, binary.LittleEndian, i)
}

// AddBool method adds the given bool for hashing.
// Returns an error if encountered.
func (m *Manager) AddBool(b bool) error {
	if m.Hash == nil {
		return errHashNil
	}
	return binary.Write(m.Hash, binary.LittleEndian, b)
}

// Generate method returns a hash value for added content.
func (m *Manager) Generate() (string, error) {
	if m.Hash == nil {
		return "", errHashNil
	}
	return hex.EncodeToString(m.Hash.Sum(nil)), nil
}

// ResetHash resets the underneath Hash Manager.
// This method should be called before generating Hash if the Manager is used for more than once.
func (m *Manager) ResetHash() {
	m.Hash.Reset()
}
