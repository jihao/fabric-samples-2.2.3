#!/bin/bash
export FABRIC_CFG_PATH=/root/fabric/fabric-samples/config
export VERBOSE=false
export MAX_RETRY=5
export CLI_DELAY=3
export DELAY=3
export CHANNEL_NAME="mychannel"
export CC_NAME="daniu_1"
export CC_SRC_PATH="./daniu"
export CC_END_POLICY="NA"
export CC_COLL_CONFIG="NA"
export CC_INIT_FCN="NA"
export CC_SRC_LANGUAGE="go"
export CC_VERSION="1.0"
export CC_SEQUENCE=1


chaincodeInvokeAddCompany() {
  parsePeerConnectionParameters $@
  res=$?
  verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
  fcn_call='{"function":"AddCompany", "Args":[ "{\"CompanyName\":\"上海达牛信息DEMO有限公司\",\"RegAdd\":\"曲阳路街道\",\"AreaCode\":110101,\"RegCapital\":10,\"RegDate\":20190101, \"BizScope\":\"网络信息、计算机科技领域内的技术开发、\",\"TransportNumber\":\"989845878484\",\"RegTel\":\"010-12345678\",\"CompanyCode\":\"999956789999956789\",\"LARName\":\"赵先生\",\"LARMobile\":18222931234,\"ContactsName\":\"赵先生\",\"ContactsTel\":13835551111,\"EstablishedTime\":20190101,\"BizLicensePhotoURL\":\"https://www.a.com/营业执照图片.jpg\",\"TransportPhotoURL\":\"https://www.a.com/道路运输经营许可证.jpg\",\"ICPNumber\":\"48484844848\",\"SecurityNumber\":\"122448475876\",\"PlatformCompanyName\":\"上海达牛信息DEMO有限公司\",\"PlatformCompanyCode\":\"3243432432432432432F\",\"Ext\":\"扩展信息\"}","100",""]}'
  infoln "invoke fcn call:${fcn_call}"
  set -x
  peer chaincode invoke -o orderer.niuinfo.com:7050 --ordererTLSHostnameOverride orderer.niuinfo.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c "${fcn_call}" >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Invoke execution on $PEERS failed "
  successln "Invoke transaction successful on $PEERS on channel '$CHANNEL_NAME'"
}

chaincodeQueryCompany() {
  ORG=$1
  setGlobals $ORG
  infoln "Querying on peer0.org${ORG} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  fcn_call='{"function":"QueryCompany", "Args":[ "{\"CompanyName\":\"上海达牛信息DEMO有限公司\",\"CompanyCode\":\"999956789999956789\"}","100",""]}'
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to Query peer0.org${ORG}, Retry after $DELAY seconds."
    set -x
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c "${fcn_call}" >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
    successln "Query successful on peer0.org${ORG} on channel '$CHANNEL_NAME'"
  else
    fatalln "After $MAX_RETRY attempts, Query result on peer0.org${ORG} is INVALID!"
  fi
}