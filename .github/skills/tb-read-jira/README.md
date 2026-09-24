# tb-read-jira

ดึงและสรุป Jira ticket จาก `digital-and-technology.atlassian.net` โดยอัตโนมัติ

## ทำอะไร

- ดึง metadata (Type, Status, Priority, Reporter, Assignee)
- อ่าน **Key details** (`customfield_10119`) — Background Story + Acceptance Criteria
- ดึงและอ่าน image attachments
- สรุปเป็นภาษามนุษย์

## Trigger เมื่อไหร่

Claude จะ activate skill นี้อัตโนมัติเมื่อ:

- พิมพ์ key เช่น `RTRED-2612`
- วาง link เช่น `https://digital-and-technology.atlassian.net/browse/RTRED-2612`
- ถามว่า "card นี้ทำอะไร", "assign ให้ใคร", "สรุป ticket ให้หน่อย"

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-read-jira` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-read-jira`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-read-jira`

```bash
# macOS
cp -R tb-read-jira ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-read-jira ~\.copilot\skills\
```

จากนั้นตั้ง env vars (ใส่ใน `~/.zshrc` หรือ `~/.bashrc`):

```bash
export JIRA_EMAIL="your-email@example.com"
export JIRA_TOKEN="your-jira-api-token"
```

### ขอ Jira API Token

1. ไปที่ https://id.atlassian.com/manage-profile/security/api-tokens
2. กด **Create API token**
3. Copy token มาใส่ `JIRA_TOKEN`

## Requirements

- Python 3 (built-in บน macOS)
- `JIRA_EMAIL` และ `JIRA_TOKEN` env vars
- GitHub Copilot (VS Code)
