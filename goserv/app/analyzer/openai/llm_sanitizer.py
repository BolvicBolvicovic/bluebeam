import openai
from openai import OpenAI
import sys
import json
import time

question_json = sys.argv[1]
with open(question_json, 'r') as f:
    question_data = json.load(f)

system_message = question_data.get('systemmessage')
crawledwebsites = question_data.get('crawledwebsites')
ainame = question_data.get('ainame')
sanitizer = question_data.get('sanitizer')

content = f"data: {crawledwebsites}\nsanitizer specification: {sanitizer}\n"
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
            model=ainame,
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
                    "name": "audit_output",
                    "schema": {
                        "type": "array",
                        "items": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "pageurl": { "type": "string" },
                                    "pagecontent": {
                                        "type": "object",
                                        "properties": {
                                            "links": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "href": { "type": "string" },
                                                        "text": { "type": "string" }
                                                    }
                                                }
                                            },
                                            "buttons": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "text":     { "type": "string" },
                                                        "onclick":  { "type": "string" },
                                                        "id":       { "type": "string" },
                                                        "classes":  { "type": "string" }
                                                    }
                                                }
                                            },
                                            "images": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "alt":      { "type": "string" },
                                                        "src":      { "type": "string" },
                                                        "classes":  { "type": "string" }
                                                    }
                                                }
                                            },
                                            "formInputs": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "type":     { "type": "string" },
                                                        "name":     { "type": "string" },
                                                        "value":    { "type": "string" }
                                                    }
                                                }
                                            },
                                            "metaTags": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "name":     { "type": "string" },
                                                        "content":  { "type": "string" }
                                                    }
                                                }
                                            },
                                            "headers": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "tag":      { "type": "string" },
                                                        "text":     { "type": "string" }
                                                    }
                                                }
                                            },
                                            "bodyText": { "type": "string" }
                                        }
                                    },
                                    "required": [ "pageurl", "pagecontent" ]
                                }
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

getResponse()
