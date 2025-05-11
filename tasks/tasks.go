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

func loadFile(filepath string) (*os.File, error) {
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
		return nil, err
	}

	var tasks []Task
	for _, record := range records {
		if len(record) < 4 {
			return nil, fmt.Errorf("malformed record: %v", record)
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

func ListTasks(filename string, all bool) ([]Task, error) {
	// Check if file exists
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("'%s': Unable to locate file", filename)
	}

	// Load and syslock file
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
