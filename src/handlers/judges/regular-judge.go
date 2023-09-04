package judge_service

import (
	enums "app/src/enums"
	models "app/src/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RegularJudge(test models.Test, cli *client.Client, containerId string, testNumber int, taskInfo models.Problem, runCmd []string, params models.Submission) models.Test {
	ctx := context.Background()

	runStr := []string{
		"su",
		"non-root",
		"-c",
		"/usr/bin/time -f \"{\\\"Memory\\\": %M, \\\"Time\\\": %e, \\\"Exit\\\": %x}\" -q -o /usr/src/app/" + strconv.Itoa(params.Task) + "/" + fmt.Sprintf(taskInfo.InputFile, testNumber) + ".info " +
			"/usr/bin/timeout " +
			strconv.Itoa(int(math.Ceil(float64(taskInfo.TimeLimit)/1000))) + "s " +
			fmt.Sprintf("%s < /usr/src/app/%d/%s 1> %s", strings.Join(runCmd, " "),
				params.Task, fmt.Sprintf(taskInfo.InputFile, testNumber),
				fmt.Sprintf("/usr/src/app/%d/%s.out",
					params.Task,
					fmt.Sprintf(taskInfo.InputFile, testNumber))),
	}

	execConfig := types.ExecConfig{
		Cmd:          runStr,
		AttachStderr: true,
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
	_, err = io.ReadAll(execResp.Reader)
	if err != nil && err != io.EOF {
		panic(err)
	}

	execConfig = types.ExecConfig{
		Cmd:          []string{"cat", "/usr/src/app/" + strconv.Itoa(params.Task) + "/" + fmt.Sprintf(taskInfo.InputFile, testNumber) + ".info"},
		AttachStdout: true,
	}

	execID, err = cli.ContainerExecCreate(ctx, containerId, execConfig)
	if err != nil {
		panic(err)
	}

	execResp, err = cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	defer execResp.Close()

	output, err := io.ReadAll(execResp.Reader)
	if err != nil && err != io.EOF {
		panic(err)
	}

	re := regexp.MustCompile(`[[:^print:]]`)
	sanitizedOuput := re.ReplaceAllString(string(output), "")

	Data := []byte(string(sanitizedOuput[1:]))

	err = json.Unmarshal(Data, &test)
	if err != nil {
		panic(err)
	}

	if test.Time*1000 >= float32(taskInfo.TimeLimit) {
		test.Verdict = enums.VERDICT_TLE.Verdict
	} else if test.Memory > taskInfo.MemoryLimit {
		test.Verdict = enums.VERDICT_MLE.Verdict
	} else if test.Exit != 0 {
		test.Verdict = enums.VERDICT_RE.Verdict
	}

	return test
}
