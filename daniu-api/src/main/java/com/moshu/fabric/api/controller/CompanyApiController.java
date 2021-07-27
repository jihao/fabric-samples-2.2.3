package com.moshu.fabric.api.controller;

import com.moshu.fabric.api.config.PathConfig;
import com.moshu.fabric.api.service.ChaincodeService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

@Controller
@RequestMapping("/")
public class CompanyApiController {

    @Autowired
    ChaincodeService chaincodeService;

    @PostMapping("{field}/{operation}")
    @ResponseBody
    public Map<String, Object> commonApi(
            @PathVariable("field")String field,
            @PathVariable("operation")String operation,
            @RequestBody Map<String, Object> dataMap ) throws Exception {
        System.out.println(dataMap);
        String channel = dataMap.get("channel").toString();
        String chaincode = dataMap.get("chaincode").toString();

        Map<String, Object> data = (Map<String, Object>)dataMap.get("content");
        System.out.println(channel+""+chaincode+""+data);
        String transactionFunctionName = PathConfig.getTransactionFunctionName(field + "/" + operation);
        System.out.println(field + "/" + operation);
        String queryFunctionName = PathConfig.getQueryFunctionName(field + "/" + operation);
        if (transactionFunctionName != null) {
            String Uuid  = UUID.randomUUID().toString().toUpperCase();
            data.put("Uuid",Uuid);
            return chaincodeService.transaction(channel, chaincode, transactionFunctionName, data);
        } else if(queryFunctionName != null){
//            String limitPageSize = "";
//            String limitOffect = "";
//            if (dataMap.containsKey("limitPageSize")) {
//                limitPageSize = dataMap.get("limitPageSize").toString();
//                if("".equals(limitPageSize)){
//                    limitPageSize="100";
//                }
//            } else {
//                limitPageSize="0";
//            }
//            if (dataMap.containsKey("limitOffect")){
//                 limitOffect = dataMap.get("limitOffect").toString();                                                                                                                                                                                                                                                                                   
//            }
            int limitPageSize = 0;
            int limitOffect = 0;
            if (dataMap.containsKey("limitPageSize")){
               limitPageSize = Integer.parseInt(dataMap.get("limitPageSize").toString());
            }
             if (dataMap.containsKey("limitOffect")){
               limitOffect = Integer.parseInt(dataMap.get("limitOffect").toString());
             }

            if (limitPageSize==0){
                 return chaincodeService.query(channel, chaincode, "0","", queryFunctionName, data);
            }else{
                if(limitOffect==0){
                    return chaincodeService.query(channel, chaincode, String.valueOf(limitPageSize), "", queryFunctionName, data);
                }
                 Map<String, Object> resultMap = chaincodeService.query(channel, chaincode,  String.valueOf(limitOffect), "", queryFunctionName, data);
                 Map<String, String> bookmarMap = (Map<String, String>)resultMap.get("responseMetadata");
                 return chaincodeService.query(channel, chaincode, String.valueOf(limitPageSize),  bookmarMap.get("bookmark"), queryFunctionName, data);
            }

        }else {
            Map errMap =  new HashMap<String,Object>();
            errMap.put("code",500);
            errMap.put("message",field + "/" + operation+"路径有误");
            return errMap;
        }
    }
}
