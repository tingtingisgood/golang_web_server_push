#!/bin/bash 
GOLANG_WEB_SERVER_PUSH_VERSION="0.1.0" 
wget https://github.com/yangtinngting/golang_web_server_push/releases/download/v${GOLANG_WEB_SERVER_PUSH_VERSION}/golangWebServerPush-${GOLANG_WEB_SERVER_PUSH_VERSION}.linux.amd64.tar.gz 
tar -xzvf golangWebServerPush-${GOLANG_WEB_SERVER_PUSH_VERSION}.linux.amd64.tar.gz 
cd golang_web_server_push/ 

# if you just want to start prometheus as root 
# ./prometheus --config.file=prometheus.yml  

# create user 
useradd --no-create-home --shell /bin/false golang_web_server_push 

# copy binaries   
cp golangWebServerPush /usr/local/bin 
chown golang_web_server_push:golang_web_server_push  /usr/local/bin/golangWebServerPush 

# setup systemd 
echo '[Unit]
Description=Golang Web Server Push
Wants=network-online.target
After=network-online.target

[Service]
User=golang_web_server_push
Group=golang_web_server_push
Type=simple
ExecStart=/usr/local/bin/golangWebServerPush 

[Install]
WantedBy=multi-user.target' > /etc/systemd/system/golang_web_server_push.service

systemctl daemon-reload
systemctl enable golang_web_server_push
systemctl start golang_web_server_push




