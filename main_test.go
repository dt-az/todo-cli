package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadTasks(t *testing.T) {
	// Test case 1: Empty file
	err := os.WriteFile(tasksFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tasksFile)

	tasks, err := loadTasks()
	if err != nil {
		t.Errorf("loadTasks() returned an error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("loadTasks() returned %d tasks, expected 0", len(tasks))
	}

	// Test case 2: File doesn't exist
	os.Remove(tasksFile)
	tasks, err = loadTasks()
	if err != nil {
		t.Errorf("loadTasks() returned an error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("loadTasks() returned %d tasks, expected 0", len(tasks))
	}

	// Test case 3: Valid JSON data
	data := []byte(`[{"text": "Task 1", "completed": false}, {"text": "Task 2", "completed": true}]`)
	err = os.WriteFile(tasksFile, data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	tasks, err = loadTasks()
	if err != nil {
		t.Errorf("loadTasks() returned an error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("loadTasks() returned %d tasks, expected 2", len(tasks))
	}
	if tasks[0].Text != "Task 1" || tasks[0].Completed != false {
		t.Errorf("Task 1 is incorrect: %+v", tasks[0])
	}
	if tasks[1].Text != "Task 2" || tasks[1].Completed != true {
		t.Errorf("Task 2 is incorrect: %+v", tasks[1])
	}

	// Test case 4: Invalid JSON data
	err = os.WriteFile(tasksFile, []byte(`[{"text": "Task 1", "completed": false}`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = loadTasks()
	if err == nil {
		t.Errorf("loadTasks() should return error for invalid JSON")
	}
}

func TestSaveTasks(t *testing.T) {
	tasks := []Task{{Text: "Task 1", Completed: false}, {Text: "Task 2", Completed: true}}
	err := saveTasks(tasks)
	if err != nil {
		t.Errorf("saveTasks() returned an error: %v", err)
	}
	defer os.Remove(tasksFile)

	data, err := os.ReadFile(tasksFile)
	if err != nil {
		t.Fatal(err)
	}

	var loadedTasks []Task
	err = json.Unmarshal(data, &loadedTasks)
	if err != nil {
		t.Fatal(err)
	}

	if len(loadedTasks) != 2 {
		t.Errorf("Saved file contains %d tasks, expected 2", len(loadedTasks))
	}

	if loadedTasks[0].Text != "Task 1" || loadedTasks[0].Completed != false {
		t.Errorf("Loaded Task 1 is incorrect: %+v", loadedTasks[0])
	}
	if loadedTasks[1].Text != "Task 2" || loadedTasks[1].Completed != true {
		t.Errorf("Loaded Task 2 is incorrect: %+v", loadedTasks[1])
	}
}

func TestAddTask(t *testing.T) {
	defer os.Remove(tasksFile)

	err := addTask("New Task")
	if err != nil {
		t.Errorf("addTask() returned an error: %v", err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Errorf("addTask() should add one task. Found %d", len(tasks))
	}

	if tasks[0].Text != "New Task" || tasks[0].Completed != false {
		t.Errorf("Added task is incorrect: %+v", tasks[0])
	}

	// Test adding a second task
	err = addTask("Another Task")
	if err != nil {
		t.Errorf("addTask() returned an error: %v", err)
	}

	tasks, err = loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Errorf("addTask() should add one task. Found %d", len(tasks))
	}

	if tasks[1].Text != "Another Task" || tasks[1].Completed != false {
		t.Errorf("Added task is incorrect: %+v", tasks[1])
	}
}

func TestCompleteTask(t *testing.T) {
	defer os.Remove(tasksFile)
	tasks := []Task{{Text: "Task 1", Completed: false}, {Text: "Task 2", Completed: false}}
	saveTasks(tasks)

	err := completeTask(1)
	if err != nil {
		t.Errorf("completeTask() returned an error: %v", err)
	}

	loadedTasks, _ := loadTasks()
	if !loadedTasks[0].Completed {
		t.Error("Task 1 should be completed")
	}

	if loadedTasks[1].Completed {
		t.Error("Task 2 should not be completed")
	}

	err = completeTask(3) // Invalid task number
	if err == nil {
		t.Error("completeTask() should return error for invalid task number")
	}
}

func TestClearTasks(t *testing.T) {
	defer os.Remove(tasksFile)

	tasks := []Task{{Text: "Task 1", Completed: false}}
	saveTasks(tasks)

	err := clearTasks()
	if err != nil {
		t.Errorf("clearTasks() returned an error: %v", err)
	}

	loadedTasks, _ := loadTasks()
	if len(loadedTasks) != 0 {
		t.Error("Tasks should be cleared")
	}
}
