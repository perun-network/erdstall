// SPDX-License-Identifier: Apache-2.0

package prototype

import "github.com/perun-network/erdstall/tee"

func (e *Enclave) blockProcessor(
	newBlocks <-chan *tee.Block,
	verifiedBlocks chan<- *tee.Block,
) error {
	// read blocks from newBlocks
	// verify block chain
	// write verified blocks (up to PoW-depth) to verifiedBlocks
}
