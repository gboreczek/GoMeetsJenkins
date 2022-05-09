package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/bndr/gojenkins"
)

type jobInfo struct {
	JobName     string
	BuildNumber int64
	StartDate   time.Time
	StartedBy   interface{}
	BuildResult string
}

func main() {

	fmt.Println("-----------------Jenkins Job Info Fetcher v1.1")
	host := os.Args[1]
	user := os.Args[2]
	pass := os.Args[3]
	ctx := context.Background()
	jenkins := gojenkins.CreateJenkins(nil, host, user, pass)
	// Provide CA certificate if server is using self-signed certificate
	// caCert, _ := ioutil.ReadFile("/tmp/ca.crt")
	// jenkins.Requester.CACert = caCert
	_, err := jenkins.Init(ctx)

	if err != nil {
		panic("Something Went Wrong")
	}
	fmt.Println("User Grzegorz B logged to Jenkins v.2.332.2")
	fmt.Println("Fetching jobs data")
	//nodes, _ := jenkins.GetAllNodes(ctx)

	/*for _, node := range nodes {
		// Fetch Node Data
		node.Poll(ctx)
		flag, _ := node.IsOnline(ctx)
		if flag {
			fmt.Println("Node is Online")
		}
	}*/

	//jobName := "Declarative Pipeline"
	//builds, err := jenkins.GetAllBuildIds(ctx, jobName)

	/*if err != nil {
		panic(err)
	}*/
	/*
		for _, build := range builds {
			buildId := build.Number
			data, err := jenkins.GetBuild(ctx, jobName, buildId)

			if err != nil {
				panic(err)
			}
			fmt.Println(build)
			if "SUCCESS" == data.GetResult() {
				fmt.Println("This build succeeded")
			}
		}
	*/
	// Get Last Successful/Failed/Stable Build for a Job
	/*job, err := jenkins.GetJob(ctx, jobName)

	if err != nil {
		panic(err)
	}

	job.GetLastSuccessfulBuild(ctx)
	job.GetLastStableBuild(ctx)
	*/

	jobsInfo := []jobInfo{}

	jobs, err := jenkins.GetAllJobs(ctx)

	if err != nil {
		panic(err)
	}
	counter := 0
	for _, job := range jobs {

		builds, err := jenkins.GetAllBuildIds(ctx, job.GetName())

		if err != nil {
			panic(err)
		}

		for _, build := range builds {
			buildId := build.Number
			data, err := jenkins.GetBuild(ctx, job.GetName(), buildId)
			if err != nil {
				panic(err)
			}

			userId := data.GetActions()[0].Causes[0]["userId"]

			jI := jobInfo{
				JobName: job.GetName(), BuildNumber: build.Number,
				StartDate: data.GetTimestamp(), StartedBy: userId,
				BuildResult: data.GetResult(),
			}
			jobsInfo = append(jobsInfo, jI)
		}
		counter++
		if counter == 5 {
			break
		}
	}

	input := ""
	input1 := ""
	input2 := ""

	fmt.Println("Which path of destiny will you choose young devops apprentice?\nType:\ninfo for printing Jenkins jobs information\ncount for counting failed/passed jobs (it needs fetched jobs info - it can be \nfetched from local file)\nexit for exit ;)")

	flag := true

	for flag {
		fmt.Scanln(&input)
		if input == "info" {
			fmt.Println("type b/d for sorting the data by build number or start date - b is default")
			input = ""
			fmt.Scanln(&input)
			if len(input) < 1 {
				input = "b"
			}
			fmt.Println("if you want to save the data to local file, type it's path,\notherwise the data will be shown on the screen")
			fmt.Scanln(&input1)
			if len(input1) < 1 {
				input2 = "s"
			} else {
				input2 = "f"
			}
			if input2 == "s" {
				jobDisplayer(jobSorter(jobsInfo, input), input2, input1)
			} else {
				jobDisplayer(jobSorter(jobsInfo, input), input2, input1)
			}
			fmt.Println("\nHit enter to continue")
			fmt.Scanln(&input)
		}
		if input == "exit" {
			print("★─▄█▀▀║░▄█▀▄║▄█▀▄║██▀▄║─★\n★─██║▀█║██║█║██║█║██║█║─★\n★─▀███▀║▀██▀║▀██▀║███▀║─★\n★───────────────────────★\n★───▐█▀▄─ ▀▄─▄▀ █▀▀──█───★\n★───▐█▀▀▄ ──█── █▀▀──▀───★\n★───▐█▄▄▀ ──▀── ▀▀▀──▄───★")
			return
		}

	}
	//_ = jobDisplayer(jobSorter(jobsInfo, "d"), "s", "")
	//fmt.Println(jobDisplayer(jobSorter(jobsInfo, "d"), "f", "json.txt"))
	//jobSorter(jobsInfo, "b")

}

func jobSorter(j []jobInfo, s string) []jobInfo {
	jr := []jobInfo{}
	if s == "d" {
		dates := make(map[int]int)
		for d := range j {
			dates[int(j[d].StartDate.Unix())] = d
		}
		datesSorted := make([]int, 0, len(dates))
		for k := range dates {
			datesSorted = append(datesSorted, k)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(datesSorted)))
		for _, d := range datesSorted {
			jr = append(jr, j[dates[d]])
		}
		return jr
	}
	if s == "b" {
		sort.Slice(j, func(x, y int) bool {
			return j[x].BuildNumber > j[y].BuildNumber
		})
		for k := range j {
			jr = append(jr, j[k])
		}
		return jr
	}
	fmt.Println("Wrong sorting parameter")
	return nil
}

func jobDisplayer(j []jobInfo, d string, path string) error {
	if d == "s" {
		j, _ := json.MarshalIndent(j, "", "  ")
		fmt.Println(string(j))
	}
	if d == "f" && len(path) > 0 {
		j, _ := json.MarshalIndent(j, "", "  ")
		return ioutil.WriteFile(path, []byte(j), 0666)
	}
	return nil
}
