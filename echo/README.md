# Echo

## Build

```
GOOS=linux go build -a --ldflags '-extldflags "-static"' \
  -tags netgo -installsuffix netgo -o echo .
```

### Container

```
export PROJECT_ID=$(gcloud config get-value core/project)
```

```
gcloud container builds submit --tag "gcr.io/${PROJECT_ID}/echo" .
```
