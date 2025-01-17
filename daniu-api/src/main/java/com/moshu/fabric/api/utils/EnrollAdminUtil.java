package com.moshu.fabric.api.utils;

import org.hyperledger.fabric.gateway.Wallet;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric.sdk.security.CryptoSuiteFactory;
import org.hyperledger.fabric_ca.sdk.EnrollmentRequest;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.nio.file.Paths;
import java.util.Properties;


public class EnrollAdminUtil {

    static {
        System.setProperty("org.hyperledger.fabric.sdk.service_discovery.as_localhost", "false"); //需要本地修改host
    }
    public static void enrollAdmin() throws Exception{
        // Create a CA client for interacting with the CA.
        Properties props = new Properties();
        //props.put("pemFile", "/data/fabric/fabric/fabric-samples/first-network/crypto-config/peerOrganizations/org1.niuinfo.com/ca/ca.org1.niuinfo.com-cert.pem");
//        props.put("pemFile", "/root/fabric/fabric-samples/first-network/crypto-config/peerOrganizations/org1.niuinfo.com/ca/ca.org1.niuinfo.com-cert.pem");
        props.put("pemFile", "/root/fabric/fabric-samples/daniu-network-1/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem");
        props.put("allowAllHostNames", "true");
        HFCAClient caClient = HFCAClient.createNewInstance("https://localhost:7054", props);
        CryptoSuite cryptoSuite = CryptoSuiteFactory.getDefault().getCryptoSuite();
        caClient.setCryptoSuite(cryptoSuite);

        // Create a wallet for managing identities
        Wallet wallet = Wallet.createFileSystemWallet(Paths.get("wallet"));

        // Check to see if we've already enrolled the admin user.
        boolean adminExists = wallet.exists("admin");
        if (adminExists) {
            System.out.println("An identity for the admin user \"admin\" already exists in the wallet");
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        final EnrollmentRequest enrollmentRequestTLS = new EnrollmentRequest();
        enrollmentRequestTLS.addHost("localhost");
        enrollmentRequestTLS.setProfile("tls");
        Enrollment enrollment = caClient.enroll("admin", "adminpw", enrollmentRequestTLS);
        Wallet.Identity user = Wallet.Identity.createIdentity("Org1MSP", enrollment.getCert(), enrollment.getKey());
        wallet.put("admin", user);
        System.out.println("Successfully enrolled user \"admin\" and imported it into the wallet");

    }

}
