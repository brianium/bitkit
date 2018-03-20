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

def extract_and_transform(current_mempool, sent_txids):
    """Given mempool snapshot and already sent transactions,
    construct json-serializable object of new transactions
    WEIGHT CALCULATION INCOMPLETE
    """
    new_txs = []
    for txid, info in current_mempool.items():
        if txid not in sent_txids:
            fee = info['ancestorfees'] + info['descendantfees'] - float(info['fee']) * 1e8
            size = int(info['ancestorsize'] + info['descendantsize'] - float(info['size']))
            fee_rate = fee / size
            new_txs.append({'txid': txid, 'fee_rate': fee_rate, 'weight': size})
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
        sent_txids = set()
    new_txs = extract_and_transform(current_mempool, sent_txids)
    method = "append" if sent_txids else "reset"
    if call_api(new_txs, method):
        sent_txids = sent_txids.union(t['txid'] for t in new_txs)
    infostr = "height {}, method {}, new_txs {}, sent_txids {}"
    print(infostr.format(height, method, len(new_txs), len(sent_txids)))
    return sent_txids, max_height

if __name__ == "__main__":
    SENT_TXIDS = set()
    MAX_HEIGHT = 0
    while True:
        SENT_TXIDS, MAX_HEIGHT = main(SENT_TXIDS, MAX_HEIGHT)
        sleep(5)
