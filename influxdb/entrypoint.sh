#!/bin/sh
set -e

# Run the original InfluxDB entrypoint in the background so the DB becomes
# available, then create the additional buckets once it is up.

/entrypoint.sh "$@" &
ENTRYPOINT_PID=$!

# Wait for the InfluxDB HTTP API to be ready
until curl -sf http://localhost:8086/health >/dev/null 2>&1; do
  sleep 1
done

# Give the setup a moment to complete before creating buckets
sleep 2

INFLUX_HOST=http://localhost:8086
INFLUX_ORG="${DOCKER_INFLUXDB_INIT_ORG:-icmon_org}"
INFLUX_TOKEN="${DOCKER_INFLUXDB_INIT_ADMIN_TOKEN}"

BUCKET_NAMES="AIRCOM2 AIRCOM3 AIRCOM4 AIRCOM5 AIRCOM6 AIRCOM7 AIRCOM8 AIRCOM9 AIRCOM10 AIRCOM11 AIRCOM12 AIRCOM13 AIRCOM14 AIRCOM15 AIRCOM16 AIRCOM17 AIRCOM18 AIRCOM19 AIRCOM20 BAACTW01 BAACTW02 BAACTW03 BAACTW04 BAACTW05 BAACTW06 BAACTW07 BAACTW08 BAACTW09 BAACTW10 BAACTW11 BAACTW12 BAACTW13 BAACTW14 BAACTW15 BAACTW16 BAACTW17 BAACTW18 BAACTW19 BAACTW20 CMONBUGKET CMONBUGKET1 CMONBUGKET2 CMONBUGKET3 CMONBUGKET4 CMONBUGKET5 CMONBUGKET6 CMONBUGKET7 CMONBUGKET8 CMONBUGKET9 CMONBUGKET10 SORACELL1 SORACELL2 SORACELL3 SORACELL4 SORACELL5 SORACELL6 SORACELL7 SORACELL8 SORACELL9 SORACELL10"

for bucket in $BUCKET_NAMES; do
  echo "Creating bucket: ${bucket}"
  influx bucket create \
    --host "${INFLUX_HOST}" \
    --org "${INFLUX_ORG}" \
    --token "${INFLUX_TOKEN}" \
    --name "${bucket}" \
    --retention 0 \
    --skip-verify 2>/dev/null || true
done

echo "All buckets created."

wait "$ENTRYPOINT_PID"
