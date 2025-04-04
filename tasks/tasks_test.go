package tasks

import (
	"reflect"
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

	t.Run("Read_Data_From_File", func(t *testing.T) {

	})

}
