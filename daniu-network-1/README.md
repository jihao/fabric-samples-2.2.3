## 在第一台上10.18.188.177

将`daniu_network-1`目录上传至`/root/fabric/fabric-samples/`目录

./network.sh up createChannel -ca -c mychannel -s couchdb

scp -r /root/fabric/fabric-samples/daniu-network-1 root@10.18.188.178:/root/fabric/fabric-samples/daniu-network-1


## 在第二台上10.18.188.178

cd /root/fabric/fabric-samples/daniu_network-1/addOrg3

./addOrg3.sh up -c mychannel -s couchdb


