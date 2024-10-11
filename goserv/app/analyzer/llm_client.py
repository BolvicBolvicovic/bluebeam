import openai
from openai import OpenAI
import sys
import json

question_json = sys.argv[1]
question_data = json.loads(question_json)

system_message = question_data.get('systemmessage')
data = question_data.get('data')
feature = question_data.get('feature')
client = OpenAI()
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
                "content": f"Data: {data}, Feature: {feature}"
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
                        "insight": {
                            "description": "Give a very short (1 line) insight on the feature and its presence or no presence, if it has a positve or negative impact on the webpage",
                            "type": "string"
                        }
                    }
                }
            }
        }
    )
    
    print(response.choices[0].message.content);
except openai.APIConnectionError as e:
    print("The server could not be reached")
    print(e.__cause__)  # an underlying Exception, likely raised within httpx.
except openai.RateLimitError as e:
    print("A 429 status code was received; we should back off a bit.")
except openai.APIStatusError as e:
    print("Another non-200-range status code was received")
    print(e.status_code)
    print(e.response)
