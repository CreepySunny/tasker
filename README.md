# Tasker

Tasker is a simple command-line Todo List application written in Go. It allows you to manage your tasks directly from the terminal, supporting adding, listing, completing, and deleting tasks, all stored in a CSV file with safe file locking.

## Features
- Add new tasks
- List uncompleted or all tasks
- Mark tasks as complete
- Delete tasks
- Data stored in a CSV file with file locking for safety
- Friendly time display (e.g., "a minute ago")

## Installation

Clone the repository and build the binary:

```sh
git clone https://github.com/CreepySunny/tasker.git
cd tasker
go build -o tasker
```

## Usage

### Add a Task
```
$ tasks add "Tidy my desk"
```

### List Tasks
List only uncompleted tasks:
```
$ tasks list
```
List all tasks (including completed):
```
$ tasks list --all
$ tasks list -a
```

### Complete a Task
```
$ tasks complete <taskid>
```

### Delete a Task
```
$ tasks delete <taskid>
```

## Example Data File

A sample `tasks.csv` file:
```
ID,Description,CreatedAt,IsComplete
1,My new task,2024-07-27T16:45:19-05:00,true
2,Finish this video,2024-07-27T16:45:26-05:00,true
3,Find a video editor,2024-07-27T16:45:31-05:00,false
```

## Notable Packages Used
- [`encoding/csv`](https://pkg.go.dev/encoding/csv) for CSV file operations
- [`strconv`](https://pkg.go.dev/strconv) for string conversions
- [`text/tabwriter`](https://pkg.go.dev/text/tabwriter) for tabular output
- [`os`](https://pkg.go.dev/os) for file operations
- [`github.com/spf13/cobra`](https://github.com/spf13/cobra) for CLI
- [`github.com/mergestat/timediff`](https://github.com/mergestat/timediff) for friendly time differences

## Technical Considerations
- **File Locking:** Uses `syscall.Flock` to prevent concurrent read/writes to the data file.
- **Error Handling:** Errors and diagnostics are written to stderr; output is written to stdout.

## License

MIT License © 2024 Sidney Thiel
