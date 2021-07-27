package com.moshu.fabric.api;

import com.moshu.fabric.api.utils.EnrollAdminUtil;
import com.moshu.fabric.api.utils.RegisterUserUtils;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {

	static {
		System.setProperty("org.hyperledger.fabric.sdk.service_discovery.as_localhost", "true");
	}
	public static void main(String[] args) {
		try {
			EnrollAdminUtil.enrollAdmin();
			RegisterUserUtils.registerUser();
		} catch (Exception ex) {
			ex.printStackTrace();
		}
		SpringApplication.run(Application.class, args);
	}

}
