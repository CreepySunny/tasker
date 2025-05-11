package tasks

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseTasksFromData(t *testing.T) {
	t.Run("Read_Tasks_From_Data", func(t *testing.T) {
		data := []byte(`ID,Description,CreatedAt,IsComplete
1,My new task,2024-07-27T16:45:19-05:00,true
2,Finish this video,2024-07-27T16:45:26-05:00,true
3,Find a video editor,2024-07-27T16:45:31-05:00,false`)

		got, err := readTasksFromCSVData(data)
		want := []Task{
			{ID: 1, Description: "My new task", CreatedAt: time.Date(2024, 7, 27, 16, 45, 19, 0, time.FixedZone("", -5*3600)), IsCompleted: true},
			{ID: 2, Description: "Finish this video", CreatedAt: time.Date(2024, 7, 27, 16, 45, 26, 0, time.FixedZone("", -5*3600)), IsCompleted: true},
			{ID: 3, Description: "Find a video editor", CreatedAt: time.Date(2024, 7, 27, 16, 45, 31, 0, time.FixedZone("", -5*3600)), IsCompleted: false},
		}

		if err != nil {
			t.Errorf("[!] error [!]: %v", err)
			return
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})
}

func TestTaskString(t *testing.T) {
	task := Task{
		ID:          42,
		Description: "Test task",
		CreatedAt:   time.Date(2025, 5, 12, 10, 0, 0, 0, time.UTC),
		IsCompleted: true,
	}
	expected := "42\tTest task\t2025-05-12T10:00:00Z\ttrue"
	if got := task.String(); got != expected {
		t.Errorf("Task.String() = %q, want %q", got, expected)
	}
}

func TestLoadAndCloseFile(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test_loadfile.csv")
	defer os.Remove(tmpFile)

	f, err := loadFile(tmpFile)
	if err != nil {
		t.Fatalf("loadFile() error = %v", err)
	}
	if f == nil {
		t.Fatal("loadFile() returned nil file")
	}

	err = closeFile(f)
	if err != nil {
		t.Errorf("closeFile() error = %v", err)
	}
}

func TestReadTasksFromCSVData(t *testing.T) {
	cases := []struct {
		name    string
		data    []byte
		wantErr string // substring to match in error, or "" for no error
		wantLen int
	}{
		{"Valid data", []byte("ID,Description,CreatedAt,IsComplete\n1,Test,2025-05-12T10:00:00Z,true"), "", 1},
		{"Empty data", []byte{}, "", 0},
		{"Malformed record", []byte("ID,Description,CreatedAt,IsComplete\n1,Test,2025-05-12T10:00:00Z"), "malformed record", 0},
		{"Invalid ID", []byte("ID,Description,CreatedAt,IsComplete\nfoo,Test,2025-05-12T10:00:00Z,true"), "failed to parse ID", 0},
		{"Invalid CreatedAt", []byte("ID,Description,CreatedAt,IsComplete\n1,Test,notatime,true"), "failed to parse CreatedAt", 0},
		{"Invalid IsCompleted", []byte("ID,Description,CreatedAt,IsComplete\n1,Test,2025-05-12T10:00:00Z,notabool"), "failed to parse IsCompleted", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readTasksFromCSVData(tc.data)
			if tc.wantErr == "" && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Errorf("expected error containing %q, got %v", tc.wantErr, err)
			}
			if len(got) != tc.wantLen {
				t.Errorf("expected %d tasks, got %d", tc.wantLen, len(got))
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test_listtasks.csv")
	defer os.Remove(tmpFile)

	content := "ID,Description,CreatedAt,IsComplete\n1,Task1,2025-05-12T10:00:00Z,true\n2,Task2,2025-05-12T11:00:00Z,false\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	t.Run("all_true_returns_all_tasks", func(t *testing.T) {
		tasks, err := ListTasks(tmpFile, true)
		if err != nil {
			t.Errorf("ListTasks(all=true) error: %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("all_false_returns_uncompleted_tasks", func(t *testing.T) {
		tasks, err := ListTasks(tmpFile, false)
		if err != nil {
			t.Errorf("ListTasks(all=false) error: %v", err)
		}
		if len(tasks) != 1 || tasks[0].ID != 2 {
			t.Errorf("expected only uncompleted task, got %+v", tasks)
		}
	})

	t.Run("file_does_not_exist_returns_error", func(t *testing.T) {
		nonexistent := filepath.Join(os.TempDir(), "doesnotexist.csv")
		_, err := ListTasks(nonexistent, true)
		if err == nil || !strings.Contains(err.Error(), "Unable to locate file") {
			t.Errorf("expected file not found error, got %v", err)
		}
	})
}
