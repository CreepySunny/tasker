package tasks

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestReadTasksFromCSVData(t *testing.T) {
	cases := []struct {
		name    string
		data    []byte
		wantErr string // substring to match in error, or "" for no error
		want    []Task // expected parsed tasks
	}{
		{
			name:    "Valid data",
			data:    []byte("ID,Description,CreatedAt,IsComplete\n1,Test,2025-05-12T10:00:00Z,true"),
			wantErr: "",
			want: []Task{{
				ID:          1,
				Description: "Test",
				CreatedAt:   time.Date(2025, 5, 12, 10, 0, 0, 0, time.UTC),
				IsCompleted: true,
			}},
		},
		{
			name:    "Empty data",
			data:    []byte{},
			wantErr: "",
			want:    nil,
		},
		{
			name:    "Malformed record",
			data:    []byte("ID,Description,CreatedAt,IsComplete\n1,Test,2025-05-12T10:00:00Z"),
			wantErr: "wrong number of fields",
			want:    nil,
		},
		{
			name:    "Invalid ID",
			data:    []byte("ID,Description,CreatedAt,IsComplete\nfoo,Test,2025-05-12T10:00:00Z,true"),
			wantErr: "failed to parse ID",
			want:    nil,
		},
		{
			name:    "Invalid CreatedAt",
			data:    []byte("ID,Description,CreatedAt,IsComplete\n1,Test,notatime,true"),
			wantErr: "failed to parse CreatedAt",
			want:    nil,
		},
		{
			name:    "Invalid IsCompleted",
			data:    []byte("ID,Description,CreatedAt,IsComplete\n1,Test,2025-05-12T10:00:00Z,notabool"),
			wantErr: "failed to parse IsCompleted",
			want:    nil,
		},
		{
			name:    "Multiple valid tasks with timezone",
			data:    []byte("ID,Description,CreatedAt,IsComplete\n1,My new task,2024-07-27T16:45:19-05:00,true\n2,Finish this video,2024-07-27T16:45:26-05:00,true\n3,Find a video editor,2024-07-27T16:45:31-05:00,false"),
			wantErr: "",
			want: []Task{
				{ID: 1, Description: "My new task", CreatedAt: time.Date(2024, 7, 27, 16, 45, 19, 0, time.FixedZone("", -5*3600)), IsCompleted: true},
				{ID: 2, Description: "Finish this video", CreatedAt: time.Date(2024, 7, 27, 16, 45, 26, 0, time.FixedZone("", -5*3600)), IsCompleted: true},
				{ID: 3, Description: "Find a video editor", CreatedAt: time.Date(2024, 7, 27, 16, 45, 31, 0, time.FixedZone("", -5*3600)), IsCompleted: false},
			},
		},
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
			if tc.want != nil && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parsed tasks mismatch:\n got:  %#v\n want: %#v", got, tc.want)
			}
			if tc.want == nil && len(got) != 0 {
				t.Errorf("expected no tasks, got: %#v", got)
			}
		})
	}
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
}

func TestAddTask(t *testing.T) {
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_addtask.csv")
	defer os.Remove(tmpFile)

	testCases := []struct {
		name        string
		fileSetup   func() string // returns filename
		description string
		wantErr     string
	}{
		{
			name: "Add to new file (creates header)",
			fileSetup: func() string {
				os.Remove(tmpFile)
				return tmpFile
			},
			description: "Test Add",
			wantErr:     "",
		},
		{
			name: "Add to file with only header",
			fileSetup: func() string {
				os.WriteFile(tmpFile, []byte("ID,Description,CreatedAt,IsComplete\n"), 0644)
				return tmpFile
			},
			description: "Header Only",
			wantErr:     "",
		},
		{
			name: "Add to file with existing tasks",
			fileSetup: func() string {
				os.WriteFile(tmpFile, []byte("ID,Description,CreatedAt,IsComplete\n1,Old Task,2025-05-12T09:00:00Z,false\n"), 0644)
				return tmpFile
			},
			description: "New Task",
			wantErr:     "",
		},
		{
			name: "Unwritable directory (should error)",
			fileSetup: func() string {
				return "/root/should_not_exist.csv"
			},
			description: "Fail",
			wantErr:     "failed to open datasource for appending",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := tc.fileSetup()
			err := AddTask(filename, tc.description)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
			if tc.wantErr == "" {
				data, err := os.ReadFile(filename)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")
				if len(lines) < 2 {
					t.Fatalf("expected at least 2 lines (header + task), got %d", len(lines))
				}
				// Check header
				if lines[0] != "ID,Description,CreatedAt,IsComplete" {
					t.Errorf("header: got %q, want %q", lines[0], "ID,Description,CreatedAt,IsComplete")
				}
				// Check that a line contains the description
				found := false
				for _, line := range lines[1:] {
					if strings.Contains(line, tc.description) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected to find description %q in file, but did not. Lines: %v", tc.description, lines)
				}
			}
		})
	}
}
