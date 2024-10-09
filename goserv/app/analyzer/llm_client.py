from openai import OpenAI
import sys
import json

question_json = sys.argv[1]
question_data = json.loads(question_json)

system_message = question_data.get('systemmessage')
data = question_data.get('data')
feature = question_data.get('feature')
client = OpenAI()

response = client.chat.completions.create(
    model="gpt-3.5-turbo",
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
                        "description": "If the feature is present then extract the concise part of the data that shows the feature is present else leave an empty string",
                        "type": "string"
                    },
                    "insight": {
                        "description": "Give a very short (1 line) insight on the feature and its presence or no presence, if it has a positve or negative impact",
                        "type": "string"
                    }
                }
            }
        }
    }
)

print(response.choices[0].message.content);
