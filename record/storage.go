package record

import (
	"encoding/xml"
	"time"
)

// Storage represents parsed storage record.
type Storage struct {
	RecordID                  string    `xml:"RECORD_ID"`
	CreateTime                time.Time `xml:"CREATE_TIME"`
	StorageSystem             string    `xml:"STORAGE_SYSTEM"`
	Site                      *string   `xml:"SITE"`
	StorageShare              *string   `xml:"STORAGE_SHARE"`
	StorageMedia              *string   `xml:"STORAGE_MEDIA"`
	StorageClass              *string   `xml:"STORAGE_CLASS"`
	FileCount                 *string   `xml:"FILE_COUNT"`
	DirectoryPath             *string   `xml:"DIRECTORY_PATH"`
	LocalUser                 *string   `xml:"LOCAL_USER"`
	LocalGroup                *string   `xml:"LOCAL_GROUP"`
	UserIdentity              *string   `xml:"USER_IDENTITY"`
	Group                     *string   `xml:"GROUP"`
	GroupAttribute            *string   `xml:"GROUP_ATTRIBUTE"`
	GroupAttributeType        *string   `xml:"GROUP_ATTRIBUTE_TYPE"`
	StartTime                 time.Time `xml:"START_TIME"`
	EndTime                   time.Time `xml:"END_TIME"`
	ResourceCapacityUsed      uint64    `xml:"RESOURCE_CAPACITY_USED"`
	LogicalCapacityUsed       *uint64   `xml:"LOGICAL_CAPACITY_USED"`
	ResourceCapacityAllocated *uint64   `xml:"RESOURCE_CAPACITY_ALLOCATED"`
}

// Storages represents storages structure parsed from XML where storage records are wrapped.
type Storages struct {
	XMLName  xml.Name  `xml:"STORAGES"`
	Storages []Storage `xml:"STORAGE"`
}
