TLS_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/niuinfo.com/orderers/orderer2.niuinfo.com/tls/server.crt
peer channel fetch config config_block.pb -o orderer.niuinfo.com:7050 -c mychannel --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem
configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json
echo "{\"client_tls_cert\":\"$(cat $TLS_FILE | base64)\",\"host\":\"orderer2.niuinfo.com\",\"port\":8050,\"server_tls_cert\":\"$(cat $TLS_FILE | base64)\"}" > $PWD/orderer2.json
jq ".channel_group.groups.Orderer.values.ConsensusType.value.metadata.consenters += [$(cat orderer2.json)]" config.json > modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id mychannel --original config.pb --updated modified_config.pb --output config_update.pb
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo "{\"payload\":{\"header\":{\"channel_header\":{\"channel_id\":\"mychannel\", \"type\":2}},\"data\":{\"config_update\":"$(cat config_update.json)"}}}" | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
peer channel update -f config_update_in_envelope.pb -c mychannel -o orderer.niuinfo.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/niuinfo.com/orderers/orderer.niuinfo.com/msp/tlscacerts/tlsca.niuinfo.com-cert.pem