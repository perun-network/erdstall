// SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GenericJSONMarshallingTest tests that an object can be marshalled and
// unmarshalled correctly with the golang json marshaller.
// `orig` should be an instance of the type that you want to test, filled
// with values. `empty` should be a default contructed object, e.g. &T{}.
// `empty` must be passed as pointer.
func GenericJSONMarshallingTest(t *testing.T, orig, empty interface{}) {
	data, err := json.Marshal(orig)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, empty))
	assert.Equal(t, orig, reflect.ValueOf(empty).Elem().Interface())
}
