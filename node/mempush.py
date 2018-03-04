"""
This module is responsible for pushing mempool data
to the memcool API
"""
import os
from time import sleep
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
    """Returns dictionary of form {txid1:{<info>}, txid2:{<info>}}
    info includes:
    size, fee, modifiedfee, time, height, descendantcount, descendantsize,
    descendantfees, ancestorcount, ancestorsize, ancestorfees, depends
    """
    return call_node([["getrawmempool", verbose]])[0]

def extract_and_transform(current_mempool, sent_txids):
    """Given mempool snapshot and already sent transactions,
    construct json-serializable object of new transactions
    WEIGHT CALCULATION INCOMPLETE
    """
    new_txs = []
    for txid, info in current_mempool.items():
        if txid not in sent_txids:
            fee = info['ancestorfees'] + info['descendantfees'] - float(info['fee']) * 1e8
            size = info['ancestorsize'] + info['descendantsize'] - float(info['size'])
            fee_rate = fee / size
            new_txs.append({'txid': txid, 'fee_rate': fee_rate, 'weight': size})
            sent_txids.add(txid)
    return new_txs, sent_txids

def get_max_height(current_mempool):
    """Checks for new block
    """
    return max(info['height'] for _, info in current_mempool.items())

def call_api(new_txs):
    """Calls memcool api
    """
    print(len(new_txs)) # replace with call to api
    successful_post = True
    return successful_post

def main(sent_txids, max_height):
    """Main function in loop
    """
    current_mempool = get_mempool()
    height = get_max_height(current_mempool)
    if height > max_height:
        max_height = height
        sent_txids = set()
    new_txs, potention_sent_txids = extract_and_transform(current_mempool, sent_txids)
    if call_api(new_txs):
        sent_txids = potention_sent_txids
    return sent_txids, max_height

if __name__ == "__main__":
    SENT_TXIDS = set()
    MAX_HEIGHT = 0
    while True:
        SENT_TXIDS, MAX_HEIGHT = main(SENT_TXIDS, MAX_HEIGHT)
        sleep(5)
