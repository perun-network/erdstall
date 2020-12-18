// SPDX-License-Identifier: Apache-2.0

package eth

import (
	"context"
	"time"
)

var (
	contextNodeReqTimeout   = 2 * time.Second
	contextWaitMinedTimeout = 10 * time.Second
)

// SetNodeReqTimeout sets the timeout used when creating contexts with
// ContextNodeReq.
func SetNodeReqTimeout(duration time.Duration) {
	contextNodeReqTimeout = duration
}

// SetWaitMinedTimeout sets the timeout used when creating contexts with
// ContextWaitMined.
func SetWaitMinedTimeout(duration time.Duration) {
	contextWaitMinedTimeout = duration
}

// ContextNodeReq creates a node request timeout context (default: 2 sec).
// Use this context when performing a node request that doesn't wait for any
// blocks to be mined, e.g., setting up a subscription, sending a transaction
// (but not waiting for it to be mined).
//
// Can be changed using SetNodeReqTimeout.
func ContextNodeReq() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextNodeReqTimeout)
}

// ContextWaitMined creates a mining timeout context (default: 10 sec).
// Use this context when waiting for a transaction to be mined or similar.
//
// Can be changed using SetWaitMinedTimeout.
func ContextWaitMined() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextWaitMinedTimeout)
}
