#!/bin/bash
DIR=.

exist_file()
{
    if [ -e "$1" ]
    then
        return 1
    else
        return 2
    fi
}

for k in $(ls $DIR)
do
  echo $k
  [ -d $k ] && cd $k 

  exist_file *.go
  value=$?
  if [ $value -eq 1 ]
  then
    echo "build:$k"
    go build
  fi
 
  cd ..
done





