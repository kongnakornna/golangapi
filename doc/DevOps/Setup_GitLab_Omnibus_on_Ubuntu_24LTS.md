## Setup GitLab Omnibus บน Ubuntu 24.04 LTS

### สรุปภาพรวมสถาปัตยกรรม GitLab Self-Managed (Omnibus)

┌──────────────────────────────────────────┐
│                  Nginx                   │
│        (reverse proxy / TLS terminator)  │
└───────────────────┬──────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────┐
│          GitLab Workhorse (Go)           │
│ - Proxy for Rails / Gitaly / Registry    │
└───────────────────┬──────────────────────┘
                    │
                    ▼
┌──────────────────────────────────────────┐
│          Ruby on Rails (Main App)        │
│ - API / Web / Sidekiq                    │
└──────┬───────────────────────────────────┘
       │                         │
       ▼                         ▼
┌───────────────────┐   ┌──────────────────┐
│ PostgreSQL        │   │ Redis            │
│ (database)        │   │ (cache/jobs)     │
└───────────────────┘   └──────────────────┘
       │
       ▼
┌──────────────┐
│ Gitaly (Go)  │ ←→ Git repositories
└──────────────┘

> 📦 ทั้งหมดนี้รวมอยู่ในแพ็กเกจ “GitLab Omnibus”
ซึ่งติดตั้งผ่าน apt ได้ตัวเดียว (gitlab-ce หรือ gitlab-ee)
และภายในมีทุก service ข้างต้นครบหมด

> “GitLab ถูกพัฒนาด้วย Ruby on Rails เป็นแกนหลักของระบบ web application
ใช้ PostgreSQL เป็นฐานข้อมูลหลัก และ Redis สำหรับ cache กับ job queue
โดยใช้ Gitaly (ภาษา Go) เป็น backend สำหรับจัดการ repository
และ GitLab Runner (ภาษา Go) สำหรับระบบ CI/CD”

### Stack หลักที่ GitLab ใช้งานจริง
| Layer                                 | เทคโนโลยีที่ใช้                           | หน้าที่หลัก                                                                             |
| :------------------------------------ | :---------------------------------------- | :-------------------------------------------------------------------------------------- |
| **Frontend (Web UI)**                 | Vue.js (ต่อยอดจาก HAML/ERB เดิมของ Rails) | แสดงหน้าเว็บ, Project Dashboard, Merge Request UI ฯลฯ                                   |
| **Backend (Core API)**                | **Ruby on Rails 7.x**                     | ระบบหลักทั้งหมด — users, projects, issues, pipelines, webhooks, permissions             |
| **Database**                          | **PostgreSQL 16.x**                       | เก็บข้อมูลหลักทั้งหมด (users, repos metadata, merge requests, pipelines logs index ฯลฯ) |
| **Cache / Queue**                     | **Redis 7.x**                             | ใช้เก็บ session, background job queue (Sidekiq)                                         |
| **Background Worker**                 | **Sidekiq (Ruby)**                        | ประมวลผลงานเบื้องหลัง เช่น email, CI job scheduling, indexing                           |
| **Git Repository Storage**            | **Gitaly (Go)**                           | จัดการ Git operations (clone, fetch, push) แบบ distributed                              |
| **CI/CD Execution**                   | **GitLab Runner (Go)**                    | รัน pipeline jobs ในเครื่องหรือ Docker                                                  |
| **Proxy / API Gateway**               | **GitLab Workhorse (Go)**                 | จัดการ HTTP requests, file uploads, Git over HTTP                                       |
| **Metrics / Monitoring**              | Prometheus + node_exporter                | เก็บและดู performance metrics                                                           |
| **Mail Delivery (optional)**          | SMTP (เช่น Gmail, Postfix)                | ส่งอีเมลแจ้งเตือนจากระบบ                                                                |
| **Registry (optional)**               | Docker Registry (Go-based)                | เก็บ container images ภายในองค์กร                                                       |
| **Pages / Static Hosting (optional)** | GitLab Pages Daemon (Go)                  | โฮสต์ static websites จาก repo                                                          |


### ข้อกำหนดเบื้องต้น:
- เครื่อง Ubuntu 24.04 LTS (สามารถใช้ VM หรือ Cloud VM ได้)
- สิทธิ์ root หรือ sudo
- ทรัพยากรเครื่องขั้นต่ำ:
  - CPU: 4 cores
  - RAM: 8 GB
  - Storage: 20 GB
- Network: 1 Gbps
- Pupblic IP Address: 128.199.230.251
- port ที่ต้องเปิด: 80 (HTTP), 443 (HTTPS), 22 (SSH), 5050 (GitLab Container Registry)
- domain name (optional): เช่น gitlab.aibisec.com, registry.aibisec.com

### สเปกเซิร์ฟเวอร์แนะนำสำหรับทีมขนาด ~100+ คน ที่มี repo จำนวนมาก, CI/CD, Container Registry ค่อนข้างหนัก

- CPU: 8 – 16 คอร์ (16 threads) เพื่อรองรับผู้ใช้พร้อมกันและงาน CI/CD หนัก
- RAM: 16 – 32 GB
- Storage: SSD/NVMe สำหรับ OS + GitLab data + Registry data
- Network: Bandwidth ดีพอสำหรับ push/pull images และ CI/CD artifacts
- Back-up / Redundancy: ควรวางแผน backup และกรณีเครื่องเสีย (อาจใช้ secondary mirror/Georeplicationในอนาคต)

> หมายเหตุ: เอกสารบอกว่า “สำหรับ ~1000 ผู้ใช้ ถ้า 8 vCPU + 16 GB น่าจะใช้ได้” ดังนั้นสำหรับ ~100 + ผู้ใช้ แนะนำให้มากกว่านั้นเพื่อความเสถียร

### 0. จดโดเมนเนมที่ name.com และ ชี้ DNS
- จดโดเมนเนมที่ name.com (เช่น aibisec.com)
- เพิ่ม A Record ชื่อ gitlab ชี้มาที่ 128.199.230.251
- เพิ่ม A Record ชื่อ registry ชี้มาที่ 128.199.230.251

### 1. เตรียมเครื่องและ OS

- สร้าง droplet Ubuntu 24.04 LTS บน DigitalOcean หรือใช้ VM/Cloud VM ที่มีอยู่
- ssh เข้าไปยังเครื่องด้วยสิทธิ์ root หรือ sudo user

```bash
ssh root@128.199.230.251
```
- อัพเดตระบบและติดตั้งแพ็กเกจที่จำเป็น:

```bash
sudo apt update && sudo apt upgrade -y
```
- reboot เครื่องถ้ามีการอัพเดตเคอร์เนล:

```bash
sudo reboot
```
- แสดงข้อมูล Kernel/สถาปัตยกรรมอย่างย่อ
```bash
uname -a
```

- แสดงรายละเอียดเวอร์ชัน Ubuntu (ต้องมีแพ็กเกจ lsb-release)
```bash
lsb_release -a
```

- แสดงรายละเอียดเครื่อง โฮสต์เนม OS และ Kernel ในมุมมองที่อ่านง่าย
```bash
hostnamectl
```

- แสดงรายละเอียด OS จากไฟล์มาตรฐาน
```bash
cat /etc/os-release
```

### 2. ตรวจสอบก่อนเริ่ม
- ยืนยันว่า DNS ชี้มาถูกแล้ว (ควรได้ 128.199.230.251)

```bash
# Linux/Mac
host gitlab.aibisec.com
host registry.aibisec.com

# Windows (PowerShell)
Resolve-DnsName gitlab.aibisec.com
Resolve-DnsName registry.aibisec.com

# หรือ CMD
nslookup gitlab.aibisec.com
nslookup registry.aibisec.com
```
ถ้ายังไม่ชี้ถูกต้อง รอ DNS propagate สักพัก (อาจใช้เวลาถึง 24-48 ชั่วโมง)

- ตั้งชื่อเครื่องและ /etc/hosts ให้แมตช์โดเมน
```bash
sudo hostnamectl set-hostname gitlab.aibisec.com
```
แก้ไขไฟล์ /etc/hosts:
```bash
sudo nano /etc/hosts
```
เพิ่มบรรทัด:
```
128.199.230.251 gitlab.aibisec.com gitlab
128.199.230.251 registry.aibisec.com
```
- เปิดพอร์ต (ทั้งบนเครื่อง และใน DigitalOcean Firewall ถ้าใช้)
```bash
sudo apt update && sudo apt install -y ufw
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP (สำหรับ Let's Encrypt)
sudo ufw allow 443/tcp  # HTTPS
sudo ufw allow 5050/tcp # Registry (ถ้าเปิดใช้งาน)
sudo ufw enable
sudo ufw status
```

### 3. ติดตั้ง GitLab CE (Omnibus)
> เริ่มแบบ HTTP ก่อน เพื่อเลี่ยงปัญหา SSL ตอนแรก แล้วค่อยเปิด HTTPS + Let’s Encrypt ภายหลัง
```bash
# พื้นฐาน
sudo apt install -y curl openssh-server ca-certificates tzdata perl

# เพิ่ม repo GitLab และติดตั้ง
curl -sS https://packages.gitlab.com/install/repositories/gitlab/gitlab-ce/script.deb.sh | sudo bash
sudo EXTERNAL_URL="http://gitlab.aibisec.com" apt install -y gitlab-ce

# คอนฟิกครั้งแรก
sudo gitlab-ctl reconfigure
sudo gitlab-ctl status
```

ทดสอบเปิดเว็บ: http://gitlab.aibisec.com
เข้าสู่ระบบด้วยผู้ใช้ root แล้วดูรหัสผ่านเริ่มต้น:
```bash
sudo cat /etc/gitlab/initial_root_password
```

**คำแนะนำ: เปลี่ยนรหัสผ่าน root ทันทีหลังล็อกอินครั้งแรก**
**เปลี่ยนรหัสผ่าน root ใน GitLab**
- เปิดเว็บเบราว์เซอร์ไปที่ http://gitlab.aibisec.com
- ล็อกอินด้วยผู้ใช้ `root` และรหัสผ่านที่ได้จากคำสั่งข้างต้น
- เข้าเมนูเปลี่ยนรหัสผ่าน และตั้งรหัสผ่านใหม่
- ลบไฟล์รหัสผ่านทิ้งเองเลย ไม่ต้องรอ 24 ชม
```bash
sudo rm /etc/gitlab/initial_root_password
```

**ตั้งค่า SMTP เพื่อเปิดระบบ reset password ได้ในอนาคต (กรณีลืมรหัส root)**
แก้ไขไฟล์คอนฟิก GitLab:
```bash
sudo nano /etc/gitlab/gitlab.rb
```
เพิ่มบรรทัดเหล่านี้ (แก้ไขตาม SMTP ที่ใช้งาน):
```
gitlab_rails['smtp_enable'] = true
gitlab_rails['smtp_address'] = "smtp.gmail.com"
gitlab_rails['smtp_port'] = 587
gitlab_rails['smtp_user_name'] = "no-reply@aibisec.com"
gitlab_rails['smtp_password'] = "รหัสผ่านแอปพลิเคชัน"
gitlab_rails['smtp_domain'] = "aibisec.com"
gitlab_rails['smtp_authentication'] = "login"
gitlab_rails['smtp_enable_starttls_auto'] = true
gitlab_rails['gitlab_email_from'] = "no-reply@aibisec.com"
gitlab_rails['gitlab_email_display_name'] = "AIBISEC GitLab"
gitlab_rails['gitlab_email_reply_to'] = "no-reply@aibisec.com"
```
บันทึกไฟล์แล้วรันคำสั่ง:
```bash
sudo gitlab-ctl reconfigure
```

**สร้างผู้ใช้ admin เพิ่มอย่างน้อย 1 บัญชี (กันลืมรหัส root)**
- ไปที่เมนู Admin Area (ไอคอนรูปประแจมุมขวาบน)
- เลือก "Users" แล้วคลิก "New User"
- กรอกข้อมูลผู้ใช้ใหม่ ตั้งเป็น Admin ด้วย

**หากพบปัญหาการใช้งาน**
```bash
sudo gitlab-ctl tail        # ดู log แบบ real-time
sudo gitlab-ctl restart     # รีสตาร์ท GitLab services
sudo gitlab-ctl reconfigure # รันคอนฟิกใหม่
sudo gitlab-ctl status      # ตรวจสอบสถานะ services
```


### 4. ตั้งค่า HTTPS ด้วย Let’s Encrypt
เมื่อยืนยันว่าเว็บเข้าได้แล้ว (และ DNS/พอร์ต 80, 443 เปิด)

แก้ไฟล์ config ให้เป็น HTTPS:
เปิดไฟล์ config
```bash
sudo nano /etc/gitlab/gitlab.rb
```
เพิ่มบรรทัดเหล่านี้ (หรือแก้ไขบรรทัดที่มีอยู่แล้ว)
```
external_url 'https://gitlab.aibisec.com'
```
> หมายเหตุ: ต้องใช้ https:// จริง ๆ (มี s) เพราะระบบจะตรวจชนิด protocol จากตรงนี้


```bash
# ปรับค่าคอนฟิก GitLab ให้ใช้ HTTPS และเปิด Let’s Encrypt
sudo sed -i 's|external_url "http://gitlab.aibisec.com"|external_url "https://gitlab.aibisec.com"|' /etc/gitlab/gitlab.rb

sudo bash -lc 'printf "\nletsencrypt[\"enable\"] = true\nletsencrypt[\"auto_renew\"] = true\nletsencrypt[\"contact_emails\"] = [\"admin@aibisec.com\"]\n" >> /etc/gitlab/gitlab.rb'


# ใช้งานคอนฟิก
sudo gitlab-ctl reconfigure
sudo gitlab-ctl restart

# ตรวจสอบพอร์ต 443 ว่าเปิดใช้งานแล้ว
sudo ss -tulpen | grep -E ':443'
```
ทดสอบเปิดเว็บ: https://gitlab.aibisec.com
ตรวจสอบว่า HTTPS ใช้งานได้และมีใบรับรองถูกต้อง


### 5. ตั้งค่า GitLab Container Registry (ถ้าต้องการใช้)
แก้ไขไฟล์ config:
```bash
sudo vi /etc/gitlab/gitlab.rb
```
เพิ่มบรรทัดเหล่านี้:
```
registry_external_url 'https://registry.aibisec.com'
gitlab_rails['registry_enabled'] = true
gitlab_rails['registry_host'] = "registry.aibisec.com"
gitlab_rails['registry_port'] = "5050"
gitlab_rails['registry_api_url'] = "http://localhost:5050"
```

รันคำสั่ง:
```bash
sudo gitlab-ctl reconfigure
sudo gitlab-ctl restart
```

- ทดสอบเปิดเว็บ: https://registry.aibisec.com จะได้ 404 (ถือว่าปกติ เพราะไม่มี UI)
- ใน GitLab (UI) ไปที่ Project → Deploy → Container Registry ควรเห็น endpoint ของโปรเจ็กต์

**ทดสอบ login/push จากเครื่อง developer**
```bash
docker login registry.aibisec.com
# ใส่ username + Personal Access Token/Deploy Token (มี scope write_registry)

# ตัวอย่าง push
docker build -t registry.aibisec.com/<group>/<project>:test .
docker push registry.aibisec.com/<group>/<project>:test
```
> ถ้าจะล้างพื้นที่อัตโนมัติ แนะนำเปิด Cleanup Policies ในหน้า Project → Settings → Packages & Registries → Cleanup policies

### 6. ติดตั้ง GitLab Runner (เครื่องเดิมหรือเครื่องแยกก็ได้)
> แนะนำแยกเครื่อง/VM ถ้าโหลดงานเยอะ แต่เริ่มที่เครื่องเดียวกันได้
```bash
# ติดตั้ง runner
curl -L --output gitlab-runner.deb https://gitlab-runner-downloads.s3.amazonaws.com/latest/deb/gitlab-runner_amd64.deb
sudo dpkg -i gitlab-runner.deb

# ลงทะเบียน runner (เอา token จาก GitLab: Admin → Runners หรือ Group/Project → Runners)
sudo gitlab-runner register \
  --url https://gitlab.aibisec.com/ \
  --registration-token <YOUR_TOKEN> \
  --executor docker \
  --docker-image docker:25 \
  --description "docker-runner" \
  --tag-list "docker,linux" \
  --run-untagged="true" \
  --locked="false"

sudo systemctl enable gitlab-runner && sudo systemctl restart gitlab-runner
sudo gitlab-runner status
```
### 7. ตัวอย่างไฟล์ .gitlab-ci.yml
> ตัวอย่างนี้ build แล้ว push image ไปที่ Registry ภายใน GitLab
```yaml
stages:
  - lint
  - test
  - build
  - deploy

variables:
  GIT_DEPTH: "1"
  DOCKER_BUILDKIT: "1"

lint:
  stage: lint
  image: node:22
  script:
    - npm ci
    - npm run lint
  rules:
    - changes:
        - "**/*.js"
        - "**/*.ts"

test:
  stage: test
  image: node:22
  needs: [lint]
  script:
    - npm ci
    - npm test
  artifacts:
    when: always
    reports:
      junit: reports/junit.xml

build_image:
  stage: build
  needs: [test]
  image: docker:25
  services:
    - name: docker:25-dind
      command: ["--tls=false"]
  variables:
    DOCKER_HOST: tcp://docker:2375
    DOCKER_DRIVER: overlay2
  script:
    - echo $CI_REGISTRY
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA
  rules:
    - if: $CI_COMMIT_BRANCH

deploy_prod:
  stage: deploy
  needs: [build_image]
  script:
    - ./scripts/deploy.sh
  environment:
    name: production
    url: https://app.aibisec.com
  rules:
    - if: $CI_COMMIT_TAG
```

### 8. สำรองข้อมูล GitLab
ตั้งค่าการสำรองข้อมูลอัตโนมัติในไฟล์ config:
```bash
sudo vi /etc/gitlab/gitlab.rb
```
เพิ่มบรรทัดเหล่านี้:
```
gitlab_rails['backup_path'] = "/var/opt/gitlab/backups"
gitlab_rails['backup_archive_permissions'] = 0644
gitlab_rails['backup_keep_time'] = 604800  # เก็บสำรองข้อมูล 7 วัน
```
รันคำสั่ง:
```bash
sudo gitlab-ctl reconfigure
```
สร้างสคริปต์สำรองข้อมูลรายวัน (เพิ่มใน cron)
```bash
sudo vi /etc/cron.d/gitlab-backup
```
เพิ่มบรรทัดนี้ (สำรองทุกวันตี 1)
```
0 1 * * * root /usr/bin/gitlab-backup create CRON=1
```
บันทึกไฟล์แล้วออก

ทดสอบรันสำรองข้อมูลด้วยตนเอง:
```bash
sudo gitlab-backup create
```

### 9. การบำรุงรักษาและอัพเกรด GitLab
- ตรวจสอบสถานะ GitLab:
```bash
sudo gitlab-ctl status
```
- ดู log เมื่อมีปัญหา:
```bash
sudo gitlab-ctl tail
```
- อัพเกรด GitLab:
```bash
sudo apt update
sudo apt install -y gitlab-ce
sudo gitlab-ctl reconfigure
```
### 10. การตั้งค่าเพิ่มเติมที่แนะนำ
- ตั้งค่า 2FA สำหรับผู้ใช้ทุกคน
- ตั้งค่า SSO (LDAP/Active Directory) ถ้ามี
- ตั้งค่า Webhooks สำหรับการแจ้งเตือน
- ตั้งค่า CI/CD Templates สำหรับโปรเจ็กต์ต่าง ๆ

### Appendix: คำสั่งพื้นฐานสำหรับจัดการ GitLab Omnibus
```bash
# ตรวจสอบสถานะ services
sudo gitlab-ctl status
# รีสตาร์ท GitLab services
sudo gitlab-ctl restart
# ดู log แบบ real-time
sudo gitlab-ctl tail
# รันคอนฟิกใหม่
sudo gitlab-ctl reconfigure
# ดูเวอร์ชัน GitLab
sudo gitlab-rake gitlab:env:info
# สร้างสำรองข้อมูล
sudo gitlab-backup create
# คืนค่าจากสำรองข้อมูล
sudo gitlab-backup restore BACKUP=timestamp_of_backup
```
