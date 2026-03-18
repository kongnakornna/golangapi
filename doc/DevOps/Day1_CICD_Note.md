## DevOps Jenkins & GitLab Actions & N8N - Day 1

### ดาวน์โหลดเอกสารประกอบการอบรม

[คลิกที่นี่เพื่อดาวน์โหลดเอกสารประกอบการอบรม](https://bit.ly/devops_easybuy)

### 📋 สารบัญ

1. [ภาพรวมหลักสูตร](#ภาพรวมหลักสูตร)
2. [การติดตั้งเครื่องมือ DevOps CI/CD](#การติดตั้งเครื่องมือ-devops-cicd)
3. [พื้นฐานการใช้งาน Git และ GitLab](#พื้นฐานการใช้งาน-git-และ-gitlab)
4. [พื้นฐานการใช้งาน Docker](#พื้นฐานการใช้งาน-docker)
5. [พื้นฐานการใช้งาน Jenkins](#พื้นฐานการใช้งาน-jenkins)
6. [การสร้าง Pipeline ด้วย Jenkins](#การสร้าง-pipeline-ด้วย-jenkins)
7. [การเชื่อมต่อ Jenkins กับ GitLab](#การเชื่อมต่อ-jenkins-กับ-gitlab)
8. [การตั้งค่า Webhook ใน GitLab เพื่อ Trigger Jenkins Job](#การตั้งค่า-webhook-ใน-gitlab-เพื่อ-trigger-jenkins-job)

## โปรแกรม (Tool and Editor) ที่ใช้อบรม

1. **Visual Studio Code**
2. **Java JDK 21.x**
3. **Docker Desktop**
4. **Git**
5. **GitLab Account**

> แนะนำ 
> หลักสูตรนี้ใช้ Java JDK เวอร์ชั่น 17 หรือ 21 เท่านั้น (เก่าหรือใหม่กว่านี้ไม่รองรับ)
---

## การตรวจสอบความเรียบร้อยของเครื่องมือที่ติดตั้งบน Windows / Mac OS / Linux

เปิด Command Prompt บน Windows หรือ Terminal บน Mac ขึ้นมาป้อนคำสั่งดังนี้

### Visual Studio Code
```bash
code --version
```

### Node JS
```bash
node -v
npm -v
npx -v
```

### Java JDK
```bash
java -version
set JAVA_HOME
where java
```

### Docker
```bash
docker --version
docker compose version
docker info
```

### Git
```bash
git version
```

---

## ภาพรวมหลักสูตร
ในหลักสูตรนี้ เราจะได้เรียนรู้เกี่ยวกับการใช้งานเครื่องมือต่าง ๆ ที่เกี่ยวข้องกับ DevOps, CI/CD, Jenkins, GitLab Actions และ N8N ผ่านการปฏิบัติจริง โดยมีเนื้อหาหลัก ๆ ดังนี้:
1. การติดตั้งและตั้งค่าเครื่องมือ DevOps CI/CD
2. การใช้งาน Git และ GitLab สำหรับการจัดการเวอร์ชันของโค้ด
3. การสร้างและจัดการคอนเทนเนอร์ด้วย Docker
4. การตั้งค่าและใช้งาน Jenkins สำหรับการทำ CI/CD
5. การสร้างและจัดการเวิร์กโฟลว์ด้วย GitLab Actions
6. การใช้งาน N8N สำหรับการทำงานอัตโนมัติต่าง ๆ

---

## การติดตั้งเครื่องมือ DevOps CI/CD
ในส่วนนี้ เราจะทำการติดตั้งเครื่องมือต่าง ๆ ที่จำเป็นสำหรับการทำงานในหลักสูตรนี้ ได้แก่ Visual Studio Code, Java JDK, Docker Desktop, และ Git

### 1. การติดตั้ง Visual Studio Code
- ไปที่เว็บไซต์ [Visual Studio Code](https://code.visualstudio.com/)
- ดาวน์โหลดตัวติดตั้งที่เหมาะสมกับระบบปฏิบัติการของคุณ (Windows, Mac OS, Linux)
- ทำการติดตั้งตามขั้นตอนที่แนะนำในเว็บไซต์

### 2. การติดตั้ง Java JDK
- ไปที่เว็บไซต์ [Oracle](https://www.oracle.com/java/technologies/javase/jdk21-archive-downloads.html)
- เลือก Java JDK เวอร์ชัน 21.x และดาวน์โหลดตัวติดตั้งที่เหมาะสมกับระบบปฏิบัติการของคุณ
- ทำการติดตั้งตามขั้นตอนที่แนะนำในเว็บไซต์
- ตั้งค่า JAVA_HOME ในระบบปฏิบัติการของคุณ

### 3. การติดตั้ง Docker Desktop
- ไปที่เว็บไซต์ [Docker](https://www.docker.com/products/docker-desktop/)
- ดาวน์โหลดตัวติดตั้งที่เหมาะสมกับระบบปฏิบัติการของคุณ
- ทำการติดตั้งตามขั้นตอนที่แนะนำในเว็บไซต์
- หลังการติดตั้ง ให้ทำการล็อกอินด้วยบัญชี Docker Hub ของคุณ

### 4. การติดตั้ง Git
- ไปที่เว็บไซต์ [Git](https://git-scm.com/)
- ดาวน์โหลดตัวติดตั้งที่เหมาะสมกับระบบปฏิบัติการของคุณ
- ทำการติดตั้งตามขั้นตอนที่แนะนำในเว็บไซต์
- ตั้งค่า Git ด้วยคำสั่ง:
```bash
git config --global user.name "Your Name"
git config --global user.email "you@example.com"
```
- ดูการตั้งค่าโดยใช้คำสั่ง:
```bash
git config --list
```

---

## พื้นฐานการใช้งาน Git และ GitLab
ในส่วนนี้ เราจะเรียนรู้พื้นฐานการใช้งาน Git และ GitLab สำหรับการจัดการเวอร์ชันของโค้ด โดยจะมีเนื้อหาหลัก ๆ ดังนี้:
1. การสร้าง Repository บน GitLab
2. การโคลน Repository ลงในเครื่องของเรา
3. การเพิ่มไฟล์และการ commit การเปลี่ยนแปลง
4. การ push การเปลี่ยนแปลงกลับไปยัง GitLab
5. การสร้างและจัดการ Branches ใน Git

## พื้นฐานการใช้งาน Git

- เริ่มต้นใช้งาน Git ในโปรเจ็กต์
- การตั้งค่า , การ clone , คำสั่งพื้นฐาน

### 1. Git First time setup
```bash
git config --global user.name "Your Name"
git config --global user.email "you@example.com"
```
> คำสั่งนี้จะทำให้ Git รู้ว่า ใครเป็นคนทำการ commit โค้ด และทำเพียงครั้งเดียวเท่านั้น

คำสั่งเช็คข้อมูลที่กำหนดไว้

```bash
git config  --list
git config  --global --list
```

### 2. Set Default branch "main"
```bash
# ตั้งค่า default branch เป็น main
git config --global init.defaultBranch main

# ตรวจสอบหลังตั้งค่า
git config --list
git config --list --show-origin
```
> Git มีการตั้งค่าได้ 3 ระดับ ซึ่งจะถูกอ่านเรียงต่อกันตามลำดับ
> 1. System: การตั้งค่าสำหรับทุก User บนเครื่องคอมพิวเตอร์นี้
> 2. Global: การตั้งค่าสำหรับ User ของคุณคนเดียว (ไฟล์ ~/.gitconfig)
> 3. Local: การตั้งค่าสำหรับโปรเจกต์นั้นๆ โปรเจกต์เดียว (ไฟล์ .git/config ในโฟลเดอร์โปรเจกต์)

### 3. Git Workflow

##### 3.1 สร้างโฟลเดอร์โปรเจ็กต์ใหม่
```bash
mkdir jenkins-gitlab-n8n-easybuy/basic-git
cd jenkins-gitlab-n8n-easybuy/basic-git
```

##### 3.2 เริ่มต้น git ในโฟลเดอร์
```bash
git init
```

##### 3.3 สร้างไฟล์ greeting.txt
```bash
echo "Greeting message line 1" >> greeting.txt
```

##### 3.4 ตรวจสอบสถานะ
```bash
git status
```

##### 3.5 เพิ่มไฟล์เข้า staging area
```bash
git add greeting.txt
```

##### 3.6 ตรวจสอบสถานะ
```bash
git status
```

##### 3.7 commit ไฟล์
```bash
git commit -m "Initial commit"
```

##### 3.8 ตรวจสอบสถานะ
```bash
git status
```

##### 3.9 แก้ไขไฟล์ greeting.txt
```bash
echo "Greeting message line 2" >> greeting.txt
```

##### 3.10 ตรวจสอบสถานะ
```bash
git status
```

##### 3.11 เพิ่มไฟล์เข้า staging area
```bash
git add greeting.txt
or
git add .
```

##### 3.12 ตรวจสอบสถานะ
```bash
git status
```

##### 3.13 commit ไฟล์
```bash
git commit -m "Add line 2 to greeting.txt"
```

##### 3.14 ตรวจสอบสถานะ
```bash
git status
```

##### 3.13 commit ไฟล์
```bash
git commit -m "Add line 2 to greeting.txt"
```

##### 3.14 ตรวจสอบสถานะ
```bash
git status
```

##### 3.15 ดูประวัติการ commit
```bash
git log
```

##### 3.16 ดูประวัติการ commit แบบย่อ
```bash
git log --oneline
git log --oneline --graph --decorate --all
```

##### 3.17 ดูความแตกต่างของไฟล์
```bash
# ปกติ
git diff HEAD~1 HEAD greeting.txt

# สีสันสวยงาม
git diff --color HEAD~1 HEAD greeting.txt

# แบบมีบริเวณบรรทัดเพิ่ม-ลบ
git diff --color --unified=5 HEAD~1 HEAD greeting.txt
```

##### 3.18 เพิ่มไฟล์ readme.txt
```bash
echo "# Basic Git" >> readme.txt
git add readme.txt
git commit -m "Add readme.txt"
git status
git log --oneline
```

##### 3.19 ย้อนกลับและแก้ไขการเปลี่ยนแปลง (Undo/Repair)

##### วิธีที่ 1: Git Reset --hard (การย้อนเวลาแบบทำลายล้าง) 💣
`git reset --hard <commit>` ทำ 3 อย่างพร้อมกัน:

1. ย้าย HEAD: ย้ายตัวชี้ (HEAD) และ Branch ปัจจุบันกลับไปที่ <commit> ที่ระบุ
2. ล้าง Staging Area: ทำให้ Staging Area (index) ตรงกับสถานะของ <commit> นั้น
3. ล้าง Working Directory: ลบการเปลี่ยนแปลงทั้งหมด ในไฟล์ที่คุณกำลังทำงานอยู่ (ที่ยังไม่ได้ commit) ให้กลับไปเหมือนกับ <commit> นั้น

แสดงภาพกราฟิก
```plaintext
HEAD -> main
 o---o---o---o---o  main (HEAD)
          ^
          |
        <commit>
          |
          +-- Staging Area (index) ถูกล้าง
          |
          +-- Working Directory ถูกล้าง
```

> สรุป: เป็นคำสั่งสำหรับ "ทิ้งทุกอย่าง" ที่ทำหลังจาก <commit> นั้นไป แล้วย้อนเวลากลับไปจุดนั้นอย่างถาวร

```bash
git reset --hard <commit>
git status
git log --oneline
```

##### วิธีที่ 2: Git Reset --soft (การย้อนเวลาแบบเก็บงาน) 🛠️
`git reset --soft <commit>` ทำ 2 อย่าง:

1. ย้าย HEAD: ย้ายตัวชี้ (HEAD) และ Branch ปัจจุบันกลับไปที่ <commit> ที่ระบุ
2. เก็บการเปลี่ยนแปลง: ทำให้การเปลี่ยนแปลงทั้งหมดใน Staging Area (index) ยังคงอยู่

แสดงภาพกราฟิก
```plaintext
HEAD -> main
 o---o---o---o---o  main (HEAD)
          ^
          |
        <commit>
          |
          +-- Staging Area (index) ยังคงอยู่
          |
          +-- Working Directory ยังคงอยู่
```

> สรุป: เป็นคำสั่งสำหรับ "ย้อนกลับ" ไปยัง <commit> ที่ระบุ แต่ยังคงเก็บการเปลี่ยนแปลงทั้งหมดใน Staging Area (index) ไว้

```bash
git reset --soft <commit>
git status
git log --oneline
```

##### วิธีที่ 3: Git Revert (การย้อนกลับแบบสร้าง commit ใหม่) 🔄
`git revert <commit>` สร้าง commit ใหม่ที่ย้อนกลับการเปลี่ยนแปลงที่ทำใน <commit> ที่ระบุ โดยไม่เปลี่ยนแปลงประวัติของ commit

แสดงภาพกราฟิก
```plaintext
HEAD -> main
 o---o---o---o---o  main (HEAD)
          ^
          |
        <commit>
          |
          +-- Staging Area (index) ถูกล้าง
          |
          +-- Working Directory ถูกล้าง
          |
        +---o  main (HEAD) (commit ใหม่ที่ย้อนกลับ)
        |
      <commit>
        |
        +---o  main (HEAD)
```
> สรุป: เป็นคำสั่งสำหรับ "ย้อนกลับ" การเปลี่ยนแปลงที่ทำใน <commit> ที่ระบุ โดยสร้าง commit ใหม่ที่ทำการย้อนกลับนั้น

```bash
git revert <commit>
git status
git log --oneline
```

##### วิธีที่ 4: Git Checkout (การแวะดูอดีต) 🕵️‍♂️

`git checkout <commit>` ทำหน้าที่แตกต่างออกไป:

1. ย้าย HEAD: ย้ายตัวชี้ (HEAD) ไปยัง <commit> ที่ระบุ แต่ ไม่ได้ย้าย Branch ตามไปด้วย สิ่งนี้จะทำให้คุณเข้าสู่สถานะที่เรียกว่า "detached HEAD"

2. อัปเดต Working Directory: ทำให้ไฟล์ใน Working Directory ของคุณตรงกับสถานะของ <commit> นั้น เพื่อให้คุณสามารถดูหรือทดสอบโค้ด ณ จุดนั้นได้

แสดงภาพกราฟิก
```plaintext
HEAD -> main
 o---o---o---o---o  main (HEAD)
          ^
          |
        <commit>
          |
          +-- Staging Area (index) ถูกล้าง
          |
          +-- Working Directory ถูกล้าง
          |
        +---o  main (HEAD) (commit ใหม่ที่ย้อนกลับ)
        |
      <commit>
        |
        +---o  main (HEAD)
```
> สรุป: เป็นคำสั่งสำหรับ "สลับไปดู" โค้ดในอดีตชั่วคราว เมื่อคุณดูเสร็จแล้ว สามารถใช้ git checkout <branch-name> (เช่น git checkout main) เพื่อกลับมาที่ปัจจุบันได้อย่างปลอดภัย โดยที่ประวัติ commit และงานล่าสุดของคุณยังอยู่ครบ

> สรุปถ้า commit ปัจจุบันมี error ในโค้ดและอยากย้ายกลับไปเริ่มใหม่ใน commit ก่อนหน้าแนะนำให้ใช้ git reset --hard <commit> เพื่อย้อนกลับไปยัง commit ที่ต้องการ

```bash
git reset --hard HEAD~1
```

เพราะเป้าหมายของคุณคือ "ทิ้ง commit ล่าสุดที่มีปัญหา และย้อนกลับไปเริ่มต้นใหม่ที่ commit ก่อนหน้านั้น"

`git reset --hard HEAD~1` จะทำสิ่งที่คุณต้องการทุกอย่าง:

- ลบ commit ล่าสุด ที่มี error ออกจากประวัติของ Branch
- ล้างการเปลี่ยนแปลงทั้งหมด ที่มาจาก commit ที่ผิดพลาดนั้นออกจากโค้ดของคุณ
- ทำให้ไฟล์ทั้งหมดในโปรเจกต์ของคุณกลับไปอยู่ในสถานะที่สมบูรณ์ ณ commit ก่อนหน้า พร้อมให้คุณเริ่มทำงานต่อได้ทันที

### 4. เชื่อมต่อและจัดการ Remote (Remotes)
##### 4.1 สร้าง Remote Repository บน GitLab
- ไปที่เว็บไซต์ [GitLab](https://gitlab.com/)
- ลงชื่อเข้าใช้ด้วยบัญชีของคุณ
- คลิกที่ปุ่ม "New Project/repository"
- เลือก "Create blank project"
- กรอกชื่อโปรเจ็กต์ เช่น "basic-git"
- ตั้งค่า Visibility Level เป็น "Private" หรือ "Public" ตามต้องการ
- คลิกที่ปุ่ม "Create project"

##### 4.2 เชื่อมต่อ Local Repository กับ Remote Repository
```bash
git remote add origin https://gitlab.com/your-username/basic-git.git
```

##### 4.3 ตรวจสอบ Remote Repository
```bash
git remote -v
```

##### 4.4 Push การเปลี่ยนแปลงไปยัง Remote Repository
```bash
git branch -M main
git push -u origin main
```

### 5. การจัดการ Branches
##### 5.1 สร้าง Branch ใหม่
```bash
git branch <branch-name>

# หรือ
git checkout -b <branch-name>

# เช่น
git branch develop

# หรือ
git checkout -b develop
```

##### 5.2 สลับไปยัง Branch ที่ต้องการ
```bash
git switch <branch-name>

# หรือ
git checkout <branch-name>

# เช่น
git switch develop
# หรือ
git checkout develop
```

>  git switch (ผู้เชี่ยวชาญด้านการสลับ 🚂)
> git switch ถูกสร้างขึ้นมาใหม่ (ใน Git เวอร์ชัน 2.23) โดยมีวัตถุประสงค์เดียวคือ การสลับ Branch
> หน้าที่ชัดเจน: ใช้สำหรับเปลี่ยนจาก Branch หนึ่งไปยังอีก Branch หนึ่งเท่านั้น
> ปลอดภัยกว่า: git switch จะ ป้องกัน ไม่ให้คุณสลับ Branch หากการสลับนั้นจะทำให้งานที่คุณทำค้างไว้ (แต่ยังไม่ได้ commit) สูญหาย มันจะเตือนให้คุณ commit หรือ stash งานก่อน ซึ่งช่วยลดความผิดพลาดได้มาก

> git checkout (มีดพกสวิส)
> git checkout เป็นคำสั่งที่มีความสามารถหลากหลายมากกว่า มันสามารถใช้ได้ทั้งการสลับ Branch, การสร้าง Branch ใหม่, การแวะดู commit ในอดีต, และการกู้คืนไฟล์จาก commit ก่อนหน้า
> ความสามารถหลากหลาย: git checkout สามารถทำได้หลายอย่างในคำสั่งเดียว ซึ่งทำให้มันมีความซับซ้อนและอาจทำให้เกิดความสับสนได้
> เสี่ยงต่อความผิดพลาด: เนื่องจาก git checkout มีหลายหน้าที่ มันอาจทำให้ผู้ใช้เผลอทำสิ่งที่ไม่ต้องการ เช่น การสลับ Branch โดยไม่ตั้งใจและสูญเสียงานที่ยังไม่ได้ commit

##### 5.3 สร้างไฟล์ใหม่และ commit ใน Branch ใหม่
```bash
touch develop_file.txt
git add develop_file.txt
git commit -m "เพิ่มไฟล์ใหม่ใน Branch develop"
git push -u origin develop
```

##### 5.4 รวม Branch กลับไปยัง main
```bash
git checkout main
git merge develop
```

##### 5.5 ลบ Branch ที่ไม่ต้องการ
```bash
git branch -d develop
```

#### 6. การแก้ไขข้อขัดแย้ง (Merge Conflicts)
##### 6.1 สร้างข้อขัดแย้ง
```bash
git checkout develop
echo "Hello from develop" > conflict.txt
git add conflict.txt
git commit -m "เพิ่มไฟล์ conflict.txt ใน Branch develop"
git push -u origin develop
```

```bash
git checkout main
echo "Hello from main" > conflict.txt
git add conflict.txt
git commit -m "เพิ่มไฟล์ conflict.txt ใน Branch main"
git push origin main
```

##### 6.2 พยายามรวม Branch
```bash
git merge develop
# จะเกิดข้อขัดแย้ง
git status
# แก้ไขไฟล์ conflict.txt
git add conflict.txt
git commit -m "แก้ไขข้อขัดแย้งในไฟล์ conflict.txt"
git push origin main
```
> แนวทางการแก้ไขข้อขัดแย้ง
> 1. เปิดไฟล์ที่มีข้อขัดแย้งในโปรแกรมแก้ไขข้อความ (เช่น VS Code)
> 2. ค้นหาส่วนที่มีข้อขัดแย้ง ซึ่งจะถูกทำเครื่องหมายด้วย `<<<<<<<`, `=======`, และ `>>>>>>>`
> 3. ตัดสินใจเลือกว่าจะเก็บส่วนใด หรือผสมผสานกัน
> 4. ลบเครื่องหมายข้อขัดแย้งทั้งหมดออกจากไฟล์
> 5. บันทึกไฟล์และปิดโปรแกรมแก้ไขข้อความ
> 6. ใช้คำสั่ง `git add <file>` เพื่อทำเครื่องหมายว่าได้แก้ไขข้อขัดแย้งแล้ว
> 7. ใช้คำสั่ง `git commit` เพื่อสร้าง commit ใหม่ที่บันทึกการแก้ไขข้อขัดแย้ง

#### 7. การใช้งาน .gitignore
> ไฟล์ .gitignore ใช้เพื่อบอก Git ว่าไม่ต้องติดตาม (track) หรือไม่ต้องรวมไฟล์หรือโฟลเดอร์บางอย่างในระบบควบคุมเวอร์ชัน
##### 7.1 สร้างไฟล์ .gitignore
```bash
echo "node_modules/" >> .gitignore
echo "*.log" >> .gitignore
echo "dist/" >> .gitignore
echo ".env" >> .gitignore
```
---

## พื้นฐานการใช้งาน Docker

- การติดตั้ง Docker Desktop
- การใช้งานคำสั่งพื้นฐานของ Docker
- การสร้างและจัดการ Container
- การสร้าง Dockerfile และ Docker Image
- การใช้งาน Docker Compose

#### 1. การติดตั้ง Docker Desktop
- ดาวน์โหลดและติดตั้ง Docker Desktop จาก [https://www.docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop)
- ตรวจสอบการติดตั้งโดยรันคำสั่ง:
```bash
docker --version
docker compose version
docker info
```

#### 2. การใช้งานคำสั่งพื้นฐานของ Docker
##### 2.1 รัน hello-world container
```bash
docker run hello-world
```
> คำสั่งนี้จะดาวน์โหลดและรัน container ที่แสดงข้อความต้อนรับจาก Docker

##### 2.2 ดูรายการ container ที่กำลังรันอยู่
```bash
docker ps
docker ps -a  # ดู container ทั้งหมดรวมถึงที่หยุดแล้ว
```

##### 2.3 หยุด container
```bash
docker stop <container_id>
```

##### 2.4 ลบ container
```bash
docker rm <container_id>
```

##### 2.5 ดูรายการ image ที่มีอยู่
```bash
docker images
docker image ls
```
##### 2.6 ลบ image
```bash
docker rmi <image_id>
```

#### 3. การสร้างและจัดการ Container
##### 3.1 รัน container จาก image
```bash
docker run -d -p 8880:80 --name mynginx nginx
```
> คำสั่งนี้จะรัน Nginx container ในโหมด detached และแมปพอร์ต 80 ของ container ไปยังพอร์ต 8880 ของโฮสต์

##### 3.2 เข้าสู่ shell ของ container
```bash
docker exec -it mynginx /bin/bash
# หรือ
docker exec -it mynginx /bin/sh
```

##### 3.3 ดู log ของ container
```bash
docker logs mynginx
```

#### 4. การสร้าง Dockerfile และ Docker Image

##### 4.1 สร้างโฟลเดอร์โปรเจ็กต์
```bash
mkdir basic-docker/docker-node-app
cd basic-docker/docker-node-app
```

##### 4.2 สร้างไฟล์ package.json
```bash
npm init -y
```

##### 4.3 กำหนดสคริปต์ใน package.json
```json
{
  "name": "docker-node-app",
  "version": "1.0.0",
  "description": "A simple Node.js app for Docker",
  "main": "index.js",
  "scripts": {
    "dev": "nodemon index.js",
    "start": "node index.js"
  },
  "author": "Your Name",
  "license": "ISC",
  "dependencies": {
    "express": "^4.18.2",
    "nodemon": "^2.0.22"
  }
}
```
> nodemon เป็นเครื่องมือที่ช่วยในการพัฒนา Node.js โดยจะทำการรีสตาร์ทเซิร์ฟเวอร์อัตโนมัติเมื่อมีการเปลี่ยนแปลงในโค้ด

##### 4.4 สร้างไฟล์ index.js
```javascript
const express = require('express')
const app = express()

// ทำ url ให้สามารถเข้าถึงได้
app.get('/', (req, res) => {
  res.send('Hello World!')
})

// run the server
app.listen(3000, () => {
  console.log('Server is running on http://localhost:3000')
})
```

##### 4.5 สร้างไฟล์ Dockerfile

> ไฟล์ Dockerfile เป็นไฟล์ข้อความธรรมดาที่ใช้กำหนดขั้นตอนการสร้าง Docker Image โดยระบุฐานของ image, การติดตั้ง dependencies, การคัดลอกไฟล์, การตั้งค่าพอร์ต และคำสั่งที่ต้องรันเมื่อ container เริ่มทำงาน

```Dockerfile
# โหลด image ของ node จาก docker hub
FROM node:alpine

# กำหนด directory ที่จะใช้เก็บไฟล์ของโปรเจค
WORKDIR /app

# คัดลอกไฟล์ package.json และ package-lock.json ไปยัง directory ที่กำหนดไว้
COPY package*.json ./

# ติดตั้ง package ที่ระบุในไฟล์ package.json
RUN npm install

# ติดตั้ง nodemon เพื่อใช้ในการรันโปรเจค
RUN npm install -g nodemon

# คัดลอกไฟล์ทั้งหมดไปยัง directory ที่กำหนดไว้
COPY . .

# ระบุ port ที่จะใช้
EXPOSE 3000

# รันคำสั่ง npm run dev เมื่อ container ถูกสร้างขึ้น
CMD ["npm", "run", "dev"]
```

##### 4.6 สร้าง Docker Image
```bash
docker build -t docker-node-app .
```

##### 4.7 รัน Docker Container
```bash
docker run -d -p 3300:3000 --name mydockerapp docker-node-app
```

#### 5. การใช้งาน Docker Compose
##### 5.1 สร้างไฟล์ docker-compose.yml
```yaml
networks:
  nodejs_network:
    name: nodejs_network
    driver: bridge

services:

  # NodeJS App
  nodejs:
    build:
      context: .
      dockerfile: Dockerfile
      tags:
        - "mynodeapp:1.0"
    container_name: mynodeapp
    volumes:
      - .:/app
      - /app/node_modules
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=development
      - CHOKIDAR_USEPOLLING=true # สำหรับ Windows เพื่อให้ nodemon ทำงานได้
    networks:
      - nodejs_network
    restart: always

  # MongoDB
  mongodb:
    image: mongo
    container_name: mongodb
    ports:
      - "28017:27017"
    networks:
      - nodejs_network
    restart: always
    volumes:
      - mongo_data:/data/db

volumes:
  mongo_data:
    name: mongo_data
    driver: local
```

##### 5.2 รัน Docker Compose
```bash
# เช็คความถูกต้องของไฟล์ docker-compose.yml
docker compose config

# รัน Docker Compose
docker compose up -d

# หรือถ้ามีการเปลี่ยนแปลงไฟล์ Dockerfile หรือ docker-compose.yml
docker compose up -d --build
```

##### 5.3 ตรวจสอบสถานะของ Container
```bash
docker compose ps
```

##### 5.4 หยุดและลบ Container
```bash
docker compose down

# หรือถ้าต้องการลบข้อมูลทั้งหมดรวมถึง volume
docker compose down -v

# หรือถ้าต้องการลบ image ด้วย
docker compose down --rmi all -v
```

#### 6. DockerHub (การใช้งาน Docker Hub)
##### 6.1 สร้างบัญชีผู้ใช้บน Docker Hub
- ไปที่ [Docker Hub](https://hub.docker.com/) และสมัครสมาชิก
- สร้าง Repository ใหม่บน Docker Hub
  - คลิกที่ปุ่ม "Create Repository"
  - กรอกชื่อ Repository เช่น `mynodeapp`
  - เลือก Public หรือ Private ตามต้องการ
  - คลิกที่ปุ่ม "Create"

> สำหรับการใช้งาน Docker Hub แนะนำให้ใช้บัญชีแบบ Public เพราะจะง่ายต่อการเข้าถึงและใช้งานร่วมกับผู้อื่น
> เวอ์ร์ชันฟรีของ Docker Hub มีข้อจำกัดในการสร้าง Repository แบบ Private และมีข้อจำกัดในการดึง (pull) image จาก Repository ในระยะเวลาหนึ่ง

##### 6.2 เข้าสู่ระบบ Docker Hub ผ่าน Command Line
```bash
docker login
```

##### 6.3 การ Tag Image และ Push ขึ้น Docker Hub

| คำสั่ง | คำอธิบาย |
|--------|-----------|
| `docker tag <image>:<old_tag> <username>/<repo>:<new_tag>` | สร้างแท็กใหม่ให้ image เพื่อเตรียม push |


```bash
docker tag mynodeapp:1.0 <your-dockerhub-username>/mynodeapp:1.0
```

##### 6.4 ดัน (Push) Docker Image ไปยัง Docker Hub

| คำสั่ง | คำอธิบาย |
|--------|-----------|
| `docker push <username>/<repo>:<tag>` | อัปโหลด image พร้อมแท็กขึ้น Docker Hub |

```bash
docker push <your-dockerhub-username>/mynodeapp:1.0
```

##### 6.5 ดึง (Pull) Docker Image จาก Docker Hub

| คำสั่ง | คำอธิบาย |
|--------|-----------|
| `docker pull <username>/<repo>:<tag>` | ดาวน์โหลด image จาก Docker Hub |

```bash
docker pull <your-dockerhub-username>/mynodeapp:1.0
```
---

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
    image: jenkins/jenkins:jdk21
    container_name: jenkins
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
docker-compose up -d
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

## Workshop Express Docker Application

<img src="https://miro.medium.com/v2/1*Jr3NFSKTfQWRUyjblBSKeg.png" width="200">

โปรเจ็กต์ Express.js + TypeScript REST API ที่ทำงานใน Docker Container พร้อมด้วย CI/CD Pipeline ผ่าน Jenkins, Github Actions และการแจ้งเตือนผ่าน N8N

- [Express Docker Application](#express-docker-application)
  - [📋 Table of Contents](#-table-of-contents)
  - [🚀 Features](#-features)
  - [🏗️ Project Structure](#️-project-structure)
  - [🛠️ Prerequisites](#️-prerequisites)
  - [⚡ Quick Start](#-quick-start)
    - [1. Clone Repository](#1-clone-repository)
    - [2. Local Development](#2-local-development)
    - [3. Docker Development](#3-docker-development)
  - [🐳 Docker Commands](#-docker-commands)
  - [🧪 Testing](#-testing)
  - [🔄 CI/CD Pipeline](#-cicd-pipeline)
  - [⚡ GitHub Actions](#-github-actions)
  - [📡 API Endpoints](#-api-endpoints)
  - [🔧 Configuration](#-configuration)
  - [📝 Environment Variables](#-environment-variables)
  - [📊 Test Coverage](#-test-coverage)
  - [🤝 Contributing](#-contributing)
  - [📄 License](#-license)

## 🚀 Features

- ✅ **Express.js REST API** - Fast, minimal Node.js web framework
- 🔷 **TypeScript Support** - Type-safe JavaScript development
- 🐳 **Docker Containerization** - Lightweight Alpine-based container
- 🧪 **Jest Testing** - Comprehensive test suite with Supertest
- 🔄 **CI/CD Pipeline** - Automated Jenkins pipeline with deployment
- ⚡ **GitHub Actions** - Modern CI/CD with GitHub-native automation
- 📊 **Test Coverage** - Code coverage reports and analysis
- 📦 **Docker Hub Integration** - Automated image publishing
- 🔔 **Server Deployment** - Automated deployment to remote servers
- 🏷️ **Semantic Versioning** - Build number based tagging
- ⚡ **Hot Reload** - Development with ts-node

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

#### 1. เตรียมโปรเจ็กต์ nodejs express docker application

1.1 สร้างโฟลเดอร์โปรเจ็กต์ใหม่
```bash
mkdir express-docker-app
cd express-docker-app
```
1.2 สร้างไฟล์ `package.json`
```bash
npm init -y
```
1.3 กำหนดค่าใน `package.json`
```json
{
  "name": "express-docker-app",
  "version": "1.0.0",
  "description": ",
  "main": "index.js",
  "scripts": {
    "start:ts": "ts-node src/app.ts",
    "start": "node dist/app.js",
    "build": "tsc",
    "test": "jest"
  },
  "keywords": [],
  "author": ",
  "license": "ISC",
  "type": "commonjs",
  "dependencies": {
    "express": "^5.1.0"
  },
  "devDependencies": {
    "@types/express": "^5.0.3",
    "@types/jest": "^30.0.0",
    "@types/node": "^24.5.2",
    "@types/supertest": "^6.0.3",
    "jest": "^30.1.3",
    "supertest": "^7.1.4",
    "ts-jest": "^29.4.4",
    "ts-node": "^10.9.2",
    "typescript": "^5.9.2"
  }
}
```

1.4 กำหนดค่า TypeScript
```bash
npx tsc --init
```

1.5 แก้ไข `tsconfig.json`
```json
{
  "compilerOptions": {
    "target": "es2016",
    "module": "commonjs",
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "strict": true,
    "skipLibCheck": true,
    "outDir": "dist",
    "rootDir": "src",
    "types": ["jest", "node"]
  },
  // คอมไพล์เฉพาะโค้ดแอป; การทดสอบจะถูกแปลงด้วย ts-jest แยกต่างหาก
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

1.6 สร้างโฟลเดอร์ `src` และไฟล์ `app.ts`
```bash
mkdir src
touch src/app.ts
```

1.7 เพิ่มโค้ดตัวอย่างใน `src/app.ts`
```typescript
import express, { type Express, type Request, type Response } from 'express'

const app: Express = express()

const port: number = 3000

// Routes
// GET /
app.get('/', (_: Request, res: Response) => {
  res.json({
    message: 'Hello Express + TypeScript!'
  })
})

// GET /api/hello
app.get('/api/hello', (_: Request, res: Response) => {
  res.json({
    message: 'Hello from Express API!'
  })
})

// GET /api/health
app.get('/api/health', (_: Request, res: Response) => {
  res.json({
    status: 'UP'
  })
})

// Start server
app.listen(port, () => console.log(`Application is running on port ${port}`))
```

1.8 สร้าง tests folder และไฟล์ `app.test.ts`
```bash
mkdir tests
touch tests/app.test.ts
```

1.9 เพิ่มโค้ดตัวอย่างใน `tests/app.test.ts`
```typescript
import request from 'supertest'
import express from 'express'

const app = express()
app.get('/api/hello', (_, res) => res.json({ message: 'Hello from Express API!' }))

test('GET /api/hello', async () => {
  const res = await request(app).get('/api/hello')
  expect(res.statusCode).toBe(200)
  expect(res.body.message).toBe('Hello from Express API!')
})
```
1.10 กำหนดค่า Jest โดยสร้างไฟล์ `jest.config.js`
```javascript
module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src', '<rootDir>/tests'],
  testMatch: ['**/__tests__/**/*.ts', '**/?(*.)+(spec|test).ts'],
  transform: {
    '^.+\\.ts$': 'ts-jest',
  },
  collectCoverageFrom: [
    'src/**/*.ts',
    '!src/**/*.d.ts',
  ],
}
```
1.11 สร้างไฟล์ `.gitignore`
```bash
touch .gitignore
```
1.12 เพิ่มโค้ดใน `.gitignore`
```
node_modules
dist
.env
```


#### 2. เตรียม Dockerfile สำหรับโปรเจ็กต์

2.1 สร้างไฟล์ `Dockerfile`
```bash
touch Dockerfile
```

2.2 เพิ่มโค้ดใน `Dockerfile`
```Dockerfile
# Build stage - สำหรับ development และ testing
FROM node:22-alpine AS builder

# กำหนด Working Directory ภายใน Container
WORKDIR /app

# Copy ไฟล์ package.json และ package-lock.json เข้าไปก่อน
# เพื่อใช้ประโยชน์จาก Docker cache layer ทำให้ไม่ต้อง install dependencies ใหม่ทุกครั้งที่แก้โค้ด
COPY package*.json ./

# ติดตั้ง Dependencies (รวม dev dependencies สำหรับ testing)
RUN npm install

# Copy โค้ดทั้งหมดในโปรเจกต์เข้าไปใน container
COPY . .

# Compile TypeScript เป็น JavaScript
RUN npm run build

# Production stage - สำหรับ production deployment
FROM node:22-alpine AS production

# กำหนด Working Directory ภายใน Container
WORKDIR /app

# Copy package files
COPY package*.json ./

# ติดตั้งเฉพาะ production dependencies
RUN npm ci --only=production && npm cache clean --force

# Copy โค้ดที่ compiled แล้วจาก builder stage
COPY --from=builder /app/dist ./dist
# COPY --from=builder /app/src ./src

# กำหนด Port ที่ Container จะทำงาน
EXPOSE 3000

# คำสั่งสำหรับรัน Express Application (ใช้ compiled JavaScript)
CMD ["npm", "start"]
```

2.3 สร้างไฟล์ `.dockerignore`
```bash
# Dependencies
node_modules
npm-debug.log*

# Build outputs
dist
build

# Environment files
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# Testing
coverage
*.lcov

# Git
.git
.gitignore

# Docker
Dockerfile
.dockerignore

# Documentation
README.md
*.md

# IDE
.vscode
.idea
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
logs
*.log

# Temporary files
.tmp
.temp
```

#### 3. สร้าง Jenkinsfile สำหรับ CI/CD Pipeline
```bash
touch Jenkinsfile
```

3.1 เพิ่มโค้ดใน `Jenkinsfile`
```groovy
pipeline {
  
    // ใช้ any agent เพื่อหลีกเลี่ยงปัญหา Docker path mounting บน Windows
    agent any

    // กำหนด environment variables
    environment {
        // ใช้ค่าเป็น "credentialsId" ของ Jenkins โดยตรงสำหรับ docker.withRegistry
        DOCKER_HUB_CREDENTIALS_ID = 'dockerhub-cred'
        DOCKER_REPO = "iamsamitdev/express-docker-app"
        APP_NAME = "express-docker-app"
    }

    // กำหนด stages ของ Pipeline
    stages {

        // Stage 1: ดึงโค้ดล่าสุดจาก Git
        stage('Checkout') {
            steps {
                echo "Checking out code..."
                checkout scm
            }
        }

        // Stage 2: ติดตั้ง dependencies และรันเทสต์ (รองรับทุก Platform)
        stage('Install & Test') {
            steps {
                script {
                    // ตรวจสอบว่ามี Node.js บน host หรือไม่
                    def hasNodeJS = false
                    def isWindows = isUnix() ? false : true
                    
                    try {
                        if (isWindows) {
                            bat 'node --version && npm --version'
                        } else {
                            sh 'node --version && npm --version'
                        }
                        hasNodeJS = true
                        echo "Using Node.js installed on ${isWindows ? 'Windows' : 'Unix'}"
                    } catch (Exception e) {
                        echo "Node.js not found on host, using Docker"
                        hasNodeJS = false
                    }
                    
                    if (hasNodeJS) {
                        // ใช้ Node.js บน host
                        if (isWindows) {
                            bat '''
                                npm install
                                npm test
                            '''
                        } else {
                            sh '''
                                npm install
                                npm test
                            '''
                        }
                    } else {
                        // ใช้ Docker run command (รองรับทุก platform)
                        if (isWindows) {
                            bat '''
                                docker run --rm ^
                                -v "%cd%":/workspace ^
                                -w /workspace ^
                                node:22-alpine sh -c "npm install && npm test"
                            '''
                        } else {
                            sh '''
                                docker run --rm \\
                                -v "$(pwd)":/workspace \\
                                -w /workspace \\
                                node:22-alpine sh -c "npm install && npm test"
                            '''
                        }
                    }
                }
            }
        }

        // Stage 3: สร้าง Docker Image สำหรับ production
        stage('Build Docker Image') {
            steps {
                script {
                    echo "Building Docker image: ${DOCKER_REPO}:${BUILD_NUMBER}"
                    docker.build("${DOCKER_REPO}:${BUILD_NUMBER}", "--target production .")
                }
            }
        }

        // Stage 4: Push Image ไปยัง Docker Hub
        stage('Push Docker Image') {
            steps {
                script {
                    // ต้องส่งค่าเป็น credentialsId เท่านั้น ไม่ใช่ค่าที่ mask ของ credentials()
                    docker.withRegistry('https://index.docker.io/v1/', env.DOCKER_HUB_CREDENTIALS_ID) {
                        echo "Pushing image to Docker Hub..."
                        def image = docker.image("${DOCKER_REPO}:${BUILD_NUMBER}")
                        image.push()
                        image.push('latest')
                    }
                }
            }
        }

        // Stage 5: เคลียร์ Docker images และ cache บน agent
        stage('Cleanup Docker') {
            steps {
                script {
                    def isWindows = isUnix() ? false : true
                    echo "Cleaning up local Docker images/cache on agent..."
                    if (isWindows) {
                        bat ""
                            docker image rm -f ${DOCKER_REPO}:${BUILD_NUMBER} || echo ignore
                            docker image rm -f ${DOCKER_REPO}:latest || echo ignore
                            docker image prune -af -f
                            docker builder prune -af -f
                        ""
                    } else {
                        sh ""
                            docker image rm -f ${DOCKER_REPO}:${BUILD_NUMBER} || true
                            docker image rm -f ${DOCKER_REPO}:latest || true
                            docker image prune -af -f
                            docker builder prune -af -f
                        ""
                    }
                }
            }
        }

        // Stage 6: Deploy ไปยังเครื่อง local (รองรับทุก Platform)
        stage('Deploy Local') {
            steps {
                script {
                    def isWindows = isUnix() ? false : true
                    echo "Deploying container ${APP_NAME} from latest image..."
                    if (isWindows) {
                        bat ""
                            docker pull ${DOCKER_REPO}:latest
                            docker stop ${APP_NAME} || echo ignore
                            docker rm ${APP_NAME} || echo ignore
                            docker run -d --name ${APP_NAME} -p 3000:3000 ${DOCKER_REPO}:latest
                            docker ps --filter name=${APP_NAME} --format \"table {{.Names}}\t{{.Image}}\t{{.Status}}\"
                        ""
                    } else {
                        sh ""
                            docker pull ${DOCKER_REPO}:latest
                            docker stop ${APP_NAME} || true
                            docker rm ${APP_NAME} || true
                            docker run -d --name ${APP_NAME} -p 3000:3000 ${DOCKER_REPO}:latest
                            docker ps --filter name=${APP_NAME} --format "table {{.Names}}\t{{.Image}}\t{{.Status}}"
                        ""
                    }
                }
            }
        }
    }
}
```

#### 4. Push โค้ดขึ้น GitLab

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://gitlab.com/your-username/your-repo.git
git push -u origin main
```

#### 5. สร้าง Jenkins Job แบบ Pipeline
5.1 ไปที่ Jenkins Dashboard
5.2 คลิกที่ "New Item"
5.3 ตั้งชื่อ Job และเลือก "Pipeline" จากนั้นคลิก "OK"
5.4 ในส่วน "Pipeline" ให้เลือก "Pipeline script from SCM"
5.5 ตั้งค่า SCM เป็น "Git" และกรอก URL ของ GitLab Repo
5.6 ตั้งค่า "Branch Specifier" เป็น "*/main"
5.7 ในส่วน "Script Path" ให้ระบุที่อยู่ของ Jenkinsfile (เช่น `Jenkinsfile`)
5.8 คลิก "Save"
5.9 คลิก "Build Now" เพื่อทดสอบการทำงานของ Pipeline

#### 6. การตั้งค่า Webhook ใน GitLab
**จำเป็นต้องติดตั้ง plugin "GitLab Plugin, Git Plugin, GitLab API Plugin" ใน Jenkins ก่อน**

6.1 ไปที่หน้า GitLab Repository ของคุณ
6.2 คลิกที่ "Settings" > "Webhooks" > "Add webhook"
6.3 ในช่อง "URL" ให้ใส่ URL ของ Jenkins Server (เช่น `https://sweet-dolls-rush.loca.lt/project/test-jenkinsfile-gitlab`)
6.4 Trigger: เลือก "Push events" หรือ event อื่น ๆ ตามต้องการ
6.5 คลิก "Add webhook"

## สิ่งที่เรียนรู้ใน Day 1
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