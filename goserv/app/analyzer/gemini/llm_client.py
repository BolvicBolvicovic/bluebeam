import google.generativeai as genai
import google.api_core.exceptions as exceptions
import os
import sys
import json
import datetime
import time
from typing_extensions import TypedDict

class ResponseFormat(TypedDict):
    feature_name: str
    is_feature_present: bool
    text_if_feature_present: str
    url_of_page_text_if_feature_present: str
    thought_process: str

question_json = sys.argv[1]
with open(question_json, 'r') as f:
    question_data = json.load(f)

system_message = question_data.get('systemmessage')
data = question_data.get('data')
feature = question_data.get('feature')
content = json.dumps(data)

genai.configure(api_key=os.environ["GEMINI_API_KEY"])

model = genai.GenerativeModel(
        model_name=sys.argv[2],
        system_instruction=system_message,
)

try:
    result = model.generate_content(
        [f"data: {data}", f"feature: {feature}"],
        generation_config=genai.GenerationConfig(
            temperature=0.1,
            top_k=10,
            top_p=0.1,
            max_output_tokens=4096,
            response_mime_type="application/json",
            response_schema=ResponseFormat
        ),
    )
    
    if result and result.candidates[0]:
        parts_text = result.candidates[0].content.parts[0].text
        print(parts_text)
    else:
        print("error: No content part [0] text for candidates [0]")
except exceptions.ResourceExhausted:
    print("error: Resource has been exhausted (e.g. check quota)")
