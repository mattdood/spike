package run

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"
)

var (
	OutputBaseDirectory = getUserHome()
	OutputDirectory     = path.Join(OutputBaseDirectory, InstallBaseDirectory)
	OutputPath          = path.Join(OutputDirectory, TasksPath)
)

const (
	InstallBaseDirectory string      = "spikes/"
	TasksPath            string      = "tasks.json"
	FolderPermission     fs.FileMode = 00775
	FilePermission       fs.FileMode = 00644
	OpenStatus           string      = "O"
	ClosedStatus         string      = "C"
)

func getUserHome() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	return homedir
}

// Nullable values should be pointers
type Task struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	ID          int    `json:"id"`
}

type TaskList struct {
	Open   []Task `json:"O"`
	Closed []Task `json:"C"`
}

// Add a task to the task list
func (tl *TaskList) AddTask(t Task, status string) error {
	if status == OpenStatus {
		tl.Open = append(tl.Open, t)
		return nil
	}

	if status == ClosedStatus {
		tl.Closed = append(tl.Closed, t)
		return nil
	}

	return errors.New("status should be one of: 'O', 'C'")
}

// Swap a task status between closed an open
func (tl *TaskList) UpdateTaskStatus(id int, status string) error {
	today := time.Now().Format("2006-01-02")

	if status == ClosedStatus {
		for i, task := range tl.Open {
			if task.ID == id {
				task.Updated = today
				// Move from open to closed
				tl.Closed = append(tl.Closed, task)
				// Remove task fro mopen
				tl.Open = append(tl.Open[:i], tl.Open[i+1:]...)

				return nil
			}
		}

		return errors.New("task not found in status 'O' to be changed to 'C'")
	}

	if status == OpenStatus {
		for i, task := range tl.Closed {
			if task.ID == id {
				task.Updated = today
				// Move from open to closed
				tl.Open = append(tl.Open, task)
				// Remove task fro mopen
				tl.Closed = append(tl.Closed[:i], tl.Closed[i+1:]...)

				return nil
			}
		}

		return errors.New("task not found in status 'C' to be changed to 'O'")
	}

	return errors.New("status must be one of 'O' or 'C'")
}

func (tl *TaskList) UpdateTaskName(id int, name string) error {
	var foundTask bool

	today := time.Now().Format("2006-01-02")

	for i, task := range tl.Open {
		if task.ID == id {
			tl.Open[i].Name = name

			tl.Open[i].Updated = today

			return nil
		}
	}

	if !foundTask {
		for i, task := range tl.Closed {
			if task.ID == id {
				tl.Closed[i].Name = name

				tl.Closed[i].Updated = today

				return nil
			}
		}
	}

	return fmt.Errorf("task not found for ID: %d", id)
}

func (tl *TaskList) UpdateTaskDescription(id int, desc string) error {
	var foundTask bool

	today := time.Now().Format("2006-01-02")

	for i, task := range tl.Open {
		if task.ID == id {
			tl.Open[i].Name = desc
			tl.Open[i].Updated = today

			return nil
		}
	}

	if !foundTask {
		for i, task := range tl.Closed {
			if task.ID == id {
				tl.Closed[i].Name = desc
				tl.Closed[i].Updated = today

				return nil
			}
		}
	}

	return fmt.Errorf("task not found for ID: %d", id)
}

// Read tasks.json file and return representation
// of data in sorted form. Each "O" and "C" slice is
// sorted in descending order.
func NewTaskList() *TaskList {
	fileData, err := os.ReadFile(OutputPath)
	if err != nil {
		fmt.Println(err)
	}

	var taskList TaskList

	err = json.Unmarshal(fileData, &taskList)
	if err != nil {
		fmt.Println(err)
	}

	sort.Slice(taskList.Open, func(i, j int) bool { return taskList.Open[i].ID > taskList.Open[j].ID })
	sort.Slice(taskList.Closed, func(i, j int) bool { return taskList.Closed[i].ID > taskList.Closed[j].ID })

	return &taskList
}

// Save the task list back to JSON
func SaveTaskList(tasks *TaskList) error {
	file, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %s", err)
	}

	err = os.WriteFile(OutputPath, file, FilePermission)
	if err != nil {
		return fmt.Errorf("error writing JSON: %s", err)
	}

	return nil
}

// Find the highest ID in the tasks list then
// increment by 1
func NewTaskID(taskList *TaskList) int {
	// Iterate ID field
	var nextID int

	if len(taskList.Open) > 0 && len(taskList.Closed) > 0 {
		// Both have tasks
		if taskList.Open[0].ID > taskList.Closed[0].ID {
			nextID = taskList.Open[0].ID + 1
		} else {
			nextID = taskList.Closed[0].ID + 1
		}
	} else if len(taskList.Open) > 0 && len(taskList.Closed) == 0 {
		// Only open has tasks
		nextID = taskList.Open[0].ID + 1
	} else if len(taskList.Closed) > 0 && len(taskList.Open) == 0 {
		// Only closed has tasks
		nextID = taskList.Closed[0].ID + 1
	} else {
		// Empty file
		nextID = 0
	}

	return nextID
}

// Create a new task based on the data input via the CLI
func Create(name string, description string) int {
	tasks := NewTaskList()

	today := time.Now().Format("2006-01-02")

	nextID := NewTaskID(tasks)

	err := tasks.AddTask(
		Task{
			ID:          nextID,
			Name:        name,
			Description: description,
			Created:     today,
			Updated:     today,
		},
		OpenStatus,
	)
	if err != nil {
		fmt.Println("error adding task: ", err)
		return 1
	}

	err = SaveTaskList(tasks)
	if err != nil {
		fmt.Println("error saving JSON: ", err)
		return 1
	}

	return 0
}

// Update a task's name, description, or status
func Update(id int, key string, value string, status string) int {
	tasks := NewTaskList()

	if key == "" && value == "" && status != "" {
		// Status update
		err := tasks.UpdateTaskStatus(id, status)
		if err != nil {
			fmt.Println("updating task status failed: ", err)
			return 1
		}
	}

	if key != "" && value != "" && status == "" {
		if key == "name" {
			// Updating Task.Name
			if err := tasks.UpdateTaskName(id, value); err != nil {
				fmt.Println(err)
				return 1
			}
		} else if key == "desc" {
			// Update Task.Description
			if err := tasks.UpdateTaskDescription(id, value); err != nil {
				fmt.Println(err)
				return 1
			}
		} else {
			// invalid key
			fmt.Println("key must be one of 'name', 'desc': ", key)
			return 1
		}
	}

	err := SaveTaskList(tasks)
	if err != nil {
		fmt.Println("error saving JSON: ", err)
		return 1
	}

	return 0
}

// List tasks
func List(status string) int {
	tasks := NewTaskList()

	fmt.Printf("Available tasks for the '%s' status\n", status)

	if status == OpenStatus {
		for _, task := range tasks.Open {
			fmt.Printf("ID: %d\nName: %s\nDesc.: %s\nCreated: %s\nUpdated: %s\n\n", task.ID, task.Name, task.Description, task.Created, task.Updated)
		}
	}

	if status == ClosedStatus {
		for _, task := range tasks.Closed {
			fmt.Printf("ID: %d\nName: %s\nDesc.: %s\nCreated: %s\nUpdated: %s\n\n", task.ID, task.Name, task.Description, task.Created, task.Updated)
		}
	}

	return 0
}

// Git command wrapper for `git add`
func Add(files []string) {
	out, err := exec.Command("git", "-C", OutputDirectory, "add", strings.Join(files, " ")).Output()
	if err != nil {
		fmt.Println("`git add` exited abnormally")
		fmt.Println(err)
	}

	output := string(out)

	fmt.Print(output)
}

// Git command wrapper for `git init`
func Init() {
	if _, err := os.Stat(OutputDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(OutputDirectory, FolderPermission)
		if err != nil {
			fmt.Println("Failed to create directories")
			fmt.Println(err)
		}

		// Create empty task list
		tasks := TaskList{
			Open:   []Task{},
			Closed: []Task{},
		}

		err = SaveTaskList(&tasks)
		if err != nil {
			fmt.Println("error saving JSON: ", err)
		}
	}

	out, err := exec.Command("git", "-C", OutputDirectory, "init").Output()
	if err != nil {
		fmt.Println("`git init` exited abnormally")
		fmt.Println(err)
	}

	output := string(out)

	fmt.Print(output)
}

// Git command wrapper for `git push`
func Push() {
	out, err := exec.Command("git", "-C", OutputDirectory, "push").Output()
	if err != nil {
		fmt.Println("`git push` exited abnormally")
		fmt.Println(err)
	}

	output := string(out)

	fmt.Print(output)
}

// Git command wrapper for `git pull`
func Pull() {
	out, err := exec.Command("git", "-C", OutputDirectory, "pull").Output()
	if err != nil {
		fmt.Println("`git push` exited abnormally")
		fmt.Println(err)
	}

	output := string(out[:])

	fmt.Print(output)
}
