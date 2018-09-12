# Google Cloud Service Accounts with Google Container Engine (GKE) - Tutorial

Applications running on [Google Container Engine](https://cloud.google.com/container-engine) have access to other [Google Cloud Platform](https://cloud.google.com) services such as [Stackdriver Trace](https://cloud.google.com/trace) and [Cloud Pub/Sub](https://cloud.google.com/pubsub). In order to access these services a [Service Account](https://cloud.google.com/compute/docs/access/service-accounts) must be created and used by client applications.

This tutorial will walk you through deploying the `echo` application which creates Pub/Sub messages from HTTP requests and sends trace data to Stackdriver Trace.  

## Create a Service Account

The `echo` application requires the following permissions:

* The ability to publish messages to a pub/sub topic.
* The ability to write trace data to Stackdriver Trace.

Create a service account for the `echo` application:

```
export PROJECT_ID=$(gcloud config get-value core/project)
```

```
export SERVICE_ACCOUNT_NAME="echo-service-account"
```

```
gcloud iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
  --display-name "echo service account"
```

Add the `pubsub.editor` and `cloudtrace.agent` IAM permissions to the echo service account:

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

### Generate and download the `echo` service account:

```
gcloud iam service-accounts keys create \
  --iam-account "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  service-account.json
```

## Setup Google Pub/Sub Topics and Subscriptions

The `echo` application publishes messages to the `echo` topic. Create the `echo` topic:

```
gcloud pubsub topics create echo
```

Once messages have been pushed to the `echo` topic, they can be fetch using a subscription. Create the `echo` subscription:

```
gcloud pubsub subscriptions create echo --topic echo
```

Test the `echo` subscription:

```
gcloud pubsub subscriptions pull echo --auto-ack
```

```
Listed 0 items.
```

## Deploy to Google Container Engine

The `echo` application needs access to the echo service account created earlier. Create a Kubernetes secret from the `service-account.json` file:

```
kubectl create configmap echo --from-literal "project-id=${PROJECT_ID}"
```

```
kubectl create secret generic echo --from-file service-account.json
```

Deploy the `echo` container image using a replicaset:

```
kubectl create -f deployments/echo.yaml
```

At this point the `gcr.io/hightowerlabs/echo` container image should be running:

```
kubectl get pods
```
```
NAME                    READY     STATUS    RESTARTS   AGE
echo-6f5964c9fb-xh5zn   1/1       Running   0          1m
```

### Publishing Messages

The `echo` service is now running in the cluster on a private IP address. In a seperate terminal create a proxy to the `echo` pod:

```
kubectl port-forward \
  $(kubectl get pods -l app=echo -o jsonpath='{.items[0].metadata.name}') \
  8080:8080
```

At this point the `echo` service is available at `http://127.0.0.1:8080/pubsub`

#### Submit a request to the echo service

```
curl http://127.0.0.1:8080/pubsub -d 'Hello GKE!'
```

Fetch a message from the echo subscription:

```
gcloud pubsub subscriptions pull echo --auto-ack
```

```
┌────────────┬────────────────┬────────────┐
│    DATA    │   MESSAGE_ID   │ ATTRIBUTES │
├────────────┼────────────────┼────────────┤
│ Hello GKE! │ 26699545805948 │            │
└────────────┴────────────────┴────────────┘
```

### Stackdriver Trace

The `echo` application is configured to send 1 out of 10 request to Stackdriver. Once a trace has been submitted it will be viewable via the Stackdriver Trace dashboard.

![Image of Stackdriver Trace Dashboard](stackdriver-trace.png)

## Cleanup

```
kubectl delete configmap echo
```

```
kubectl delete secret echo
```

```
kubectl delete deployment echo
```

```
gcloud pubsub subscriptions delete echo
```

```
gcloud pubsub topics delete echo
```

```
gcloud iam service-accounts delete "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
```
