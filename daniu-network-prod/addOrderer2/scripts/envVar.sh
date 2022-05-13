#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

# imports
. scripts/utils.sh

export BASEDIR=/root/fabric/fabric-samples/daniu-network-prod
export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${BASEDIR}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem
export PEER0_ORG1_CA=${BASEDIR}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/ca.crt
export PEER0_ORG2_CA=${BASEDIR}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/ca.crt
export PEER0_ORG3_CA=${BASEDIR}/organizations/peerOrganizations/org1.dayunban.com/peers/peer0.org1.dayunban.com/tls/ca.crt
export PEER0_ORG4_CA=${BASEDIR}/organizations/peerOrganizations/org4.niuinfo.com/peers/peer0.org4.niuinfo.com/tls/ca.crt
export PEER0_ORG5_CA=${BASEDIR}/organizations/peerOrganizations/org5.niuinfo.com/peers/peer0.org5.niuinfo.com/tls/ca.crt

# Set environment variables for the peer org
setGlobals() {
  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  infoln "Using organization ${USING_ORG}"
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${BASEDIR}/organizations/peerOrganizations/org1.niuinfo.com/users/Admin@org1.niuinfo.com/msp
    # export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_ADDRESS=peer0.org1.niuinfo.com:7051
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${BASEDIR}/organizations/peerOrganizations/org2.niuinfo.com/users/Admin@org2.niuinfo.com/msp
    # export CORE_PEER_ADDRESS=localhost:9051
    export CORE_PEER_ADDRESS=peer0.org2.niuinfo.com:9051
  elif [ $USING_ORG -eq 3 ]; then
    export CORE_PEER_LOCALMSPID="Org3MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
    export CORE_PEER_MSPCONFIGPATH=${BASEDIR}/organizations/peerOrganizations/org1.dayunban.com/users/Admin@org1.dayunban.com/msp
    # export CORE_PEER_ADDRESS=localhost:11051
    export CORE_PEER_ADDRESS=peer0.org1.dayunban.com:11051
  elif [ $USING_ORG -eq 4 ]; then
    export CORE_PEER_LOCALMSPID="Org4MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG4_CA
    export CORE_PEER_MSPCONFIGPATH=${BASEDIR}/organizations/peerOrganizations/org4.niuinfo.com/users/Admin@org4.niuinfo.com/msp
    # export CORE_PEER_ADDRESS=localhost:12051
    export CORE_PEER_ADDRESS=peer0.org4.niuinfo.com:12051
  elif [ $USING_ORG -eq 5 ]; then
    export CORE_PEER_LOCALMSPID="Org5MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG5_CA
    export CORE_PEER_MSPCONFIGPATH=${BASEDIR}/organizations/peerOrganizations/org5.niuinfo.com/users/Admin@org5.niuinfo.com/msp
    # export CORE_PEER_ADDRESS=localhost:13051
    export CORE_PEER_ADDRESS=peer0.org5.niuinfo.com:13051
  else
    errorln "ORG Unknown"
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

# Set environment variables for use in the CLI container 
setGlobalsCLI() {
  setGlobals $1

  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_ADDRESS=peer0.org1.niuinfo.com:7051
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_ADDRESS=peer0.org2.niuinfo.com:9051
  elif [ $USING_ORG -eq 3 ]; then
    export CORE_PEER_ADDRESS=peer0.org1.dayunban.com:11051
  elif [ $USING_ORG -eq 4 ]; then
    export CORE_PEER_ADDRESS=peer0.org4.niuinfo.com:12051
  elif [ $USING_ORG -eq 5 ]; then
    export CORE_PEER_ADDRESS=peer0.org5.niuinfo.com:13051
  else
    errorln "ORG Unknown"
  fi
}

# parsePeerConnectionParameters $@
# Helper function that sets the peer connection parameters for a chaincode
# operation
parsePeerConnectionParameters() {
  PEER_CONN_PARMS=""
  PEERS=""
  while [ "$#" -gt 0 ]; do
    setGlobals $1
    PEER="peer0.org$1"
    ## Set peer addresses
    PEERS="$PEERS $PEER"
    PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
    ## Set path to TLS certificate
    TLSINFO=$(eval echo "--tlsRootCertFiles \$PEER0_ORG$1_CA")
    PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
    # shift by one to get to the next organization
    shift
  done
  # remove leading space for output
  PEERS="$(echo -e "$PEERS" | sed -e 's/^[[:space:]]*//')"
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}