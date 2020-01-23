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

// TODO comment
package model

// Entity represents a table in the database.
type Entity interface {
	TableName() string
	PrimaryKey() string
}

// ServiceInstance represents the ServiceInstance model in the Database.
type ServiceInstance struct {
	ID              string `gorm:"primary_key;type:varchar(100)"`
	ApplicationID   string `gorm:"type:varchar(100);not null;unique;column:application_id"`
	ApplicationName string `gorm:"type:varchar(100);not null"`
	SpaceID         string `gorm:"type:varchar(100);not null"`
	OrgID           string `gorm:"type:varchar(100);not null"`
	ConsumerKey     string `gorm:"type:varchar(100);not null"`
	ConsumerSecret  string `gorm:"type:varchar(100);not null"`
	ParameterHash   string `gorm:"type:varchar(100);not null"`
}

// Subscription represents the Subscription model in the database.
type Subscription struct {
	ID            string `gorm:"primary_key;type:varchar(100);not null;unique"`
	ApplicationID string `gorm:"type:varchar(100);not null"`
	APIName       string `gorm:"type:varchar(100);not null"`
	APIVersion    string `gorm:"type:varchar(100);not null"`
	User          string `gorm:"type:varchar(100);not null"`
	SVCInstanceID string `gorm:"type:varchar(100);not null;column:svc_instance_id"`
}

// Bind represents the Bind model in the Database
type Bind struct {
	ID            string `gorm:"primary_key;type:varchar(100)"`
	SVCInstanceID string `gorm:"type:varchar(100);not null;column:svc_instance_id"`
	PlatformAppID string `gorm:"type:varchar(100)"`
}

func (ServiceInstance) TableName() string {
	return TableServiceInstance
}

func (s ServiceInstance) PrimaryKey() string {
	return s.ID
}

func (Bind) TableName() string {
	return TableBind
}

func (b Bind) PrimaryKey() string {
	return b.ID
}

func (s Subscription) PrimaryKey() string {
	return s.ID
}

func (Subscription) TableName() string {
	return TableSubscriptions
}

const TableServiceInstance = "service_instances"

const TableBind = "binds"

const TableSubscriptions = "subscriptions"

const ServiceInstanceIDFieldName = "svc_instance_id"

const ForeignKeyDestAppID = TableServiceInstance + "(" + ServiceInstanceIDFieldName + ")"

const ForeignKeyDestSVCInstanceID = TableServiceInstance + "(id)"
