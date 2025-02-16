// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build e2e

package integration

import (
	"context"
	"testing"

	"github.com/gotomicro/eorm"
	"github.com/gotomicro/eorm/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SelectTestSuite struct {
	Suite
	data *test.SimpleStruct
}

func (s *SelectTestSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.data = test.NewSimpleStruct(1)
	res := eorm.NewInserter[test.SimpleStruct](s.orm).Values(s.data).Exec(context.Background())
	if res.Err() != nil {
		s.T().Fatal(res.Err())
	}
}

func (s *SelectTestSuite) TearDownSuite() {
	res := eorm.RawQuery[any](s.orm, "DELETE FROM `simple_struct`").Exec(context.Background())
	if res.Err() != nil {
		s.T().Fatal(res.Err())
	}
}

func (s *SelectTestSuite) TestGet() {
	testCases := []struct {
		name    string
		s       *eorm.Selector[test.SimpleStruct]
		wantErr error
		wantRes *test.SimpleStruct
	}{
		{
			name: "not found",
			s: eorm.NewSelector[test.SimpleStruct](s.orm).
				From(&test.SimpleStruct{}).
				Where(eorm.C("Id").EQ(9)),
			wantErr: eorm.ErrNoRows,
		},
		{
			name: "found",
			s: eorm.NewSelector[test.SimpleStruct](s.orm).
				From(&test.SimpleStruct{}).
				Where(eorm.C("Id").EQ(1)),
			wantRes: s.data,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestMySQL8Select(t *testing.T) {
	suite.Run(t, &SelectTestSuite{
		Suite: Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
	suite.Run(t, &SelectTestSuiteGetMulti{
		Suite: Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

type SelectTestSuiteGetMulti struct {
	Suite
	data []*test.SimpleStruct
}

func (s *SelectTestSuiteGetMulti) SetupSuite() {
	s.Suite.SetupSuite()
	s.data = append(s.data, &test.SimpleStruct{Id: 1})
	s.data = append(s.data, &test.SimpleStruct{Id: 2})
	res := eorm.NewInserter[test.SimpleStruct](s.orm).Values(s.data...).Exec(context.Background())
	if res.Err() != nil {
		s.T().Fatal(res.Err())
	}
}

func (s *SelectTestSuiteGetMulti) TearDownSuite() {
	res := eorm.RawQuery[any](s.orm, "DELETE FROM `simple_struct`").Exec(context.Background())
	if res.Err() != nil {
		s.T().Fatal(res.Err())
	}
}

func (s *SelectTestSuiteGetMulti) TestGetMulti() {
	testCases := []struct {
		name    string
		s       *eorm.Selector[test.SimpleStruct]
		wantErr error
		wantRes []*test.SimpleStruct
	}{
		{
			name: "not found",
			s: eorm.NewSelector[test.SimpleStruct](s.orm).
				From(&test.SimpleStruct{}).
				Where(eorm.C("Id").EQ(9)),
			wantRes: []*test.SimpleStruct{},
		},
		{
			name: "found",
			s: eorm.NewSelector[test.SimpleStruct](s.orm).
				From(&test.SimpleStruct{}).
				Where(eorm.C("Id").LT(3)),
			wantRes: s.data,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := tc.s.GetMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
