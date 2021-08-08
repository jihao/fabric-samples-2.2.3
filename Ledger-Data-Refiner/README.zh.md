# Ledger-Data-Refiner 安装步骤

1. 使用第一台服务器 10.18.188.177  
2. 上传 Ledger-Data-Refiner 至 /root/fabric/fabric-samples/Ledger-Data-Refiner  
3. make docker_all  
4. docker-compose -f docker-compose.yaml up -d  
5. 访问 http://47.98.138.84:30052 (177的外网IP)  


## 相关配置

    docker-compose.yaml
    config/config.ini
    config/connection-config-docker.yaml




