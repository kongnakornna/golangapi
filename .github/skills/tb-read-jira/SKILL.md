---
name: tb-read-jira
description: Use this whenever a message contains an "RTRED-" key (e.g. RTRED-2612), a digital-and-technology.atlassian.net link, or a bare issue number, AND the user wants to know what it says — reading, opening, summarizing, checking the assignee/status, or pulling the background story or acceptance criteria. This is the default skill for any "what is this card/link about", "read/อ่าน/เปิด/ดึงข้อมูล/สรุป RTRED-xxxx", "ลิงก์นี้เกี่ยวกับอะไร", or "assign ให้ใคร" request, even when the words "Jira" or "ticket" never appear. The user's intent is simply to understand the contents of one of these issues; that is enough to trigger. Skip only when the user wants to create, edit, update, or change the status of an issue, when the link points to a different Atlassian site, or for Trello boards and local files.
---

# Read Jira

Fetch a Jira issue from `digital-and-technology.atlassian.net` and summarize it.

## Why this skill exists

Two traps make naive Jira reads fail on this instance:

1. **Private URLs** — fetching the `/browse/...` page just returns a login wall. You must use the REST API with auth.
2. **The `description` field is empty.** This org puts the real ticket content in a custom field called **"Key details"** (`customfield_10119`), which holds the Background Story and Acceptance Criteria. If you read `description` you'll wrongly conclude the card is blank. **Read Key details, not description.**

## Auth

Auth uses the user's Jira email + API token via env vars `JIRA_EMAIL` and `JIRA_TOKEN`.

- If they're already set in the environment, use them.
- If not, ask the user for their Jira email and API token. They can set them in one line in the prompt, e.g.:
  ```
  ! export JIRA_EMAIL="you@example.com" && export JIRA_TOKEN="<token>"
  ```
  Note that env vars do not persist across separate Bash calls, so run the export together with the fetch in a single command.

## How to fetch

Run the bundled script. It accepts a full URL, an issue key, or just a number:

```bash
export JIRA_EMAIL="..." && export JIRA_TOKEN="..." && \
python3 ~/.copilot/skills/tb-read-jira/scripts/fetch_jira.py RTRED-2612
```

Accepted argument forms (all resolve to the issue key):
- Full URL: `https://digital-and-technology.atlassian.net/browse/RTRED-2612`
- Issue key: `RTRED-2612`
- Bare number: `2612` → defaults to project prefix `RTRED-`

The script reads from the hardcoded instance `digital-and-technology.atlassian.net`, pulls `expand=renderedFields`, converts the HTML to readable text, and prints core metadata plus **Key details**.

If you ever need a field the script doesn't cover, fall back to a direct call and inspect `renderedFields`:
```bash
curl -s -u "$JIRA_EMAIL:$JIRA_TOKEN" -H "Accept: application/json" \
  "https://digital-and-technology.atlassian.net/rest/api/3/issue/<KEY>?expand=renderedFields"
```

## Image Attachments

After fetching Key details, always check for image attachments and read them:

```bash
# 1. List attachments
curl -s -u "$JIRA_EMAIL:$JIRA_TOKEN" -H "Accept: application/json" \
  "https://digital-and-technology.atlassian.net/rest/api/3/issue/<KEY>?fields=attachment" | \
  python3 -c "import json,sys; data=json.load(sys.stdin); attachments=data.get('fields',{}).get('attachment',[]); [print(f\"name: {a['filename']}\nurl: {a['content']}\nmimeType: {a['mimeType']}\n\") for a in attachments] if attachments else print('No attachments found')"

# 2. Download each image attachment
curl -s -u "$JIRA_EMAIL:$JIRA_TOKEN" -L "<attachment-content-url>" -o /tmp/jira_<KEY>_<n>.png
```

Then use the Read tool on each downloaded image file — Claude can read images visually and extract context from screenshots, tables, and diagrams.

Include what you see in the images as part of your summary. Images often contain the real problem description when Key details is empty.

## Output

Always present the output in **two parts, in this order**:

### Part 1 — Raw content (show everything verbatim)

Show all raw text exactly as received, before any interpretation:

1. **Metadata block** — Key, Type, Status, Priority, Reporter, Assignee, Link (as a code block or table)
2. **Key details** — the full text of `customfield_10119` verbatim, including every line, table row, and bullet
3. **Image attachments** — for each image, state the filename and describe every piece of text, data, and visual content you can read from it (column headers, row values, error messages, chat bubbles, etc.)

Do not skip, truncate, or paraphrase anything in Part 1. The user needs to see the raw source before your interpretation.

### Part 2 — Summary

After the raw content, provide a plain-language summary with:
- A short metadata table (Key, Type, Status, Priority, Reporter, Assignee)
- The **Background Story** and **Acceptance Criteria** from Key details
- Key observations from any image attachments
- A plain-language interpretation of what the card is asking for

If Key details is genuinely empty, say so — but only after confirming the field itself was empty, never because `description` was null.
