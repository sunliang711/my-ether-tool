#!/bin/bash
if [ -z "${BASH_SOURCE}" ]; then
    this=${PWD}
else
    rpath="$(readlink ${BASH_SOURCE})"
    if [ -z "$rpath" ]; then
        rpath=${BASH_SOURCE}
    elif echo "$rpath" | grep -q '^/'; then
        # absolute path
        echo
    else
        # relative path
        rpath="$(dirname ${BASH_SOURCE})/$rpath"
    fi
    this="$(cd $(dirname $rpath) && pwd)"
fi

if [ -r ${SHELLRC_ROOT}/shelllib ]; then
    source ${SHELLRC_ROOT}/shelllib
elif [ -r /tmp/shelllib ]; then
    source /tmp/shelllib
else
    # download shelllib then source
    shelllibURL=https://cdn.jsdelivr.net/gh/sunliang711/init/shellConfigs/shelllib
    (cd /tmp && curl -s -LO ${shelllibURL})
    if [ -r /tmp/shelllib ]; then
        source /tmp/shelllib
    fi
fi

# available VARs: user, home, rootID
# available functions:
#    _err(): print "$*" to stderror
#    _command_exists(): check command "$1" existence
#    _require_command(): exit when command "$1" not exist
#    _runAsRoot():
#                  -x (trace)
#                  -s (run in subshell)
#                  --nostdout (discard stdout)
#                  --nostderr (discard stderr)
#    _insert_path(): insert "$1" to PATH
#    _run():
#                  -x (trace)
#                  -s (run in subshell)
#                  --no-stdout (discard stdout)
#                  --no-stderr (discard stderr)
#    _ensureDir(): mkdir if $@ not exist
#    _root(): check if it is run as root
#    _require_root(): exit when not run as root
#    _linux(): check if it is on Linux
#    _require_linux(): exit when not on Linux
#    _wait(): wait $i seconds in script
#    _must_ok(): exit when $? not zero
#    _info(): info log
#    _infoln(): info log with \n
#    _error(): error log
#    _errorln(): error log with \n
#    _checkService(): check $1 exist in systemd

###############################################################################
# write your code below (just define function[s])
# function is hidden when begin with '_'
function _parseOptions() {
    if [ $(uname) != "Linux" ]; then
        echo "getopt only on Linux"
        exit 1
    fi

    options=$(getopt -o dv --long debug --long name: -- "$@")
    [ $? -eq 0 ] || {
        echo "Incorrect option provided"
        exit 1
    }
    eval set -- "$options"
    while true; do
        case "$1" in
        -v)
            VERBOSE=1
            ;;
        -d)
            DEBUG=1
            ;;
        --debug)
            DEBUG=1
            ;;
        --name)
            shift # The arg is next in position args
            NAME=$1
            ;;
        --)
            shift
            break
            ;;
        esac
        shift
    done
}

_example() {
    _parseOptions "$0" "$@"
    # TODO
}
add() {
    if [ -z "$ID" ]; then
        echo "Missing env var ID" 1>&2
        exit 1
    fi

    ./met network add --name eth --rpc "https://mainnet.infura.io/v3/${ID}" --explorer https://etherscan.io --symbol ETH
    ./met network add --name  polygon --rpc https://polygon-rpc.com --explorer https://polygonscan.com --symbol MATIC
    ./met network add --name  bsc --rpc https://bsc-dataseed1.binance.org --explorer https://bscscan.com --symbol BNB
    ./met network add --name  op --rpc https://mainnet.optimism.io --explorer https://optimistic.etherscan.io --symbol ETH
    ./met network add --name  arbitrum --rpc https://arb1.arbitrum.io/rpc --explorer https://arbiscan.io --symbol ETH
    ./met network add --name  arbi --rpc https://arb1.arbitrum.io/rpc --explorer https://arbiscan.io --symbol ETH
    ./met network add --name  goerli --rpc "https://goerli.infura.io/v3/${ID}" --explorer https://goerli.etherscan.io --symbol GETH
    ./met network add --name  sepolia --rpc "https://sepolia.infura.io/v3/${ID}" --explorer https://sepolia.etherscan.io --symbol ETH
    ./met network add --name  ftm --rpc https://1rpc.io/ftm --explorer https://ftmscan.com --symbol FTM
    ./met network add --name  ftmTest --rpc https://rpc.ankr.com/fantom_testnet --explorer https://testnet.ftmscan.com --symbol FTM
    ./met network add --name  avax --rpc https://1rpc.io/avax/c --explorer https://snowtrace.io --symbol AVAX
    ./met network add --name  bscTest --rpc https://data-seed-prebsc-1-s2.binance.org:8545 --explorer https://testnet.bscscan.com --symbol tBNB

}
# write your code above
###############################################################################

em() {
    $ed $0
}

function _help() {
    cd "${this}"
    cat <<EOF2
Usage: $(basename $0) ${bold}CMD${reset}

${bold}CMD${reset}:
EOF2
    perl -lne 'print "\t$2" if /^\s*(function)?\s*(\S+)\s*\(\)\s*\{$/' $(basename ${BASH_SOURCE}) | perl -lne "print if /^\t[^_]/"
}

case "$1" in
"" | -h | --help | help)
    _help
    ;;
*)
    "$@"
    ;;
esac
