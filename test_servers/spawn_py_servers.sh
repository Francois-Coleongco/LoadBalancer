#!/bin/bash

for i in $(cat ./port_list)
do
    python3 ./main.py -p ${i} &
done
