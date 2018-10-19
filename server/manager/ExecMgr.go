package managers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	//	"time"

	. "github.com/bigwind/goWeb/common"
)

type ExecMgr struct {
	Prog string
	Args []string
}

func (this *ExecMgr) ExecCmd() error {
	locker := GetGlobeLocker()
	locker.Lock()
	defer locker.Unlock()

	//	time.Sleep(time.Second * 10)

	cmd := exec.Command(this.Prog, this.Args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	serr, _ := ioutil.ReadAll(stderr)
	sout, _ := ioutil.ReadAll(stdout)

	fmt.Println(string(serr))
	fmt.Println("----------------------")
	fmt.Println(string(sout))
	errMsg := string(serr)
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}
