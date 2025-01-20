import time
import requests

def get_task_xiongmao():
    url = "http://112.74.176.127:8020/studio/api/task/get?key=AKID5af7bb633ff38158e8587d0b8cfcbfa6&type=dz&uid=4358060915864643&sec_uid=MS4wLjABAAAAIL3rVTcL76atlBoKdLiExLSD6qLxEaq0xbN2jGy8u6vysn1IPW7qgR5gKxmGf1lA&filter=video&platform=dy"
    bgResponse = requests.get(url)
    task_data = bgResponse.json()
    print(task_data)
    if task_data["success"] and task_data["code"] == 0:
        return task_data["data"]
    else:
        return None

def submit_order_xiongmao(order_no, share_url):
    params = {
        "orderNo": order_no,
        "totalNum": 100,
        "businessKey": share_url,
        "encryptionKey": "00400B1DCFCA851CB76DE97AE9D51321",
        "shopKey": "F6AF352FF95B6C8E1C8BBBB7652F43F2",
        "userName": "test_admin"
    }
    url = "http://111.180.188.251:20030/kakrolot_web/orders/submit"
    bgResponse = requests.post(url, json=params)
    task_data = bgResponse.json()
    print("submit order result ", task_data)

for i in range(100):
    task = get_task_xiongmao()
    if task:
        submit_order_xiongmao(task["studiotask_id"], task["params"]["share_url"])
    else:
        print("get task failed")
    time.sleep(2)
