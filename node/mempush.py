"""
This module is responsible for pushing mempool data
to the bitkit API
"""
import os
from time import sleep
import requests
from requests.auth import HTTPBasicAuth
from bitcoinrpc.authproxy import AuthServiceProxy

def make_connection_string():
    """Creates a connection string for connecting to a full bitcoin node
    """
    node_ip = os.environ['NODE_IP']
    rpc_user = os.environ['RPC_USER']
    rpc_password = os.environ['RPC_PASSWORD']
    return "http://{}:{}@{}".format(rpc_user, rpc_password, node_ip)

CONNECTION_STRING = make_connection_string()
def call_node(commands):
    """Creates a proxy for communicating with a full bitcoin node
    """
    return AuthServiceProxy(CONNECTION_STRING).batch_(commands)

def get_mempool(verbose=True):
    """Returns dictionary of form {txid1:{<info>}, txid2:{<info>}} and current block height
    info includes:
    size, fee, modifiedfee, time, height, descendantcount, descendantsize,
    descendantfees, ancestorcount, ancestorsize, ancestorfees, depends
    """
    data = call_node([["getrawmempool", verbose], ["getblockcount"]])
    return data[0], data[1]

def calculate_mining_fee_rate(mempool):
    """Calculates effective fee rate to mine a transaction.
    Parent transactions must be mined with or before child transactions.
    This means that the effective fee rate for a transaction
    is the maximum of [its average fee rate (including its parents)]
    and [the average fee rate including child transactions]
    """
    mining_fee_rate = {txid:info['ancestorfees'] / info['ancestorsize']
                       for txid, info in mempool.items()}
    current_and_depends = [[k, k] for k in mining_fee_rate.keys()]
    while current_and_depends:
        temp = current_and_depends.copy()
        current_and_depends = []
        for current, depends in temp:
            current_fee_rate = mining_fee_rate.get(current, 0)
            depends_fee_rate = mining_fee_rate.get(depends, 0)
            mining_fee_rate[depends] = max(current_fee_rate, depends_fee_rate)
            for new_depends in mempool[depends]['depends']:
                current_and_depends.append([depends, new_depends])
    return mining_fee_rate

def extract_and_transform(current_mempool, sent_txids):
    """Given mempool snapshot and already sent transactions,
    construct json-serializable object of new transactions
    WEIGHT CALCULATION INCOMPLETE
    """
    mining_fee_rate = calculate_mining_fee_rate(current_mempool)
    new_txs = []
    for txid, info in current_mempool.items():
        fee_rate = mining_fee_rate[txid]
        if txid not in sent_txids or sent_txids[txid] != fee_rate:
            vsize = int(float(info['size'])) # float(info['fee']) * 1e8
            sent_txids[txid] = fee_rate
            new_txs.append({'txid': txid, 'fee_rate': fee_rate, 'weight': vsize})
    return new_txs

def call_api(new_txs, method):
    """Calls bitkit api
    """
    uri = os.environ['BITKIT']
    data = {"method": method, "data":new_txs}
    authy = HTTPBasicAuth(os.environ['AUTH_USER'], os.environ['AUTH_PASSWORD'])
    response = requests.post(uri, json=data, verify=False, auth=authy)
    return response.ok

def main(sent_txids, max_height):
    """Main function in loop
    """
    current_mempool, height = get_mempool()
    if height > max_height:
        max_height = height
        sent_txids = {}
    method = "append" if sent_txids else "reset"
    new_txs = extract_and_transform(current_mempool, sent_txids)
    call_api(new_txs, method)
    infostr = "height {}, method {}, new_txs {}, sent_txids {}"
    print(infostr.format(height, method, len(new_txs), len(sent_txids)))
    return sent_txids, max_height

if __name__ == "__main__":
    SENT_TXIDS = {}
    MAX_HEIGHT = 0
    while True:
        SENT_TXIDS, MAX_HEIGHT = main(SENT_TXIDS, MAX_HEIGHT)
        sleep(5)
