/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"net"
	"time"

	"github.com/sony/sonyflake"
	"github.com/vogo/vogo/vnet"
)

var (
	snowflakeStartTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	machineIDFetcher   = localIPMachineIDFetcher
)

// MachineIDFetcher function for fetching machine id
// NOTE: for a distributed system, it's better to register machine into a general center, which make sure the machine id is unique.
type MachineIDFetcher func() (uint16, error)

// SetMachineIDFetcher set MachineIDFetcher
func SetMachineIDFetcher(fetcher MachineIDFetcher) {
	machineIDFetcher = fetcher
}

var machineIDFromIP uint16

// localIPMachineIDFetcher get machine id from local ip address.
func localIPMachineIDFetcher() (uint16, error) {
	if machineIDFromIP > 0 {
		return machineIDFromIP, nil
	}

	ipString, err := vnet.LocalIP()
	if err != nil {
		return 0, err
	}

	ip := net.ParseIP(ipString)

	machineIDFromIP = binary.BigEndian.Uint16(ip[len(ip)-2:])

	return machineIDFromIP, nil
}

func newSnowflake() *sonyflake.Sonyflake {
	return sonyflake.NewSonyflake(sonyflake.Settings{
		// StartTime is the time since which the Sonyflake time is defined as the elapsed time.
		// If StartTime is 0, the start time of the Sonyflake is set to "2014-09-01 00:00:00 +0000 UTC".
		// If StartTime is ahead of the current time, Sonyflake is not created.
		StartTime: snowflakeStartTime,

		// MachineID returns the unique ID of the Sonyflake instance.
		// If MachineID returns an error, Sonyflake is not created.
		// If MachineID is nil, default MachineID is used.
		// Default MachineID returns the lower 16 bits of the private IP address.
		MachineID: machineIDFetcher,

		// CheckMachineID validates the uniqueness of the machine ID.
		// If CheckMachineID returns false, Sonyflake is not created.
		// If CheckMachineID is nil, no validation is done.
		CheckMachineID: nil,
	})
}

// Snowflake snow flake id interface
type Snowflake interface {
	NextID() uint64
	NextSID() string
}

type sonySnowflake struct {
	flake *sonyflake.Sonyflake
}

// NextID can continue to generate IDs for about 174 years from StartTime.
// But after the Sonyflake time is over the limit, NextID panics.
func (s *sonySnowflake) NextID() uint64 {
	id, err := s.flake.NextID()
	if err != nil {
		panic(err)
	}

	return id
}

// NextSID generate string format id
func (s *sonySnowflake) NextSID() string {
	id, err := s.flake.NextID()
	if err != nil {
		panic(err)
	}

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return base64.RawURLEncoding.EncodeToString(b)
}

// New Snowflake
func New() Snowflake {
	return &sonySnowflake{flake: newSnowflake()}
}
