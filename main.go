package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	/* Open the file, defer close. (os.O_RDWR is opening the file for read and write,
	O_CREATE creates if does not exist, '0666 gives read and write permissions) */
	tasks, err := os.OpenFile("tasks.txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer tasks.Close()

	// Use the command line arguments to decide what to do
	if len(os.Args) <= 1 {
		log.Fatal("No arguments provided.")
	}

	action := os.Args[1]

	switch action {
	case "add-task":
		CheckArgLength(4, os.Args)
		taskToAdd := os.Args[2]
		dueDate := os.Args[3]
		newId := GenerateTaskID(tasks)
		AddTask(tasks, newId, taskToAdd, dueDate)
	case "remove":
		CheckArgLength(3, os.Args)
		taskToRemove := os.Args[2]
		RemoveTask(tasks, taskToRemove)
	case "list":
		ViewTasks(tasks)
	default:
		log.Fatal("Invalid action provided.")
	}
}

// Gets the largest ID in the file and returns a new id incremented by one
func GenerateTaskID(tasks *os.File) int {
	var (
		biggestID int
		line      string
	)
	scanner := bufio.NewScanner(tasks)
	// Finds the biggest ID in the file
	for scanner.Scan() {
		line = scanner.Text()
		taskId := string(line[0])
		if taskId != "" {
			num, err := strconv.Atoi(taskId)
			if err != nil {
				log.Fatal(err)
			} else if num > biggestID {
				biggestID = num
			}
		}
	}
	return biggestID + 1
}

func RemoveTask(tasks *os.File, taskToRemove string) {
	ResetFilePointer(tasks)

	// Creates a slice with every line except the one to remove
	scanner := bufio.NewScanner(tasks)
	re := regexp.MustCompile(`^\d+`)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindString(line)
		if match != taskToRemove {
			lines = append(lines, line)
		}
	}

	// Empty the file
	err := tasks.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	ResetFilePointer(tasks)

	// Write the lines back to the file
	_, err = tasks.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func ResetFilePointer(tasks *os.File) {
	// Put the file pointer back to the start of the file
	_, err := tasks.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func AddTask(tasks *os.File, id int, taskToAdd string, dueDate string) {
	lineToAdd := fmt.Sprintf("%d. %s - %s\n", id, taskToAdd, dueDate)
	_, err := tasks.WriteString(lineToAdd)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckArgLength(expectedLength int, args []string) {
	if len(args) != expectedLength {
		log.Fatal("Incorrect number of arguments provided.")
	}
}

func ViewTasks(tasks *os.File) {
	allTasks, err := io.ReadAll(tasks)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(allTasks))
}
