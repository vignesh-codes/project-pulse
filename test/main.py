import requests
import json
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

def send_request(user_id, event_name, unix_ts):
    url = 'http://localhost:8080/v1/internal/log/'
    headers = {
        'x-client-id': 'id',
        'x-client-secret': 'secret',
        'Content-Type': 'application/json'
    }
    payload = {
        'user_id': user_id,
        'event_name': event_name,
        'unix_ts': unix_ts
    }
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    print(f"Response: {response.status_code} - User ID: {user_id} - Event: {event_name}")

def generate_data():
    user_ids = range(1, 5001)  
    event_names = ['login', 'logout']
    current_timestamp = int(time.time())

    data = []
    for user_id in user_ids:
        for event_name in event_names:
            unix_ts = current_timestamp + user_id 
            data.append((user_id, event_name, unix_ts))
    
    return data

def send_concurrent_requests():
    data = generate_data()
    workers = 1000
    entries_per_batch = 10000
    delay = 1 / workers

    with ThreadPoolExecutor(max_workers=workers) as executor:
        for i in range(0, entries_per_batch, workers):
            batch = data[i:i + workers]
            futures = [executor.submit(send_request, user_id, event_name, unix_ts) for user_id, event_name, unix_ts in batch]

            for future in as_completed(futures):
                try:
                    future.result()
                except Exception as e:
                    print(f"An error occurred: {str(e)}")

            time.sleep(delay)

send_concurrent_requests()
