from posting import Posting, Header, RequestModel

def on_request(request: RequestModel, posting: Posting) -> None:
    auth_token = posting.get_variable("auth_token")
    if auth_token:
        request.headers.append(Header(name="Authorization", value=f"Bearer {auth_token}"))
        posting.notify("Added auth token to request")
    else:
        posting.notify("No auth token found")
