/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package peer

import (
	"fmt"

	pb "github.com/idci/protos"
)

// DuplicateHandlerError returned if attempt to register same chaincodeID while a stream already exists.
type DuplicateHandlerError struct {
	To pb.PeerEndpoint
}

func (d *DuplicateHandlerError) Error() string {
	return fmt.Sprintf("Duplicate Handler error: %s", d.To)
}

func newDuplicateHandlerError(msgHandler MessageHandler) error {
	to, err := msgHandler.To()
	if err != nil {
		return fmt.Errorf("Error creating Duplicate Handler error: %s", err)
	}
	return &DuplicateHandlerError{To: to}
}
