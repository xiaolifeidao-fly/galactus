

import time
import requests


def get_task():
    url = "http://task.06km.com/task-api/s/task/get?uid=7437422433697317946&businessId=1&secUid=MS4wLjABAAAAkJxYlT9MWT1Zdpc_2btsI6EYbcY9sJ5Ns5H6QOQALbVuU2lXmeNb_QbXmJm7Qq1P&shieldSoStatus=0&onlyTw=2"
    headers = {
        "u_key":"o13bv4tepkoilskg",
        "s_key":"0cebqqn0ft7w6fh8"
    }
    bgResponse = requests.get(url, headers=headers)
    task_data = bgResponse.json()
    print(task_data)
    if task_data["code"] == '200':
        return task_data["data"]
    else:
        return None


def submit_order(order_no, task_url):
    params = {
        "orderNo": order_no,
        "totalNum": 100,
        "businessKey": task_url,
        "encryptionKey": "00400B1DCFCA851CB76DE97AE9D51321",
        "shopKey": "F6AF352FF95B6C8E1C8BBBB7652F43F2",
        "userName": "test_admin"
    }
    url = "http://111.180.188.251:20030/kakrolot_web/orders/submit"
    bgResponse = requests.post(url, json=params)
    task_data = bgResponse.json()
    print("submit order result ", task_data)
for i in range(100):
    task = get_task()
    if task:
        submit_order(task["orderId"], task["taskUrl"])
    else:
        print("get task failed")
    time.sleep(2)
