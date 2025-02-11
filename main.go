package main

import (
        "encoding/json"
        "flag"
        "fmt"
        "os"
)

type Task struct {
        Text      string `json:"text"`
        Completed bool   `json:"completed"`
}

const tasksFile = "tasks.json"

func loadTasks() ([]Task, error) {
        data, err := os.ReadFile(tasksFile)
        if err != nil {
                if os.IsNotExist(err) {
                        return []Task{}, nil
                }
                return nil, fmt.Errorf("loading tasks: %w", err)
        }

        var tasks []Task
        err = json.Unmarshal(data, &tasks)
        if err != nil {
                return nil, fmt.Errorf("unmarshaling tasks: %w", err)
        }
        return tasks, nil
}

func saveTasks(tasks []Task) error {
        data, err := json.MarshalIndent(tasks, "", "  ")
        if err != nil {
                return fmt.Errorf("marshaling tasks: %w", err)
        }
        err = os.WriteFile(tasksFile, data, 0644)
        if err != nil {
                return fmt.Errorf("saving tasks: %w", err)
        }
        return nil
}

func addTask(text string) error {
        tasks, err := loadTasks()
        if err != nil {
                return err
        }

        newTask := Task{Text: text, Completed: false}
        tasks = append(tasks, newTask)

        return saveTasks(tasks)
}

func listTasks() error {
        tasks, err := loadTasks()
        if err != nil {
                return err
        }

        if len(tasks) == 0 {
                fmt.Println("No tasks found.")
                return nil
        }

        for i, task := range tasks {
                status := "[ ]"
                if task.Completed {
                        status = "[x]"
                }
                fmt.Printf("%d. %s %s\n", i+1, status, task.Text)
        }
        return nil
}

func completeTask(taskNumber int) error {
        tasks, err := loadTasks()
        if err != nil {
                return err
        }

        if taskNumber <= 0 || taskNumber > len(tasks) {
                return fmt.Errorf("invalid task number")
        }

        tasks[taskNumber-1].Completed = true
        return saveTasks(tasks)
}

func clearTasks() error {
	return saveTasks([]Task{})
}

func printHelp() {
	fmt.Println("Usage: todo <command> [options]")
	fmt.Println("Available commands:")
	fmt.Println("  add   - Add a new task")
	fmt.Println("    Options:")
	fmt.Println("      -text <task_text> (required): The text of the task")
	fmt.Println("  list  - List all tasks")
	fmt.Println("  complete - Mark a task as complete")
	fmt.Println("    Options:")
	fmt.Println("      -number <task_number> (required): The number of the task to complete")
	fmt.Println("  clear - Clear all tasks")
	fmt.Println("  help  - Display this help message")
}

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addText := addCmd.String("text", "", "Task text")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	completeNumber := completeCmd.Int("number", 0, "Task number to complete")

	if len(os.Args) < 2 {
					printHelp()
					os.Exit(1)
	}

	command := os.Args[1] // Store the command

	switch command {
	case "add":
					addCmd.Parse(os.Args[2:])
					if *addText == "" {
									fmt.Println("Task text is required.")
									addCmd.PrintDefaults()
									os.Exit(1)
					}
					if err := addTask(*addText); err != nil {
									fmt.Println("Error adding task:", err)
									os.Exit(1)
					}
					fmt.Println("Task added:", *addText)

	case "list":
					listCmd.Parse(os.Args[2:])
					if err := listTasks(); err != nil {
									fmt.Println("Error listing tasks:", err)
									os.Exit(1)
					}

	case "complete":
					completeCmd.Parse(os.Args[2:])
					if *completeNumber <= 0 {
									fmt.Println("Task number is required and must be greater than 0.")
									completeCmd.PrintDefaults()
									os.Exit(1)
					}
					if err := completeTask(*completeNumber); err != nil {
									fmt.Println("Error completing task:", err)
									os.Exit(1)
					}
					fmt.Printf("Task %d completed.\n", *completeNumber)
	case "clear":
					if err := clearTasks(); err != nil {
									fmt.Println("Error clearing tasks:", err)
									os.Exit(1)
					}
					fmt.Println("All tasks cleared.")

	case "help":
					printHelp()
					os.Exit(0) // Exit cleanly after printing help

	default:
					fmt.Println("Unknown command:", command) // Use the stored command
					printHelp() // Print help on unknown command
					os.Exit(1)
	}
}
