package utils

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
)

func ParseArgs() (bool, bool) {
	tasks := flag.String("task", "", "Check and Claim Task (Y/n)")
	reffs := flag.String("reff", "", "Do you want to claim referrals? (Y/n)")
	flag.Parse()

	var tasksEnable bool
	var reffsEnable bool

	if *tasks == "" {
		fmt.Print("Do you want to check and claim tasks? (Y/n): ")
		var taskInput string
		fmt.Scanln(&taskInput)
		taskInput = strings.TrimSpace(strings.ToLower(taskInput))

		switch taskInput {
		case "y":
			tasksEnable = true
		case "n":
			tasksEnable = false
		default:
			tasksEnable = true
		}
	}

	if *reffs == "" {
		fmt.Print("Do you want to claim Referrals? (Y/n): ")
		var reffInput string
		fmt.Scanln(&reffInput)
		reffInput = strings.TrimSpace(strings.ToLower(reffInput))

		switch reffInput {
		case "y":
			reffsEnable = true
		case "no":
			reffsEnable = false
		default:
			reffsEnable = true
		}
	}

	return tasksEnable, reffsEnable
}

func PrintLogo() {

	fmt.Printf(yellow(" ____  _    _   _ __  __ ____   ___ _____ \n| __ )| |  | | | |  \\/  | __ ) / _ \\_   _|\n|  _ \\| |  | | | | |\\/| |  _ \\| | | || |  \n| |_) | |__| |_| | |  | | |_) | |_| || |  \n|____/|_____\\___/|_|  |_|____/ \\___/ |_|  \n------------------------------------------\n"))
}

func ClearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func loadListFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var parseList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			parseList = append(parseList, strings.TrimSpace(scanner.Text()))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(parseList) < 1 {
		return nil, errors.New(fmt.Sprintf("\"%v\" is empty!", fileName))
	}

	return parseList, nil

}

func ParseQueries() ([]string, error) {
	return loadListFile("./configs/query_list.conf")
}

func FormatUpTime(d time.Duration) string {
	totalSeconds := int(d.Seconds())

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
}

func ConvertStrTimestamp(timestampStr string) (string, error) {
	// Convert the string to an integer
	epochTime, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("error converting timestamp: %w", err)
	}

	// Convert milliseconds to seconds and nanoseconds
	seconds := epochTime / 1000
	nanoseconds := (epochTime % 1000) * 1e6

	t := time.Unix(seconds, nanoseconds)
	const layout = "Monday, 3:04 PM"

	humanReadableTime := t.Format(layout)

	return humanReadableTime, nil
}
