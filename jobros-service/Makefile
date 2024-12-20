VERSION=$(shell cat version.txt)

update-helm:
    sed -i "s/^version:.*$$/version: ${VERSION}/" deploy/helm/jobros/Chart.yaml
    sed -i "s/^appVersion:.*$$/appVersion: ${VERSION}/" deploy/helm/jobros/Chart.yaml

update-go:
    sed -i "s/^const Version = \".*\"/const Version = \"${VERSION}\"/" backend/version.go

version-update: update-helm update-go
