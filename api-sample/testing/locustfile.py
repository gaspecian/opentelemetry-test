from locust import HttpUser, task, between
import random
import string

class APIUser(HttpUser):
    wait_time = between(1, 3)
    
    def on_start(self):
        self.user_ids = []
    
    @task(3)
    def create_user(self):
        name = ''.join(random.choices(string.ascii_letters, k=8))
        email = f"{name.lower()}@example.com"
        
        response = self.client.post("/users", json={
            "name": name,
            "email": email
        })
        
        if response.status_code == 201:
            user_id = response.json().get("id")
            if user_id:
                self.user_ids.append(user_id)
    
    @task(5)
    def list_users(self):
        self.client.get("/users")
    
    @task(4)
    def get_user(self):
        if self.user_ids:
            user_id = random.choice(self.user_ids)
            self.client.get(f"/users/{user_id}")
    
    @task(2)
    def update_user(self):
        if self.user_ids:
            user_id = random.choice(self.user_ids)
            name = ''.join(random.choices(string.ascii_letters, k=8))
            email = f"{name.lower()}@example.com"
            
            self.client.put(f"/users/{user_id}", json={
                "name": name,
                "email": email
            })
    
    @task(1)
    def delete_user(self):
        if self.user_ids:
            user_id = self.user_ids.pop(random.randrange(len(self.user_ids)))
            self.client.delete(f"/users/{user_id}")
