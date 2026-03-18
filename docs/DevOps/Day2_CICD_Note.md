## DevOps Jenkins & GitLab Actions & N8N - Day 2

[คลิกที่นี่เพื่อดาวน์โหลดเอกสารประกอบการอบรม](https://bit.ly/devops_easybuy)

### 📋 สารบัญ
1. [พื้นฐานการใช้งาน Jenkins](#พื้นฐานการใช้งาน-jenkins)
2. [การสร้าง Pipeline ด้วย Jenkins](#การสร้าง-pipeline-ด้วย-jenkins)
3. [การเชื่อมต่อ Jenkins กับ GitLab](#การเชื่อมต่อ-jenkins-กับ-gitlab)
4. [การตั้งค่า Webhook ใน GitLab เพื่อ Trigger Jenkins Job](#การตั้งค่า-webhook-ใน-gitlab-เพื่อ-trigger-jenkins-job)
5. [Jenkins multibranch pipeline](#jenkins-multibranch-pipeline)
6. [Jenkins on Ubuntu Server with Docker](#jenkins-on-ubuntu-server-with-docker)
7. [N8N on Ubuntu Server with Docker](#n8n-on-ubuntu-server-with-docker)
8. [Jenkins CI/CD deployed to Server with Docker and SSH](#jenkins-cicd-deployed-to-server-with-docker-and-ssh)


## พื้นฐาน CI/CD ด้วย Jenkins
- การติดตั้ง Jenkins แบบ Bare Metal Installation
- การติดตั้ง Jenkins แบบ Containerized Deployment
- การตั้งค่า Jenkins เบื้องต้น
- การสร้าง Pipeline แบบพื้นฐาน

#### 1. การติดตั้ง Jenkins แบบ Bare Metal Installation
- ดาวน์โหลดและติดตั้ง Jenkins จาก [https://www.jenkins.io/download/](https://www.jenkins.io/download/)
- ติดตั้ง Java JDK 17 หรือ 21 (แนะนำให้ใช้เวอร์ชันล่าสุด)

- เปิดเว็บเบราว์เซอร์และเข้าไปที่ [http://localhost:8080](http://localhost:8080)
- ทำตามขั้นตอนการตั้งค่าเบื้องต้น

#### 2. การติดตั้ง Jenkins แบบ Containerized Deployment
- สร้าง Dockerfile สำหรับ Jenkins
```Dockerfile
# เริ่มต้นจาก Image Jenkins ที่เราต้องการ
FROM jenkins/jenkins:jdk21

USER root

# ติดตั้งเฉพาะ Docker CLI เท่านั้น
RUN apt-get update && \
    apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && \
    apt-get install -y docker-ce-cli && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# (ส่วน Entrypoint script ยังคงเดิม เพื่อจัดการสิทธิ์ docker.sock)
RUN echo '#!/bin/bash\n\
DOCKER_SOCK="/var/run/docker.sock"\n\
if [ -S "$DOCKER_SOCK" ]; then\n\
    DOCKER_GID=$(stat -c "%g" $DOCKER_SOCK)\n\
    if ! getent group $DOCKER_GID > /dev/null 2>&1; then\n\
        groupadd -g $DOCKER_GID docker\n\
    fi\n\
    usermod -aG $DOCKER_GID jenkins\n\
fi\n\
exec /usr/bin/tini -- /usr/local/bin/jenkins.sh "$@"' > /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh

USER jenkins

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
```

- สร้าง docker-compose.yml สำหรับ Jenkins
```yaml
# Define Network
networks:
  jenkins_network:
    name: jenkins_network
    driver: bridge

# Define Services
services:
  jenkins:
    build: .
    image: jenkins-with-docker:jdk21
    container_name: jenkins
    user: root
    volumes:
      - ./jenkins_home:/var/jenkins_home # เพื่อให้ Jenkins สามารถเก็บข้อมูลไว้ใน host
      - /var/run/docker.sock:/var/run/docker.sock # เพื่อให้ Jenkins สามารถใช้งาน Docker daemon ที่รันบน host ได้
    environment:
      - JENKINS_OPTS=--httpPort=8800 # กำหนด Port สำหรับ Jenkins UI
    ports:
      - "8800:8800" # สำหรับ Jenkins UI
    restart: always
    networks:
      - jenkins_network
```
- รัน Jenkins ด้วย Docker Compose
```bash
docker-compose up -d --build
```
- เปิดเว็บเบราว์เซอร์และเข้าไปที่ [http://localhost:8800](http://localhost:8800)
- ทำตามขั้นตอนการตั้งค่าเบื้องต้น

#### 3. รู้จัก Jenkins เบื้องต้น

- Jenkins เป็นเครื่องมือที่ช่วยในการทำ CI/CD โดยอัตโนมัติ
- CI/CD คือแนวทางการพัฒนาซอฟต์แวร์ที่ช่วยให้การพัฒนาและการนำส่งซอฟต์แวร์เป็นไปอย่างรวดเร็วและมีประสิทธิภาพ
- Jenkins รองรับการทำงานร่วมกับเครื่องมือและเทคโนโลยีต่าง ๆ เช่น Git, Docker, Kubernetes เป็นต้น
- Jenkins มีระบบ Plugin ที่ช่วยเพิ่มความสามารถและฟีเจอร์ต่าง ๆ ให้กับ Jenkins
- Jenkins สามารถทำงานร่วมกับระบบ Version Control เช่น GitHub, GitLab, Bitbucket เป็นต้น
- Jenkins รองรับการสร้าง Pipeline ที่ช่วยในการจัดการกระบวนการ CI/CD ได้อย่างมีประสิทธิภาพ
- Jenkins มีระบบการแจ้งเตือนผ่านทาง Email, Slack, Microsoft Teams เป็นต้น
- Jenkins รองรับการทำงานร่วมกับ Docker เพื่อสร้างสภาพแวดล้อมที่เหมาะสมสำหรับการทดสอบและการนำส่งซอฟต์แวร์
- Jenkins มีระบบการจัดการผู้ใช้งานและสิทธิ์การเข้าถึงที่ช่วยให้การบริหารจัดการ Jenkins เป็นไปอย่างมีประสิทธิภาพ
- Jenkins รองรับการทำงานร่วมกับ Cloud Providers เช่น AWS, Azure, Google Cloud เป็นต้น
- Jenkins มีระบบการสำรองข้อมูลและการกู้คืนข้อมูลที่ช่วยให้การบริหารจัดการ Jenkins เป็นไปอย่างมีประสิทธิภาพ

---

## What is Jenkins Pipeline ?

![Jenkins Pipeline](https://kubedemy.io/wp-content/uploads/2023/06/4418c3cd93a28e984510f8d25a6fd815.png)

Jenkins Pipeline (หรือเรียกสั้นๆ ว่า "Pipeline" โดยใช้ตัวพิมพ์ใหญ่ "P") เป็นชุดของปลั๊กอินที่สนับสนุนการนำเสนอกระบวนการและการผนวกรวม pipeline สำหรับ continuous delivery เข้าไปใน Jenkins

Continuous delivery (CD) pipeline เป็นการแสดงกระบวนการอัตโนมัติที่ช่วยนำซอฟต์แวร์จากระบบควบคุมเวอร์ชันไปจนถึงมือผู้ใช้และลูกค้าของคุณ ทุกครั้งที่มีการเปลี่ยนแปลงซอฟต์แวร์ (ที่ถูก commit ในระบบควบคุมเวอร์ชัน) จะต้องผ่านกระบวนการที่ซับซ้อนเพื่อให้สามารถนำออกใช้งานได้ กระบวนการนี้รวมถึงการสร้างซอฟต์แวร์ในลักษณะที่เชื่อถือได้และสามารถทำซ้ำได้ รวมถึงการนำซอฟต์แวร์ที่สร้างเสร็จแล้ว (เรียกว่า "build") ผ่านหลายขั้นตอนของการทดสอบและการปรับใช้

Pipeline ช่วยให้คุณสามารถกำหนดกระบวนการเหล่านี้ในรูปแบบที่เป็นโค้ด ซึ่งทำให้สามารถจัดการและปรับปรุงกระบวนการได้ง่ายขึ้น นอกจากนี้ยังช่วยให้สามารถทำงานร่วมกับทีมพัฒนาและทีมปฏิบัติการได้อย่างมีประสิทธิภาพมากขึ้น

## ประเภทของ Jenkins Pipeline

![Jenkins Pipeline Types](https://kouzie.github.io/assets/cicd/jenkins1.png)

การเขียน Pipeline ใน Jenkins สามารถทำได้หลัก ๆ 2 แบบ คือ Declarative Pipeline และ Scripted Pipeline โดยแต่ละแบบมีลักษณะและการใช้งานที่แตกต่างกัน ดังนี้

![Jenkins Pipeline Types](https://www.lambdatest.com/blog/wp-content/uploads/2021/04/Screenshot-2021-02-13-at-10.05.27-AM.png)

1. Declarative Pipeline: เป็นรูปแบบที่มีโครงสร้างชัดเจนและง่ายต่อการอ่านและเขียน โดยใช้คำสั่งที่กำหนดไว้ล่วงหน้า เช่น `pipeline`, `stage`, `steps` เป็นต้น

![Jenkins Declarative Pipeline](https://devops.com/wp-content/uploads/2018/07/Jenkinspic4-1.png)

2. Scripted Pipeline: เป็นรูปแบบที่มีความยืดหยุ่นมากกว่า Declarative Pipeline โดยใช้ Groovy เป็นภาษาหลักในการเขียนโค้ด ซึ่งช่วยให้สามารถเขียนโค้ดที่ซับซ้อนและมีความยืดหยุ่นมากขึ้น

## การสร้าง Pipeline ด้วย Jenkins
1. ติดตั้ง Jenkins และ Plugins ที่จำเป็น
2. สร้าง Jenkins Job ใหม่
3. เลือกประเภทของ Job เป็น "Pipeline"
4. กำหนดชื่อ Job และตั้งค่าอื่น ๆ ตามต้องการ

## ตัวอย่างการสร้าง Jenkinsfile แบบ Declarative Pipeline

```groovy
pipeline {
    agent any // ใช้ agent ใดก็ได้
    stages {
        stage('Build') {
            steps {
                echo 'Building...'
                // คำสั่งสำหรับการ build เช่น การ compile โค้ด
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
                // คำสั่งสำหรับการทดสอบ เช่น การรัน unit tests
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying...'
                // คำสั่งสำหรับการ deploy เช่น การส่งโค้ดไปยังเซิร์ฟเวอร์
            }
        }
    }
}

```
อธิบายโค้ด:
- `pipeline { ... }`: กำหนดว่าเป็น Jenkins Pipeline
- `agent any`: กำหนดให้ Jenkins ใช้ agent ใดก็ได้ในการรัน Pipeline
- `stages { ... }`: กำหนดขั้นตอนต่าง ๆ ใน Pipeline
- `stage('Build') { ... }`: กำหนดขั้นตอน "Build"
- `steps { ... }`: กำหนดคำสั่งที่ต้องการรันในแต่ละขั้นตอน

> agent คือเครื่องมือหรือสภาพแวดล้อมที่ Jenkins ใช้ในการรัน Pipeline เช่น Docker, Kubernetes, หรือเครื่องเซิร์ฟเวอร์ที่ติดตั้ง Jenkins

> ภาษา Groovy ที่ใช้ใน Jenkinsfile สามารถใช้ได้ทั้งแบบ Declarative และ Scripted Pipeline ขึ้นอยู่กับความต้องการและความซับซ้อนของกระบวนการ CI/CD ที่ต้องการสร้าง

> Jenkinsfile สามารถเก็บไว้ในระบบควบคุมเวอร์ชัน เช่น GitHub เพื่อให้สามารถจัดการและติดตามการเปลี่ยนแปลงได้ง่ายขึ้น

## ตัวอย่างการสร้าง Jenkinsfile แบบ Scripted Pipeline

```groovy
node('built-in') { // ระบุ agent ที่ต้องการใช้เป็น built-in
    stage('Build') {
        echo 'Building...'
        // คำสั่งสำหรับการ build เช่น การ compile โค้ด
    }
    stage('Test') {
        echo 'Testing...'
        // คำสั่งสำหรับการทดสอบ เช่น การรัน unit tests
    }
    stage('Deploy') {
        echo 'Deploying...'
        // คำสั่งสำหรับการ deploy เช่น การส่งโค้ดไปยังเซิร์ฟเวอร์
    }
}
```
อธิบายโค้ด:
- `node { ... }`: กำหนดว่าเป็น Scripted Pipeline
- `stage('Build') { ... }`: กำหนดขั้นตอน "Build"
- `echo 'Building...'`: คำสั่งสำหรับการแสดงข้อความใน Console Output
- คำสั่งสำหรับการ build, test, และ deploy สามารถปรับเปลี่ยนได้ตามความต้องการ

> แบบ Scripted Pipeline มีความยืดหยุ่นมากกว่า Declarative Pipeline แต่โค้ดอาจจะซับซ้อนและอ่านยากกว่า
> กำหนด agent ได้โดยใช้คำสั่ง `node('label') { ... }` เพื่อระบุ agent ที่ต้องการใช้
> built-in คือ label ของ agent ที่ต้องการใช้ในการรัน Pipeline

#### วิธีการตรวจสอบ Available Nodes/Agents:
1. ไปที่หน้า Jenkins Dashboard
2. คลิกที่ "Manage Jenkins"
3. คลิกที่ "Manage Nodes and Clouds"
4. จะเห็นรายการ Nodes/Agents ที่มีอยู่ในระบบ พร้อมกับ Labels ที่สามารถใช้ในการระบุ agent ใน Jenkinsfile

## การเลือกใช้ Declarative หรือ Scripted Pipeline
- Declarative Pipeline เหมาะสำหรับผู้เริ่มต้นและการสร้าง Pipeline ที่ไม่ซับซ้อน
- Scripted Pipeline เหมาะสำหรับผู้ที่มีประสบการณ์และต้องการความยืดหยุ่นมากขึ้นในการเขียนโค้ด
- สามารถผสมผสานการใช้ทั้งสองแบบในโปรเจคเดียวกันได้ตามความเหมาะสม
---

## การเชื่อมต่อ Jenkins กับ GitLab
- การติดตั้ง GitLab Plugin ใน Jenkins
- การสร้าง Jenkins Job ที่เชื่อมต่อกับ GitLab Repository
- การตั้งค่า Webhook ใน GitLab เพื่อ Trigger Jenkins Job

#### 1. การติดตั้ง GitLab Plugin ใน Jenkins
- ไปที่หน้า Jenkins Dashboard
- คลิกที่ "Manage Jenkins"
- คลิกที่ "Manage Plugins"
- ไปที่แท็บ "Available"
- ค้นหา "GitLab Plugin"
- เลือก Plugin และคลิกที่ "Install without restart"

#### 2. สร้าง Personal Access Token (PAT) ใน GitLab
- ไปที่หน้า GitLab
- คลิกที่รูปโปรไฟล์ของคุณที่มุมขวาบน
- เลือก "Preferences"
- ในเมนูด้านซ้าย เลือก "Personal access tokens"
- กรอกชื่อสำหรับ Token ชื่อ "Jenkins PAT"
- กรอกรายละเอียด Description (optional) "Token for Jenkins integration"
- ตั้งค่า Expiration date ตามต้องการ (แนะนำให้ตั้งเป็น 30 วัน หรือ 60 วัน)
- เลือก Scopes ที่ต้องการ: 
  - read_repository
  - write_repository
  - api
- คลิกที่ "Create token"
- คัดลอก Token ที่สร้างขึ้นมาเก็บไว้ใช้ใน Jenkins

#### 3. เพิ่ม Credentials ใน Jenkins โดยใช้ PAT ที่สร้างขึ้น
- กลับไปที่ Jenkins Dashboard
- คลิก **"Manage Jenkins"**
- เลือก **"Manage Credentials"**
- คลิกที่ **"(global)"** domain
- หรือ **"Global credentials (unrestricted)"**
- คลิก **"Add Credentials"** ที่มุมซ้าย
**Kind:** `Username with password`
**Scope:** `Global (Jenkins, nodes, items, all child items, etc)`
**Username:** ใส่ username ของ GitHub ของคุณ
**Password:** วาง Personal Access Token ที่สร้างไว้
**ID:** ตั้งชื่อที่จดจำง่าย เช่น `gitlab-pat-credentials`
**Description:** คำอธิบายเพิ่มเติม (optional) เช่น `GitLab Personal Access Token for Jenkins`
- คลิก **"OK"** เพื่อบันทึก

#### 4. การสร้าง Jenkins Job ที่เชื่อมต่อกับ GitLab Repository
- กลับไปที่ Jenkins Dashboard
- คลิกที่ **"New Item"**
- กรอกชื่อ Job เช่น `test-gitlab-connnection`
- เลือกประเภทเป็น **"Freestyle project"**
- คลิก **"OK"**
- ในหน้าการตั้งค่า Job:
  - ในส่วน **"Source Code Management"** เลือก **"Git"**
  - ในช่อง **"Repository URL"** ใส่ URL ของ GitLab Repository เช่น `https://gitlab.com/your-username/your-repo.git`
  - ในช่อง **"Credentials"** เลือก Credentials ที่เพิ่มไว้ก่อนหน้านี้ (เช่น `gitlab-pat-credentials`)
- ในส่วน **"Build Triggers"** เลือก **"Build when a change is pushed to GitLab"**
- ในส่วน **"Build"** เพิ่มขั้นตอนการ build ตามต้องการ เช่น การรันสคริปต์ shell
- คลิก **"Save"** เพื่อบันทึกการตั้งค่า Job

#### 5.การทำงานกับ Jenkinsfile

![Jenkinsfile](https://intellij-support.jetbrains.com/hc/user_images/7IueXvigXEGikkwJ-RGd0A.png)

> Jenkinsfile คือไฟล์ที่ใช้ในการกำหนดขั้นตอนการทำงานของ Jenkins Pipeline โดยใช้ภาษา Groovy ซึ่งช่วยให้สามารถจัดการกระบวนการ CI/CD ได้อย่างมีประสิทธิภาพ

### 1. โครงสร้างพื้นฐานของ Jenkinsfile

```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                echo 'Building...'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying...'
            }
        }
    }
}
```

### ตัวอย่าง Jenkinsfile สำหรับโปรเจ็กต์ทดสอบใน GitLab ก่อนหน้าที่สร้างไว้
สามารถสร้างไฟล์ชื่อ `Jenkinsfile` ใน root directory ของโปรเจ็กต์ โดยมีเนื้อหาดังนี้:

```groovy
pipeline {
    agent any
    stages {
        stage('Checkout') {
            steps {
                echo "Checking out code..."
                checkout scm // ดึงโค้ดจาก repository ที่เชื่อมต่อกับ Jenkins Job
            }
        }
        stage('Build') {
            steps {
                echo 'Building...'
                // คำสั่งสำหรับการ build เช่น การ compile โค้ด
            }
        }
    }
}
```

### 2. การใช้งาน Jenkinsfile ใน Jenkins Job
1. สร้าง Jenkins Job ใหม่ ชื่อ "test-jenkinsfile-gitlab"
2. เลือกประเภทของ Job เป็น "Pipeline"
3. ในส่วนของ Pipeline Definition เลือก "Pipeline script from SCM"
4. เลือก "Git" เป็น SCM
5. ใส่ Repository URL และเลือก Credentials ที่สร้างไว้
6. กำหนด Branch ที่ต้องการใช้ (เช่น `main` หรือ `master`)
7. กำหนด Script Path เป็น `Jenkinsfile` (หรือชื่อไฟล์อื่น ๆ ที่ใช้)
8. บันทึกการตั้งค่าและรัน Job

## การเพิ่ม Docker Credentials ใน Jenkins Secrets
1. ไปที่หน้า Jenkins Dashboard
2. คลิกที่ "Manage Jenkins"
3. เลือก "Manage Credentials"
4. เลือกโดเมนที่ต้องการเพิ่ม Credentials (เช่น Global)
5. คลิก "Add Credentials"
6. เลือกชนิดเป็น "Username with password"
7. กรอก Username ด้วย Docker Hub Username
8. กรอก Password ด้วย Docker Hub Access Token
9. ตั้ง ID เช่น `dockerhub-cred`
10. คลิก "OK" เพื่อบันทึก

## Workshop Express Docker Application

<img src="https://miro.medium.com/v2/1*Jr3NFSKTfQWRUyjblBSKeg.png" width="200">

โปรเจ็กต์ Express.js + TypeScript REST API ที่ทำงานใน Docker Container พร้อมด้วย CI/CD Pipeline ผ่าน Jenkins, Github Actions และการแจ้งเตือนผ่าน N8N


## 🏗️ Project Structure

```
express-docker-app/
├── 📁 src/
│   └── 📄 app.ts                   # Main Express application (TypeScript)
├── 📁 tests/
│   └── 📄 app.test.ts              # Jest test suite with Supertest
├── 📁 dist/                       # Compiled JavaScript output
│   └── 📄 app.js                  # Compiled application
├── 📁 node_modules/               # Node.js dependencies
├── 📁 .github/
│   └── 📁 workflows/
│       └── 📄 main.yml               # GitHub Actions workflow
├── 🐳 Dockerfile                  # Docker build configuration
├── 🔧 Jenkinsfile                 # Jenkins CI/CD pipeline
├── ⚙️ jest.config.js              # Jest testing configuration
├── 📄 package.json                # Node.js project configuration
├── 📄 package-lock.json           # Dependency lock file
├── 📄 tsconfig.json               # TypeScript configuration
└── 📖 README.md                   # Project documentation
```
## 🛠️ Prerequisites

- **Node.js 22+** (LTS recommended)
- **npm 10+** or **yarn**
- **TypeScript** (installed globally or via npx)
- **Docker & Docker Compose**
- **Git**
- **Jenkins** (for CI/CD)
- **SSH Access** (for deployment)

## ⚡ Quick Start

#### ขั้นตอนที่ 1: ดาวน์โหลดโปรเจ็กต์

ดาวน์โหลดโค้ดตัวอย่างได้ที่ลิงก์นี้:
[https://drive.google.com/file/d/1NPlqrDK0d9pX1egWKaq__3vybVHqZaW-/view?usp=sharing](https://drive.google.com/file/d/1NPlqrDK0d9pX1egWKaq__3vybVHqZaW-/view?usp=sharing)

#### ขั้นตอนที่ 2. สร้าง Jenkinsfile สำหรับ CI/CD Pipeline

```groovy
pipeline {
    // agent: กำหนด agent ที่จะใช้รัน Pipeline
    // any หมายถึง ใช้ agent ใดก็ได้ที่มีอยู่ใน Jenkins
    agent any

    // กัน “เช็คเอาต์ซ้ำซ้อน”
    // ถ้า job เป็นแบบ Pipeline from SCM / Multibranch แนะนำเพิ่ม options { skipDefaultCheckout(true) }
    // เพื่อปิดการ checkout อัตโนมัติก่อนเข้า stages (เพราะเรามี checkout scm อยู่แล้ว)
    options { 
        skipDefaultCheckout(true)   // ถ้าเป็น Pipeline from SCM/Multi-branch
    }

    // กำหนด environment variables สำหรับ Docker Hub credentials และ Docker repository
    environment {
        DOCKER_HUB_CREDENTIALS_ID = 'dockerhub-cred'
        DOCKER_REPO               = "your-dockerhub-username/express-docker-app"
        APP_NAME                  = "express-docker-app"
    }

    // กำหนด stages ของ Pipeline
    stages {

        // Stage 1: ดึงโค้ดล่าสุดจาก Git
        // ใช้ checkout scm หากใช้ Pipeline from SCM
        // หรือใช้ git url: 'https://gitlab.com/your-username/your-repo.git', branch: 'main', credentialsId: 'gitlab-pat-credentials'
        stage('Checkout') {
            steps {
                echo "Checking out code..."
                checkout scm
                // หรือใช้แบบกำหนดเอง หากไม่ใช้ Pipeline from SCM:
                // git url: 'https://gitlab.com/your-username/your-repo.git', branch: 'main', credentialsId: 'gitlab-pat-credentials'
            }
        }

        // Stage 2: ติดตั้ง dependencies และ Run test
        // ใช้ Node.js plugin (ต้องติดตั้ง NodeJS plugin ก่อน) ใน Jenkins หรือ Node.js ใน Docker 
        // ถ้ามี package-lock.json ให้ใช้ npm ci แทน npm install จะเร็วและล็อกเวอร์ชันชัดเจนกว่า
        stage('Install & Test') {
            steps {
                sh '''
                    if [ -f package-lock.json ]; then npm ci; else npm install; fi
                    npm test
                '''
            }
        }

        // Stage 3: สร้าง Docker Image
        // ใช้ Docker ที่ติดตั้งบน Jenkins agent (ต้องติดตั้ง Docker plugin ก่อน) ใน Jenkins หรือ Docker ใน Docker
        stage('Build Docker Image') {
            steps {
                sh ""
                    echo "Building Docker image: ${DOCKER_REPO}:${BUILD_NUMBER}"
                    docker build --target production -t ${DOCKER_REPO}:${BUILD_NUMBER} -t ${DOCKER_REPO}:latest .
                ""
            }
        }

        // Stage 4: Push Image ไปยัง Docker Hub
        // ใช้ docker.withRegistry() เพื่อความปลอดภัยและเรียบง่าย
        stage('Push Docker Image') {
            steps {
                script {
                    docker.withRegistry('https://index.docker.io/v1/', DOCKER_HUB_CREDENTIALS_ID) {
                        sh ""
                            echo "Pushing Docker image: ${DOCKER_REPO}:${BUILD_NUMBER} and ${DOCKER_REPO}:latest"
                            docker push ${DOCKER_REPO}:${BUILD_NUMBER}
                            docker push ${DOCKER_REPO}:latest
                        ""
                    }
                }
            }
        }

        // Stage 5: เคลียร์ Docker images บน agent
        // เพื่อประหยัดพื้นที่บน Jenkins agent หลังจาก push image ขึ้น Docker Hub แล้ว
        // ไม่จำเป็นต้องเก็บ image ไว้บน agent อีกต่อไป
        // หลักการทำงานคือ ลบ image ที่สร้างขึ้น (ทั้งแบบมี tag build number และ latest)
        // และลบ cache ที่ไม่จำเป็นออกไป
        stage('Cleanup Docker') {
            steps {
                sh ""
                    echo "Cleaning up local Docker images/cache on agent..."
                    docker image rm -f ${DOCKER_REPO}:${BUILD_NUMBER} || true
                    docker image rm -f ${DOCKER_REPO}:latest || true
                    docker image prune -af || true
                    docker builder prune -af || true
                ""
            }
        }

        // Stage 6: Deploy ไปยังเครื่อง local
        // ดึง image ล่าสุดจาก Docker Hub มาใช้งาน
        // หยุดและลบ container เก่าที่ชื่อ ${APP_NAME} (ถ้ามี)
        // สร้างและรัน container ใหม่จาก image ล่าสุด
        stage('Deploy Local') {
            steps {
                sh ""
                    echo "Deploying container ${APP_NAME} from latest image..."
                    docker pull ${DOCKER_REPO}:latest
                    docker stop ${APP_NAME} || true
                    docker rm ${APP_NAME} || true
                    docker run -d --name ${APP_NAME} -p 3000:3000 ${DOCKER_REPO}:latest
                    docker ps --filter name=${APP_NAME} --format "table {{.Names}}\\t{{.Image}}\\t{{.Status}}"
                ""
            }
        }
    }

    // กำหนด post actions
    // เช่น การแจ้งเตือนเมื่อ pipeline เสร็จสิ้น
    // สามารถเพิ่มการแจ้งเตือนผ่าน email, Slack, หรืออื่นๆ ได้ตามต้องการ
    post {
        always {
            echo "Pipeline finished with status: ${currentBuild.currentResult}"
        }
        success {
            echo "Pipeline succeeded!"
        }
        failure {
            echo "Pipeline failed!"
        }
    }

}
```

#### ขั้นตอนที่ 3. Push โค้ดขึ้น GitLab

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://gitlab.com/your-username/your-repo.git
git push -u origin main
```

#### ขั้นตอนที่ 4. สร้าง Jenkins Job แบบ Pipeline
4.1 ไปที่ Jenkins Dashboard
4.2 คลิกที่ "New Item"
4.3 ตั้งชื่อ Job และเลือก "Pipeline" จากนั้นคลิก "OK"
4.4 ในส่วน "Pipeline" ให้เลือก "Pipeline script from SCM"
4.5 ตั้งค่า SCM เป็น "Git" และกรอก URL ของ GitLab Repo
4.6 ตั้งค่า "Branch Specifier" เป็น "*/main"
4.7 ในส่วน "Script Path" ให้ระบุที่อยู่ของ Jenkinsfile (เช่น `Jenkinsfile`)
4.8 คลิก "Save"
4.9 คลิก "Build Now" เพื่อทดสอบการทำงานของ Pipeline

#### ขั้นตอนที่ 5. ติดตั้ง localtunnel และเปิดใช้งาน (ถ้ายังไม่มี)
- ติดตั้ง localtunnel
```bash
npm install -g localtunnel
```
- เปิดใช้งาน localtunnel เพื่อสร้าง public URL สำหรับ Jenkins Server
```bash
lt --port 8800
```

#### ขั้นตอนที่ 6. การตั้งค่า Webhook ใน GitLab
**จำเป็นต้องติดตั้ง plugin "GitLab Plugin, Git Plugin, GitLab API Plugin" ใน Jenkins ก่อน**

6.1 ไปที่หน้า GitLab Repository ของคุณ
6.2 คลิกที่ "Settings" > "Webhooks" > "Add webhook"
6.3 ในช่อง "URL" ให้ใส่ URL ของ Jenkins Server (เช่น `https://sweet-dolls-rush.loca.lt/project/test-jenkinsfile-gitlab`)
6.4 Trigger: เลือก "Push events" หรือ event อื่น ๆ ตามต้องการ
6.5 คลิก "Add webhook"

#### ขั้นตอนที่ 7. ทดสอบการทำงาน
7.1 ทำการแก้ไขโค้ดในโปรเจ็กต์ของคุณ (เช่น แก้ไขไฟล์ `src/app.ts`)
7.2 Commit และ Push การเปลี่ยนแปลงขึ้น GitLab
```bash
git add .
git commit -m "Test Jenkins CI/CD"
git push origin main
```

## Setup N8N for Notification

1. สมัครบัญชี ngrok ฟรีที่ https://ngrok.com/ และคัดลอก Authtoken ของคุณมาเก็บไว้

2. สร้างโฟลเดอร์ใหม่สำหรับ n8n
```bash
mkdir n8n-postgres-ngrok
```

3. สร้างไฟล์ `.env` ในโฟลเดอร์ `n8n-postgres-ngrok`
```bash
cd n8n-postgres-ngrok
touch .env
```
4. เพิ่มโค้ดในไฟล์ `.env`

```env
# PostgreSQL Credentials
POSTGRES_DB=n8n
POSTGRES_USER=admin
POSTGRES_PASSWORD=your_password_here

# n8n Encryption Key (สำคัญมาก ห้ามทำหาย)
# สร้างคีย์สุ่มยาวๆ ได้จาก: openssl rand -hex 32
N8N_ENCRYPTION_KEY=0123456789abcdef0123456789abcdef

# Timezone Settings
GENERIC_TIMEZONE=Asia/Bangkok
TZ=Asia/Bangkok

# ngrok Settings
#  สมัคร ngrok ฟรีได้ที่: https://dashboard.ngrok.com/signup
NGROK_AUTHTOKEN=your_ngrok_authtoken_here
```

5. สร้างไฟล์ `docker-compose.yml` ในโฟลเดอร์ `n8n-postgres-ngrok`
```bash
touch docker-compose.yml
```

6. เพิ่มโค้ดในไฟล์ `docker-compose.yml`
```yaml
networks:
  n8n_network:
    name: n8n_network
    driver: bridge

services:

  # service สำหรับ PostgreSQL
  postgres:
    image: postgres:16
    container_name: n8n_postgres
    restart: always
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    volumes:
      - ./postgres_data:/var/lib/postgresql/data
    networks:
      - n8n_network

  # service สำหรับ n8n
  n8n:
    image: docker.n8n.io/n8nio/n8n
    container_name: n8n_main
    restart: always
    ports:
      - "5678:5678"
    environment:
      - DB_TYPE=postgresdb
      - DB_POSTGRESDB_HOST=postgres
      - DB_POSTGRESDB_PORT=5432
      - DB_POSTGRESDB_DATABASE=${POSTGRES_DB}
      - DB_POSTGRESDB_USER=${POSTGRES_USER}
      - DB_POSTGRESDB_PASSWORD=${POSTGRES_PASSWORD}
      - N8N_ENCRYPTION_KEY=${N8N_ENCRYPTION_KEY} # ต้องตั้งค่านี้เพื่อความปลอดภัย
      - GENERIC_TIMEZONE=${GENERIC_TIMEZONE}
      - TZ=${TZ}
      - N8N_HOST=localhost # จำเป็นเพื่อให้ Tunnel ทำงานถูกต้อง
    volumes:
      - ./n8n_data:/home/node/.n8n
    networks:
      - n8n_network
    depends_on:
      - postgres

  # service สำหรับ ngrok
  ngrok:
    image: ngrok/ngrok:latest
    container_name: n8n_ngrok_tunnel
    restart: unless-stopped
    environment:
      - NGROK_AUTHTOKEN=${NGROK_AUTHTOKEN}
    command: http n8n:5678
    ports:
      - "4040:4040" # สำหรับเข้าดูหน้า Web UI
    networks:
      - n8n_network
    depends_on:
      - n8n
```

7. รัน n8n, PostgreSQL และ ngrok ด้วย Docker Compose
```bash
docker-compose up -d --build
```

8. ตรวจสอบสถานะ container ทั้งหมด
```bash
docker-compose ps
```

9. เข้าใช้งาน n8n ผ่าน URL ที่ ngrok สร้างให้
- เปิดเบราว์เซอร์แล้วไปที่ http://127.0.0.1:4040
- ดูที่ช่อง "Forwarding" จะเห็น URL ที่ ngrok สร้างให้ (เช่น https://abcd1234.ngrok.io)
- คลิกที่ URL นั้นเพื่อเข้าใช้งาน n8n (เช่น https://abcd1234.ngrok.io)
10. ตั้งค่า Webhook ใน n8n
- สร้าง Workflow ใหม่ใน n8n
- เพิ่ม Node "Webhook" และตั้งค่า Method เป็น "POST"
- คัดลอก URL ของ Webhook (เช่น https://abcd1234.ngrok.io/webhook/your-webhook-id)
11. ตั้งค่า Jenkins ให้ส่ง Notification ไปยัง n8n
- ใน Jenkinsfile ของคุณ เพิ่ม stage ใหม่หลังจาก stage 'Deploy Local' ดังนี้

```groovy
stage('Deploy Local') {
    steps {
        ...
    }
    // ส่งข้อมูลไปยัง n8n webhook เมื่อ deploy สำเร็จ
    // ใช้ Jenkins HTTP Request Plugin (ต้องติดตั้งก่อน)
    // หรือใช้ Java URLConnection แทน (fallback) ถ้า httpRequest ไม่ได้ติดตั้ง
    // n8n-webhook คือ Jenkins Secret Text Credential ที่เก็บ URL ของ n8n webhook
    // ต้องสร้าง Credential นี้ใน Jenkins ก่อน ใช้งาน
    // โดยใช้ ID ว่า n8n-webhook

    post {
        success {
            script {
                withCredentials([string(credentialsId: 'n8n-webhook', variable: 'N8N_WEBHOOK_URL')]) {
                    def payload = [
                        project  : env.JOB_NAME,
                        stage    : 'Deploy Local',
                        status   : 'success',
                        build    : env.BUILD_NUMBER,
                        image    : "${env.DOCKER_REPO}:latest",
                        container: env.APP_NAME,
                        url      : 'http://localhost:3000/',
                        timestamp: new Date().format("yyyy-MM-dd'T'HH:mm:ssXXX")
                    ]
                    def body = groovy.json.JsonOutput.toJson(payload)
                    try {
                        httpRequest acceptType: 'APPLICATION_JSON',
                                    contentType: 'APPLICATION_JSON',
                                    httpMode: 'POST',
                                    requestBody: body,
                                    url: N8N_WEBHOOK_URL,
                                    validResponseCodes: '100:599'
                        echo 'n8n webhook (success) sent via httpRequest.'
                    } catch (err) {
                        echo "httpRequest failed or not available: ${err}. Falling back to Java URLConnection..."
                        try {
                            def conn = new java.net.URL(N8N_WEBHOOK_URL).openConnection()
                            conn.setRequestMethod('POST')
                            conn.setDoOutput(true)
                            conn.setRequestProperty('Content-Type', 'application/json')
                            conn.getOutputStream().withWriter('UTF-8') { it << body }
                            int rc = conn.getResponseCode()
                            echo "n8n webhook (success) via URLConnection, response code: ${rc}"
                        } catch (e2) {
                            echo "Failed to notify n8n (success): ${e2}"
                        }
                    }
                }
            }
        }
    }
}

 post {
    always {
        echo "Pipeline finished with status: ${currentBuild.currentResult}"
    }
    success {
        echo "Pipeline succeeded!"
    }
    failure {
        // ส่งข้อมูลไปยัง n8n webhook เมื่อ pipeline ล้มเหลว
        // ใช้ Jenkins HTTP Request Plugin (ต้องติดตั้งก่อน)
        // หรือใช้ Java URLConnection แทน (fallback) ถ้า httpRequest ไม่ได้ติดตั้ง
        // n8n-webhook คือ Jenkins Secret Text Credential ที่เก็บ URL ของ n8
        // ต้องสร้าง Credential นี้ใน Jenkins ก่อน ใช้งาน
        // โดยใช้ ID ว่า n8n-webhook
        script {
            withCredentials([string(credentialsId: 'n8n-webhook', variable: 'N8N_WEBHOOK_URL')]) {
                def payload = [
                    project  : env.JOB_NAME,
                    stage    : 'Pipeline',
                    status   : 'failed',
                    build    : env.BUILD_NUMBER,
                    image    : "${env.DOCKER_REPO}:latest",
                    container: env.APP_NAME,
                    url      : 'http://localhost:3000/',
                    timestamp: new Date().format("yyyy-MM-dd'T'HH:mm:ssXXX")
                ]
                def body = groovy.json.JsonOutput.toJson(payload)
                try {
                    httpRequest acceptType: 'APPLICATION_JSON',
                                contentType: 'APPLICATION_JSON',
                                httpMode: 'POST',
                                requestBody: body,
                                url: N8N_WEBHOOK_URL,
                                validResponseCodes: '100:599'
                    echo 'n8n webhook (failure) sent via httpRequest.'
                } catch (err) {
                    echo "httpRequest failed or not available: ${err}. Falling back to Java URLConnection..."
                    try {
                        def conn = new java.net.URL(N8N_WEBHOOK_URL).openConnection()
                        conn.setRequestMethod('POST')
                        conn.setDoOutput(true)
                        conn.setRequestProperty('Content-Type', 'application/json')
                        conn.getOutputStream().withWriter('UTF-8') { it << body }
                        int rc = conn.getResponseCode()
                        echo "n8n webhook (failure) via URLConnection, response code: ${rc}"
                    } catch (e2) {
                        echo "Failed to notify n8n (failure): ${e2}"
                    }
                }
            }
        }
    }
}
```

> ข้อควรระวัง new java.net.URL ใน Jenkinsfile อาจไม่ทำงานในบางสภาพแวดล้อมของ Jenkins ที่มีการจำกัด Security Sandbox หรือไม่มี Java standard library ที่สมบูรณ์ อาจต้องใช้ Jenkins HTTP Request Plugin เป็นหลัก

11. กำหนดตัวแปร Jenkins Credential
- ไปที่ Jenkins Dashboard > Manage Jenkins > Manage Credentials
- คลิกที่ (global) > Add Credentials
- เลือก Kind เป็น Secret text
- กรอกข้อมูลในช่อง Secret เป็น URL ของ n8n webhook
- กำหนด ID ว่า n8n-webhook
- คลิก OK เพื่อบันทึก

12. ติดตั้ง Jenkins HTTP Request Plugin
- ไปที่ Jenkins Dashboard > Manage Jenkins > Manage Plugins
- ในแท็บ Available ค้นหา "HTTP Request"
- เลือกและติดตั้ง จากนั้นรีสตาร์ท Jenkins

13. สร้าง node code ใน n8n
- ใน n8n Workflow ที่สร้างไว้ เพิ่ม Node "Set" เพื่อกำหนดค่าข้อมูลที่ต้องการส่งไปยัง Slack
- เพิ่ม Node "Slack" เพื่อส่งข้อความแจ้งเตือน

```javascript
// Normalize input from Webhook node
// n8n Webhook node โดยปกติจะให้ { body, headers, query, params }
// แต่เผื่อกรณี payload ถูก map ขึ้นมาบน root ก็รองรับทั้งสองแบบ

const items = $input.all();
if (items.length === 0) {
  return [{ json: { error: 'No input items from Webhook' } }];
}

const raw = items[0].json || {};
const payload = (raw.body && typeof raw.body === 'object') ? raw.body : raw;

// Extract fields with sane defaults
const project   = String(payload.project ?? payload.job ?? 'unknown-project');
const stage     = String(payload.stage   ?? 'unknown-stage');
const status    = String(payload.status  ?? 'unknown');
const build     = String(payload.build   ?? payload.buildNumber ?? 'n/a');
const image     = String(payload.image   ?? 'n/a');
const container = String(payload.container ?? 'n/a');
const url       = String(payload.url     ?? 'http://localhost:3000/');
const timestamp = payload.timestamp ? new Date(payload.timestamp).toISOString() : new Date().toISOString();

// Small helpers
const emoji = status.toLowerCase() === 'success' ? '✅'
            : status.toLowerCase() === 'failed'  ? '❌'
            : 'ℹ️';

const lines = [
  `${emoji} Deploy ${status.toUpperCase()}: ${project} (${stage})`,
  `Build: ${build}`,
  `Image: ${image}`,
  `Container: ${container}`,
  `URL: ${url}`,
  `Time: ${timestamp}`
];
const slackText = lines.join('\n');

// Optional: Slack Block Kit (ถ้าคุณจะ map ไปใช้กับ Slack node แบบ Blocks)
const slackBlocks = [
  {
    type: 'header',
    text: { type: 'plain_text', text: `${emoji} ${project} – ${stage}` }
  },
  { type: 'divider' },
  {
    type: 'section',
    fields: [
      { type: 'mrkdwn', text: `*Status:*\n${status.toUpperCase()}` },
      { type: 'mrkdwn', text: `*Build:*\n${build}` },
      { type: 'mrkdwn', text: `*Image:*\n${image}` },
      { type: 'mrkdwn', text: `*Container:*\n${container}` },
      { type: 'mrkdwn', text: `*URL:*\n${url}` },
      { type: 'mrkdwn', text: `*Time:*\n${timestamp}` }
    ]
  }
];

// Return a single normalized item
return [{
  json: {
    // raw webhook data (for debugging)
    _webhook: raw,

    // normalized
    project, stage, status, build, image, container, url, timestamp,

    // for Slack node
    slack: {
      text: slackText,
      blocks: slackBlocks
    }
  }
}];
```

14. ตั้งค่า Slack ให้กับ node “Send a message” ใน n8n 
- ไปที่ https://api.slack.com/apps และสร้างแอปใหม่
- ตั้งค่า OAuth & Permissions โดยเพิ่ม Scopes ที่ต้องการ เช่น chat:write, incoming-webhook
- ติดตั้งแอปใน workspace ของคุณและคัดลอก OAuth Access Token
- ใน n8n เพิ่ม Node "Slack" และตั้งค่า Credentials โดยใช้ OAuth Access Token ที่ได้มา
- เชื่อมต่อ Node "Webhook" กับ Node "Slack" เพื่อส่งข้อความแจ้งเตือนเมื่อมีการเรียก Webhook

15. ตั้งค่า ส่งผ่าน n8n ไป Discord
- ใน Discord สร้าง Webhook URL สำหรับช่องที่ต้องการส่งข้อความ
- ใน n8n เพิ่ม Node "HTTP Request" หลังจาก Node "Set"
- ตั้งค่า HTTP Request ดังนี้
  - Method: POST
  - URL: [Webhook URL ที่สร้างใน Discord]
  - Body: JSON
  - JSON Body:
    ```json
    {
      "content": "ข้อความที่ต้องการส่งไปยัง Discord"
    }
    ```
- เชื่อมต่อ Node "Set" กับ Node "HTTP Request" เพื่อส่งข้อความไปยัง Discord เมื่อมีการเรียก Webhook

16. ตั้งค่า ส่งผ่าน Line Messaging API
- สร้าง Channel ใน LINE Developers Console และคัดลอก Channel Access Token
- ใน n8n เพิ่ม Node "HTTP Request" หลังจาก Node "Set"
- ตั้งค่า HTTP Request ดังนี้
  - Method: POST
  - URL: https://api.line.me/v2/bot/message/push
  - Headers:
    - Authorization: Bearer [Channel Access Token]
    - Content-Type: application/json
    - Body: JSON
    - JSON Body:
      ```json
      {
        "to": "[User ID หรือ Group ID]",
        "messages": [
          {
            "type": "text",
            "text": "ข้อความที่ต้องการส่งไปยัง LINE"
          }
        ]
      }
      ```
- เชื่อมต่อ Node "Set" กับ Node "HTTP Request" เพื่อส่งข้อความไปยัง LINE เมื่อมีการเรียก Webhook

17. ตั้งค่าแจ้งเตือนผ่าน Email
- ใน n8n เพิ่ม Node "Email" หลังจาก Node "Set"
- ตั้งค่า Email ดังนี้
  - To: [ที่อยู่อีเมลผู้รับ]
  - Subject: [หัวข้ออีเมล]
  - Body: [เนื้อหาอีเมล]
- เชื่อมต่อ Node "Set" กับ Node "Email" เพื่อส่งอีเมลแจ้งเตือนเมื่อมีการเรียก Webhook

## Jenkins multibranch pipeline

> Jenkins Multibranch Pipeline คือฟีเจอร์ที่ช่วยให้เราสามารถสร้าง Pipeline ที่สามารถทำงานกับหลายๆ สาขา (branches) ของโค้ดในระบบควบคุมเวอร์ชัน เช่น Git ได้อย่างง่ายดาย

### ข้อดีของ Jenkins Multibranch Pipeline
1. **การจัดการหลายสาขาได้ง่าย**: สามารถสร้าง Pipeline สำหรับแต่ละสาขาได้โดยอัตโนมัติ
2. **การตรวจสอบคุณภาพโค้ด**: สามารถรันการทดสอบและการตรวจสอบคุณภาพโค้ดสำหรับแต่ละสาขาได้
3. **การปรับปรุงอย่างต่อเนื่อง**: สนับสนุนการพัฒนาแบบ Agile และ CI/CD ได้ดี
4. **การแยกสภาพแวดล้อม**: สามารถแยกสภาพแวดล้อมการพัฒนา การทดสอบ และการผลิตได้อย่างชัดเจน

### การตั้งค่า Jenkins Multibranch Pipeline
1. ติดตั้ง Jenkins และ Plugins ที่จำเป็น เช่น Git, GitHub, GitLab, Pipeline, Pipeline Utility Steps, HTML Publisher
2. สร้าง Multibranch Pipeline Job ใหม่
3. ตั้งค่า Repository URL และ Credentials
4. กำหนด Branch Sources และ Strategies
5. บันทึกและรัน Pipeline

## Workshop Jenkins multibranch pipeline
1. สร้าง Repository ใหม่ใน GitHub สำหรับโค้ดตัวอย่าง
2. สร้าง Jenkins Multibranch Pipeline Job ใหม่
3. ตั้งค่า Repository URL และ Credentials
4. กำหนด Branch Sources และ Strategies
5. สร้าง Jenkinsfile สำหรับแต่ละสาขา
6. บันทึกและรัน Pipeline
7. ตรวจสอบผลลัพธ์และแก้ไขปัญหาที่เกิดขึ้น

### .NET Core Jenkins multibranch pipeline

### 🏗️ Project Structure

```
dotnet-docker-app/
├── 📄 Program.cs                           # Main application entry point
├── 📄 dotnet-docker-app.csproj            # .NET project file
├── 📄 appsettings.json                    # Application settings
├── 📄 appsettings.Development.json        # Development settings
├── 📄 dotnet-docker-app.http              # HTTP requests for testing
├── 📁 Properties/
│   └── 📄 launchSettings.json             # Launch settings
├── 📁 bin/                                 # Binary output (compiled)
│   └── 📁 Debug/
│       └── 📁 net9.0/
│           ├── 📄 dotnet-docker-app.dll
│           ├── 📄 dotnet-docker-app.exe
│           ├── 📄 dotnet-docker-app.pdb
│           └── 📄 *.json
├── 📁 obj/                                 # Object files (intermediate)
│   ├── 📄 project.assets.json
│   └── 📁 Debug/
├── 📄 .dockerignore                        # Files to ignore in Docker build
├── 📄 .gitignore                           # Files to ignore in Git
├── 🐳 Dockerfile                            # Docker build configuration (Multi-stage)
├── 🐳 docker-compose.dev.yml                # Docker Compose for development
├── 🔧 Jenkinsfile                           # Jenkins CI/CD pipeline
└── 📄 README.md                            # Project documentation
```

#### ขั้นตอนที่ 1: ดาวน์โหลดโค้ดตัวอย่าง
[https://drive.google.com/file/d/1U21ZrjAaeJVdgvq46wSOuk9qB5SlUUGA/view?usp=sharing](https://drive.google.com/file/d/1U21ZrjAaeJVdgvq46wSOuk9qB5SlUUGA/view?usp=sharing)

#### ขั้นตอนที่ 2: สร้าง Jenkinsfile สำหรับ CI/CD Pipeline
```groovy
// =================================================================
// HELPER FUNCTION: สร้างฟังก์ชันสำหรับส่ง Notification ไปยัง n8n
// การสร้างฟังก์ชันช่วยลดการเขียนโค้ดซ้ำซ้อน (DRY Principle)
// =================================================================
def sendNotificationToN8n(String status, String stageName, String imageTag, String containerName, String hostPort) {
    script {
        withCredentials([string(credentialsId: 'n8n-webhook', variable: 'N8N_WEBHOOK_URL')]) {
            def payload = [
                project  : env.JOB_NAME,
                stage    : stageName,
                status   : status,
                build    : env.BUILD_NUMBER,
                image    : "${env.DOCKER_REPO}:${imageTag}",
                container: containerName,
                url      : "http://localhost:${hostPort}/weatherforecast",
                timestamp: new Date().format("yyyy-MM-dd'T'HH:mm:ssXXX")
            ]
            def body = groovy.json.JsonOutput.toJson(payload)
            try {
                httpRequest acceptType: 'APPLICATION_JSON',
                            contentType: 'APPLICATION_JSON',
                            httpMode: 'POST',
                            requestBody: body,
                            url: N8N_WEBHOOK_URL,
                            validResponseCodes: '200:299'
                echo "n8n webhook (${status}) sent successfully."
            } catch (err) {
                echo "Failed to send n8n webhook (${status}): ${err}"
            }
        }
    }
}

pipeline {
    agent any

    options { 
        skipDefaultCheckout(true)
    }

    environment {
        DOCKER_HUB_CREDENTIALS_ID = 'dockerhub-cred'
        DOCKER_REPO               = "your-dockerhub-username/dotnet-docker-app"

        // DEV environment
        DEV_APP_NAME              = "dotnet-app-dev"
        DEV_HOST_PORT             = "6001"

        // PROD environment
        PROD_APP_NAME             = "dotnet-app-prod"
        PROD_HOST_PORT            = "6000"
    }

    parameters {
        choice(name: 'ACTION', choices: ['Build & Deploy', 'Rollback'], description: 'เลือก Action ที่ต้องการ')
        string(name: 'ROLLBACK_TAG', defaultValue: '', description: 'สำหรับ Rollback: ใส่ Image Tag ที่ต้องการ (เช่น Git Hash หรือ dev-123)')
        choice(name: 'ROLLBACK_TARGET', choices: ['dev', 'prod'], description: 'สำหรับ Rollback: เลือกว่าจะ Rollback ที่ Environment ไหน')
    }

    stages {

        // Stage 1: Checkout
        stage('Checkout') {
            when { expression { params.ACTION == 'Build & Deploy' } }
            steps {
                echo "Checking out code..."
                checkout scm
            }
        }

        // Stage 2: Restore & Test
        stage('Restore & Test') {
            when { expression { params.ACTION == 'Build & Deploy' } }
            steps {
                echo "Running tests inside a consistent Docker environment..."
                script {
                    docker.image('mcr.microsoft.com/dotnet/sdk:10.0').inside {
                        sh '''
                            dotnet restore
                            dotnet build --no-restore -c Release
                            dotnet test --no-build --verbosity normal -c Release
                        '''
                    }
                }
            }
        }

        // Stage 3: Build & Push Docker Image
        stage('Build & Push Docker Image') {
            when { expression { params.ACTION == 'Build & Deploy' } }
            steps {
                script {
                    def imageTag = (env.BRANCH_NAME == 'main') ? sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim() : "dev-${env.BUILD_NUMBER}"
                    env.IMAGE_TAG = imageTag
                    
                    docker.withRegistry('https://index.docker.io/v1/', DOCKER_HUB_CREDENTIALS_ID) {
                        echo "Building image: ${DOCKER_REPO}:${env.IMAGE_TAG}"
                        def customImage = docker.build("${DOCKER_REPO}:${env.IMAGE_TAG}", "--target final .")
                        
                        echo "Pushing images to Docker Hub..."
                        customImage.push()
                        if (env.BRANCH_NAME == 'main') {
                            customImage.push('latest')
                        }
                    }
                }
            }
        }

        // Deploy to DEV
        stage('Deploy to DEV (Local Docker)') {
            when {
                expression { params.ACTION == 'Build & Deploy' }
                branch 'develop'
            } 
            steps {
                script {
                    def deployCmd = ""
                        echo "Deploying container ${DEV_APP_NAME} from latest image..."
                        docker pull ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker stop ${DEV_APP_NAME} || true
                        docker rm ${DEV_APP_NAME} || true
                        docker run -d --name ${DEV_APP_NAME} \
                            -p ${DEV_HOST_PORT}:8080 \
                            -e ASPNETCORE_ENVIRONMENT=Development \
                            -e ASPNETCORE_URLS=http://+:8080 \
                            ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker ps --filter name=${DEV_APP_NAME} --format "table {{.Names}}\\t{{.Image}}\\t{{.Status}}"
                    ""
                    sh deployCmd
                }
            }
            post {
                success {
                    sendNotificationToN8n('success', 'Deploy to DEV (Local Docker)', env.IMAGE_TAG, env.DEV_APP_NAME, env.DEV_HOST_PORT)
                }
            }
        }

        // Approval for Production
        stage('Approval for Production') {
            when {
                expression { params.ACTION == 'Build & Deploy' }
                branch 'main'
            }
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    input message: "Deploy image tag '${env.IMAGE_TAG}' to PRODUCTION (Local Docker on port ${PROD_HOST_PORT})?"
                }
            }
        }

        // Deploy to PROD
        stage('Deploy to PRODUCTION (Local Docker)') {
            when {
                expression { params.ACTION == 'Build & Deploy' }
                branch 'main'
            } 
            steps {
                script {
                    def deployCmd = ""
                        echo "Deploying container ${PROD_APP_NAME} from latest image..."
                        docker pull ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker stop ${PROD_APP_NAME} || true
                        docker rm ${PROD_APP_NAME} || true
                        docker run -d --name ${PROD_APP_NAME} \
                            -p ${PROD_HOST_PORT}:8080 \
                            -e ASPNETCORE_ENVIRONMENT=Production \
                            -e ASPNETCORE_URLS=http://+:8080 \
                            ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker ps --filter name=${PROD_APP_NAME} --format "table {{.Names}}\\t{{.Image}}\\t{{.Status}}"
                    ""
                    sh deployCmd
                }
            }
            post {
                success {
                    sendNotificationToN8n('success', 'Deploy to PRODUCTION (Local Docker)', env.IMAGE_TAG, env.PROD_APP_NAME, env.PROD_HOST_PORT)
                }
            }
        }

        // Rollback
        stage('Execute Rollback') {
            when { expression { params.ACTION == 'Rollback' } }
            steps {
                script {
                    if (params.ROLLBACK_TAG.trim().isEmpty()) {
                        error "เมื่อเลือก Rollback กรุณาระบุ 'ROLLBACK_TAG'"
                    }

                    def targetAppName = (params.ROLLBACK_TARGET == 'dev') ? DEV_APP_NAME : PROD_APP_NAME
                    def targetHostPort = (params.ROLLBACK_TARGET == 'dev') ? DEV_HOST_PORT : PROD_HOST_PORT
                    def targetEnv = (params.ROLLBACK_TARGET == 'dev') ? 'Development' : 'Production'
                    def imageToDeploy = "${DOCKER_REPO}:${params.ROLLBACK_TAG.trim()}"
                    
                    echo "ROLLING BACK ${params.ROLLBACK_TARGET.toUpperCase()} to image: ${imageToDeploy}"
                    
                    def deployCmd = ""
                        docker pull ${imageToDeploy}
                        docker stop ${targetAppName} || true
                        docker rm ${targetAppName} || true
                        docker run -d --name ${targetAppName} \
                            -p ${targetHostPort}:8080 \
                            -e ASPNETCORE_ENVIRONMENT=${targetEnv} \
                            -e ASPNETCORE_URLS=http://+:8080 \
                            ${imageToDeploy}
                    ""
                    sh(deployCmd)
                }
            }
            post {
                success {
                    script {
                        def targetAppName = (params.ROLLBACK_TARGET == 'dev') ? DEV_APP_NAME : PROD_APP_NAME
                        def targetHostPort = (params.ROLLBACK_TARGET == 'dev') ? DEV_HOST_PORT : PROD_HOST_PORT
                        sendNotificationToN8n('success', "Rollback ${params.ROLLBACK_TARGET.toUpperCase()}", params.ROLLBACK_TAG, targetAppName, targetHostPort)
                    }
                }
            }
        }
    }

    post {
        always {
            script {
                if (params.ACTION == 'Build & Deploy') {
                    echo "Cleaning up Docker images on agent..."
                    try {
                        sh ""
                            docker image rm -f ${DOCKER_REPO}:${env.IMAGE_TAG} || true
                            docker image rm -f ${DOCKER_REPO}:latest || true
                        ""
                    } catch (err) {
                        echo "Could not clean up images, but continuing..."
                    }
                }
            }
        }
        failure {
            sendNotificationToN8n('failed', "Pipeline Failed", 'N/A', 'N/A', 'N/A')
        }
    }
}
```

####  ขั้นตอนที่ 3: Push โค้ดขึ้น GitLab

```bash
git init
git add .
git commit -m "Initial commit for Jenkins multibranch pipeline"
git branch -M main
git remote add origin https://gitlab.com/your-username/your-repo.git
git push -u origin main
```

#### ขั้นตอนที่ 4: แยก Branch สำหรับ DEV และ PROD
```bash
git checkout -b develop
git push origin develop
```

#### ขั้นตอนที่ 5: แก้ไข Program.cs เพื่อทดสอบ
```csharp
.
.
// เพิ่ม endpoint ใหม่สำหรับทดสอบ
app.MapGet("/api/products", () =>
{
    var products = new[]
    {
        new { Id = 1, Name = "Laptop", Price = 25000.00M, InStock = true },
        new { Id = 2, Name = "Mouse", Price = 500.00M, InStock = true },
        new { Id = 3, Name = "Keyboard", Price = 1200.00M, InStock = false },
        new { Id = 4, Name = "Monitor", Price = 8000.00M, InStock = true }
    };
    return Results.Ok(products);
})
.WithName("GetProducts");
.
.
```
#### ขั้นตอนที่ 6: สร้าง Jenkins Multibranch Pipeline Job ใหม่
- เปิด Jenkins Dashboard
- คลิก "New Item"
- ตั้งชื่อ Job เช่น "DotNet-Docker-App-Multibranch"
- เลือก "Multibranch Pipeline" แล้วคลิก "OK"
- ในส่วน "Branch Sources" คลิก "Add Source" แล้วเลือก "Github"
- กรอก Repository URL และเลือก Credentials ที่ตั้งค่าไว้
- ในส่วน Behavior เลือก ตามภาพ

![Jenkins Multibranch Pipeline Config](https://www.itgenius.co.th/assets/frondend/images/course_detail/devopsjenkins/itgn-1186.jpg)

- **Discover branches:** Exclude branches that are also filed as PRs (ป้องกันการสร้าง job ซ้ำซ้อนระหว่าง branch กับ PR)
- **Discover pull requests from origin:** The current pull request revision (สร้าง job สำหรับ PR จาก origin)
- **Discover pull requests from forks:** The current pull request revision (สร้าง job สำหรับ PR จาก forks)
- กำหนด "strategy" เป็น The current pull request revision 
- กำหนด "Trust" เป็น From users with Admin or Write permission
- กำหนด "Property strategy" เป็น All branches get the same properties
- ในส่วน "Build Configuration" เลือก "by Jenkinsfile"
- คลิก "Save" เพื่อบันทึกการตั้งค่า

#### ขั้นตอนที่ 7: ทดสอบการทำงานของ Jenkins Multibranch Pipeline
- Push โค้ดไปที่ branch develop เพื่อทดสอบการ deploy ไปยัง DEV
- Push โค้ดไปที่ branch main เพื่อทดสอบการ deploy ไปยัง PROD
- ทดสอบการ Rollback โดยการเลือก Action เป็น Rollback และระบุ ROLLBACK_TAG กับ ROLLBACK_TARGET
- ตรวจสอบผลลัพธ์และสถานะของ Jenkins Job

#### ขั้นตอนที่ 8: ทดสอบ API endpoints
```bash
# ทดสอบ DEV environment (port 6001)
curl http://localhost:6001/
curl http://localhost:6001/ping
curl http://localhost:6001/weatherforecast
curl http://localhost:6001/api/products

# ทดสอบ PROD environment (port 6000)
curl http://localhost:6000/
curl http://localhost:6000/ping
curl http://localhost:6000/weatherforecast
```

## Jenkins on Ubuntu Server with Docker

> Spec ของ Server (VM/Droplet) ที่จะใช้รัน Jenkins และ Docker — RAM และ CPU เป็นปัจจัยที่สำคัญที่สุด เพื่อให้การทำงานราบรื่น โดยเฉพาะตอนที่ Jenkins ทำงานพร้อมกับ Docker build และ npm install

### สรุป Spec ที่แนะนำ (TL;DR)

- 🚀 แนะนำเป็นอย่างยิ่ง (Recommended):
    - CPU: 2 vCPUs
    - RAM: 4 GB
    - Disk: 80 GB SSD
    - DigitalOcean Plan: Basic Regular Droplet, $24/เดือน (หรือ Premium Intel/AMD ที่ราคาใกล้เคียงกัน)

- 👌 ขั้นต่ำสุด (Bare Minimum):
    - CPU: 2 vCPUs
    - RAM: 2 GB
    - Disk: 50 GB SSD
    - DigitalOcean Plan: Basic Regular Droplet, $12/เดือน (อาจจะช้าตอน Build)

![Jenkins on Docker](https://www.itgenius.co.th/assets/frondend/images/course_detail/devopsjenkins/itgn-1045.jpg)

---

### ตารางเปรียบเทียบและเหตุผล

| ส่วนประกอบ | ขั้นต่ำสุด (Bare Minimum) | แนะนำ (Recommended) | เหตุผล |
| --- | --- | --- | --- |
| CPU | 2 vCPUs | 2 vCPUs | Build Time: การ Build Docker Image และรัน npm install ใช้พลังประมวลผลสูง กรณี 2 Cores จะช่วยให้ Build เสร็จเร็วขึ้นอย่างชัดเจน (1 Core อาจจะช้ามาก) |
| RAM | 2 GB | 4 GB | Performance: มีผลอย่างสำคัญมาก! Jenkins เป็นแอปที่ใช้ Java ซึ่งกิน RAM ค่อนข้างเยอะ — 2 GB: พอให้ Jenkins ทำงานได้ แต่ตอน Build อาจจะช้าหนักเพราะระบบอาจสลับ Memory ไปใช้ Disk (Swapping) — 4 GB: เพียงพอสำหรับ OS, Jenkins, Docker และตอนที่ Build โปรแกรมพร้อมกัน ทำให้ระบบโดยรวมลื่นไหลกว่ามาก |
| Disk | 50 GB SSD | 80 GB SSD | Speed & Space: ต้องเป็น SSD เท่านั้น เพื่อความเร็วในการอ่าน/เขียนไฟล์ตอน Build และสิ่งที่ Docker Image + Build History ใช้พื้นที่เยอะ การมี 80 GB จะช่วยให้คุณไม่ต้องลบภาพ (image) บ่อย ๆ |

### ทำไมถึงไม่ควรใช้ Spec ต่ำกว่านี้?

- RAM 1 GB: ไม่เพียงพออย่างแน่นอน Jenkins อาจจะล้มเหลว (Out of Memory) หรือระบบจะช้ามากจนแทบใช้งานไม่ได้
- CPU 1 vCPU: จะทำให้ทุกขั้นตอนตั้งแต่ git clone, npm install ไปจนถึง docker build ช้ามาก ทำให้วงจรการพัฒนา (Development Cycle) ของคุณช้าตามไปด้วย
- HDD (Hard Disk จานหมุน): ไม่ควรใช้เด็ดขาด เพราะความเร็ว I/O ต่ำเกินไปสำหรับงาน CI/CD จะทำให้เกิดคอขวดอย่างรุนแรง

สรุป: การลงทุนเพิ่มเติมเล็กน้อยเพื่อใช้ Plan ที่มี RAM 4 GB และ CPU 2 vCPUs จะให้ประสบการณ์การใช้งานที่ดี คุ้มค่า และมีเสถียรภาพมากกว่าอย่างชัดเจน 🚀

---
## 🗺️ แผนการย้าย (Migration Plan)

1. หยุดการทำงาน ของ Jenkins บนเครื่อง Local
2. บีบอัดและสำรองข้อมูล ทั้งหมดในโฟลเดอร์ `jenkins_home`
3. ย้ายไฟล์ ที่สำรองไว้ไปยัง Server ใหม่
4. เตรียมสภาพแวดล้อม บน Server ใหม่ (ติดตั้ง Docker, Docker Compose)
5. แตกไฟล์และคืนค่า ข้อมูล `jenkins_home` บน Server
6. เริ่มการทำงาน Jenkins บน Server ใหม่ด้วย `docker-compose.yml` เดิม
7. ตรวจสอบ ความเรียบร้อย

### หยุดและสำรองข้อมูลบนเครื่อง Local

1. หยุด Jenkins:

   ```bash
   docker-compose down
   ```

2. บีบอัดโฟลเดอร์ `jenkins_home`:

   ```bash
   tar -czvf jenkins_backup.tar.gz ./jenkins_home

   # ถ้าพบปัญหา
   # Exclude the Problematic Directory (Recommended)
   tar -czvf jenkins_backup.tar.gz \
    --exclude='./jenkins_home/tools/jenkins.plugins.nodejs.tools.NodeJSInstallation' \
    ./jenkins_home
   ```

## Jenkins Setup on Ubuntu Server with Docker (IP: 152.42.162.142)

**🚀 ขั้นตอนที่ 1: เตรียม Server (Update & Upgrade)**: ขั้นตอนแรกหลังจากสร้าง Droplet และเชื่อมต่อผ่าน SSH เข้าไปแล้ว คือการอัปเดตระบบให้เป็นเวอร์ชันล่าสุดเสมอ

1. เชื่อมต่อไปยัง Droplet ผ่าน SSH:
   
   ```bash
   ssh root@152.42.162.142
   ```

2. อัปเดต Package List และ Upgrade ระบบ:
   
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

3. ตรวจสอบข้อมูลระบบพื้นฐาน (เช็ค Kernel/OS/Hostname/Release):

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

**🐳 ขั้นตอนที่ 2: ติดตั้ง Docker Engine**:

1. ติดตั้ง Package ที่จำเป็น:
   
   ```bash
   sudo apt-get install -y ca-certificates curl gnupg
   ```

2. เพิ่ม Docker’s Official GPG Key:
   
   ```bash
   sudo install -m 0755 -d /etc/apt/keyrings
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
   sudo chmod a+r /etc/apt/keyrings/docker.gpg
   ```

3. ตั้งค่า Docker Repository:
   
   ```bash
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```

4. ติดตั้ง Docker Engine:
   
   ```bash
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin
    ```

5. (แนะนำ) เพิ่ม User ปัจจุบันของคุณไปยัง Group `docker`:
   
   ```bash
   sudo usermod -aG docker ${USER}
   exit
   ```
    - **สำคัญ**: หลังจากรันคำสั่งนี้ ให้ออกจาก SSH แล้วเชื่อมต่อเข้ามาใหม่ เพื่อให้สิทธิ์มีผล

6. ตรวจสอบการติดตั้ง Docker:
    
    ```bash
    docker --version
    docker run hello-world
    ```

**📦 ขั้นตอนที่ 3: ติดตั้ง Docker Compose**
> บน Ubuntu เวอร์ชันใหม่ๆ Docker Compose จะถูกติดตั้งเป็น Plugin ของ Docker CLI ซึ่งง่ายมาก

1. ติดตั้ง Docker Compose Plugin:
   
   ```bash
   sudo apt-get install -y docker-compose-plugin
   ```

2. ตรวจสอบการติดตั้ง:
   
   ```bash
   docker compose version
   ```

**🧱 ขั้นตอนที่ 4: ตั้งค่าและรัน Jenkins**

> ตอนนี้ระบบพร้อมแล้ว เราจะสร้างไฟล์ docker-compose.yml และรัน Jenkins container

1. สร้าง Directory สำหรับโปรเจกต์:
   
   ```bash
   mkdir jenkins-server && cd jenkins-server
   ```
2. สร้างไฟล์ `Dockerfile` สำหรับการปรับแต่ง Jenkins (ถ้าต้องการ):
   
   ```bash
   nano Dockerfile
   ```

    ```Dockerfile
    # เริ่มต้นจาก Image Jenkins ที่เราต้องการ
    FROM jenkins/jenkins:jdk21

    USER root

    # Install Docker CLI and Node.js
    RUN apt-get update && \
        apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release && \
        curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.co>
        apt-get update && \
        apt-get install -y docker-ce-cli && \
        curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
        apt-get install -y nodejs && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*

    # สร้าง entrypoint script เพื่อแก้สิทธิ์ docker socket ตอน runtime
    RUN echo '#!/bin/bash\n\
    DOCKER_SOCK="/var/run/docker.sock"\n\
    if [ -S "$DOCKER_SOCK" ]; then\n\
        DOCKER_GID=$(stat -c "%g" $DOCKER_SOCK)\n\
        if ! getent group $DOCKER_GID > /dev/null 2>&1; then\n\
            groupadd -g $DOCKER_GID docker\n\
        fi\n\
        usermod -aG $DOCKER_GID jenkins\n\
    fi\n\
    exec /usr/bin/tini -- /usr/local/bin/jenkins.sh "$@"' > /usr/local/bin/entrypoint.sh && \
        chmod +x /usr/local/bin/entrypoint.sh

    USER jenkins

    ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
    ```

   - กด `CTRL + X` แล้วกด `Y` เพื่อบันทึกไฟล์ จากนั้นกด `ENTER` เพื่อยืนยันชื่อไฟล์
   
3. สร้างไฟล์ `docker-compose.yml` ด้วยเนื้อหาดังนี้:

    ```bash
    nano docker-compose.yml
    ```
    > สำคัญ: Network ต้องเป็นแบบ external ที่สร้างไว้ล่วงหน้าแล้ว (n8n-server_n8n-network) มาจาก n8n setup
    ```yaml
    # Define Network
    networks:
        n8n-server_n8n-network:
            external: true

    # Define Services
    services:
        jenkins:
            build: .
            image: jenkins-with-docker:jdk21
            container_name: jenkins
            user: root # Run as root to allow entrypoint script to set permissions
            volumes:
                - ./jenkins_home:/var/jenkins_home # เพื่อให้ Jenkins สามารถเก็บข้อมูลไว้ใน host
                - /var/run/docker.sock:/var/run/docker.sock # เพื่อให้ Jenkins สามารถใช้งาน Docker daemon ที่รันบน host ได้
            environment:
                - JENKINS_OPTS=--httpPort=8800 # กำหนด Port สำหรับ Jenkins UI
            ports:
                - "8800:8800" # สำหรับ Jenkins UI
            restart: always
            networks:
                - n8n-server_n8n-network
    ```
    - กด `CTRL + X` แล้วกด `Y` เพื่อบันทึกไฟล์ จากนั้นกด `ENTER` เพื่อยืนยันชื่อไฟล์

4. 🚚 ย้ายไฟล์ Backup ไปยัง Server ใหม่

    ```bash
    exit
    # รูปแบบ: scp [ไฟล์ต้นทาง] [user@server_ip:ปลายทาง]
    scp ./jenkins_backup.tar.gz user@your-server-ip:~/

    # ตัวอย่างเช่น
    # scp C:\TrainingWorkshop\devops-jenkins-github-actions-n8n\jenkins-cicd-workshop\jenkins-server\jenkins_backup.tar.gz root@128.199.133.80:~/
    ```
    > คำสั่งนี้จะคัดลอกไฟล์ jenkins_backup.tar.gz ไปไว้ที่ Home directory (~/) ของ user บน Server ใหม่

5. ✨ คืนค่าข้อมูลและเริ่ม Jenkins

   ```bash
   # บน Server ใหม่, ภายในโฟลเดอร์ jenkins-server
   cd jenkins-server
   mv ~/jenkins_backup.tar.gz .
   tar -xzvf jenkins_backup.tar.gz
   ```
   > คำสั่งนี้จะสร้างโฟลเดอร์ ./jenkins_home ขึ้นมาใหม่ ซึ่งมีข้อมูลครบถ้วนเหมือนบนเครื่อง Local

6. กำหนดสิทธิ์ (สำคัญมาก):
   
   ```bash
   sudo chown -R 1000:1000 jenkins_home
   ```

7. รัน Jenkins Container:
   
   ```bash
   docker compose up -d
   ```
   > คำสั่งนี้จะดาวน์โหลด Jenkins image และเริ่ม container ในพื้นหลัง (-d)

6. ตรวจสอบสถานะ Container:
   
   ```bash
   docker compose ps
   docker compose logs -f jenkins
   ```
   > รอจนกว่าจะเห็นข้อความ Jenkins is fully up and running

7. เข้าใช้งาน Jenkins ผ่านเว็บเบราว์เซอร์:
   - เปิดเบราว์เซอร์และไปที่ `http://your-server-ip:8800`

**🛡️ ขั้นตอนที่ 5: ตั้งค่า Firewall**
> DigitalOcean Droplets จะมี Firewall (ufw) ป้องกันอยู่ เราต้องเปิด Port ที่จำเป็น

1. อนุญาตการเชื่อมต่อ SSH (สำคัญมาก!):
   ```bash
   sudo ufw allow OpenSSH
   ```

2. อนุญาตการเชื่อมต่อไปยัง Port ของ Jenkins (8800):
   ```bash
   sudo ufw allow 8800
   ```

3. เปิดใช้งาน Firewall:
   ```bash
   sudo ufw enable
   ```

4. ตรวจสอบสถานะของ Firewall:
   ```bash
   sudo ufw status
   ```
   > คุณควรจะเห็นว่า Port 22 (OpenSSH) และ 8800 ได้รับการอนุญาต (ALLOW)

**✨ ขั้นตอนที่ 6: ตั้งค่า Jenkins ครั้งแรก (กรณีการติดตั้งใหม่)**

1. ดูรหัสผ่านเริ่มต้น (Initial Admin Password):
   ```bash
   sudo cat /var/jenkins_home/secrets/initialAdminPassword
   ```
   - คัดลอกรหัสผ่านนี้ไปใช้ในขั้นตอนถัดไป

2. เปิดเบราว์เซอร์และไปที่ `http://your-server-ip:8800`
3. ใตั้งค่า Jenkins:
   - ใส่รหัสผ่านที่ได้จากขั้นตอนก่อนหน้า
   - เลือก "Install suggested plugins"
   - สร้างผู้ใช้ Admin ใหม่
   - ตั้งค่า URL ของ Jenkins (ถ้าจำเป็น)

---

## N8N on Ubuntu Server with Docker (IP: 152.42.162.142)
> การติดตั้ง N8N บน Ubuntu Server ด้วย Docker นั้นง่ายและรวดเร็ว โดยใช้ Docker Compose เพื่อจัดการ Container

1. สร้างโฟลเดอร์สำหรับ N8N:
   ```bash
   mkdir n8n-server && cd n8n-server
   ```

2. สร้างไฟล์ `.env` ด้วยเนื้อหาดังนี้:
    ```env
    # PostgreSQL Credentials
    POSTGRES_DB=n8n
    POSTGRES_USER=admin
    POSTGRES_PASSWORD=yourpassword

    # n8n Encryption Key (สำคัญมาก ห้ามทำหาย)
    # สร้างคีย์สุ่มยาวๆ ได้จาก: openssl rand -hex 32
    N8N_ENCRYPTION_KEY=0123456789abcdef0123456789abcdef

    # Timezone Settings
    GENERIC_TIMEZONE=Asia/Bangkok
    TZ=Asia/Bangkok

    # กรณีไม่ได้ใช้ https ให้ตั้งค่า Secure Cookie เป็น false
    # Disable secure cookie for HTTP access
    N8N_SECURE_COOKIE=false
    ```

3. สร้างไฟล์ `docker-compose.yml` ด้วยเนื้อหาดังนี้:
    ```yaml
    networks:
    n8n-network:
        driver: bridge

    services:

    # service สำหรับ PostgreSQL
    postgres:
        image: postgres:16
        container_name: n8n_postgres
        restart: always
        environment:
            - POSTGRES_USER=${POSTGRES_USER}
            - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
            - POSTGRES_DB=${POSTGRES_DB}
        volumes:
            - ./postgres_data:/var/lib/postgresql/data
        networks:
            - n8n-network
    # service สำหรับ n8n
    n8n:
        image: docker.n8n.io/n8nio/n8n
        container_name: n8n_main
        restart: always
        ports:
            - "152.42.162.142:5678:5678"
        environment:
            - DB_TYPE=postgresdb
            - DB_POSTGRESDB_HOST=postgres
            - DB_POSTGRESDB_PORT=5432
            - DB_POSTGRESDB_DATABASE=${POSTGRES_DB}
            - DB_POSTGRESDB_USER=${POSTGRES_USER}
            - DB_POSTGRESDB_PASSWORD=${POSTGRES_PASSWORD}
            - N8N_ENCRYPTION_KEY=${N8N_ENCRYPTION_KEY}
            - GENERIC_TIMEZONE=${GENERIC_TIMEZONE}
            - TZ=${TZ}
            - N8N_HOST=152.42.162.142
            - N8N_SECURE_COOKIE=${N8N_SECURE_COOKIE}
        volumes:
            - ./n8n_data:/home/node/.n8n
        networks:
            - n8n-network
        depends_on:
            - postgres
    ```

4. รัน N8N ด้วย Docker Compose:
   ```bash
   docker-compose up -d
   ```

5. ตรวจสอบสถานะ Container:
   ```bash
   docker-compose ps
   ```

6. กำหนดสิทธิ์ (สำคัญมาก):
   ```bash
   sudo chown -R 1000:1000 n8n_data
   sudo chown -R 999:999 postgres_data
   ```
   > 1000:1000 คือ UID:GID ของ user node (n8n) และ 999:999 คือ UID:GID ของ user postgres

7. ตั้งค่า Firewall (ถ้ายังไม่ได้ตั้งค่า):
   ```bash
   sudo ufw allow 5678
   sudo ufw reload
   ```
8. เข้าใช้งาน N8N ผ่านเว็บเบราว์เซอร์:
   - เปิดเบราว์เซอร์และไปที่ `http://152.42.162.142:5678`
   - คุณควรเห็นหน้าเข้าลงงทะเบียน/เข้าสู่ระบบของ N8N

---
## Jenkins CI/CD deployed to Server with Docker and SSH

ต้องปรับปรุงหลัก ๆ 3 ส่วนดังนี้

1. ติดตั้ง Plugin และตั้งค่า Credentials: บน Jenkins Server
2. เตรียม Server ปลายทาง: ตั้งค่าให้ Jenkins เชื่อมต่อผ่าน SSH ได้
3. แก้ไข `Jenkinsfile`: เปลี่ยนจากการรันคำสั่ง `sh 'docker ...'` ตรงๆ ไปเป็นการรันผ่าน sshagent

#### ขั้นตอนที่ 1: ติดตั้ง Plugin และตั้งค่า Credentials (บน Jenkins Server IP: 152.42.162.142)

1. ติดตั้ง Plugin:
   - ไปที่ "Manage Jenkins" > "Manage Plugins"
   - ติดตั้ง Plugin ที่จำเป็น เช่น "SSH Agent" และ "Docker Pipeline"
   - Restart Jenkins ถ้าจำเป็น

2. สร้าง SSH Credentials:
   - คุณต้องมีคู่คีย์ SSH (Private และ Public Key) ก่อน หากยังไม่มี ให้สร้างด้วย `ssh-keygen`
   - ไปที่ Manage Jenkins > Credentials > System > Global credentials (unrestricted).
   - คลิก Add Credentials:
     - Kind: เลือก `SSH Username with private key`
     - ID: ตั้ง ID ที่จำง่าย (เช่น `remote-deploy-key`)
     - Description: (Optional) e.g., "Key for deploying to remote server"
     - Username: ชื่อ user บนเครื่องปลายทางที่คุณจะใช้ login (เช่น `ubuntu` หรือ `jenkins_deploy`)
     - Private Key: เลือก `Enter directly` แล้ววางเนื้อหาของ Private Key (ไฟล์ `id_rsa`) ลงไป
     - ตัวอย่าง Private Key:
       ```
       -----BEGIN OPENSSH PRIVATE KEY-----
       b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
       ...
       -----END OPENSSH PRIVATE KEY-----
       ```
     - Passphrase: (ถ้ามี) ใส่ passphrase ของ Private Key (ตัวอย่าง 123456)
   - คลิก OK เพื่อบันทึก

#### ขั้นตอนที่ 2: เตรียม Server ปลายทาง (Remote Target Server IP: 128.199.242.92)

1. ติดตั้ง Docker และ Docker Compose บน Server ปลายทาง (ถ้ายังไม่ได้ติดตั้ง)
2. เพิ่ม Public Key:
   - นำ Public Key (ไฟล์ `id_rsa.pub`) ที่คู่กับ Private Key ในขั้นตอนที่ 1
   - เพิ่ม Public Key นี้ลงในไฟล์ `~/.ssh/authorized_keys` ของ user ที่คุณระบุใน Credential (เช่น ubuntu หรือ jenkins_deploy)
   - ใช้ nano หรือ echo เพื่อเพิ่มคีย์:
   ```
    nano ~/.ssh/authorized_keys
    ```
   - หรือ
   ```
    echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGW40E43GWkuAiSoPKD1tLTdYZFEIbaLr6J6G6bEJwfW root@ubuntu-jenkin-n8n-server" >> ~/.ssh/authorized_keys
   ```
3. ให้สิทธิ์ Docker (สำคัญมาก):
   - User ที่คุณใช้ SSH เข้าไป จะต้องรันคำสั่ง docker ได้โดยไม่ต้องใช้ sudo 
   - คำสั่ง: `sudo usermod -aG docker <your-ssh-username>` (เช่น `sudo usermod -aG docker ubuntu`)
   - ออกจากระบบแล้ว login ใหม่ หรือ restart server เพื่อให้การเปลี่ยนแปลงมีผล

#### ขั้นตอนที่ 3: แก้ไข Jenkinsfile

แก้ไข Jenkinsfile ของคุณให้เรียกใช้ sshagent เพื่อห่อ (wrap) คำสั่ง sh ที่จะไปรันบนเครื่องปลายทาง

การเปลี่ยนแปลงหลัก:

1. เพิ่ม environment variables:

- `SSH_CREDENTIALS_ID`: ID ของ Credential ที่คุณสร้างในขั้นตอนที่ 1
- `REMOTE_USER`: Username ที่ใช้ SSH (เช่น ubuntu)
- `REMOTE_HOST_IP`: IP Address หรือ Domain name ของ Server ปลายทาง

2. แก้ไข sendNotificationToN8n:

- เพิ่มพารามิเตอร์ `hostIp`
- เปลี่ยน `url` จาก `localhost` เป็น `hostIp` ที่ได้รับมา เพื่อให้ URL ใน Notification ถูกต้อง

3. แก้ไข Deploy และ Rollback stages:

- ลบ `def deployCmd = ...` และ `sh deployCmd` แบบเดิม
- ใช้ `script { ... }`
- ห่อคำสั่ง `sh` ด้วย `sshagent(credentials: [env.SSH_CREDENTIALS_ID]) { ... }`
- เปลี่ยนคำสั่ง `sh` ให้เป็น `sh "ssh -o StrictHostKeyChecking=no ${env.REMOTE_USER}@${env.REMOTE_HOST_IP} '...your commands... ' "`
  `-o StrictHostKeyChecking=no`: เพื่อป้องกันไม่ให้ Pipeline ค้างรอถาม "Are you sure you want to continue connecting (yes/no)?" ในครั้งแรก

นี่คือ `Jenkinsfile` ฉบับปรับปรุง:

```groovy
// =================================================================
// HELPER FUNCTION: สร้างฟังก์ชันสำหรับส่ง Notification ไปยัง n8n
// [MODIFIED] เพิ่ม hostIp และเปลี่ยน localhost
// =================================================================

def sendNotificationToN8n(String status, String stageName, String imageTag, String containerName, String hostPort, String hostIp) { // [MODIFIED] Added hostIp
    script {
        withCredentials([string(credentialsId: 'n8n-webhook', variable: 'N8N_WEBHOOK_URL')]) {
            def payload = [
                project  : env.JOB_NAME,
                stage    : stageName,
                status   : status,
                build    : env.BUILD_NUMBER,
                image    : "${env.DOCKER_REPO}:${imageTag}",
                container: containerName,
                url      : "http://${hostIp}:${hostPort}/", // [MODIFIED] Use hostIp instead of localhost
                timestamp: new Date().format("yyyy-MM-dd'T'HH:mm:ssXXX")
            ]
            def body = groovy.json.JsonOutput.toJson(payload)
            try {
                httpRequest acceptType: 'APPLICATION_JSON',
                            contentType: 'APPLICATION_JSON',
                            httpMode: 'POST',
                            requestBody: body,
                            url: N8N_WEBHOOK_URL,
                            validResponseCodes: '200:299'
                echo "n8n webhook (${status}) sent successfully."
            } catch (err) {
                echo "Failed to send n8n webhook (${status}): ${err}"
            }
        }
    }
}

pipeline {
    agent any

    options { 
        skipDefaultCheckout(true) 
    }

    environment {
        DOCKER_HUB_CREDENTIALS_ID = 'dockerhub-cred'
        DOCKER_REPO               = "your-dockerhub-username/express-app"

        DEV_APP_NAME              = "express-app-dev"
        DEV_HOST_PORT             = "3001"

        PROD_APP_NAME             = "express-app-prod"
        PROD_HOST_PORT            = "3000"

        // ========================================================
        // [NEW] เพิ่มตัวแปรสำหรับ Remote Server
        // ========================================================
        SSH_CREDENTIALS_ID        = 'remote-deploy-key' // ID ของ SSH Credential ที่สร้างใน Jenkins
        REMOTE_USER               = 'root'            // User บนเครื่องปลายทาง
        REMOTE_HOST_IP            = '128.199.242.92' // IP ของเครื่องปลายทาง
    }

    parameters {
        choice(name: 'ACTION', choices: ['Build & Deploy', 'Rollback'], description: 'เลือก Action ที่ต้องการ')
        string(name: 'ROLLBACK_TAG', defaultValue: '', description: 'สำหรับ Rollback: ใส่ Image Tag ที่ต้องการ (เช่น Git Hash หรือ dev-123)')
        choice(name: 'ROLLBACK_TARGET', choices: ['dev', 'prod'], description: 'สำหรับ Rollback: เลือกว่าจะ Rollback ที่ Environment ไหน')
    }

    stages {

        // =================================================================
        // BUILD STAGES: (ไม่มีการเปลี่ยนแปลง)
        // =================================================================

        stage('Checkout') {
            when { expression { params.ACTION == 'Build & Deploy' } }
            steps {
                echo "Checking out code..."
                checkout scm
            }
        }

        stage('Install & Test') {
            when { expression { params.ACTION == 'Build & Deploy' } }
            steps {
                echo "Running tests inside a consistent Docker environment..."
                script {
                    docker.image('node:22-alpine').inside {
                        sh '''
                            if [ -f package-lock.json ]; then npm ci; else npm install; fi
                            npm test
                        '''
                    }
                }
            }
        }

        stage('Build & Push Docker Image') {
            when { expression { params.ACTION == 'Build & Deploy' } }
            steps {
                script {
                    def imageTag = (env.BRANCH_NAME == 'main') ? sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim() : "dev-${env.BUILD_NUMBER}"
                    env.IMAGE_TAG = imageTag
                    
                    docker.withRegistry('https://index.docker.io/v1/', DOCKER_HUB_CREDENTIALS_ID) {
                        echo "Building image: ${DOCKER_REPO}:${env.IMAGE_TAG}"
                        def customImage = docker.build("${DOCKER_REPO}:${env.IMAGE_TAG}", "--target production .")
                        
                        echo "Pushing images to Docker Hub..."
                        customImage.push()
                        if (env.BRANCH_NAME == 'main') {
                            customImage.push('latest')
                        }
                    }
                }
            }
        }

        // =================================================================
        // DEPLOY STAGES: [MODIFIED]
        // =================================================================

        stage('Deploy to DEV (Remote Server)') { // [MODIFIED] Changed name
            when {
                expression { params.ACTION == 'Build & Deploy' }
                branch 'develop'
            } 
            steps {
                script {
                    // [MODIFIED] สร้าง command string
                    def deployCmd = ""
                        echo 'Deploying container ${DEV_APP_NAME} on REMOTE server...'
                        docker pull ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker stop ${DEV_APP_NAME} || true
                        docker rm ${DEV_APP_NAME} || true
                        docker run -d --name ${DEV_APP_NAME} -p ${DEV_HOST_PORT}:3000 ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker ps --filter name=${DEV_APP_NAME} --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
                    ""
                    
                    // [MODIFIED] ใช้ sshagent ห่อ sh command
                    sshagent(credentials: [env.SSH_CREDENTIALS_ID]) {
                        sh "ssh -o StrictHostKeyChecking=no ${env.REMOTE_USER}@${env.REMOTE_HOST_IP} \"${deployCmd}\"
                    }
                }
            }
            post {
                success {
                    // [MODIFIED] ส่ง IP ของ Remote Server ไปด้วย
                    sendNotificationToN8n('success', 'Deploy to DEV (Remote Server)', env.IMAGE_TAG, env.DEV_APP_NAME, env.DEV_HOST_PORT, env.REMOTE_HOST_IP)
                }
            }
        }

        stage('Approval for Production') {
            when {
                expression { params.ACTION == 'Build & Deploy' }
                branch 'main'
            }
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    // [MODIFIED] อัปเดตข้อความให้ชัดเจนว่าไป Remote
                    input message: "Deploy image tag '${env.IMAGE_TAG}' to PRODUCTION (Remote Server: ${REMOTE_HOST_IP} on port ${PROD_HOST_PORT})?"
                }
            }
        }

        stage('Deploy to PRODUCTION (Remote Server)') { // [MODIFIED] Changed name
            when {
                expression { params.ACTION == 'Build & Deploy' }
                branch 'main'
            } 
            steps {
                script {
                    // [MODIFIED] สร้าง command string
                    def deployCmd = ""
                        echo 'Deploying container ${PROD_APP_NAME} on REMOTE server...'
                        docker pull ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker stop ${PROD_APP_NAME} || true
                        docker rm ${PROD_APP_NAME} || true
                        docker run -d --name ${PROD_APP_NAME} -p ${PROD_HOST_PORT}:3000 ${DOCKER_REPO}:${env.IMAGE_TAG}
                        docker ps --filter name=${PROD_APP_NAME} --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
                    ""
                    
                    // [MODIFIED] ใช้ sshagent ห่อ sh command
                    sshagent(credentials: [env.SSH_CREDENTIALS_ID]) {
                        sh "ssh -o StrictHostKeyChecking=no ${env.REMOTE_USER}@${env.REMOTE_HOST_IP} \"${deployCmd}\"
                    }
                }
            }
            post {
                success {
                    // [MODIFIED] ส่ง IP ของ Remote Server ไปด้วย
                    sendNotificationToN8n('success', 'Deploy to PRODUCTION (Remote Server)', env.IMAGE_TAG, env.PROD_APP_NAME, env.PROD_HOST_PORT, env.REMOTE_HOST_IP)
                }
            }
        }

        // =================================================================
        // ROLLBACK STAGE: [MODIFIED]
        // =================================================================
        stage('Execute Rollback (Remote)') { // [MODIFIED] Changed name
            when { expression { params.ACTION == 'Rollback' } }
            steps {
                script {
                    if (params.ROLLBACK_TAG.trim().isEmpty()) {
                        error "เมื่อเลือก Rollback กรุณาระบุ 'ROLLBACK_TAG'"
                    }

                    def targetAppName = (params.ROLLBACK_TARGET == 'dev') ? DEV_APP_NAME : PROD_APP_NAME
                    def targetHostPort = (params.ROLLBACK_TARGET == 'dev') ? DEV_HOST_PORT : PROD_HOST_PORT
                    def imageToDeploy = "${DOCKER_REPO}:${params.ROLLBACK_TAG.trim()}"
                    
                    echo "ROLLING BACK ${params.ROLLBACK_TARGET.toUpperCase()} on REMOTE server to image: ${imageToDeploy}"
                    
                    // [MODIFIED] สร้าง command string
                    def deployCmd = ""
                        docker pull ${imageToDeploy}
                        docker stop ${targetAppName} || true
                        docker rm ${targetAppName} || true
                        docker run -d --name ${targetAppName} -p ${targetHostPort}:3000 ${imageToDeploy}
                    ""

                    // [MODIFIED] ใช้ sshagent ห่อ sh command
                    sshagent(credentials: [env.SSH_CREDENTIALS_ID]) {
                        sh "ssh -o StrictHostKeyChecking=no ${env.REMOTE_USER}@${env.REMOTE_HOST_IP} \"${deployCmd}\"
                    }
                }
            }
            post {
                success { 
                    // [MODIFIED] ส่ง IP ของ Remote Server ไปด้วย และต้องดึงตัวแปรมาไว้นอก script block
                    script {
                        def targetAppName = (params.ROLLBACK_TARGET == 'dev') ? env.DEV_APP_NAME : env.PROD_APP_NAME
                        def targetHostPort = (params.ROLLBACK_TARGET == 'dev') ? env.DEV_HOST_PORT : env.PROD_HOST_PORT
                        sendNotificationToN8n('success', "Rollback ${params.ROLLBACK_TARGET.toUpperCase()}", params.ROLLBACK_TAG, targetAppName, targetHostPort, env.REMOTE_HOST_IP)
                    }
                }
            }
        }
    }

    post {
        always {
            script {
                // [MODIFIED] การลบ image บน agent (Jenkins server) ยังคงเหมือนเดิม
                if (params.ACTION == 'Build & Deploy') {
                    echo "Cleaning up Docker images on agent..."
                    try {
                        sh ""
                            docker image rm -f ${DOCKER_REPO}:${env.IMAGE_TAG} || true
                            docker image rm -f ${DOCKER_REPO}:latest || true
                        ""
                    } catch (err) {
                        echo "Could not clean up images, but continuing..."
                    }
                }
                echo "Cleaning up workspace..."
                cleanWs()
            }
        }
        failure {
            sendNotificationToN8n('failed', "Pipeline Failed", 'N/A', 'N/A', 'N/A', 'N/A') // [MODIFIED]
        }
    }
}
```

## สิ่งที่เรียนรู้ใน Day 2
- เข้าใจแนวคิดของ CI/CD และความสำคัญในการพัฒนา Software
- รู้จักกับเครื่องมือที่ใช้ในกระบวนการ CI/CD เช่น GitLab, Jenkins, Docker
- สามารถสร้าง Pipeline บน Jenkins เพื่อทำการ Build, Test, และ Deploy Application
- ตั้งค่า Webhook ใน GitLab เพื่อเชื่อมต่อกับ Jenkins
- เข้าใจการทำงานของ Docker และสามารถสร้าง Docker Image และ Push ขึ้น Docker Hub
- เข้าใจการใช้งาน Jenkinsfile ในการกำหนดขั้นตอนของ Pipeline
- เรียนรู้การจัดการ Credentials ใน Jenkins เพื่อเชื่อมต่อกับ GitLab และ Docker Hub
- เข้าใจการทำงานของ Declarative Pipeline และ Scripted Pipeline ใน Jenkins
- สามารถสร้างและจัดการ Jenkins Job เพื่อทำงานร่วมกับ GitLab Repository
- เรียนรู้การใช้งาน Node.js, Express.js, TypeScript, และ Jest ในการพัฒนาและทดสอบ Application
- การเลือกสเปค Server ที่เหมาะสมสำหรับ Jenkins CI/CD
- การติดตั้ง Jenkins บน Ubuntu Server ด้วย Docker และ Docker Compose
- การตั้งค่า Jenkins Pipeline สำหรับ CI/CD
- การใช้งาน Jenkinsfile และ Declarative Pipeline
- การติดตั้งและใช้งาน N8N สำหรับ Automation