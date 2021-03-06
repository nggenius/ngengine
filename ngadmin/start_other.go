// +build !windows

package master

import (
	"encoding/json"
	"github.com/nggenius/ngengine/logger"
	"os/exec"
	"server/libs/log"
	"syscall"
)

// start 启动进程
func start(startPath string, startPara *ServiceLink, l *logger.Log) error {
	b, err := json.Marshal(startPara.CoreOption)
	if err != nil {
		l.LogErr(err)
		return err
	}

	cmd := exec.Command(startPath, "-p", string(b))

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	} else {
		cmd.SysProcAttr.Setpgid = true
	}

	err = cmd.Start()
	if err != nil {
		log.LogErr(err)
		return err
	}

	l.LogInfo("master", " app start ", startPara.ServType, ",", startPara.ServId)

	return nil
}
