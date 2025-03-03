package src

import (
	"bytes"
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsCompleted bool
}

func readDataFromFile(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func readTasksFromData(data []byte) ([]Task, error) {
	csvReader := csv.NewReader(bytes.NewReader(data))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for i, record := range records {
		if i == 0 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		createdAt, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			return nil, err
		}

		completed, err := strconv.ParseBool(record[3])
		if err != nil {
			return nil, err
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

func ListTasks(fileName string) ([]Task, error) {

	return nil, nil
}
