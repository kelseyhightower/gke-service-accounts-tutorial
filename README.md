# Google Container Engine (GKE) and Service Accounts

Applications running in GKE have access to other GCP services such as Stack Drive Trace and Google PubSub. In order to access these services a service account can be used.

## Create a Service Account

The example application requires the following permissions:

* The ability to publish messages to a Pubsub topic.
* The ability to write trace data to StackDriver Trace.

Capture the GCP project ID:

```
export PROJECT_ID=$(gcloud config get-value core/project)
```

Set the service account name:

```
export SERVICE_ACCOUNT_NAME="service-account-example"
```

Create the example service account:

```
gcloud beta iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
  --display-name "Google Container Engine example service account"
```

Add the `pubsub.editor` and `cloudtrace.agent` IAM permissions:

```
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role='roles/pubsub.editor'
```

```
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role='roles/cloudtrace.agent'
```

Generate and download the example service account configuration:

```
gcloud beta iam service-accounts keys create \
  --iam-account "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  service-account.json
```

### Setup Google PubSub

```
gcloud beta pubsub topics create echo
``` 

```
gcloud beta pubsub subscriptions create echo --topic echo
```

```
gcloud beta pubsub subscriptions pull echo
```

```
Listed 0 items.
```

### Kubernetes

```
kubectl create secret generic echo --from-file service-account.json
```

```
kubectl create -f replicasets/echo.yaml
```

In a seperate terminal create a proxy to the `echo` pod:

```
kubectl port-forward \
  $(kubectl get pods -l app=echo -o jsonpath='{.items[0].metadata.name}') \
  8080:8080
```

Submit a request to the echo pod:

```
curl http://127.0.0.1:8080/pubsub -d 'Hello GKE!'
```

Fetch a message from the echo subscription:

```
gcloud beta pubsub subscriptions pull echo
```