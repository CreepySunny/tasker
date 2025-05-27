package tasks

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"time"
)

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsCompleted bool
}

func (t Task) String() string {
	return fmt.Sprintf("%d\t%s\t%s\t%t", t.ID, t.Description, t.CreatedAt.Format(time.RFC3339), t.IsCompleted)
}

// ensureDataSource checks if the file exists, and if not, creates it with the correct header.
func ensureDataSource(filepath string) error {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to create datasource: %w", err)
		}
		if _, err = file.WriteString("ID,Description,CreatedAt,IsComplete\n"); err != nil {
			file.Close()
			return fmt.Errorf("failed to write header: %w", err)
		}
		file.Close()
	}
	return nil
}

func loadFile(filepath string) (*os.File, error) {
	if err := ensureDataSource(filepath); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}
	return f, nil
}

func closeFile(f *os.File) error {
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		return fmt.Errorf("failed to unlock file: %w", err)
	}
	return f.Close()
}

func readTasksFromCSVData(data []byte) ([]Task, error) {
	if len(data) == 0 {
		return []Task{}, nil
	}

	csvReader := csv.NewReader(bytes.NewReader(data))
	_, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("malformed record: %w", err)
	}

	var tasks []Task
	for _, record := range records {
		if len(record) < 4 {
			return nil, fmt.Errorf("malformed record: expected at least 4 fields but got %d in record: %v", len(record), record)
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse ID: %w", err)
		}

		createdAt, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse CreatedAt: %w", err)
		}

		completed, err := strconv.ParseBool(record[3])
		if err != nil {
			return nil, fmt.Errorf("failed to parse IsCompleted: %w", err)
		}

		task := Task{
			ID:          id,
			Description: record[1],
			CreatedAt:   createdAt,
			IsCompleted: completed,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// AddTask appends a new task to the datasource (CSV file).
func AddTask(filename string, description string) error {
	file, err := loadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open datasource for appending: %w", err)
	}
	defer func() {
		if err := closeFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()

	// Move to end of file for appending
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("failed to seek to end of file: %w", err)
	}

	// Determine next ID
	file.Seek(0, io.SeekStart)
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file for ID: %w", err)
	}
	tasks, err := readTasksFromCSVData(data)
	if err != nil {
		return fmt.Errorf("failed to parse tasks for ID: %w", err)
	}
	nextID := 1
	if len(tasks) > 0 {
		nextID = tasks[len(tasks)-1].ID + 1
	}

	csvWriter := csv.NewWriter(file)
	record := []string{
		strconv.Itoa(nextID),
		description,
		time.Now().Format(time.RFC3339),
		"false",
	}
	if err := csvWriter.Write(record); err != nil {
		return fmt.Errorf("failed to write task: %w", err)
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

func ListTasks(filename string, all bool) ([]Task, error) {
	// Load and syslock file (this will create the file if it doesn't exist)
	file, err := loadFile(filename)
	if err != nil {
		return nil, err
	}

	// Defer file closing
	defer func() {
		if err := closeFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()

	// Read data from file
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Read tasks from CSV Data
	tasks, err := readTasksFromCSVData(data)
	if err != nil {
		return nil, err
	}

	// If all flag
	if all {
		return tasks, nil
	}

	// Reparate uncompleted tasks
	uncompletedTasks := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if !task.IsCompleted {
			uncompletedTasks = append(uncompletedTasks, task)
		}
	}

	// Return uncompleted list
	return uncompletedTasks, nil
}

func CompleteTask(filename string, taskID string) error {
	// Load and syslock file
	file, err := loadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open datasource for updating: %w", err)
	}
	defer func() {
		if err := closeFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()

	// Read data from file
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	tasks, err := readTasksFromCSVData(data)
	if err != nil {
		return fmt.Errorf("failed to parse tasks: %w", err)
	}
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("failed to parse task ID: %w", err)
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].IsCompleted = true
			tasks[i].CreatedAt = time.Now() // Update the CreatedAt to now
			break
		}
	}
	// Write updated tasks back to file
	csvWriter := csv.NewWriter(file)
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Description,
			task.CreatedAt.Format(time.RFC3339),
			strconv.FormatBool(task.IsCompleted),
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write task: %w", err)
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("failed to flush csv writer: %w", err)
	}
	// If we reach here, the task was successfully marked as completed
	return nil
}
