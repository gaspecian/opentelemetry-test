from locust import HttpUser, task, between
import json
import random

class TodoUser(HttpUser):
    wait_time = between(1, 3)
    
    def on_start(self):
        """Initialize with some todos"""
        self.created_todos = []
        for i in range(3):
            response = self.client.post("/todos", json={
                "title": f"Initial todo {i}",
                "done": False
            })
            if response.status_code == 201:
                self.created_todos.append(response.json()["id"])

    @task(3)
    def get_all_todos(self):
        """Get all todos - most frequent operation"""
        self.client.get("/todos")

    @task(2)
    def create_todo(self):
        """Create a new todo"""
        response = self.client.post("/todos", json={
            "title": f"Todo {random.randint(1, 1000)}",
            "done": random.choice([True, False])
        })
        if response.status_code == 201:
            self.created_todos.append(response.json()["id"])

    @task(2)
    def get_specific_todo(self):
        """Get a specific todo"""
        if self.created_todos:
            todo_id = random.choice(self.created_todos)
            self.client.get(f"/todos/{todo_id}")

    @task(1)
    def update_todo(self):
        """Update an existing todo"""
        if self.created_todos:
            todo_id = random.choice(self.created_todos)
            self.client.put(f"/todos/{todo_id}", json={
                "title": f"Updated todo {random.randint(1, 1000)}",
                "done": random.choice([True, False])
            })

    @task(1)
    def delete_todo(self):
        """Delete a todo"""
        if self.created_todos:
            todo_id = self.created_todos.pop(random.randint(0, len(self.created_todos) - 1))
            self.client.delete(f"/todos/{todo_id}")
