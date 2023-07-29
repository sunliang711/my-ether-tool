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
