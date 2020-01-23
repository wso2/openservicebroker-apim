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

// Package utils holds a common set of Util functions
package utils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
	"path"
)

var (
	ErrCannotParseURL = errors.New("unable to parse URL")
	ErrNoPaths        = errors.New("no paths found")
)

// IsValidParams returns false if one of the arguments are empty or argument array is nil
func IsValidParams(params ...string) bool {
	if params == nil {
		return false
	}
	for _, val := range params {
		if val == "" {
			return false
		}
	}
	return true
}

// RawMsgToString converts json.RawMessage into String
// Returns the string representation of json.RawMessage and any error occurred
func RawMsgToString(msg *json.RawMessage) (string, error) {
	j, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// ConstructURL construct URL by joining the paths provided.
// First param will be treated as the base and the rest of params configured as paths.
// An error will be thrown if the number of paths is equal to zero.
// Returns constructed path and any error occurred.
func ConstructURL(paths ...string) (string, error) {
	if len(paths) == 0 {
		return "", ErrNoPaths
	}
	u, err := url.Parse(paths[0])
	if err != nil {
		return "", ErrCannotParseURL
	}
	for i := 1; i < len(paths); i++ {
		u.Path = path.Join(u.Path, paths[i])
	}
	return u.String(), nil
}

// GetJSONSchema returns the schema struct representation of the given JSON string and any error encountered.
func GetJSONSchema(rawSchema string) (map[string]interface{}, error) {
	var schema map[string]interface{}
	err := json.Unmarshal([]byte(rawSchema), &schema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}