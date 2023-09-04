package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	run_service "app/src/handlers/run"
	setup_services "app/src/handlers/setup"
	models "app/src/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	submission, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	lang := os.Args[2]
	taskId, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	unparsedCode, err := rClient.Get(ctx, "submissions:"+strconv.Itoa(submission)+":code").Result()
	if err != nil {
		panic(err)
	}

	code, err := strconv.Unquote(unparsedCode)
	if err != nil {
		panic(err)
	}

	body := models.Submission{SubmissionId: submission, Code: code, Language: lang, Task: taskId}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	rClient.Set(ctx, "submissions:" + strconv.Itoa(submission) + ":result", "Compiling", 0)
	containerId := setup_services.StartContainer(body, cli)
	log.Println("Container ID:", containerId)

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic Occurred:", err)
			if strings.HasPrefix(fmt.Sprintf("%s", err), "Compilation Error: ") {
				rClient.Set(ctx, "submissions:"+strconv.Itoa(submission)+":result", "Compilation Error", 0)
				rClient.Set(ctx, "submissions:"+strconv.Itoa(submission)+":ce", err, 0)
			}
			cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{Force: true})
			log.Println("Container " + containerId + " forcefully removed")
		}
	}()

	setup_services.TransferFiles(body, cli, containerId)
	runCmd := setup_services.BuildProgram(cli, containerId)

	rClient.Set(ctx, "submissions:"+strconv.Itoa(body.SubmissionId)+":onTest", 1, 0)

	rClient.Set(ctx, "submissions:" + strconv.Itoa(submission) + ":result", "Running", 0)
	run_service.LoopCases(body, cli, containerId, runCmd, *rClient)
	panic("Judging successful")
}
