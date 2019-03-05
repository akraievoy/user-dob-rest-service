#!/bin/bash

trap trap_exit SIGINT SIGTERM

# LATER if you kill the test with sigterm - it still DOES PASS and trap_exit output is NOT observed at all
trap_exit() {
  echo "Caught Signal ... test has failed."
  trap - cleanup SIGINT SIGTERM # clear the trap
  kill -- -$$ # Sends SIGTERM sto child/sub processes
  exit 1
}

RETRIES=0
RETRIES_MAX=24

Q0="select 1;"

RES=$(echo "${Q1}" | PGPASSWORD=bazinga psql -tz0 -h postgres -U upserter -d upserter )
Q_STATUS=$?
while test ${Q_STATUS} -ne 0 ; do
    if test ${RETRIES} -ge ${RETRIES_MAX} ; then
      echo retries max reached
      exit 1
    fi
    RES=$(echo "${Q1}" | PGPASSWORD=bazinga psql -tz0 -h postgres -U upserter -d upserter )
    Q_STATUS=$?
    RETRIES=$((RETRIES + 1))
    echo "sleeping on query 0: status=${Q_STATUS}, retries=${RETRIES}"
    sleep 3
done

echo "query 0 PASSED"

Q1="select count(*) from users;"
E1=100

RES=$(echo "${Q1}" | PGPASSWORD=bazinga psql -tz0 -h postgres -U upserter -d upserter )
Q_STATUS=$?
while test ${Q_STATUS} -ne 0 -o ${RES} -lt ${E1} ; do
    if test ${RETRIES} -ge ${RETRIES_MAX} ; then
      echo retries max reached
      exit 1
    fi
    RES=$(echo "${Q1}" | PGPASSWORD=bazinga psql -tz0 -h postgres -U upserter -d upserter )
    Q_STATUS=$?
    RETRIES=$((RETRIES + 1))
    echo "sleeping on query 1: status=${Q_STATUS}, RES=${RES}, retries=${RETRIES}"
    sleep 3
done

echo "query 1 PASSED with RES=${RES}"

Q2="select count(*) from users where id=33 and name='Damien';"
E2=1

RES=$(echo "${Q2}" | PGPASSWORD=bazinga psql -tz0 -h postgres -U upserter -d upserter )
Q_STATUS=$?
while test ${Q_STATUS} -ne 0 -o ${RES} -lt ${E2} ; do
    if test ${RETRIES} -ge ${RETRIES_MAX} ; then
      echo retries max reached
      exit 1
    fi
    RES=$(echo "${Q2}" | PGPASSWORD=bazinga psql -tz0 -h postgres -U upserter -d upserter )
    Q_STATUS=$?
    RETRIES=$((RETRIES + 1))
    echo "sleeping on query 2: status=${Q_STATUS}, RES=${RES}, retries=${RETRIES}"
    sleep 3
done

echo "query 2 PASSED with RES=${RES}"
