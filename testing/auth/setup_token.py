from posting import Posting, RequestModel, Auth, Header
import httpx

def on_request(request: RequestModel, posting: Posting):
    token = posting.get_variable("auth_token")
    if not token:
        posting.notify(title="Error", message="No authentication token found.", severity="error")
        return False
    request.headers.append(Header(name="Authorization", value=f"Bearer {token}"))
