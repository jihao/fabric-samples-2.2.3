package com.moshu.fabric.api.service;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeoutException;

import com.alibaba.fastjson.JSON;
import com.moshu.fabric.api.gateway.FabricGatewayFactory;
import com.moshu.fabric.api.gateway.IdentityCreateException;
import com.moshu.fabric.api.gateway.WalletCreateException;
import org.hyperledger.fabric.gateway.*;
import org.springframework.stereotype.Service;
import org.springframework.web.bind.annotation.RequestBody;


@Service
public class ChaincodeService {

    public Map<String, Object> transaction(String channel, String chaincode, String functionName,  @RequestBody Map<String, Object> dataMap)  {
        Gateway gateway = null;
        try {
            gateway = FabricGatewayFactory.getGateWay();
        } catch (WalletCreateException e) {
            return buildErrorMap(e, e.getMessage());
        } catch (IdentityCreateException e) {
            return buildErrorMap(e, e.getMessage());
        }
        Map<String, Object> resultMap;
        try {
            Network network = gateway.getNetwork(channel);
            Contract contract = network.getContract(chaincode);
            byte[] result;
            result = contract.submitTransaction(functionName,JSON.toJSONString(dataMap));
            String resultStr = new String(result);
            resultMap = JSON.parseObject(resultStr);
        } catch (ContractException e) {
            resultMap = buildErrorMap(e,"链码调用失败ContractException");
        }catch (TimeoutException e) {
            resultMap = buildErrorMap(e,"链码调用失败TimeoutException");
        }catch (InterruptedException e) {
            resultMap = buildErrorMap(e,"链码调用失败InterruptedException");
        }catch (Exception e) {
            resultMap = buildErrorMap(e,"链码调用失败");
        }

        return resultMap;
    }

    public  Map<String, Object> query(String channel, String chaincode,String pageSize, String bookmark, String functionName,Map<String, Object> dataMap){
        Gateway gateway = null;
        try {
            gateway = FabricGatewayFactory.getGateWay();
        } catch (WalletCreateException e) {
            return buildErrorMap(e, e.getMessage());
        } catch (IdentityCreateException e) {
            return buildErrorMap(e, e.getMessage());
        }
        Map<String, Object> resultMap;
        try  {
            System.out.println("channel = " + channel);
            System.out.println("chaincode = " + chaincode);
            System.out.println("dataMap = " + JSON.toJSONString(dataMap));
            System.out.println("functionName = " + functionName);
            System.out.println("-----------------------------------------------");

                Network network = gateway.getNetwork(channel);
                Contract contract = network.getContract(chaincode);

            System.out.println("contract = " + contract.toString());

                byte[] result = contract.submitTransaction(functionName, JSON.toJSONString(dataMap),pageSize,bookmark);
                resultMap = JSON.parseObject(new String(result));
                System.out.println(resultMap.keySet());
            } catch (Exception e) {
                e.printStackTrace();
                resultMap = buildErrorMap(e,"链码调用失败");
            }
            return resultMap;
    }

    private Map<String, Object> buildErrorMap(Exception e,String message) {
        Map<String, Object> errorResult = new HashMap<>();
        errorResult.put("code",500);
        errorResult.put("message",message);
        errorResult.put("stackTrace",e.getStackTrace());
        return errorResult;
    }
}