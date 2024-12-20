# JoBros Consumer service platform

JoBros is a platform for users to offer their services and for users hire services from other users.

## Getting started

Build and push the Docker image:
```sh
docker build -t jobros-api:latest .
```

Deploy the service using the following command:
```sh
# Add Bitnami repo for MongoDB dependency
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Create a Kubernetes secret for MongoDB credentials
kubectl create secret generic mongodb-credentials \
  --from-literal=mongodb-root-password='<your-mongodb-root-password>' \
  --from-literal=mongodb-user-password='<your-mongodb-user-password>' \
  --from-literal=mongodb-uri='mongodb://<user>:<password>@<mongodb-service>:<port>/<database>'

# Create a Kubernetes secret for JWT secret key
kubectl create secret generic jwt-secret --from-literal=JWT_SECRET_KEY='your-secure-random-key'

# Install MongoDB
helm install jobros-mongodb bitnami/mongodb \
  --set auth.existingSecret=mongodb-credentials \
  --set auth.rootPasswordKey=mongodb-root-password \
  --set auth.passwordKey=mongodb-user-password \
  --set auth.database=jobros # set the database name

# Install the jobros service
helm install jobros ./helm/jobros \
  --set mongodb.existingSecret=mongodb-credentials \
  --set mongodb.uriKey=mongodb-uri
```

## Development

When developing on windows, configure git to not convert line endings to CRLF.
```sh
git config core.autocrlf false
```

Run unit tests
```azure
go test ./...
```