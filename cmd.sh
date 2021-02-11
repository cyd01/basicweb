#!/bin/bash

if [ $REQUEST_METHOD == "GET" ]; then
  echo "Status: 200"
  echo "Content-type: text/html"
  echo
  env | sort
else

  if [ $HTTP_CONTENT_LENGTH -gt 0 ]; then
    echo "Status: 200"
    echo "Content-type: text/html"
    echo
    echo "Here is the body:"
    #read -n $HTTP_CONTENT_LENGTH body
    #echo $body
    cat
  else
    echo "Status: 200"
    echo "Content-type: text/html"
    echo
    echo "There is no body"
  fi
fi
