package com.moshu.fabric.api.utils;

import org.hyperledger.fabric.gateway.Wallet;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric.sdk.security.CryptoSuiteFactory;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.hyperledger.fabric_ca.sdk.RegistrationRequest;

import java.nio.file.Paths;
import java.security.PrivateKey;
import java.util.Properties;
import java.util.Set;


public class RegisterUserUtils {
    static {
        System.setProperty("org.hyperledger.fabric.sdk.service_discovery.as_localhost", "true");
    }
    public static void registerUser()throws Exception{
        // Create a CA client for interacting with the CA.
        Properties props = new Properties();
        //props.put("pemFile", "/root/fabric/fabric-samples/first-network/crypto-config/peerOrganizations/org1.niuinfo.com/ca/ca.org1.niuinfo.com-cert.pem");
        props.put("pemFile", "/root/fabric/fabric-samples/daniu-network-1/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem");
        props.put("allowAllHostNames", "true");
        HFCAClient caClient = HFCAClient.createNewInstance("https://localhost:7054", props);
        CryptoSuite cryptoSuite = CryptoSuiteFactory.getDefault().getCryptoSuite();
        caClient.setCryptoSuite(cryptoSuite);

        // Create a wallet for managing identities
        Wallet wallet = Wallet.createFileSystemWallet(Paths.get("wallet"));

        // Check to see if we've already enrolled the user.
        boolean userExists = wallet.exists("appUser");
        if (userExists) {
            System.out.println("An identity for the user \"appUser\" already exists in the wallet");
            return;
        }

        userExists = wallet.exists("admin");
        if (!userExists) {
            System.out.println("\"admin\" needs to be enrolled and added to the wallet first");
            return;
        }

        Wallet.Identity adminIdentity = wallet.get("admin");
        User admin = new User() {

            @Override
            public String getName() {
                return "admin";
            }

            @Override
            public Set<String> getRoles() {
                return null;
            }

            @Override
            public String getAccount() {
                return null;
            }

            @Override
            public String getAffiliation() {
                return "org1.department1";
            }

            @Override
            public Enrollment getEnrollment() {
                return new Enrollment() {

                    @Override
                    public PrivateKey getKey() {
                        return adminIdentity.getPrivateKey();
                    }

                    @Override
                    public String getCert() {
                        return adminIdentity.getCertificate();
                    }
                };
            }

            @Override
            public String getMspId() {
                return "Org1MSP";
            }

        };

        // Register the user, enroll the user, and import the new identity into the wallet.
        RegistrationRequest registrationRequest = new RegistrationRequest("appUser");
        registrationRequest.setAffiliation("org1.department1");
        registrationRequest.setEnrollmentID("appUser");

        if (!userExists) {
            caClient.revoke(admin, "admin", "reason");
        }

        String enrollmentSecret = caClient.register(registrationRequest, admin);
        Enrollment enrollment = caClient.enroll("appUser", enrollmentSecret);
        Wallet.Identity user = Wallet.Identity.createIdentity("Org1MSP", enrollment.getCert(), enrollment.getKey());
        wallet.put("appUser", user);
        System.out.println("Successfully enrolled user \"appUser\" and imported it into the wallet");
    }

}
