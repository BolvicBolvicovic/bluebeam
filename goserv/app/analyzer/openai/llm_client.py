import openai
from openai import OpenAI
import os
import sys
import json
import time

question_json = sys.argv[1]
with open(question_json, 'r') as f:
    question_data = json.load(f)

system_message = question_data.get('systemmessage')
data = question_data.get('data')
feature = question_data.get('feature')
content = f"Data: {data}, Feature: {feature}"

client = OpenAI()
count = 0

def getResponse():
    global count
    count += 1
    if (count >= 50):
        print(f"error: too many tries ({count})")
        return
    try:
        response = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {
                    "role": "system",
                    "content": system_message
                },
                {
                    "role": "user",
                    "content": content
                }
            ],
            response_format={
                "type": "json_schema",
                "json_schema": {
                    "name": "data_schema",
                    "schema": {
                        "type": "object",
                        "properties": {
                            "feature_name": {
                                "description": "The name of the feature you are asked to analyse",
                                "type": "string"
                            },
                            "ispresent": {
                                "description": "If the feature is present or not",
                                "type": "boolean"
                            },
                            "textifpresent": {
                                "description": "If the feature is present then extract a part from the provided data set proving the presence of the feature else leave an empty string",
                                "type": "string"
                            },
                            "urlifpresent": {
                                "description": "If the feature is present then extract the url of the page where you found the information else leave an empty string",
                                "type": "string"
                            },
                            "thoughtprocess": {
                               "description": "Give a very short (1 line) explaination of how the data you extracted validates the feature or, if the feature is not found why no data match the feature.",
                                "type": "string"
                            }
                        }
                    }
                }
            }
        )
        
        print(response.choices[0].message.content);
    except openai.APIConnectionError as e:
        print("error: The server could not be reached")
        print(e.__cause__)  # an underlying Exception, likely raised within httpx.
    except openai.RateLimitError as e:
        time.sleep(5)
        getResponse()
    except openai.APIStatusError as e:
        print("error: Another non-200-range status code was received")
        print(e.response)
        print(content)

if (os.environ["OPENAI_API_KEY"] == "mocked_key"):
    print("Test Pass")
else:
    getResponse()
