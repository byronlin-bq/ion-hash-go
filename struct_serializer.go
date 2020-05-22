/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package ionhash

type structSerializer struct {
	baseSerializer

	scalarSerializer serializer
	fieldHashes      [][]byte
}

func newStructSerializer(hashFunction IonHasher, depth int, hashFunctionProvider IonHasherProvider) serializer {
	return &structSerializer{
		baseSerializer:   baseSerializer{hashFunction: hashFunction, depth: depth},
		scalarSerializer: newScalarSerializer(hashFunctionProvider.newHasher(), depth+1)}
}

func (structSerializer *structSerializer) scalar(ionValue interface{}) error {
	handleFieldNameErr := structSerializer.scalarSerializer.handleFieldName(ionValue)
	if handleFieldNameErr != nil {
		return handleFieldNameErr
	}

	scalarErr := structSerializer.scalarSerializer.scalar(ionValue)
	if scalarErr != nil {
		return scalarErr
	}

	sum := structSerializer.sum()
	structSerializer.appendFieldHash(sum)
	return nil
}

func (structSerializer *structSerializer) stepOut() {
	//compareBytes(structSerializer.fieldHashes)

	for i := 0; i < len(structSerializer.fieldHashes); i++ {
		structSerializer.update(escape(structSerializer.fieldHashes[i]))
	}
	structSerializer.scalarSerializer.stepOut()
}

func (structSerializer *structSerializer) appendFieldHash(sum []byte) {
	structSerializer.fieldHashes = append(structSerializer.fieldHashes, sum)
}

func compareBytes(bs1, bs2 []byte) []int16 {
	panic("implement me")
}
