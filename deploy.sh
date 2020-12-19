#!/bin/bash

if [[ $1 == vue ]]; then
    echo "Building and deploying the Vue application..."
    cd frontend
    yarn build
    ssh -i ~/.ssh/mcculleydj.dev.pem ubuntu@54.203.119.188 rm -rf /home/ubuntu/cosign
    scp -r -i ~/.ssh/mcculleydj.dev.pem dist ubuntu@54.203.119.188:/home/ubuntu/cosign
    ssh -i ~/.ssh/mcculleydj.dev.pem ubuntu@54.203.119.188 sudo systemctl restart nginx
    cd ..
elif [[ $1 == api ]]; then
    echo "Building and deploying the Go API..."
    cd backend
    GOOS=linux GOARCH=amd64 go build cmd/api/main.go
    ssh -i ~/.ssh/mcculleydj.dev.pem ubuntu@54.203.119.188 sudo supervisorctl stop cosign-api
    scp -i ~/.ssh/mcculleydj.dev.pem main ubuntu@54.203.119.188:~/
    ssh -i ~/.ssh/mcculleydj.dev.pem ubuntu@54.203.119.188 sudo supervisorctl start cosign-api
    rm main
    cd ..
fi