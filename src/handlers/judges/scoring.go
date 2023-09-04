package judge_service

import (
	enums "app/src/enums"
	models "app/src/models"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Scoring(test models.Test, cli *client.Client, containerId string, testNumber int, taskInfo models.Problem, params models.Submission) models.Test {
	ctx := context.Background()
	checkerStr := []string{
		"/usr/src/app/" + strconv.Itoa(params.Task) + "/checker",
		"/usr/src/app/" + strconv.Itoa(params.Task) + "/" + fmt.Sprintf(taskInfo.InputFile, testNumber),
		"/usr/src/app/" + strconv.Itoa(params.Task) + "/" + fmt.Sprintf(taskInfo.InputFile, testNumber) + ".out",
		"/usr/src/app/" + strconv.Itoa(params.Task) + "/" + fmt.Sprintf(taskInfo.OutputFile, testNumber),
	}

	execConfig := types.ExecConfig{
		Cmd:          checkerStr,
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

	output, err := io.ReadAll(execResp.Reader)
	if err != nil && err != io.EOF {
		panic(err)
	}

	re := regexp.MustCompile(`[[:^print:]]`)
	sanitizedOuput := re.ReplaceAllString(string(output), "")

	if len(sanitizedOuput) > 0 && !unicode.IsLetter(rune(sanitizedOuput[0])) {
		sanitizedOuput = sanitizedOuput[1:]
	}

	if strings.HasPrefix(sanitizedOuput, enums.VERDICT_ACCEPTED.Prefix) {
		test.Verdict = enums.VERDICT_ACCEPTED.Verdict
		test.ScoredPoints = 100
	} else if strings.HasPrefix(sanitizedOuput, enums.VERDICT_PARTIAL.Prefix) {
		test.Verdict = enums.VERDICT_PARTIAL.Verdict
		tmpPoints, err := strconv.ParseFloat(strings.TrimSpace(strings.SplitAfter(sanitizedOuput, " ")[1]), 32)
		if err != nil {
			panic(err)
		}
		test.ScoredPoints = float32(tmpPoints)
	} else if strings.HasPrefix(sanitizedOuput, enums.VERDICT_WA.Prefix) {
		test.Verdict = enums.VERDICT_WA.Verdict
	} else if strings.HasPrefix(sanitizedOuput, enums.VERDICT_WOF.Prefix) {
		test.Verdict = enums.VERDICT_WOF.Verdict
	} else if strings.HasPrefix(sanitizedOuput, enums.VERDICT_UEOF.Prefix) {
		test.Verdict = enums.VERDICT_UEOF.Verdict
	}

	return test
}
