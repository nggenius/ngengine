package logic

import (
	"ngengine/console/models"
	"time"
)

type ServerData struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	ServerId   int       `json:"serverid"`
	ServerIp   string    `json:"serverip"`
	LoginTime  time.Time `json:"logintime"`
	DeleteTime time.Time `json:"deletetime"`
}

type ServerDataList struct {
	ServerNum  int          `json:"servernum"`
	ServerData []ServerData `json:"serverdata"`
}

func GetServerList() (*ServerDataList, error) {
	// 从数据库中拉取所有服务器列表
	dataList := new(ServerDataList)

	var allData []models.NxConsole
	data := &models.NxConsole{}
	allData, err := data.ReadAll()
	if err != nil {
		return nil, err
	}

	var allNum = len(allData)
	dataList.ServerNum = allNum
	dataList.ServerData = make([]ServerData, 0)

	for i := 0; i < allNum; i++ {
		var serverData ServerData
		serverData.Id = allData[i].Id
		serverData.Name = allData[i].Name
		serverData.ServerId = allData[i].ServerId
		serverData.ServerIp = allData[i].ServerIp
		serverData.LoginTime = allData[i].LoginTime
		serverData.DeleteTime = allData[i].DeleteTime
		dataList.ServerData = append(dataList.ServerData, serverData)
	}
	return dataList, nil
}
