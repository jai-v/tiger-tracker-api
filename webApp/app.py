import hashlib
import json
from urllib.parse import quote
import os

import flask
from flask import Flask, request, redirect , make_response, render_template
import requests
import pkce

app = Flask(__name__)

CLIENT_ID = os.environ.get("HYDRA_CLIENT_ID")
CLIENT_SECRET = os.environ.get("HYDRA_CLIENT_SECRET")
WEBAPP_HOST = os.environ.get("WEBAPP_HOST")
WEBAPP_PORT = os.environ.get("WEBAPP_PORT")
REDIRECT_URI = f"http://{WEBAPP_HOST}:{WEBAPP_PORT}/callback"
HYDRA_HOST = os.environ.get("HYDRA_HOST")
HYDRA_PUBLIC_PORT = os.environ.get("HYDRA_PUBLIC_PORT")
HYDRA_ADMIN_PORT = os.environ.get("HYDRA_ADMIN_PORT")
ENCODED_REDIRECT_URI = quote(REDIRECT_URI)

code_challenge_verifier_map = {}
session = {}


@app.route('/')
def index():
    code_verifier = pkce.generate_code_verifier(length=128)
    code_challenge = pkce.get_code_challenge(code_verifier)
    code_challenge_verifier_map[code_challenge] = code_verifier
    state = hashlib.sha256(os.urandom(1024)).hexdigest()
    session[state] = {}
    payload = {
        "oauth_init_url": f'http://{HYDRA_HOST}:{HYDRA_PUBLIC_PORT}/oauth2/auth?client_id={CLIENT_ID}&redirect_uri={ENCODED_REDIRECT_URI}&response_type=code&scope=offline&state={state}',
    }
    print(HYDRA_HOST, HYDRA_PUBLIC_PORT)
    return flask.render_template("welcome.html", **payload)


@app.route('/authentication/login')
def login():
    request_args = request.args.to_dict()
    payload = {
        "login_url": "http://localhost:8080/api/tiger-tracker/v1/login",
        "error": request_args.get("error"),
        "login_challenge": request_args.get("login_challenge")
    }
    return flask.render_template("login.html", **payload)


@app.route('/authentication/consent', methods=["GET", "POST"])
def consent():
    if request.method == "POST":
        consent_challenge = request.form.get("consent_challenge").strip()
    else:
        request_args = request.args.to_dict()
        consent_challenge = request_args.get("consent_challenge").strip()
    if consent_challenge is None:
        print("error: consent challenge missing in the request param")
        return

    resp = requests.get(f"http://hydra:{HYDRA_ADMIN_PORT}/oauth2/auth/requests/consent", {"consent_challenge": consent_challenge})
    statusCode = resp.status_code
    consentDetails = resp.json()
    if statusCode != 200:
        print("error: could not fetch consent request for challenge")
        return
    if True or consentDetails.get("skip") or request.method == "POST":  # skipping the consent
        putData = {}
        if consentDetails.get("requested_access_token_audience"):
            putData["requested_access_token_audience"] = consentDetails["requested_access_token_audience"]
        if consentDetails.get("requested_scope"):
            putData["grant_scope"] = consentDetails["requested_scope"]
        putData["session"] = {"access_token": None, "id_token": None}
        resp = requests.put(url=f"http://hydra:{HYDRA_ADMIN_PORT}/oauth2/auth/requests/consent/accept",
                            data=json.dumps(putData), params={"consent_challenge": consent_challenge},
                            cookies=request.cookies)
        acceptedConsent = resp.json()
        if resp.status_code != 200:
            print("error: could not accept consent request for challenge, response:", acceptedConsent)
            return
        redirectTo = acceptedConsent.get("redirect_to")
        return redirect(redirectTo, 302)

    clientName = consentDetails.get("client", {}).get("client_name", {})
    consentMessage = f"Application {clientName} wants access resources on your behalf and to",

    payload = {
        "consent_challenge": consent_challenge,
        "consent_message": consentMessage,
        "requested_scopes": consentDetails.get("requested_scope", [])
    }
    return flask.render_template("consent.html", **payload)


@app.route('/signup')
def signup():
    payload = {
        "signup_url": "http://localhost:8080/api/tiger-tracker/v1/signup"
    }
    return flask.render_template("signup.html", **payload)


@app.route('/callback')
def callback():
    query_params = request.args.to_dict()
    code = query_params.get("code")
    state = query_params.get("state")
    data = {"grant_type": "authorization_code", "code": code, "client_id": CLIENT_ID, "client_secret": CLIENT_SECRET,
            "redirect_uri": REDIRECT_URI}
    tokenResp = requests.post(f"http://hydra:{HYDRA_PUBLIC_PORT}/oauth2/token", data=data)
    tokenDetails = tokenResp.json()
    if tokenResp.status_code != 200:
        print("failed to get token, resp:", tokenDetails)
        return

    response = make_response(render_template("callback.html", data=tokenDetails))
    response.set_cookie("access_token", tokenDetails.get("access_token"))
    return response


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5001)
