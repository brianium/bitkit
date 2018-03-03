import os
import numpy as np
from bitcoinrpc.authproxy import AuthServiceProxy, JSONRPCException

def make_connection_string():
    node_ip = os.environ['NODE_IP']
    rpc_user = os.environ['RPC_USER']
    rpc_password = os.environ['RPC_PASSWORD']
    return "http://{}:{}@{}".format(rpc_user, rpc_password, node_ip)
	
connection_string = make_connection_string()
def engine(commands):
	return AuthServiceProxy(connection_string).batch_(commands)

def get_satpb():
	memtx = engine([["getrawmempool", True]])[0]
	fees = np.array([v['ancestorfees'] + v['descendantfees'] - float(v['fee'])*1e8 for k,v in memtx.items()], dtype=np.float32)
	sizes = np.array([v['ancestorsize'] + v['descendantsize'] - float(v['size']) for k,v in memtx.items()], dtype=np.float32)
	satpb = fees/sizes
	return satpb


if __name__ == "__main__":
    print(len(get_satpb()))
