package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getHttpToken(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		c.String(200, nodeConfig.Token)
	} else {
		c.String(200, `{"status": "Login"}`)
	}
}

func getInstall(c *gin.Context) {
	token := c.Query("token")
	if token == nodeConfig.Token {
		c.Writer.Header().Set("Content-Type", "application/x-sh;charset=utf-8")
		installFile, err := ioutil.ReadFile("static/install.sh")
		if err != nil {
			fmt.Println("Read Error", err)
		}
		installFileString := strings.ReplaceAll(strings.ReplaceAll(string(installFile), "{ctrlServer}", "http://"+c.Request.Host), "{token}", nodeConfig.Token)
		c.String(200, installFileString)
	} else {
		c.String(200, ``)
	}
}

func getConfig(c *gin.Context) {
	token := c.Query("token")
	if token == nodeConfig.Token {
		c.Writer.Header().Set("Content-Type", "application/yaml;charset=utf-8")
		installFile, err := ioutil.ReadFile("static/config.yaml")
		if err != nil {
			fmt.Println("Read Error", err)
		}
		installFileString := strings.ReplaceAll(string(installFile), "{ctrlServer}", "http://"+c.Request.Host)
		c.String(200, installFileString)
	} else {
		c.String(200, ``)
	}
}

func getTls(c *gin.Context) {
	token := c.Query("token")
	if token == nodeConfig.Token {
		c.Writer.Header().Set("Content-Type", "application/zip;charset=utf-8")
		c.File("./tls.zip")
	} else {
		c.String(200, ``)
	}
}

type vhostConfig struct {
	Host       string `json:"host"`
	SourceHost string `json:"source_host"`
	Source     string `json:"source"`
	Proto      string `json:"proto"`
	Tls        bool   `json:"tls"`
	Pen        string `json:"pen"`
	Key        string `json:"key"`
}

func getVhost(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain;charset=utf-8")
	session := sessions.Default(c)
	token := c.Query("token")
	if token == nodeConfig.Token || session.Get("login") == "ok" {
		var vhostList []vhostConfig
		vhostJsonStr, err := ioutil.ReadFile("vhost.json")
		if err != nil {
			fmt.Println("Read Error", err)
		}
		err = json.Unmarshal(vhostJsonStr, &vhostList)
		if err != nil {
			fmt.Println("Json Unmarshal", err)
		}

		var configText string

		for i := range vhostList {
			if !vhostList[i].Tls {
				configText += `server {
	listen 80;
	listen [::]:80;
	server_name ` + vhostList[i].Host + `;

	error_page 404 /404.html;
	error_page 500 /500.html;
	error_page 502 /502.html;
	error_page 503 /503.html;
	error_page 504 /504.html;

	sub_filter "{when}" $date_local;
	sub_filter "{method}" $request_method;
	sub_filter "{hostname}" $hostname;
	sub_filter "{remote}" $remote_addr;

	location /404.html {
		root /node/error;
	}
	location /500.html {
		root /node/error;
	}
	location /502.html {
		root /node/error;
	}
	location /503.html {
		root /node/error;
	}
	location /504.html {
		root /node/error;
	}

	location  ~* \.(webp|gif|png|jpg|css|js|woff|woff2)$
	{
		proxy_pass ` + vhostList[i].Proto + `://` + vhostList[i].Source + `;
		proxy_set_header Host ` + vhostList[i].SourceHost + `;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header REMOTE-HOST $remote_addr;

		proxy_pass_header Server;
		add_header server CluckCDN;

		proxy_cache cache_one;
		proxy_cache_key $host$uri$is_args$args;
		proxy_cache_valid 200 304 301 302 5m;

		expires 12h;
	}

	location / 
	{
		proxy_pass ` + vhostList[i].Proto + `://` + vhostList[i].Source + `;
		proxy_set_header Host ` + vhostList[i].SourceHost + `;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header REMOTE-HOST $remote_addr;
		proxy_ssl_server_name on;

		proxy_pass_header Server;
		add_header server CluckCDN;
		
		add_header X-Cache $upstream_cache_status;

		proxy_ignore_headers Set-Cookie Cache-Control expires;
	}
}
`
			} else {
				configText += `server {
	listen 80;
	listen [::]:80;
	listen 443 ssl http2;
	listen [::]:443 ssl http2;
	server_name ` + vhostList[i].Host + `;

	error_page 404 /404.html;
	error_page 500 /500.html;
	error_page 502 /502.html;
	error_page 503 /503.html;
	error_page 504 /504.html;

	sub_filter "{when}" $date_local;
	sub_filter "{method}" $request_method;
	sub_filter "{hostname}" $hostname;
	sub_filter "{remote}" $remote_addr;

	location /404.html {
		root /node/error;
	}
	location /500.html {
		root /node/error;
	}
	location /502.html {
		root /node/error;
	}
	location /503.html {
		root /node/error;
	}
	location /504.html {
		root /node/error;
	}

	if ($server_port !~ 443){
		rewrite ^(/.*)$ https://$host$1 permanent;
	}

	ssl_certificate ` + vhostList[i].Pen + `;
	ssl_certificate_key ` + vhostList[i].Key + `;
	ssl_protocols TLSv1.2 TLSv1.3;
	ssl_ciphers EECDH+CHACHA20:EECDH+CHACHA20-draft:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;
	ssl_prefer_server_ciphers on;
	ssl_session_cache shared:SSL:10m;
	ssl_session_timeout 10m;
	add_header Strict-Transport-Security "max-age=31536000";
	error_page 497 https://$host$request_uri;

	location  ~* \.(webp|gif|png|jpg|css|js|woff|woff2)$
	{
		proxy_pass ` + vhostList[i].Proto + `://` + vhostList[i].Source + `;
		proxy_set_header Host ` + vhostList[i].SourceHost + `;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header REMOTE-HOST $remote_addr;

		proxy_pass_header Server;
		add_header server CluckCDN;

		proxy_cache cache_one;
		proxy_cache_key $host$uri$is_args$args;
		proxy_cache_valid 200 304 301 302 5m;

		expires 12h;
	}

	location / 
	{
		proxy_pass ` + vhostList[i].Proto + `://` + vhostList[i].Source + `;
		proxy_set_header Host ` + vhostList[i].SourceHost + `;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header REMOTE-HOST $remote_addr;
		proxy_ssl_server_name on;

		proxy_pass_header Server;
		add_header Server CluckCDN;
		
		add_header X-Cache $upstream_cache_status;

		proxy_ignore_headers Set-Cookie Cache-Control expires;
	}
}
`
			}
		}
		c.String(200, configText)
	} else {
		c.String(200, ``)
	}
}
