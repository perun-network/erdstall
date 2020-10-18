// SPDX-License-Identifier: Apache-2.0

package eth

import (
	"context"
	"time"
)

const defaultContextTimeout = 10 * time.Second

// NewDefaultContext creates a default timeout context (10 sec).
func NewDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultContextTimeout)
}
