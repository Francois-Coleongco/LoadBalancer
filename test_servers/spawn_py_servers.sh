#!/bin/bash

# note that these ports must be unused so flask can work with it

for i in $(cat ./port_list)
do
    python3 ./main.py -p ${i} &
done
