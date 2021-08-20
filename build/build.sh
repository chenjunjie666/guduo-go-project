#!/bin/bash

b_path="$PWD/../"
function script(){
  /bin/bash $1
}
function travFolder(){
  flist=`ls $1`
  cd $dir_
  #echo $flist
  for f in $flist
  do
    if [ -d "$b_path/$f" ]
    then
      #echo "dir:$f"
      script travFolder $f
      echo "dir:  $f"
    else
      echo "file:  $f"
      script $f
      echo $f
    fi
  done
}
travFolder "$PWD/../"