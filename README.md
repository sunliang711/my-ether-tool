subcommands

account
    add
    rm
    list
    switch

network
    add
    rm
    list
    switch

global flag for the following:
--account <>
--network <>
tx
    send
    query (query by hash)
    receipt (query receipt by hash)
    offsign

contract
    read
    write

erc20
    info --name --symbol --decimals
    transfer
    transferFrom
    approve
    allowance

codec
    abiEncode [--function <transfer(address,uint256)>] [--args <> -args <> ...]
    abiEncode [--function <transfer((address,address),(uint256,uint256))>] [--types <> --types <> ...] [--args <> -args <> ...]
    abiDecode

eip712
    sign


## examples
### abi encode
 go run main.go codec abiencode --abi "transfer(address,uint256)" --args 0x9D757Dd679bE17b4094c740fB0047fa3a7Ed6DF0 --args 1000000

 ### offsign
 go run main.go tx offsign --rpc  https://rpc.ankr.com/fantom_testnet --from 0xba536E7ce173802053435bF03d1D528f3Ff29C32 --to 0x41cbC063B4b3264F5a075012e685B9fA05e41a44 --data 0xa9059cbb0000000000000000000000009d757dd679be17b4094c740fb0047fa3a7ed6df000000000000000000000000000000000000000000000000000000000000f4240

go run main.go tx offsign --rpc  https://rpc.ankr.com/fantom_testnet --from 0xba536E7ce173802053435bF03d1D528f3Ff29C32 --to 0x41cbC063B4b3264F5a075012e685B9fA05e41a44 --abi "transfer(address,uint256)" --args 0x9D757Dd679bE17b4094c740fB0047fa3a7Ed6DF0 --args 100000

### Usage
## 发送交易
### 基本用法
met tx send --to <> --value <> --data <> --network <> < --account <> | --ledger > [-v]

### 调用合约
met tx send --to <contractAddress> --abi <abi string>|<built-in abi: erc20,erc721,erc1155> --method <methodName> --args <arg1> ... --args <argN> --network <> < --account <> | --ledger > [-v]

因为abi会很长，所以可以将abi保存到文件中，然后通过--abi "$(cat abiFile)" 来传递abi

### erc20
met erc20 transferFrom --contract <> --from <> --to <> --amount <> --network <> < --account <> | --ledger > [-v]
met erc20 transfer --contract <> --to <> --amount <> --network <> < --account <> | --ledger > [-v]
met erc20 approve --contract <> --spender <> --amount <> --network <> < --account <> | --ledger > [-v]

## abi 编码
### 普通函数
met codec abiencode --abi 'transfer(string,string)' --args "My USDT" --args "MyUSDT"
### 构造函数
met codec abiencode --abi 'constructor(string,string)' --args "My USDT" --args "MyUSDT"
