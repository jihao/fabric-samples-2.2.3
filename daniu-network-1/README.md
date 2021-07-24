## 在第一台上10.18.188.177

将`daniu-network-1`目录上传至`/root/fabric/fabric-samples/`目录

./network.sh up createChannel -ca -c mychannel -s couchdb

scp -r /root/fabric/fabric-samples/daniu-network-1/* root@10.18.188.178:/root/fabric/fabric-samples/daniu-network-1

## 在第二台上10.18.188.178

cd /root/fabric/fabric-samples/daniu-network-1/addOrg3

./addOrg3.sh up -c mychannel -s couchdb

### org3 & org4 都运行在第二台上10.18.188.178
cd /root/fabric/fabric-samples/daniu-network-1/addOrg4

./addOrg4.sh up -c mychannel -s couchdb

echo "同步peerOrganizations(org3,org4)到10.18.188.177"
scp -r /root/fabric/fabric-samples/daniu-network-1/organizations/peerOrganizations/* root@10.18.188.177:/root/fabric/fabric-samples/daniu-network-1/organizations/peerOrganizations

echo "同步节点配置到10.18.188.179"
scp -r root@10.18.188.178:/root/fabric/fabric-samples/daniu-network-1/* root@10.18.188.179:/root/fabric/fabric-samples/daniu-network-1

## 在第三台上10.18.188.179
cd /root/fabric/fabric-samples/daniu-network-1/addOrg5

./addOrg5.sh up -c mychannel -s couchdb