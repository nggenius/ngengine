// +build windows

package ngadmin

import (
	"encoding/json"
	"github.com/nggenius/ngengine/logger"
	"os/exec"
)

// Start 启动进程
func start(startPath string, startPara *ServiceLink, l *logger.Log) error {
	b, err := json.Marshal(startPara.CoreOption)
	if err != nil {
		l.LogErr(err)
		return err
	}

	cmd := exec.Command(startPath, "-p", string(b))

	err = cmd.Start()
	if err != nil {
		l.LogErr(err)
		return err
	}

	l.LogInfo("master", " app start ", startPara.ServType, ",", startPara.ServId)

	return nil
}
