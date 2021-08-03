peer channel fetch config config_block.pb -o orderer.example.com:7050 -c system-channel --tls --cafile $ORDERER_CA
configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json
jq ".channel_group.values.OrdererAddresses.value.addresses += [\"orderer2.example.com:8050\"]" config.json > modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id system-channel --original config.pb --updated modified_config.pb --output config_update.pb
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo "{\"payload\":{\"header\":{\"channel_header\":{\"channel_id\":\"system-channel\", \"type\":2}},\"data\":{\"config_update\":"$(cat config_update.json)"}}}" | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
peer channel update -f config_update_in_envelope.pb -c system-channel -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA
