#!/usr/bin/env python3
"""Fetch a Jira issue from digital-and-technology.atlassian.net and print the
"Key details" field (Background Story + Acceptance Criteria) plus core metadata.

Auth comes from env vars JIRA_EMAIL and JIRA_TOKEN.

Usage:
    JIRA_EMAIL=... JIRA_TOKEN=... python3 fetch_jira.py <issue>

<issue> can be:
    - a full URL : https://digital-and-technology.atlassian.net/browse/RTRED-2612
    - an issue key: RTRED-2612
    - just a number: 2612   (default project prefix RTRED- is added)
"""
import os
import re
import sys
import html
import json
import base64
import urllib.request

BASE_URL = "https://digital-and-technology.atlassian.net"
DEFAULT_PROJECT = "RTRED"

# Known custom field for this Jira instance. The standard `description` field
# is usually empty here — the real ticket content lives in "Key details".
KEY_DETAILS_FIELDS = ["customfield_10119", "customfield_10121"]  # Key details: Background Story + AC


def resolve_key(arg):
    """Turn a URL / key / number into an issue key like RTRED-2612."""
    m = re.search(r"/browse/([A-Z][A-Z0-9]+-\d+)", arg)
    if m:
        return m.group(1)
    if re.fullmatch(r"[A-Z][A-Z0-9]+-\d+", arg):
        return arg
    if re.fullmatch(r"\d+", arg):
        return f"{DEFAULT_PROJECT}-{arg}"
    return arg


def html_to_text(raw):
    if not raw:
        return ""
    t = raw
    t = t.replace("<li>", "\n  - ").replace("</li>", "")
    t = t.replace("<p>", "\n").replace("</p>", "\n")
    t = re.sub(r"</?(ol|ul|div|table|tbody|tr)[^>]*>", "", t)
    t = re.sub(r"<t[dh][^>]*>", "\n  | ", t)
    t = re.sub(r"<[^>]+>", "", t)
    t = html.unescape(t)
    t = re.sub(r"\n{3,}", "\n\n", t)
    return t.strip()


def main():
    if len(sys.argv) < 2:
        sys.exit("usage: fetch_jira.py <issue-url|key|number>")

    email = os.environ.get("JIRA_EMAIL")
    token = os.environ.get("JIRA_TOKEN")
    if not email or not token:
        sys.exit("error: set JIRA_EMAIL and JIRA_TOKEN environment variables")

    key = resolve_key(sys.argv[1])
    fields = ",".join(
        ["summary", "issuetype", "status", "priority", "reporter", "assignee"]
        + KEY_DETAILS_FIELDS
    )
    url = f"{BASE_URL}/rest/api/3/issue/{key}?expand=renderedFields&fields={fields}"

    cred = base64.b64encode(f"{email}:{token}".encode()).decode()
    req = urllib.request.Request(url, headers={
        "Authorization": f"Basic {cred}",
        "Accept": "application/json",
    })
    try:
        with urllib.request.urlopen(req) as resp:
            data = json.load(resp)
    except urllib.error.HTTPError as e:
        sys.exit(f"error: HTTP {e.code} fetching {key} — {e.read().decode()[:200]}")

    f = data["fields"]
    rf = data.get("renderedFields", {})

    def name(field):
        v = f.get(field)
        return v.get("name") if isinstance(v, dict) and v else "-"

    print(f"# {key} — {f.get('summary', '')}")
    print()
    print(f"- Type     : {name('issuetype')}")
    print(f"- Status   : {name('status')}")
    print(f"- Priority : {name('priority')}")
    print(f"- Reporter : {f['reporter']['displayName'] if f.get('reporter') else '-'}")
    print(f"- Assignee : {f['assignee']['displayName'] if f.get('assignee') else '-'}")
    print(f"- Link     : {BASE_URL}/browse/{key}")
    print()
    print("## Key details")
    parts = [html_to_text(rf.get(field)) for field in KEY_DETAILS_FIELDS]
    body = "\n\n".join(p for p in parts if p)
    print(body if body else "(empty)")


if __name__ == "__main__":
    main()
