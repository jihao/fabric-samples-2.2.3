package com.moshu.fabric.api.gateway;

import org.hyperledger.fabric.gateway.Contract;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.gateway.Wallet;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Date;


public final class FabricGatewayFactory {

    private static Gateway gateway = null;

    public static Contract getContract(String channel, String chaincode) throws WalletCreateException, IdentityCreateException {
        if (gateway == null) {
            gateway = connect();
        }
        Network network = gateway.getNetwork(channel);
        return network.getContract(chaincode);
    }

    public static Gateway getGateWay() throws WalletCreateException, IdentityCreateException {
        if (gateway == null)
            gateway = connect();
        return gateway;
    }

    private static Gateway connect() throws WalletCreateException, IdentityCreateException {
        long time1 = new Date().getTime();
        System.out.println(time1);
        Path walletPath = Paths.get("wallet");
        Wallet wallet = null;
        try {
            wallet = Wallet.createFileSystemWallet(walletPath);
        } catch (Exception e) {
            throw new WalletCreateException("钱包获取失败", e.getCause());
        }
        // load a CCP
//        Path networkConfigPath = Paths.get("/root/fabric/fabric-samples/daniu-network-prod/organizations/peerOrganizations/org1.niuinfo.com/connection-org1.json");
        Path networkConfigPath = Paths.get("/root/fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json");
        long time2 = new Date().getTime();
        System.out.println(time2);
        Gateway.Builder builder = Gateway.createBuilder();
        try {
            builder.identity(wallet, "appUser").networkConfig(networkConfigPath).discovery(true);
        } catch (Exception e) {
            throw new IdentityCreateException("身份创建出错", e.getCause());
        }
        return builder.connect();
    }

}
