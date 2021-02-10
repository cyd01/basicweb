#!/bin/bash

if [ $REQUEST_METHOD == "GET" ]; then
  env | sort
else

  if [ $HTTP_CONTENT_LENGTH -gt 0 ]; then
    echo "Here is the body:"
    #read -n $HTTP_CONTENT_LENGTH body
    #echo $body
    cat
  else
    echo "There is no body"
  fi
fi
