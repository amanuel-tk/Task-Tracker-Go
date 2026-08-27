package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide a command.")
		return
	}

	var task []Task

	filepath := "database.json"

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	byteValues, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("failed to read file", err)
		return
	}

	if len(byteValues) > 0 {
		err = json.Unmarshal(byteValues, &task)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
	}

	currentTime := time.Now().Format(time.RFC3339)

	switch args[0] {

	case "add":
		if len(args) < 2 {
			fmt.Println("Error: Please provide description . Example: add 'task to add'")
			return
		}
		newID := 1
		if len(task) > 0 {
			for _, t := range task {
				if t.ID > newID {
					newID = t.ID
				}
			}
			newID += 1
		}

		task = AddTask(task, Task{
			ID:          newID,
			Description: args[1],
			Status:      "todo",
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		})

	case "update":

		if len(args) < 3 {
			fmt.Println("Error: Please provide status and ID. Example: update 2 'task'")
			return
		}

		i, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error Please provide valid Id", args[1])
			return
		}

		task = UpdateTask(task, Task{
			ID:          i,
			Description: args[2],
			UpdatedAt:   currentTime,
		})
	case "delete":
		if len(args) < 2 {
			fmt.Println("Error: Please provide the right arguments. Example: delete 1")
			return
		}

		i, err := strconv.Atoi(args[1])

		if err != nil {
			fmt.Println("Error: Please provide valid id")
			return
		}

		task = DeleteTask(task, i)
	case "mark-in-progress":

		if len(args) < 2 {
			fmt.Println("Error: Please provide the right arguments. Example: mark-in-progress 1")
			return
		}

		i, err := strconv.Atoi(args[1])

		if err != nil {
			fmt.Println("Error: Please provide valid id")
			return
		}

		task = UpdateStatus(task, Task{
			ID:        i,
			Status:    "in-progress",
			UpdatedAt: currentTime,
		})
	case "mark-done":

		if len(args) < 2 {
			fmt.Println("Error: Please provide the right arguments. Example: mark-done 1")
			return
		}

		i, err := strconv.Atoi(args[1])

		if err != nil {
			fmt.Println("Error: Please provide valid id")
			return
		}

		task = UpdateStatus(task, Task{
			ID:        i,
			Status:    "done",
			UpdatedAt: currentTime,
		})

	case "list":
		if len(args) == 1 {

			list, _ := FormatJson(task)
			fmt.Println(string(list))
			return
		} else {
			value := ShowDataByStatus(task, args[1])

			list, _ := FormatJson(value)
			fmt.Println(string(list))
			return
		}
	default:
		fmt.Println("Error: please provide valid command")
		return
	}

	if err := file.Truncate(0); err != nil {
		fmt.Println("Error truncating file:", err)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(task); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

}

func AddTask(task []Task, newTask Task) []Task {

	if newTask.Description == "" {
		fmt.Println("Description cannot be empty.")
		return task
	}

	task = append(task, newTask)
	fmt.Println("Task added successfully", newTask)
	return task
}

func UpdateTask(task []Task, updateTask Task) []Task {

	if updateTask.Description == "" {
		fmt.Println("Please Provide Description, and id")
		return task
	}

	found := false

	for index, t := range task {
		if t.ID == updateTask.ID {
			task[index].Description = updateTask.Description
			task[index].UpdatedAt = updateTask.UpdatedAt
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Error: Task ID doesn't exist", updateTask.ID)
	}

	return task

}

func DeleteTask(task []Task, id int) []Task {
	for i, t := range task {
		if t.ID == id {
			task = append(task[:i], task[i+1:]...)
			fmt.Println("Deleted task:", t)
			return task
		}
	}
	fmt.Println("Error: id not found")

	return task
}

func UpdateStatus(task []Task, updateStatus Task) []Task {

	for index, t := range task {
		if t.ID == updateStatus.ID {
			task[index].Status = updateStatus.Status
			task[index].UpdatedAt = updateStatus.UpdatedAt
			return task
		}
	}
	fmt.Println("Error: Task ID doesn't exist", updateStatus.ID)
	return task
}

func FormatJson(task []Task) ([]byte, error) {
	list, err := json.MarshalIndent(task, "", " ")
	if err != nil {
		fmt.Println("Error formatting data to JSON:", err)
		return nil, err
	}
	return list, nil
}

func ShowDataByStatus(task []Task, status string) []Task {
	newTask := []Task{}
	for _, task := range task {
		if task.Status == status {
			newTask = append(newTask, task)
		}
	}
	return newTask
}
