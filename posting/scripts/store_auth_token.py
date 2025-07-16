from posting import Posting
import httpx

def on_response(response: httpx.Response, posting: Posting):
    if response.status_code == 200:
        posting.auth_token = response.json()["token"]
        posting.set_variable("auth_token", posting.auth_token)
        posting.notify("Auth token stored")
    else:
        posting.notify(f"Error storing auth token: {response.text}")
