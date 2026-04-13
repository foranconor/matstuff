#!/bin/bash
echo "Deployment started at $(date)"
cd ../
echo "Build the binary"
go build -o ./main ./cmd/
echo "Build the frontend assets"
cd web
npm run build
cd ..
echo "Replace the unit files"
sudo cp -v deploy/materials.service /etc/systemd/system/materials.service
echo "Restart the daemon"
sudo systemctl daemon-reload
sudo systemctl restart materials.service
echo "Deployment finished at $(date)"
