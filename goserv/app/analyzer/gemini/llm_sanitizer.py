import google.generativeai as genai
import google.api_core.exceptions as exceptions
import os
import sys
import json
import datetime
import time
from typing_extensions import TypedDict, List

question_json = sys.argv[1]
with open(question_json, 'r') as f:
    question_data = json.load(f)

system_message = question_data.get('systemmessage')
crawledwebsites = question_data.get('crawledwebsites')
ainame = question_data.get('ainame')
sanitizer = question_data.get('sanitizer')

content = f"data: {crawledwebsites}\nsanitizer specification: {sanitizer}\n"

genai.configure(api_key=os.environ["GEMINI_API_KEY"])

model = genai.GenerativeModel(
        model_name=ainame,
        system_instruction=system_message,
)

try:
    result = model.generate_content(
        content,
        generation_config=genai.GenerationConfig(
            temperature=0.1,
            top_k=10,
            top_p=0.1,
            response_mime_type="application/json",
            response_schema={
                "type": "ARRAY",
                "items": {
                    "type": "ARRAY",
                    "items": {
                        "type": "OBJECT",
                        "properties": {
                            "pageurl": { "type": "STRING" },
                            "pagecontent": {
                                "type": "OBJECT",
                                "properties": {
                                    "links": {
                                        "type": "ARRAY",
                                        "items": {
                                            "type": "OBJECT",
                                            "properties": {
                                                "href": { "type": "STRING" },
                                                "text": { "type": "STRING" }
                                            }
                                        }
                                    },
                                    "buttons": {
                                        "type": "ARRAY",
                                        "items": {
                                            "type": "OBJECT",
                                            "properties": {
                                                "text":     { "type": "STRING" },
                                                "onclick":  { "type": "STRING" },
                                                "id":       { "type": "STRING" },
                                                "classes":  { "type": "STRING" }
                                            }
                                        }
                                    },
                                    "images": {
                                        "type": "ARRAY",
                                        "items": {
                                            "type": "OBJECT",
                                            "properties": {
                                                "alt":      { "type": "STRING" },
                                                "src":      { "type": "STRING" },
                                                "classes":  { "type": "STRING" }
                                            }
                                        }
                                    },
                                    "formInputs": {
                                        "type": "ARRAY",
                                        "items": {
                                            "type": "OBJECT",
                                            "properties": {
                                                "type":     { "type": "STRING" },
                                                "name":     { "type": "STRING" },
                                                "value":    { "type": "STRING" }
                                            }
                                        }
                                    },
                                    "metaTags": {
                                        "type": "ARRAY",
                                        "items": {
                                            "type": "OBJECT",
                                            "properties": {
                                                "name":     { "type": "STRING" },
                                                "content":  { "type": "STRING" }
                                            }
                                        }
                                    },
                                    "headers": {
                                        "type": "ARRAY",
                                        "items": {
                                            "type": "OBJECT",
                                            "properties": {
                                                "tag":      { "type": "STRING" },
                                                "text":     { "type": "STRING" }
                                            }
                                        }
                                    },
                                    "bodyText": { "type": "STRING" }
                                }
                            }
                        }
                    }
                }
            }
        ),
    )
    
    if result and result.candidates[0]:
        parts_text = result.candidates[0].content.parts[0].text
        print(parts_text)
    else:
        print("error: No content part [0] text for candidates [0]")
except exceptions.ResourceExhausted:
    print("error: Resource has been exhausted (e.g. check quota)")
