# 部署区块链环境

此目录daniu-network-2的目的是为了验证账本数据&couchdb数据库的备份和恢复，以及升级到2.2+LTS版本的数据迁移

数据来源为 10.18.188.177/178，目标环境 173/174，恢复数据org1,org2,org3,org4的账本及数据库, 不包含org5

基本步骤为
1. 安装新环境，不安装链码
2. 备份账本 -> 迁移到新环境 -> docker重启peer，orderer
3. 备份数据库 -> 迁移到新环境 -> docker重启db
4. 验证数据是否存在

# 步骤一：安装
## 在第一台上10.18.188.173  
将`daniu-network-2`目录上传至`/root/fabric/fabric-samples/`目录  

./network.sh up createChannel -ca -c mychannel -s couchdb  
scp -r /root/fabric/fabric-samples/daniu-network-2/* root@10.18.188.174:/root/fabric/fabric-samples/daniu-network-2  

## 在第二台上10.18.188.174

### org3 & org4 都运行在第二台上10.18.188.174

cd /root/fabric/fabric-samples/daniu-network-2/addOrg3  
./addOrg3.sh up -c mychannel -s couchdb  

cd /root/fabric/fabric-samples/daniu-network-2/addOrg4   
./addOrg4.sh up -c mychannel -s couchdb  

echo "同步peerOrganizations(org3,org4)到10.18.188.173"  
scp -r /root/fabric/fabric-samples/daniu-network-2/organizations/peerOrganizations/* root@10.18.188.173:/root/fabric/fabric-samples/daniu-network-2/organizations/peerOrganizations  


----
# 步骤二：恢复

# 登录10.18.188.177 -> 10.18.188.173 
root/Niuinfo.com123!

## 复制ca数据，进行upgrade
### 10.18.188.177
scp -r /root/fabric/fabric-samples/daniu-network-1/organizations/fabric-ca/* root@10.18.188.173:/root/fabric/fabric-samples/daniu-network-2/organizations/fabric-ca/

### 10.18.188.173 
docker restart ca_orderer ca_org1 ca_org2

## 复制order节点数据, 依次进行rolling upgrade
### 10.18.188.177
scp -r  /var/lib/docker/volumes/net_orderer.example.com/_data/chains/ root@10.18.188.173:/var/lib/docker/volumes/net_orderer.example.com/_data/chains/
### 10.18.188.173 
docker restart orderer.example.com
docker logs -f orderer.example.com

### 10.18.188.177
scp -r  /var/lib/docker/volumes/net_orderer2.example.com/_data/chains/ root@10.18.188.173:/var/lib/docker/volumes/net_orderer2.example.com/_data/chains/
### 10.18.188.173 
docker restart orderer2.example.com
docker logs -f orderer2.example.com

### 10.18.188.177
scp -r  /var/lib/docker/volumes/net_orderer3.example.com/_data/chains/ root@10.18.188.173:/var/lib/docker/volumes/net_orderer3.example.com/_data/chains/
### 10.18.188.173 
docker restart orderer3.example.com
docker logs -f orderer3.example.com

### 10.18.188.177
scp -r  /var/lib/docker/volumes/net_orderer4.example.com/_data/chains/ root@10.18.188.173:/var/lib/docker/volumes/net_orderer4.example.com/_data/chains/
### 10.18.188.173 
docker restart orderer4.example.com
docker logs -f orderer4.example.com

### 10.18.188.177
scp -r  /var/lib/docker/volumes/net_orderer5.example.com/_data/chains/ root@10.18.188.173:/var/lib/docker/volumes/net_orderer5.example.com/_data/chains/
### 10.18.188.173 
docker restart orderer5.example.com
docker logs -f orderer5.example.com

## 复制peer节点数据
### 10.18.188.177
scp -r /var/lib/docker/volumes/net_peer* root@10.18.188.173:/var/lib/docker/volumes/
### 10.18.188.173 
docker restart peer0.org1.example.com
docker restart peer0.org2.example.com

docker logs -f peer0.org1.example.com


## 复制couchdb节点数据
### 10.18.188.177
export db0_volume=`docker inspect couchdb0 | grep volumes | awk -F'"' '{print $4}'`
scp -r $db0_volume root@10.18.188.173:/var/lib/docker/volumes/net_couchdb0/

export db1_volume=`docker inspect couchdb1 | grep volumes | awk -F'"' '{print $4}'`
scp -r $db1_volume root@10.18.188.173:/var/lib/docker/volumes/net_couchdb1/
### 10.18.188.173 
docker restart couchdb0
docker restart couchdb1

----
# 登录10.18.188.178 -> 10.18.188.174 
root/Niuinfo.com123!
### 10.18.188.177
export db4_volume=`docker inspect couchdb4 | grep volumes | awk -F'"' '{print $4}'`
scp -r $db4_volume root@10.18.188.174:/var/lib/docker/volumes/net_couchdb4/

export db5_volume=`docker inspect couchdb5 | grep volumes | awk -F'"' '{print $4}'`
scp -r $db5_volume root@10.18.188.174:/var/lib/docker/volumes/net_couchdb5/

scp -r /var/lib/docker/volumes/net_peer* root@10.18.188.174:/var/lib/docker/volumes/

### 10.18.188.174
docker restart couchdb4
docker restart couchdb5
docker restart peer0.org3.example.com
docker restart peer0.org4.example.com

# 步骤三：验证



---
# Notes

## 备份账本
备份账本数据
the ledger data has not been changed from the default value of /var/hyperledger/production/ (for peers) or /var/hyperledger/production/orderer (for ordering nodes)
If using CouchDB as state database, there will be no stateLeveldb directory, as the state database data would be stored within CouchDB instead.

对应的docker volumes
/var/lib/docker/volumes/net_orderer.example.com
/var/lib/docker/volumes/net_orderer2.example.com
/var/lib/docker/volumes/net_orderer3.example.com
/var/lib/docker/volumes/net_orderer4.example.com
/var/lib/docker/volumes/net_orderer5.example.com
/var/lib/docker/volumes/net_peer0.org1.example.com
/var/lib/docker/volumes/net_peer0.org2.example.com
/var/lib/docker/volumes/net_peer1.org1.example.com
/var/lib/docker/volumes/net_peer1.org2.example.com


net_peer0.org3.example.com
net_peer0.org4.example.com

## 备份couchdb

Upgrading CouchDB
If you are using CouchDB as state database, you should upgrade the peer’s CouchDB at the same time the peer is being upgraded.

To upgrade CouchDB:

Stop CouchDB.
Backup CouchDB data directory.
Install the latest CouchDB binaries or update deployment scripts to use a new Docker image.
Restart CouchDB.

访问 couchdb http://127.0.0.1:7984/_utils/
用户名/密码 admin/adminpw docker-compose-couchdb.yml 

docker inspect couchdb0 | grep volumes