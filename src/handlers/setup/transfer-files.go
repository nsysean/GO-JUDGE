package setup_services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	models "app/src/models"

	"github.com/docker/docker/client"
)

func TransferFiles(params models.Submission, cli *client.Client, containerId string) {
	taskAbsPath, _ := filepath.Abs("../../tasks/" + strconv.Itoa(params.Task))
	taskTransferStr := strings.Split(fmt.Sprintf("docker cp %s %s:/usr/src/app",
		taskAbsPath,
		containerId), " ")

	cmd := exec.Command(taskTransferStr[0], taskTransferStr[1:]...)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	log.Println("Task files successfully transfered, now creating code file")

	f, err := os.CreateTemp("./tmp/", "tmp-")
	if err != nil {
		panic(err)
	}

	defer os.Remove(f.Name())

	f.WriteString(params.Code)

	fileCreateStr := strings.Split(fmt.Sprintf("docker cp %s %s:/usr/src/app/task.cpp", f.Name(), containerId), " ")

	cmd = exec.Command(fileCreateStr[0], fileCreateStr[1:]...)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	log.Println("Code file successfully created")
}
