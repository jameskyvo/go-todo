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

/*
PROJECT: TODO APP
A todo app that takes in command line arguments and uses them to manage a to-do list in a text file.

REQUIREMENTS:
Add Task (Index, Date, Due Date, Description)
Remove Task by ID
View Tasks

Store tasks in a text file.
*/
func main() {
	/* Open the file, defer close. (os.O_RDWR is opening the file for read and write,
	O_CREATE creates if does not exist, '0666 gives read and write permissions) */
	tasks, err := os.OpenFile("tasks.txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer tasks.Close()

	// Use the command line arguments to decide what to do
	action := os.Args[1]

	if action == "add-task" {
		CheckArgLength(4, os.Args)
		taskToAdd := os.Args[2]
		dueDate := os.Args[3]
		newId := GenerateTaskID(tasks)
		AddTask(tasks, newId, taskToAdd, dueDate)
	} else if action == "remove" {
		CheckArgLength(3, os.Args)
		taskToRemove := os.Args[2]
		RemoveTask(tasks, taskToRemove)
	} else if action == "list" {
		ViewTasks(tasks)
	} else {
		log.Fatal("Invalid action provided.")
	}
}

// Gets the largest ID in the file and returns a new id incremented by one
func GenerateTaskID(tasks *os.File) int {
	biggestID := 0
	line := ""
	// Regex to match the first number in the line.
	re := regexp.MustCompile(`^\d+`)
	scanner := bufio.NewScanner(tasks)
	// Finds the biggest ID in the file
	for scanner.Scan() {
		line = scanner.Text()
		match := re.FindString(line)
		if match != "" {
			num, err := strconv.Atoi(match)
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

func CheckArgLength(i int, args []string) {
	if len(os.Args) != i {
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
