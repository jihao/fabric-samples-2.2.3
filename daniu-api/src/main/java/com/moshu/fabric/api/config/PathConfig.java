package com.moshu.fabric.api.config;

import java.util.HashMap;
import java.util.Map;

public class PathConfig {

   private static Map<String, String> modifyPathMap = null;
   private static Map<String, String> queryPathMap = null;
   public static String getTransactionFunctionName(String path) {
       if (modifyPathMap == null) {
           modifyPathMap = new HashMap<>();
           modifyPathMap.put("CompanyInfo/Add", "AddCompany");
           modifyPathMap.put("CompanyInfo/Modify", "ModifyCompany");
           modifyPathMap.put("User/Add", "AddUser");
           modifyPathMap.put("User/Modify", "ModifyUser");
           modifyPathMap.put("Driver/Add", "AddDriver");
           modifyPathMap.put("Driver/Modify", "ModifyDriver");
           modifyPathMap.put("Vehicle/Add", "AddVehicle");
           modifyPathMap.put("Vehicle/Modify", "ModifyVehicle");
           modifyPathMap.put("Contract/Add", "AddContract");
           modifyPathMap.put("Order/Add", "AddOrder");
           modifyPathMap.put("Invoice/Add", "AddInvoice");
           modifyPathMap.put("Waybill/Add", "AddWaybill");
           modifyPathMap.put("CapitalFlow/Add", "AddCapitalFlow");
           modifyPathMap.put("Complaint/Add", "AddComplaint");
           modifyPathMap.put("Comment/Add", "AddComment");

       }

       if (modifyPathMap.containsKey(path)) {
           return modifyPathMap.get(path);
       } else{
           return null;
       }
   }
    public static String getQueryFunctionName(String path) {

        if (queryPathMap == null) {
            queryPathMap = new HashMap<>();
            queryPathMap.put("CompanyInfo/Query", "QueryCompany");
            queryPathMap.put("User/Query", "QueryUser");
            queryPathMap.put("Driver/Query", "QueryDriver");
            queryPathMap.put("Vehicle/Query", "QueryVehicle");
            queryPathMap.put("Order/Query", "QueryOrder");
            queryPathMap.put("Contract/Query", "QueryContract");
            queryPathMap.put("Invoice/Query", "QueryInvoice");
            queryPathMap.put("Waybill/Query", "QueryWaybill");
            queryPathMap.put("CapitalFlow/Query", "QueryCapitalFlow");
            queryPathMap.put("Complaint/Query", "QueryComplaint");
            queryPathMap.put("Comment/Query", "QueryComment");
        }
        if (queryPathMap.containsKey(path)) {
            return queryPathMap.get(path);
        } else{
            return null;
        }
    }

}
