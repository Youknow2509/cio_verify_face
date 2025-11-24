
def getHash(value: str) -> str:
    """
        Get hash of a string value using SHA256.
    """
    import hashlib
    return hashlib.sha256(value.encode('utf-8')).hexdigest()