 
# DevOps Jenkins  & GitHub Actions +php สร้างคู่มือ แบบละเอียด Step by step

เริ่มต้นทำ CI/CD ของ PHP ด้วย Jenkins และ GitHub Actions ทำได้โดยเตรียมโปรเจกต์ PHP ให้รองรับการรันทดสอบ แล้วค่อยตั้งค่า Pipeline บน Jenkins และ Workflow บน GitHub Actions ให้ทำงานอัตโนมัติเมื่อมีการ push โค้ดขึ้น GitHub. ด้านล่างเป็นคู่มือแบบ Step by step แยกเป็น 2 ส่วน: Jenkins และ GitHub Actions สำหรับโปรเจกต์ PHP ทั่วไป (เช่น Laravel หรือ PHP ธรรมดา).[^1][^2][^3][^4][^5][^6]

***

## เตรียมโปรเจกต์ PHP + GitHub

1. สร้าง/เตรียมโปรเจกต์ PHP
    - ใช้ Git จัดการซอร์สโค้ด แล้วสร้าง repo บน GitHub.
    - แนะนำให้ใช้ Composer จัดการ dependency และใช้ PHPUnit/ Pest หรือ framework test อื่นสำหรับรันทดสอบอัตโนมัติ.[^3][^1]
2. เพิ่มไฟล์พื้นฐานในโปรเจกต์
    - `composer.json` (กำหนด dependency และสคริปต์ เช่น `"test": "phpunit"`).
    - โฟลเดอร์ `tests/` และไฟล์ทดสอบเบื้องต้น เพื่อให้ CI รันแล้วไม่ล้ม.[^5][^1]
3. จัดการ secret / ตัวแปร environment
    - ใน GitHub ให้ไปที่ Settings → Secrets and variables → Actions แล้วเพิ่มค่าที่ต้องใช้ เช่น `DB_HOST`, `DB_PASSWORD`, `SSH_KEY` (ถ้าต้อง deploy ผ่าน SSH).
    - ใน Jenkins ให้เพิ่ม Credentials (เช่น SSH username/password หรือ key, GitHub token) ที่จะใช้ใน Pipeline.[^7][^8]

***

## Jenkins: ติดตั้งและเชื่อมต่อกับ GitHub

1. ติดตั้ง Jenkins และปลั๊กอินหลัก
    - ติดตั้ง Jenkins บนเซิร์ฟเวอร์ (Linux/Windows ก็ได้) แล้วเข้าเว็บ UI.
    - แนะนำให้ติดตั้งปลั๊กอิน: Git, GitHub, Pipeline (และถ้าใช้ PHP CodeSniffer/ PHPStan ก็มีปลั๊กอินเสริมได้).[^9][^1][^7]
2. ติดตั้ง PHP และเครื่องมือบน Jenkins agent
    - บนเครื่องที่ Jenkins ใช้รัน job ให้ติดตั้ง PHP, Composer, PHPUnit/เครื่องมือทดสอบ, git ให้เรียบร้อย.
    - ถ้าใช้ Docker สามารถสร้าง image ที่มี PHP + Composer + extensions ที่ต้องใช้สำหรับโปรเจกต์ PHP แล้วให้ Jenkins ใช้ image นั้นเป็น agent.[^1][^3]
3. เชื่อม Jenkins กับ GitHub (webhook)
    - ใน GitHub repo: ไปที่ Settings → Webhooks → Add webhook ใส่ `http(s)://<JENKINS_URL>/github-webhook/` เป็น Payload URL และเลือก events สำหรับ Push/ Pull request.[^8][^10][^7]
    - ใน Jenkins:
        - ไปที่ Manage Jenkins → Configure System กำหนด GitHub server (ถ้าต้องใช้ OAuth/Token).
        - ใน job/pipeline ให้ติ๊ก “GitHub hook trigger for GITScm polling” เพื่อให้ Jenkins trigger จาก webhook ได้.[^11][^7]
4. สร้าง Pipeline Job สำหรับ PHP
    - ใน Jenkins เลือก “New Item” → ตั้งชื่อ → เลือก “Pipeline”.
    - ในส่วน Pipeline เลือก “Pipeline script from SCM” เพื่อดึง Jenkinsfile จาก GitHub repo แล้วใส่ URL repo และ branch ที่ต้องการ.[^3][^9]
5. เขียน Jenkinsfile (ตัวอย่างโครงสำหรับ PHP)
สร้างไฟล์ `Jenkinsfile` ใน root ของโปรเจกต์ เช่น:

```groovy
pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/your-org/your-php-repo.git'
            }
        }

        stage('Install dependencies') {
            steps {
                sh 'composer install --no-interaction --prefer-dist'
            }
        }

        stage('Run tests') {
            steps {
                sh './vendor/bin/phpunit'
            }
        }

        stage('Build/Prepare artifacts') {
            steps {
                sh 'php -l index.php'
            }
        }

        stage('Deploy to server') {
            when {
                branch 'main'
            }
            steps {
                sh '''
                # ตัวอย่าง deploy ผ่าน SSH
                rsync -avz --delete ./ user@your-server:/var/www/yourapp
                '''
            }
        }
    }
}
```

    - แบ่ง stage ให้ชัดเจน: Checkout → Install → Test → Build → Deploy จะช่วยให้ debug ง่าย.
    - ปรับคำสั่ง `sh` ให้ตรงกับ framework ที่ใช้ (เช่น Laravel: `php artisan test`, `php artisan migrate`, ฯลฯ).[^2][^1][^3]
6. ทดสอบ Pipeline บน Jenkins
    - กด “Build Now” ครั้งแรกเพื่อดูว่า Jenkinsfile ทำงานครบทุก stage หรือไม่.
    - แก้ไข permission, path, และคำสั่งต่างๆ ถ้ามี error จาก PHP/Composer หรือ SSH.[^12][^13][^5]

***

## GitHub Actions: CI สำหรับ PHP

1. สร้าง Workflow file
    - ในโปรเจกต์ สร้างโฟลเดอร์ `.github/workflows` แล้วสร้างไฟล์ เช่น `ci-php.yml`.
    - ไฟล์นี้คือที่กำหนดว่าเวลา push/pull request ให้รันอะไรบ้าง.[^4][^6]
2. โครงพื้นฐานของ workflow (PHP)

```yaml
name: PHP CI

on:
  push:
    branches: [ "main", "develop" ]
  pull_request:
    branches: [ "main", "develop" ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.2'
          extensions: mbstring, pdo_mysql
          coverage: none

      - name: Install dependencies
        run: composer install --no-interaction --prefer-dist

      - name: Run tests
        run: ./vendor/bin/phpunit
```

    - ใช้ action `shivammathur/setup-php` เพื่อเตรียม PHP + extensions ได้สะดวก.
    - ปรับ `php-version` และ extensions ให้ตรงตามโปรเจกต์ เช่น Laravel ต้องใช้ pdo_mysql, mbstring ฯลฯ.[^14][^15][^4]
3. เพิ่มขั้นตอนคุณภาพโค้ด (ถ้าต้องการ)
    - เพิ่ม step สำหรับ `phpcs` (PHP_CodeSniffer) หรือ `phpstan` เพื่อตรวจ style/ static analysis.
    - ตัวอย่าง: ติดตั้งด้วย Composer แล้วรัน `vendor/bin/phpcs` หรือ `vendor/bin/phpstan analyse` ใน workflow.[^1][^3]
4. เพิ่มขั้นตอน Deploy ด้วย GitHub Actions (เลือกได้)
    - ถ้าต้องการให้ deploy อัตโนมัติหลังจากเทสผ่าน ให้เพิ่ม job หรือ step ใหม่ที่รันเฉพาะบน branch main.
    - ใช้ action จาก Marketplace เช่น action ที่ deploy ผ่าน SSH หรือ rsync แล้วใช้ secret (เช่น `SSH_HOST`, `SSH_USER`, `SSH_KEY`) จาก GitHub Secrets.[^16][^6][^17]

***

## การออกแบบ Flow: Jenkins + GitHub Actions ใช้ร่วมกันยังไง

1. ใช้ Jenkins เป็น CI/CD หลัก + GitHub Actions เป็น pre-check
    - ใช้ GitHub Actions รัน test แบบเร็วๆ ทุก push/PR เพื่อ feedback ไว (unit test, lint).
    - ให้ Jenkins เน้นงาน build ใหญ่, integration test, และ deploy ไปยังหลาย environment (staging/production).[^2][^9][^14]
2. หรือใช้ GitHub Actions ทั้ง CI+CD แล้วให้ Jenkins ทำงานเฉพาะใน infra ภายใน
    - กรณีที่ใช้ Jenkins ภายในองค์กรสำหรับงานที่ต้องเข้าถึงระบบภายใน (เช่น database ภายใน, private network).
    - GitHub Actions อาจ trigger Jenkins ผ่าน webhook หรือ API เพื่อต่อ pipeline ระหว่างกัน.[^7][^8][^11]

ถ้าต้องการ ตัวอย่างจริงสำหรับ Laravel หรือโครงสร้างไฟล์ทั้ง repo (เช่น ตัวอย่าง `composer.json`, `phpunit.xml`, และ workflow ที่รองรับ MySQL) ระบุ framework (Laravel/CodeIgniter/Plain PHP) และวิธี deploy (เช่น Apache, Nginx, Docker, Kubernetes) แล้วจะช่วยออกแบบ Pipeline และไฟล์ทุกอย่างให้ครบเป็นโครงโปรเจกต์ได้.
<span style="display:none">[^18][^19][^20]</span>

<div align="center">⁂</div>

[^1]: https://www.jenkins.io/solutions/php/

[^2]: https://www.geeksforgeeks.org/devops/setting-up-a-ci-cd-pipeline-for-php-applications-using-jenkins/

[^3]: https://dev.to/kennibravo/pipeline-101-for-php-and-laravel-projects-with-jenkins-and-docker-47gh

[^4]: https://ma.ttias.be/a-github-ci-workflow-tailored-to-laravel-applications/

[^5]: https://www.somkiat.cc/php-continuous-integration-with-jenkins/

[^6]: https://docs.github.com/actions/using-workflows/workflow-syntax-for-github-actions

[^7]: https://plugins.jenkins.io/github/

[^8]: https://www.blazemeter.com/blog/how-to-integrate-your-github-repository-to-your-jenkins-project

[^9]: https://www.jenkins.io/doc/pipeline/tour/hello-world/

[^10]: https://kodekloud.com/community/t/jenkins-github-webhook-integration/292526

[^11]: https://dzone.com/articles/adding-a-github-webhook-in-your-jenkins-pipeline

[^12]: https://www.youtube.com/watch?v=5GtH-nDEEK8

[^13]: https://www.youtube.com/watch?v=64VkDEiLCY8

[^14]: https://dev.to/rafi021/step-by-step-guide-laravel-cicd-with-github-actions-i89

[^15]: https://www.ducxinh.com/en/techblog/automating-laravel-cicd-with-github-actions

[^16]: https://github.com/marketplace/actions/deploy-php-actions

[^17]: https://stackoverflow.com/questions/62447902/env-files-in-github-actions-ci-cd-workflows-how-to-provide-these-into-the-work

[^18]: https://dev.to/ehtesham_ali_abc367f36a5b/jenkins-with-php-run-your-first-pipeline-16e7

[^19]: https://gist.github.com/471ae1a98267e20530d989f64f5290ee

[^20]: https://www.youtube.com/watch?v=KxAElkZ1Hs4

