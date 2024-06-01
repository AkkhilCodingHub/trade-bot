from kiteconnect import KiteConnect
import os

def kite_init():
    api_key = os.getenv("api_key")
    api_secret = os.getenv("api_secret")
    kite = KiteConnect(api_key=api_key)

    # Generate session by visiting the URL
    print(kite.login_url())
