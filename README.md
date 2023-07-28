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
