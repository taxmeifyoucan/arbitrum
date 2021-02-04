/*
 * Copyright 2021, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package evm

import (
	"github.com/offchainlabs/arbitrum/packages/arb-util/inbox"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"
	"github.com/pkg/errors"
	"math/big"
)

type SendResult struct {
	BatchNumber *big.Int
	BatchIndex  *big.Int
	Data        []byte
}

func NewSendResultFromValue(tup *value.TupleValue) (*SendResult, error) {
	if tup.Len() != 5 {
		return nil, errors.Errorf("send result expected tuple of length 5, but recieved len %v", tup.Len())
	}

	resultKindVal, _ := tup.GetByInt64(0)
	batchNumberVal, _ := tup.GetByInt64(1)
	batchIndexVal, _ := tup.GetByInt64(2)
	dataSizeVal, _ := tup.GetByInt64(3)
	dataContentsVal, _ := tup.GetByInt64(4)

	resultKindInt, ok := resultKindVal.(value.IntValue)
	if !ok {
		return nil, errors.New("resultKind must be an int")
	}
	if resultKindInt.BigInt().Uint64() != 2 {
		return nil, errors.New("incorrect result kind for send")
	}
	batchNumberInt, ok := batchNumberVal.(value.IntValue)
	if !ok {
		return nil, errors.New("batchNumber must be an int")
	}
	batchIndexInt, ok := batchIndexVal.(value.IntValue)
	if !ok {
		return nil, errors.New("batchIndex must be an int")
	}
	dataSizeInt, ok := dataSizeVal.(value.IntValue)
	if !ok {
		return nil, errors.New("dataSize must be an int")
	}
	dataContentsBuf, ok := dataContentsVal.(*value.Buffer)
	if !ok {
		return nil, errors.New("dataContents must be an int")
	}

	data, err := inbox.BufAndLengthToBytes(dataSizeInt.BigInt(), dataContentsBuf)
	if err != nil {
		return nil, err
	}
	return &SendResult{
		BatchNumber: batchNumberInt.BigInt(),
		BatchIndex:  batchIndexInt.BigInt(),
		Data:        data,
	}, nil
}