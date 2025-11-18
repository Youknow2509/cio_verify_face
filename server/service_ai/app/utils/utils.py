
# Get token from Authorization header
from http.client import HTTPException
from typing import Optional
from urllib.request import Request


def get_token_from_header(request: Request) -> Optional[str]:
    auth_header = request.headers.get("Authorization")
    if not auth_header:
        raise HTTPException(status_code=401, detail="Missing Authorization")
    token = auth_header.split(" ", 1)[1] if " " in auth_header else auth_header
    return token