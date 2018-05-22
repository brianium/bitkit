"""check the status of the bitkit api"""
# pylint: disable=broad-except
import os
import json
import requests

def slacker(channel, text):
    """post text content to slack channel"""
    webhook_url = os.environ['SLACK']
    slack_data = {"channel": channel, "text": text}
    response = requests.post(
        webhook_url, data=json.dumps(slack_data),
        headers={'Content-Type': 'application/json'}
    )
    return response.text

def call_api(uri):
    """make a simple get request"""
    try:
        response = requests.get(uri, timeout=5)
        if response.status_code != 200:
            message = 'bitkit api returned status {}'.format(response.status_code)
            slacker('#dev', message)
    except Exception as err:
        slacker('#dev', str(err))

if __name__ == "__main__":
    URI = os.environ['BITKIT'].rstrip('/') + '/random'
    call_api(URI)
