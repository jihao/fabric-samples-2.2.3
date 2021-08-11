# 部署区块链环境
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


# 删除组织（Org5）
## 在第三台上10.18.188.179
cd /root/fabric/fabric-samples/daniu-network-1/addOrg5
./delOrg5.sh del -c mychannel -s couchdb

echo "成功删除组织Org5后，可以停止相关docker容器"
./delOrg5.sh down -c mychannel -s couchdb

相关文件2个，如果要删除其他组织，请参考
addOrg5/delOrg5.sh
scripts/org5-scripts/updateChannelConfigDel.sh


# 部署链码
root/Niuinfo.com123!

## 在第一台上10.18.188.177

scp -r root@10.18.188.177:/root/fabric/fabric-samples/config root@10.18.188.178:/root/fabric/fabric-samples/config
scp -r root@10.18.188.177:/root/fabric/fabric-samples/config root@10.18.188.179:/root/fabric/fabric-samples/config

cd /root/fabric/fabric-samples/daniu-network-1/

source ./scripts/envDaniu.sh

./scripts/deployCC12.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

## 在第二台上10.18.188.178

cd /root/fabric/fabric-samples/daniu-network-1/

source ./scripts/envDaniu.sh

./scripts/deployCC34.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

## 在第三台上10.18.188.179

cd /root/fabric/fabric-samples/daniu-network-1/

source ./scripts/envDaniu.sh

./scripts/deployCC5.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

./scripts/deployCCCommit.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

## Cli 链码调用

export FABRIC_CFG_PATH=/root/fabric/fabric-samples/config
cd /root/fabric/fabric-samples/daniu-network-1
source scripts/envVar.sh
setGlobals 1
setGlobalsCLI 1
peer chaincode query -C mychannel -n daniu_1 -c '{"function":"QueryCompany", "Args":[ "{\"CompanyName\": \"上海达牛信息技术有限公司\", \"CompanyCode\":  \"123456789123456789\"}","100",""]}'


## 升级链码
假设链码目录为 daniu_v2, 新增方法 QueryCompanyV2

## 在第一台上10.18.188.177

cd /root/fabric/fabric-samples/daniu-network-1/

source ./scripts/envDaniu_v2.sh

./scripts/deployCC12.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

## 在第二台上10.18.188.178

cd /root/fabric/fabric-samples/daniu-network-1/

source ./scripts/envDaniu_v2.sh

./scripts/deployCC34.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

## 在第三台上10.18.188.179

cd /root/fabric/fabric-samples/daniu-network-1/

source ./scripts/envDaniu_v2.sh

./scripts/deployCC5.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

./scripts/deployCCCommit.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

