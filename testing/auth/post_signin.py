from posting import Posting
import httpx


def on_response(response: httpx.Response, posting: Posting):
    if response.status_code != 200:
        posting.notify(title="Error", message=f"Failed to sign in: {response.status_code}", severity="error")
        return False

    data = response.json()
    if "error" in data:
        print(f"Error: {data['error']}")
        posting.notify(title="Error", message=data["error"], severity="error")
        return False

    posting.notify(title="Sign In Successful", message="You have successfully signed in!", severity="information")

    if "token" in data:
        posting.set_variable("auth_token", data["token"])
        posting.notify(
            title="Token Received",
            message="Your authentication token has been received and stored.",
            severity="information"
        )
    
