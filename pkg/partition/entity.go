// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package partition

import (
	"github.com/pkg/errors"

	"github.com/apache/skywalking-banyandb/api/common"
	modelv2 "github.com/apache/skywalking-banyandb/api/proto/banyandb/model/v2"
	"github.com/apache/skywalking-banyandb/banyand/tsdb"
	pbv2 "github.com/apache/skywalking-banyandb/pkg/pb/v2"
)

var (
	ErrMalformedElement = errors.New("element is malformed")
)

type EntityLocator []TagLocator

type TagLocator struct {
	FamilyOffset int
	TagOffset    int
}

func (e EntityLocator) Find(value []*modelv2.TagFamilyForWrite) (tsdb.Entity, error) {
	entity := make(tsdb.Entity, len(e))
	for i, index := range e {
		tag, err := GetTagByOffset(value, index.FamilyOffset, index.TagOffset)
		if err != nil {
			return nil, err
		}
		entry, errMarshal := pbv2.MarshalIndexFieldValue(tag)
		if errMarshal != nil {
			return nil, errMarshal
		}
		entity[i] = entry
	}
	return entity, nil
}

func (e EntityLocator) Locate(value []*modelv2.TagFamilyForWrite, shardNum uint32) (tsdb.Entity, common.ShardID, error) {
	entity, err := e.Find(value)
	if err != nil {
		return nil, 0, err
	}
	id, err := ShardID(entity.Marshal(), shardNum)
	if err != nil {
		return nil, 0, err
	}
	return entity, common.ShardID(id), nil
}

func GetTagByOffset(value []*modelv2.TagFamilyForWrite, fIndex, tIndex int) (*modelv2.TagValue, error) {
	if fIndex >= len(value) {
		return nil, errors.Wrap(ErrMalformedElement, "tag family offset is invalid")
	}
	family := value[fIndex]
	if tIndex >= len(family.GetTags()) {
		return nil, errors.Wrap(ErrMalformedElement, "tag offset is invalid")
	}
	return family.GetTags()[tIndex], nil
}
