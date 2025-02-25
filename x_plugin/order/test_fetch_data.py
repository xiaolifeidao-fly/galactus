import time
import requests


def fetch_data():
    try:
        url = "http://111.180.188.251:8091/dy/video/get?videoId=7464134195720277248"
        bgResponse = requests.get(url)
        task_data = bgResponse.json()
        data_status = task_data['data']['dataStatus']
        if data_status == 'SUCCESS':
            return True
        else:
            return False
    except Exception as e:
        print(e)
        return False
error_num = 11
success_num = 109
while True:
    result = fetch_data()
    if result:
        success_num += 1
    else:
        error_num += 1
    print(f"success_num: {success_num}, error_num: {error_num}")
    time.sleep(1)