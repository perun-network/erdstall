## Terminal 0.1: Ganache
```bash
$ ganache-cli -b 3 -a 5 -e 101 -m "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic"
```

## Terminal 0.2: Accountant
```bash
$ go run ./cmd/accountant --accounts demo/accounts.json
```

## Terminal 1: Operator
```bash
$ go run ./cmd/operator --config demo/config.json
```

## Terminal 2: Client A
```bash
$ go run ./cmd/client --contract 0x079557d7549d7D44F4b00b51d2C532674129ed51 --mnemonic "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic" --account-index 1 --username "👩 Alice"
```

## Terminal 3: Client B
```bash
$ go run ./cmd/client --contract 0x079557d7549d7D44F4b00b51d2C532674129ed51 --mnemonic "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic" --account-index 2 --username "👨 Bob"
```