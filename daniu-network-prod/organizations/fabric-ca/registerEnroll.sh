#!/bin/bash

function createOrg1() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/org1.niuinfo.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.niuinfo.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-org1 --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/org1.niuinfo.com/msp/config.yaml

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-org1 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-org1 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-org1 --id.name org1admin --id.secret org1adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/msp --csr.hosts peer0.org1.niuinfo.com --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/msp/config.yaml

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls --enrollment.profile tls --csr.hosts peer0.org1.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/server.key

  mkdir -p ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/msp/tlscacerts/ca.crt

  mkdir -p ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/tlsca
  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/tlsca/tlsca.org1.niuinfo.com-cert.pem

  mkdir -p ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/ca
  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/peers/peer0.org1.niuinfo.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/ca/ca.org1.niuinfo.com-cert.pem

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/users/User1@org1.niuinfo.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/users/User1@org1.niuinfo.com/msp/config.yaml

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://org1admin:org1adminpw@localhost:7054 --caname ca-org1 -M ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/users/Admin@org1.niuinfo.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.niuinfo.com/users/Admin@org1.niuinfo.com/msp/config.yaml
}

function createOrg2() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/org2.niuinfo.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org2.niuinfo.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca-org2 --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-org2.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/org2.niuinfo.com/msp/config.yaml

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-org2 --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-org2 --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-org2 --id.name org2admin --id.secret org2adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/msp --csr.hosts peer0.org2.niuinfo.com --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/msp/config.yaml

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls --enrollment.profile tls --csr.hosts peer0.org2.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/server.key

  mkdir -p ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/msp/tlscacerts/ca.crt

  mkdir -p ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/tlsca
  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/tlsca/tlsca.org2.niuinfo.com-cert.pem

  mkdir -p ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/ca
  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/peers/peer0.org2.niuinfo.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/ca/ca.org2.niuinfo.com-cert.pem

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/users/User1@org2.niuinfo.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/users/User1@org2.niuinfo.com/msp/config.yaml

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://org2admin:org2adminpw@localhost:8054 --caname ca-org2 -M ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/users/Admin@org2.niuinfo.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.niuinfo.com/users/Admin@org2.niuinfo.com/msp/config.yaml
}

function createOrderer() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/niuinfo.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/niuinfo.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/config.yaml

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  generateOrderer1
  generateOrderer2
  generateOrderer3
  generateOrderer4
  generateOrderer5

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/users/Admin@niuinfo.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/niuinfo.com/users/Admin@niuinfo.com/msp/config.yaml
}

function generateOrderer1() {
  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/msp --csr.hosts orderer.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/msp/config.yaml

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls --enrollment.profile tls --csr.hosts orderer.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/ca.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/server.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/server.key

  mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem

  mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem
}

function generateOrderer2() {
  infoln "Generating the orderer2 msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/msp --csr.hosts orderer2.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/msp/config.yaml

  infoln "Generating the orderer2-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls --enrollment.profile tls --csr.hosts orderer2.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/ca.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/server.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/server.key

  mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem

  #mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts
  #cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem
}

function generateOrderer3() {
  infoln "Generating the orderer3 msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/msp --csr.hosts orderer3.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/msp/config.yaml

  infoln "Generating the orderer3-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls --enrollment.profile tls --csr.hosts orderer3.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/ca.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/server.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/server.key

  mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem

  #mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts
  #cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer3.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem
}

function generateOrderer4() {
  infoln "Generating the orderer4 msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/msp --csr.hosts orderer4.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/msp/config.yaml

  infoln "Generating the orderer4-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls --enrollment.profile tls --csr.hosts orderer4.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/ca.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/server.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/server.key

  mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem

  #mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts
  #cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer4.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem
}

function generateOrderer5() {
  infoln "Generating the orderer5 msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/msp --csr.hosts orderer5.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/msp/config.yaml

  infoln "Generating the orderer5-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls --enrollment.profile tls --csr.hosts orderer5.niuinfo.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/ca.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/server.crt
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/server.key

  mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem

  #mkdir -p ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts
  #cp ${PWD}/organizations/ordererOrganizations/niuinfo.com/orderers/orderer5.niuinfo.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem
}
