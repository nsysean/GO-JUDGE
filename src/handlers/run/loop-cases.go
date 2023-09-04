package run_service

import (
	"app/src/enums"
	"app/src/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	judge_service "app/src/handlers/judges"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/redis/go-redis/v9"
)

func LoopCases(params models.Submission, cli *client.Client, containerId string, runCmd []string, rClient redis.Client) {
	ctx := context.Background()
	taskInfo := InitModes(params.Task)
	log.Println("Successfully parsed task data")

	mode := "Regular"
	if len(taskInfo.Interactor.Source) > 0 {
		mode = "Interactive"
	}

	setPermStr := []string{
		"bash",
		"-c",
		fmt.Sprintf("chmod 111 /usr/src/app/%d/%s", params.Task, strings.Replace(fmt.Sprintf(taskInfo.OutputFile, 99999), "99999", "*", -1)),
	}

	log.Println("Successfully changed permissions of output files")
	log.Println("Starting to judge tests in", mode, "mode")

	execConfig := types.ExecConfig{
		Cmd: setPermStr,
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

	groupMap := make(map[string]models.Group)
	failedGroups := make(map[string]bool)
	failed := "PASS"
	partial := false
	time := float32(0)
	mem := 0
	scr := float32(0)

	rClient.Del(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":groups:tcsDetails")
	for _, group := range taskInfo.Groups {
		if group.PointsPolicy == "complete-group" {
			group.PointStore = float32(100)
		}
		groupMap[group.Name] = group
		failedGroups[group.Name] = false
	}

	if mode == "Interactive" {
		cmd := exec.Command("docker", "exec", containerId, "mkfifo", "p1", "p2")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	for i := range taskInfo.Tests {
		if taskInfo.Tests[i].Points != 0 {
			partial = true
		}
	}

	for i := range taskInfo.Tests {
		taskInfo.Tests[i].Verdict = "Skipped"
		if failedGroups[groupMap[taskInfo.Tests[i].Group].Name] == true && groupMap[taskInfo.Tests[i].Group].PointsPolicy == "complete-group" {
			continue
		}
		for _, depedency := range groupMap[taskInfo.Tests[i].Group].Dependencies {
			if failedGroups[depedency.Group] {
				failedGroups[groupMap[taskInfo.Tests[i].Group].Name] = true
				break
			}
		}
		if failedGroups[groupMap[taskInfo.Tests[i].Group].Name] == true && groupMap[taskInfo.Tests[i].Group].PointsPolicy == "complete-group" {
			continue
		}
		taskInfo.Tests[i].Verdict = "Running"
		rClient.Set(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":onTest", i+1, 0)

		if mode == "Interactive" {
			taskInfo.Tests[i] = judge_service.InteractiveJudge(taskInfo.Tests[i], cli, containerId, i+1, taskInfo, runCmd, params)
		} else {
			taskInfo.Tests[i] = judge_service.RegularJudge(taskInfo.Tests[i], cli, containerId, i+1, taskInfo, runCmd, params)
		}

		time = max(time, taskInfo.Tests[i].Time)
		mem = max(mem, taskInfo.Tests[i].Memory)

		failedCase := false

		if taskInfo.Tests[i].Verdict == "Running" {
			taskInfo.Tests[i] = judge_service.Scoring(taskInfo.Tests[i], cli, containerId, i+1, taskInfo, params)
			data, err := json.Marshal(taskInfo.Tests[i])
			if err != nil {
				panic(err)
			}

			rClient.LPush(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":tcsDetails", string(data))
		} else {
			failedCase = true
			data, err := json.Marshal(taskInfo.Tests[i])
			if err != nil {
				panic(err)
			}

			rClient.LPush(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":tcsDetails", string(data))
			if !partial {
				failed = taskInfo.Tests[i].Verdict
				break
			}
		}

		if taskInfo.Tests[i].Verdict != enums.VERDICT_ACCEPTED.Verdict && taskInfo.Tests[i].Verdict != enums.VERDICT_PARTIAL.Verdict {
			failedCase = true
			failed = taskInfo.Tests[i].Verdict

			if !partial {
				break
			}
		}

		if taskInfo.Tests[i].Points != 0 {
			partial = true
		}

		if failedCase {
			failedGroups[groupMap[taskInfo.Tests[i].Group].Name] = true
		}

		if groupMap[taskInfo.Tests[i].Group].PointsPolicy == "complete-group" && groupMap[taskInfo.Tests[i].Group].PointStore > taskInfo.Tests[i].ScoredPoints {
			curGroup := groupMap[taskInfo.Tests[i].Group]
			curGroup.PointStore = taskInfo.Tests[i].ScoredPoints
			groupMap[taskInfo.Tests[i].Group] = curGroup
		} else if groupMap[taskInfo.Tests[i].Group].PointsPolicy == "each-test" && taskInfo.Tests[i].Points != 0 {
			curGroup := groupMap[taskInfo.Tests[i].Group]
			curGroup.PointStore += taskInfo.Tests[i].ScoredPoints / taskInfo.Tests[i].Points
			groupMap[taskInfo.Tests[i].Group] = curGroup
		} else if groupMap[taskInfo.Tests[i].Group].PointsPolicy == "" && taskInfo.Tests[i].Points != 0 {
			scr += (taskInfo.Tests[i].ScoredPoints / 100) * taskInfo.Tests[i].Points
			log.Println(scr, i + 1)
		}
	}

	if len(taskInfo.Groups) > 0 {

		scored := float32(0)
		for _, group := range groupMap {
			curScore := float32(0)
			if group.PointsPolicy == "each-test" {
				curScore = group.PointStore
			} else {
				curScore = group.PointStore / 100 * group.Points
			}
			scored += curScore
			data, err := json.Marshal(fmt.Sprintf("{\"Group\": %s, \"Score\": %f}", group.Name, curScore))
			if err != nil {
				panic(err)
			}
			groupScore, err := strconv.Unquote(string(data))
			rClient.LPush(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":groupScore", groupScore, 0)
		}
		rClient.Set(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":result", fmt.Sprintf("Partial Score %.2f", scored), 0)
	} else {
		if partial {
			rClient.Set(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":result", fmt.Sprintf("Partial Score %f", scr), 0)
		} else if failed != "PASS" {
			rClient.Set(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":result", failed, 0)
		} else {
			rClient.Set(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":result", enums.VERDICT_ACCEPTED.Verdict, 0)
		}
	}
	rClient.Set(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":time", time, 0)
	rClient.Set(ctx, "submissions:"+strconv.Itoa(params.SubmissionId)+":mem", mem, 0)
}
