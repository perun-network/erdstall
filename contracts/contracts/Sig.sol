// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.7.0;

import "./ECDSA.sol";

// Sig is a library to verify signatures.
library Sig {
    // Verify verifies whether a piece of data was signed correctly.
    function verify(bytes memory data, bytes memory signature, address signer) internal pure returns (bool) {
        bytes32 prefixedHash = ECDSA.toEthSignedMessageHash(keccak256(data));
        address recoveredAddr = ECDSA.recover(prefixedHash, signature);
        return recoveredAddr == signer;
    }
}
