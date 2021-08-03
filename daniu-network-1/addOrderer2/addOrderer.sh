#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This script extends the Hyperledger Fabric test network by adding
# adding a second orderer to the network
#

# prepending $PWD/../bin to PATH to ensure we are picking up the correct binaries
# this may be commented out to resolve installed version of tools if desired

export PATH=${PWD}/../../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}/../../config
export VERBOSE=false

. ./scripts/utils.sh
. ./scripts/envVar.sh
. ./scripts/envDaniu.sh

echo "Invoking chaincode using orderer1"
setGlobalsCLI 1
chaincodeInvokeAddCompany 1 2 3
chaincodeQueryCompany 1

echo "Starting Orderer CLI Container"
docker-compose -f ./docker/docker-compose-orderer-cli.yaml up -d 

echo "Adding new orderer TLS to the system channel (system-channel)"
docker exec orderer.cli ./scripts/orderer2-scripts/addTLSsys-channel.sh


echo "Fetch the latest configuration block"
docker exec orderer.cli ./scripts/orderer2-scripts/fetchConfigBlock.sh
docker cp orderer.cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/latest_config.block ../system-genesis-block/latest_config.block

echo "Bring Orderer2 Container"
docker-compose -f ./docker/docker-compose-orderer2.yaml up -d

echo "Adding new Orderer endpoint to the system channel (mychannel)"
docker exec orderer.cli ./scripts/orderer2-scripts/addEndPointSys-channel.sh

echo "System channel Size"
docker exec orderer.example.com ls -lh /var/hyperledger/production/orderer/chains/system-channel
docker exec orderer2.example.com ls -lh /var/hyperledger/production/orderer/chains/system-channel

echo "Application channel Size (before channel update)"
docker exec orderer.example.com ls -lh /var/hyperledger/production/orderer/chains/mychannel
docker exec orderer2.example.com ls -lh /var/hyperledger/production/orderer/chains/mychannel

echo "Add new orderer TLS to the application channel"
docker exec orderer.cli ./scripts/orderer2-scripts/addTLSapplication-channel.sh

echo "Adding new Orderer endpoint to the application channel (mychannel)"
docker exec orderer.cli ./scripts/orderer2-scripts/addEndPointapplication-channel.sh

echo "Application channel Size (after channel update)"
docker exec orderer.example.com ls -lh /var/hyperledger/production/orderer/chains/mychannel
docker exec orderer2.example.com ls -lh /var/hyperledger/production/orderer/chains/mychannel

# echo "Invoking chaincode using orderer2"
# ./scripts/transaction-ord2.sh