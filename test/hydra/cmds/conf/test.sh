#!/bin/bash

EXE_FILE=$1

echo "-----testcase1: conf_show"

cd ../testapp

echo "执行:./$EXE_FILE conf show -r lm://."
echo q | eval "./$EXE_FILE conf show -r lm://."

sleep 2
#------------------------------------------
echo "-----testcase3: conf_install"

echo "执行:./$EXE_FILE conf install -r lm://."
eval "./$EXE_FILE conf install -r lm://."
sleep 2

#------------------------------------------
echo "-----testcase3: conf_install-cover"
echo "执行:./$EXE_FILE conf install -r lm://. -cover true"
eval "./$EXE_FILE conf install -r lm://. -cover true"
sleep 2




