# 系统软件安装
## 注意在每台机器上都要安装
上传 daniu-install.zip 的内容到 /root/fabric 目录

参考 daniu-install/README.md 安装系统，设置环境变量等

# 部署区块链环境
## 在第一台上47.97.101.21

将`daniu-network-prod`目录上传至`/root/fabric/fabric-samples/`目录

./network.sh up createChannel -ca -c mychannel -s couchdb

scp -r /root/fabric/fabric-samples/daniu-network-prod/* root@106.15.180.142:/root/fabric/fabric-samples/daniu-network-prod

## 在第二台上106.15.180.142 root/aUH8Ie6o

cd /root/fabric/fabric-samples/daniu-network-prod/addOrg3

./addOrg3.sh up -c mychannel -s couchdb

echo "同步peerOrganizations(org3)到47.97.101.21"
scp -r /root/fabric/fabric-samples/daniu-network-prod/organizations/peerOrganizations/* root@47.97.101.21:/root/fabric/fabric-samples/daniu-network-prod/organizations/peerOrganizations



# 部署链码
root/Niuinfo.com123!

## 在第一台上47.97.101.21

scp -r root@47.97.101.21:/root/fabric/fabric-samples/config root@106.15.180.142:/root/fabric/fabric-samples/config

cd /root/fabric/fabric-samples/daniu-network-prod/

source ./scripts/envDaniu.sh

./scripts/deployCC12.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

## 在第二台上106.15.180.142

cd /root/fabric/fabric-samples/daniu-network-prod/

source ./scripts/envDaniu.sh

./scripts/deployCC3.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE


## Cli 链码调用

export FABRIC_CFG_PATH=/root/fabric/fabric-samples/config
cd /root/fabric/fabric-samples/daniu-network-prod
source scripts/envVar.sh
setGlobals 1
setGlobalsCLI 1
peer chaincode query -C mychannel -n daniu_1 -c '{"function":"QueryCompany", "Args":[ "{\"CompanyName\": \"上海达牛信息技术有限公司\", \"CompanyCode\":  \"123456789123456789\"}","100",""]}'


## 升级链码
假设链码目录为 daniu_v2, 新增方法 QueryCompanyV2

## 在第一台上47.97.101.21

cd /root/fabric/fabric-samples/daniu-network-prod/

source ./scripts/envDaniu_v2.sh

./scripts/deployCC12.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

## 在第二台上106.15.180.142

cd /root/fabric/fabric-samples/daniu-network-prod/

source ./scripts/envDaniu_v2.sh

./scripts/deployCC3.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE

