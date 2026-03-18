# แนวทางการติดตั้งใช้งาน

## การสร้างในเครื่อง

```bash

./scripts/build.sh

```

## อิมเมจ Docker

```bash

./scripts/docker-build.sh

```

## Docker Compose

```bash

docker compose -f deploy/docker/docker-compose.yaml up -d

```

## Kubernetes
โปรดดูที่ `deploy/k8s/deployment.yaml` การกำหนดค่าจำเป็นต้องปรับให้เหมาะสมกับสภาพแวดล้อม:

- ข้อมูลการเชื่อมต่อฐานข้อมูล

- ข้อมูลการเชื่อมต่อ Redis

- คีย์ JWT

## คำแนะนำในการกำหนดค่า

- โดยค่าเริ่มต้น จะอ่าน `configs/config.yaml` ซึ่งสามารถระบุได้ผ่าน `CONFIG_PATH`

- ในสภาพแวดล้อมการผลิต ควรให้ความสำคัญกับการแทนที่การกำหนดค่าผ่านตัวแปรสภาพแวดล้อม (ขึ้นต้นด้วย `APP_` เช่น `APP_DB_HOST`)

- Redis เป็นตัวเลือกเสริม สามารถปิดใช้งานได้โดยใช้ `APP_REDIS_ENABLED=false`