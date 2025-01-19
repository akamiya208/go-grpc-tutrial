```bash
$ grpcurl -plaintext -d '{"name": "taskName", "description": "taskDescription"}' localhost:8080 proto.TaskService/CreateTask
{
  "id": 1,
  "name": "taskName",
  "description": "taskDescription"
}

$ grpcurl -plaintext -d '{"id": 1}' localhost:8080 proto.TaskService/GetTask
{
  "id": 1,
  "name": "taskName",
  "description": "taskDescription"
}

$ grpcurl -plaintext -d '{"name": "taskName"}' localhost:8080 proto.TaskService/GetTasks
{
  "tasks": [
    {
      "id": 1,
      "name": "taskName",
      "description": "taskDescription"
    }
  ]
}

$ grpcurl -plaintext -d '{"id": 1, "name": "taskName", "description": "updateTaskDescription"}' localhost:8080 proto.TaskService/UpdateTask
{
  "id": 1,
  "name": "taskName",
  "description": "updateTaskDescription"
}

$ grpcurl -plaintext -d '{"id": 1}' localhost:8080 proto.TaskService/DeleteTask
{
  "id": 1,
  "name": "taskName",
  "description": "updateTaskDescription"
}

$ grpcurl -plaintext -d '{"id": 1}' localhost:8080 proto.TaskService/GetTask
ERROR:
  Code: Unknown
  Message: record not found
```