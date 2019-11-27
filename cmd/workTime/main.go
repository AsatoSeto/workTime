package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var exPath string

type workStruct struct {
	WorkStart  int64 `json:"work_start"`
	WorkEnd    int64 `json:"work_end"`
	BreakStart int64 `json:"break_start"`
	BreakEnd   int64 `json:"break_end"`
}

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath = filepath.Dir(ex)

	startWork := flag.Bool("wstart", false, "Start Work flag")
	startBreak := flag.Bool("bstart", false, "Start Break flag")
	stopBreak := flag.Bool("bstop", false, "Stop Break flag")
	stopWork := flag.Bool("wstop", false, "Stop Work flag")
	manualCalc := flag.Bool("manual", false, "Run dialog for calculate time")
	howMany := flag.Bool("howMany", false, "Check time")

	flag.Parse()

	if *startWork == true {
		startWorkF()
		notif("Work", "Work Day is Started")
	} else if *startBreak == true {
		startBreakF()
		notif("Lunch break", "Lunch break is Started")

	} else if *stopBreak == true {
		stopBreakF()
		notif("Lunch break", "Lunch break is Ended")

	} else if *stopWork == true {
		stopWorkF()

	} else if *manualCalc == true {
		timeCalc()
	} else if *howMany == true {
		howManyF()
	}
}

func notif(str1 string, str2 string) {
	cmd := exec.Command("bash",
		"-c", fmt.Sprintf("notify-send '%s:' '%s'", str1, str2), "1")
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func readFile() map[string]workStruct {
	workFile := make(map[string]workStruct)
	jsonFile, err := os.Open(exPath + "/workHisory.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = json.Unmarshal(bytes, &workFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	return workFile
}

func startWorkF() {
	workFile := make(map[string]workStruct)
	if _, err := os.Stat(exPath + "/workHisory.json"); err == nil {
		workFile = readFile()
	} else {
		log.Println("nope")
	}

	timeNow := time.Now()
	date := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	if _, ok := workFile[date]; ok {
		data := workFile[date]
		data.WorkStart = time.Now().Unix()
		workFile[date] = data
	} else {
		data := workStruct{}
		data.WorkStart = time.Now().Unix()
		workFile[date] = data
	}

	body, err := json.Marshal(workFile)
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(body)
}

func startBreakF() {
	workFile := make(map[string]workStruct)

	workFile = readFile()

	timeNow := time.Now()
	date := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	if _, ok := workFile[date]; ok {
		data := workFile[date]
		data.BreakStart = time.Now().Unix()
		workFile[date] = data
	} else {
		data := workStruct{}
		data.BreakStart = time.Now().Unix()
		workFile[date] = data
	}

	body, err := json.Marshal(workFile)
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(body)
}

func stopBreakF() {
	workFile := make(map[string]workStruct)

	workFile = readFile()

	timeNow := time.Now()
	date := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	if _, ok := workFile[date]; ok {
		data := workFile[date]
		data.BreakEnd = time.Now().Unix()
		workFile[date] = data
	} else {
		data := workStruct{}
		data.BreakEnd = time.Now().Unix()
		workFile[date] = data
	}

	body, err := json.Marshal(workFile)
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(body)
}

func stopWorkF() {
	workFile := make(map[string]workStruct)

	workFile = readFile()

	timeNow := time.Now()
	date := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	if _, ok := workFile[date]; ok {
		data := workFile[date]
		data.WorkEnd = time.Now().Unix()
		workFile[date] = data
	} else {
		data := workStruct{}
		data.WorkEnd = time.Now().Unix()
		workFile[date] = data
	}

	body, err := json.Marshal(workFile)
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(body)

	var (
		breakDur int64
	)

	if workFile[date].BreakEnd != 0 {
		breakDur = workFile[date].BreakEnd - workFile[date].BreakStart
	} else {
		breakDur = 0
	}

	planEnd := workFile[date].WorkStart + breakDur + (8 * 3600)

	if planEnd <= workFile[date].WorkEnd {
		cmd := exec.Command("bash",
			"-c", "notify-send 'Your overtime is:' "+fmt.Sprintf("'%d . Good Luck.'",
				((workFile[date].WorkEnd-planEnd)/60)), "1")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	} else {
		cmd := exec.Command("bash",
			"-c", "notify-send 'You not complite work:' "+fmt.Sprintf("'Reminde: %d'",
				((planEnd-workFile[date].WorkEnd)/60)), "1")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}

}

func timeCalc() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter start work: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	text = text[:len(text)-1]

	textParse := strings.Split(text, ":")

	startHour, err := strconv.Atoi(textParse[0])
	if err != nil {
		log.Println(err)
	}

	startMin, err := strconv.Atoi(textParse[1])
	if err != nil {
		log.Println(err)
	}

	startWork := (startHour * 3600) + (startMin * 60)

	fmt.Print("Enter the duration of the break: ")

	text, err = reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	text = text[:len(text)-1]

	breakMin, err := strconv.Atoi(text)
	if err != nil {
		log.Println(err)
	}

	startWork += breakMin * 60

	workDur := (startHour+8)*3600 + (startMin+breakMin)*60

	timeNow := time.Now()

	workNow := timeNow.Minute()*60 + timeNow.Hour()*3600

	endHour := workDur / 3600
	endMin := workDur % 3600
	fmt.Printf("Work End: %d:%d \n", endHour, endMin/60)
	fmt.Printf("Remained %d min \n", ((workDur - workNow) / 60))
}

func howManyF() {
	workFile := make(map[string]workStruct)
	timeNow := time.Now()
	date := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	workFile = readFile()
	var (
		breakDur int64
	)

	if workFile[date].BreakEnd != 0 {
		breakDur = workFile[date].BreakEnd - workFile[date].BreakStart
	} else {
		breakDur = 0
	}

	planEnd := workFile[date].WorkStart + breakDur + (8 * 3600)

	if planEnd <= workFile[date].WorkEnd {
		cmd := exec.Command("bash",
			"-c", "notify-send 'Your overtime is:' "+fmt.Sprintf("'%d'",
				((workFile[date].WorkEnd-planEnd)/60)), "1")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	} else {
		cmd := exec.Command("bash",
			"-c", "notify-send 'Working time left:' "+fmt.Sprintf("'%d Hour %d Min' ",
				((planEnd-time.Now().Unix())/3600), ((planEnd-time.Now().Unix())%3600/60)), "1")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}

}

func writeFile(body []byte) {
	f, err := os.Create(exPath + "/workHisory.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	_, err = f.Write(body)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
}
