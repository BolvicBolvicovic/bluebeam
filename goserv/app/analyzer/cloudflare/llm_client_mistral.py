import os
import requests
import sys
import json

question_json = sys.argv[1]
with open(question_json, 'r') as f:
    question_data = json.load(f)

system_message = question_data.get('systemmessage')
data = question_data.get('data')
feature = question_data.get('feature')
content = json.dumps(data)

ACCOUNT_ID = os.environ["CLOUDFLARE_ACCOUNT_ID"]
AUTH_TOKEN = os.environ["CLOUDFLARE_API_KEY"]

prompt = "Tell me all about PEP-8"
response = requests.post(
  f"https://api.cloudflare.com/client/v4/accounts/{ACCOUNT_ID}/ai/run/@hf/mistral/mistral-7b-instruct-v0.2",
    headers={"Authorization": f"Bearer {AUTH_TOKEN}"},
    json={
      "messages": [
        {"role": "system", "content": "You are a friendly assistant"},
        {"role": "user", "content": prompt}
      ]
    }
)
result = response.json()
print(result)
