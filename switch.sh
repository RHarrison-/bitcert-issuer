#!/bin/bash
if [ -e bind.go ]
then
    mv ./bind.go ./bind.txt 
    mv ./bind_linux_amd64.txt ./bind_linux_amd64.go
else
    mv ./bind.txt ./bind.go 
    mv ./bind_linux_amd64.go ./bind_linux_amd64.txt
fi
