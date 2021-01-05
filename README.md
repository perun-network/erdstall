# Erdstall

A 2nd Layer Plasma based on abstract trusted third parties. Build during the
ETHOnline 2020 hackathon.

## Repository Structure

The repository's root directory is a Go module. `contracts` is a truffle
project, containing the Solidity contracts. `operator` contains the operator
code and `client` the end user logic. All main programs reside inside the `cmd`
folder. All TEE related code can be found in `tee`.

## Getting started

```bash
# build
go build -o operator.bin ./cmd/operator
go build -o client.bin ./cmd/client

# gnache
ganache-cli -e 100000000000000 -b 5 -m "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic" -s 100

# op
./operator.bin
# alice
./client.bin --contract 0x4fb8637afd28492a3209017556e95dc2f8086ddb
--account-index 2
# bob
./client.bin --contract 0x4fb8637afd28492a3209017556e95dc2f8086ddb
--account-index 3
```

### Configuration file

With the command line flag `-config <file>`, you can specify a JSON
configuration file. In this file, you can override the default settings. For
reference, see `demo/config.json`. If you wish to enable TLS support for the
websocket connection, set the fields `KeyFile` to the TLS private key file path,
and `CertFile` to the TLS certificate file path.

## Description

Erdstall leverages Trusted Execution Environments (TEE) like Intel SGX (or even
MPC committees) to scale Ethereum. Similar to Plasma or Rollups, the system
consists of a smart contract, an untrusted operator running a TEE and a dynamic
group of users. Joining and leaving the system is by a single call to a smart
contract. But once assets are deposited into the system, off-chain transactions
are free and only require the exchange of signatures from users to the operator.
The TEE Enclave receives and verifies those transactions and keeps track of the
system state. The whole system evolves in epochs and at the end of each epoch,
so-called balance proofs are distributed to all users, allowing them to leave
the system at any time. Those proofs are also necessary to give all users the
possibility to exit the system shall the operator decide to cease operating.

The underlying protocols were developed and proven secure by the Chair of
Applied Cryptography research group at Technical University Darmstadt (the same
team behind the Perun generalized state channels). The related paper is
currently in submission at a cryptography conference.

## License

This work is released under the Apache 2.0 license. See LICENSE file for more
details.

_Copyright (C) 2020 - The Erdstall Authors._
