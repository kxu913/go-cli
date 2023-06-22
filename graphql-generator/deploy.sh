#!/bin/bash

echo "0 - check whether init or not"
if [-z "go.work"]; then
  echo "Already initialized."
else
  chmod +x init.sh
  ./init.sh

fi

echo "1 - Containize project"
docker build -t kevin_913/graphql-generator .
echo "Containize end..."
# cd scripts
# echo "2 - Create resources."
# kubectl apply -f service.yaml
# echo "All done"
docker run -it kevin_913/graphql-generator

