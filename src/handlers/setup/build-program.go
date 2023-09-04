package setup_services

import (
	"context"
	"encoding/json"
	"log"
	"io"
	"os"
	"path/filepath"

	models "app/src/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func BuildProgram(cli *client.Client, containerId string) []string {
	ctx := context.Background()
	absPath, _ := filepath.Abs("./assets/lang-config.json")
	jsonFile, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	jsonByteArr, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var languages models.Languages

	json.Unmarshal(jsonByteArr, &languages)

	var langConfig models.Language

	for _, language := range languages.Languages {
		if language.Name == "cpp" {
			langConfig = language
			break
		}
	}

	log.Println("Starting build commands")
	for _, buildCmd := range langConfig.Build {
		log.Println("Executing", buildCmd)
		execConfig := types.ExecConfig{
			AttachStderr: true,
			Cmd:          buildCmd,
		}

		execID, err := cli.ContainerExecCreate(ctx, containerId, execConfig)
		if err != nil {
			panic(err)
		}

		execResp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
		if err != nil {
			panic(err)
		}
		defer execResp.Close()

		output := make([]byte, 1024)
		n, err := execResp.Reader.Read(output)
		if n > 0 {
			panic("Compilation Error: " + string(output[:n]))
		}
	}
	log.Println("Build finish")
	return langConfig.Run
}
