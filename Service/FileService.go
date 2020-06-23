package Service

import (
	"encoding/json"
	"fmt"
	"os"
)

type ClientConfig struct {
	LocalAddress  string `json:"local_address"`
	LocalPort     string `json:"local_port"`
	LocalConnect  string `json:"local_connect"`
	RemoteAddress string `json:"remote_address"`
	RemotePort    string `json:"remote_port"`
	RemoteConnect string `json:"remote_connect"`
}

type ServerConfig struct {
	ServerListenPort     string `json:"server_listen_port"`
	ServiceListenConnect string `json:"service_listen_connect"`
	UserListenPort       string `json:"user_listen_port"`
	UserListenConnect    string `json:"user_listen_connect"`
}

//获取客户端配置
func GetClientConfig() ClientConfig {
	var clientConfigFile = "./client.json"
	var buf = readFile(clientConfigFile)
	var clientConfig = ClientConfig{}
	err := json.Unmarshal(buf, &clientConfig)
	if err != nil {
		fmt.Println("Unmarshal failed, ", err.Error())
	}
	println(string(buf))
	return clientConfig
}
func SetClientConfig(localAddr string, localPort string, remoteAddr string, remotePort string) {
	var clientConfig = ClientConfig{LocalAddress: localAddr, LocalPort: localPort, RemoteAddress: remoteAddr, RemotePort: remotePort}
	clientConfig.LocalConnect = "tcp"
	clientConfig.RemoteConnect = "tcp"
	clientConfigJson, _ := json.Marshal(clientConfig)
	writeFile("./client.json", clientConfigJson)
}
func GetServerConfig() ServerConfig {
	var serverConfigFile = "./server.json"
	var buf = readFile(serverConfigFile)
	var serverConfig = ServerConfig{}
	err := json.Unmarshal(buf, &serverConfig)
	if err != nil {
		fmt.Println("Unmarshal failed, ", err.Error())
	}
	println(string(buf))
	return serverConfig
}

func readFile(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return []byte("")
	}
	var buf = make([]byte, 65535)
	if file == nil {
		return []byte("")
	}
	n, err1 := file.Read(buf)
	if err1 != nil {
		return []byte("")
	}
	return buf[:n]
}

func writeFile(filePath string, data []byte) bool {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Failed, ", err.Error())
	}
	err3 := os.Chmod(filePath, 0777)
	if err3 != nil {
		fmt.Println("Failed, ", err3.Error())
	}
	n, err2 := file.Write(data)
	if err2 != nil {
		fmt.Println("Failed, ", err2.Error())
	}
	fmt.Println(n)
	return true
}
