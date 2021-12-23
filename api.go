package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loginApi(c *gin.Context) {
	user := c.PostForm("user")
	pass := c.PostForm("pass")

	adminList := adminConfig["admin"].(map[string]interface{})
	for i := range adminList {
		if i == user {
			if adminList[i] == pass {
				session := sessions.Default(c)
				session.Set("login", "ok")
				session.Save()
				c.String(200, `{"status":"ok"}`)
				return
			}
		}
	}
	c.String(200, `{"status":"uperror"}`)
}

func postNodeApi(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		nodeIp := c.PostForm("ip")
		apiName := c.PostForm("api")

		var params url.Values

		if apiName == "nodeInfo" || apiName == "updateVhost" {
			params = url.Values{"token": {nodeConfig.Token}}
		} else if apiName == "setNodeStatus" {
			status := c.PostForm("status")
			params = url.Values{"token": {nodeConfig.Token}, "status": {status}}
		}

		resp, err := http.PostForm("http://"+nodeIp+":8081/api/"+apiName, params)
		if err != nil {
			c.String(200, `{"status": "PostError"}`)
			fmt.Println("PostError", err)
			return
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		c.String(200, string(body))
	} else {
		c.String(200, `{"status":"login"}`)
	}
}

func nodeList(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		upAdminConfig() // Update

		nodeList := adminConfig["node"].([]interface{})
		var nodeJson []map[string]string
		for i := range nodeList {
			params := url.Values{"token": {nodeConfig.Token}}
			client := http.Client{
				Timeout: 2 * time.Second,
			}
			thisNode := nodeList[i].(map[string]interface{})
			resp, err := client.PostForm("http://"+fmt.Sprint(thisNode["ip"])+":8081/api/nodeInfo", params)
			if err != nil {
				nodeJson = append(nodeJson, map[string]string{"name": fmt.Sprint(thisNode["name"]), "ip": fmt.Sprint(thisNode["ip"]), "status": "error"})
				fmt.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var bodyJson map[string]interface{}
			json.Unmarshal(body, &bodyJson)
			nodeJson = append(nodeJson, map[string]string{"name": fmt.Sprint(thisNode["name"]), "ip": fmt.Sprint(thisNode["ip"]), "status": fmt.Sprint(bodyJson["status"]), "system": fmt.Sprint(bodyJson["system"])})
		}
		c.JSON(200, nodeJson)
	} else {
		c.String(200, `{"status":"login"}`)
	}
}

func updateNodeStatus(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		nodeList := c.PostForm("nodelist")
		var nodeListJson []string
		json.Unmarshal([]byte(nodeList), &nodeListJson)

		var nodeJson []string
		for i := range nodeListJson {
			params := url.Values{"token": {nodeConfig.Token}}
			client := http.Client{
				Timeout: 2 * time.Second,
			}
			resp, err := client.PostForm("http://"+nodeListJson[i]+":8081/api/nodeInfo", params)
			if err != nil {
				nodeJson = append(nodeJson, "error")
				fmt.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var bodyJson map[string]interface{}
			json.Unmarshal(body, &bodyJson)
			nodeJson = append(nodeJson, fmt.Sprint(bodyJson["status"]))
		}
		c.JSON(200, nodeJson)
	} else {
		c.String(200, `{"status":"login"}`)
	}
}

func setAllStatus(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		status := c.PostForm("status")

		nodeList := adminConfig["node"].([]interface{})
		if status != "update" {
			for i := range nodeList {
				params := url.Values{"token": {nodeConfig.Token}, "status": {status}}
				client := http.Client{
					Timeout: 1 * time.Second,
				}
				_, err := client.PostForm("http://"+fmt.Sprint(nodeList[i].(map[string]interface{})["ip"])+":8081/api/setNodeStatus", params)
				if err != nil {
					fmt.Println("PostError", err)
				}
			}
		} else {
			Zip("tls", "tls.zip")
			for i := range nodeList {
				params := url.Values{"token": {nodeConfig.Token}}
				client := http.Client{
					Timeout: 3 * time.Second,
				}
				_, err := client.PostForm("http://"+fmt.Sprint(nodeList[i].(map[string]interface{})["ip"])+":8081/api/updateVhost", params)
				if err != nil {
					fmt.Println("PostError", err)
				}
			}
		}

		c.String(200, `{"status":"`+status+`Ok"}`)
	} else {
		c.String(200, `{"status":"login"}`)
	}
}

func delNode(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		upAdminConfig() // Update

		ip := c.PostForm("ip")
		var newNodeList []interface{}
		nodeList := adminConfig["node"].([]interface{})
		delStatus := false
		for i := range nodeList {
			if nodeList[i].(map[string]interface{})["ip"] == ip {
				newNodeList = append(nodeList[:i], nodeList[i+1:]...)
				delStatus = true
				break
			}
		}
		if delStatus {
			adminConfig["node"] = newNodeList
			jsonStr, _ := json.MarshalIndent(adminConfig, "", "    ")
			WriteToFile("config.json", string(jsonStr))
			c.JSON(200, map[string]string{"status": "ok"})
		} else {
			c.JSON(200, map[string]string{"status": "no"})
		}
	} else {
		c.String(200, `{"status":"login"}`)
	}
}

func addNode(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		upAdminConfig() // Update

		ip := c.PostForm("ip")
		name := c.PostForm("name")
		nodeList := adminConfig["node"].([]interface{})
		for i := range nodeList {
			if nodeList[i].(map[string]interface{})["ip"] == ip {
				c.String(200, `{"status":"repeat"}`)
				return
			}
		}
		nodeList = append(nodeList, map[string]interface{}{"ip": ip, "name": name})
		adminConfig["node"] = nodeList
		jsonStr, _ := json.MarshalIndent(adminConfig, "", "    ")
		WriteToFile("config.json", string(jsonStr))
		c.JSON(200, map[string]string{"status": "ok"})
	} else {
		c.String(200, `{"status":"login"}`)
	}
}

func addWebSite(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		upVhostConfig() // Update

		name := c.PostForm("name")
		host := c.PostForm("host")
		source_host := c.PostForm("source_host")
		source := c.PostForm("source")
		proto := c.PostForm("proto")
		text := c.PostForm("text")

		pen := c.PostForm("pen")
		key := c.PostForm("key")

		if name != "" && host != "" && source_host != "" && source != "" && proto != "" {
			var thisConfig map[string]interface{}
			if pen != "" && key != "" {
				thisConfig = map[string]interface{}{
					"name":        name,
					"host":        host,
					"source_host": source_host,
					"source":      source,
					"proto":       proto,
					"text":        text,
					"tls":         true,
					"pen":         "/node/tls/" + host + ".pen",
					"key":         "/node/tls/" + host + ".key",
				}
				WriteTls(host, pen, key)
			} else {
				thisConfig = map[string]interface{}{
					"name":        name,
					"host":        host,
					"source_host": source_host,
					"source":      source,
					"proto":       proto,
					"text":        text,
					"tls":         false,
				}
			}
			vhostListConfig = append(vhostListConfig, thisConfig)
			jsonStr, _ := json.MarshalIndent(vhostListConfig, "", "    ")
			WriteToFile("vhost.json", string(jsonStr))
			c.JSON(200, map[string]string{"status": "ok"})
		} else {
			c.String(200, `{"status":"no"}`)
		}
	} else {
		c.String(200, `{"status":"login"}`)
	}
}

func websiteList(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		jsonStr, _ := json.Marshal(vhostListConfig)
		c.String(200, string(jsonStr))
	} else {
		c.String(200, `{"status":"error"}`)
	}
}

func delWebSite(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")

	session := sessions.Default(c)
	if session.Get("login") == "ok" {
		upVhostConfig() // Update

		host := c.PostForm("host")
		for i := range vhostListConfig {
			if host == vhostListConfig[i]["host"] {
				vhostListConfig = append(vhostListConfig[:i], vhostListConfig[i+1:]...)
				break
			}
		}
		jsonStr, _ := json.MarshalIndent(vhostListConfig, "", "    ")
		WriteToFile("vhost.json", string(jsonStr))
		c.String(200, `{"status":"ok"}`)
	} else {
		c.String(200, `{"status":"error"}`)
	}
}
