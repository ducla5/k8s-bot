# K8S-Bot

## Description
 Simple chatwork chatbot. Can be used to apply manifest file to k8s cluster.
 
## How to set up

1. Manifest store
 - S3_BUCKET_NAME:/staging/api-deployment.yaml
 - S3_BUCKET_NAME:/staging/bff-deployment.yaml
 - S3_BUCKET_NAME:/devteam/api-deployment.yaml
 - S3_BUCKET_NAME:/devteam/bff-deployment.yaml
```yaml
        - name: S3_BUCKET_NAME
          value: "xxx"
        - name: AWS_ACCESS_KEY_ID
          value: "xxx"
        - name: AWS_SECRET_ACCESS_KEY
          value: "xxx"
        - name: AWS_DEFAULT_REGION
          value: "ap-northeast-1"
```

2. Create service account
```yaml
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: k8s-bot
rules:
- apiGroups: ["", "extensions", "apps", "pods/log", "pods/exec"]
  resources: ["pods", "deployments", "replicationcontrollers","replicaset", "services", "pod", "events"]
  verbs: ["get", "watch", "list", "delete", "create"]
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8s-bot
  namespace: dev
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: k8s-bot
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: k8s-bot
subjects:
  - kind: ServiceAccount
    name: k8s-bot
    namespace: dev
```
3. Create deployment

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: k8s-bot
  namespace : dev
  annotations:
    sidecar.istio.io/inject: "false"  # Need if you deploy to istio-enabled namespace
spec:
  containers:
    - image: anhduc720/k8s-bot:latest
      imagePullPolicy: Always
      name: k8s-bot
      env:
        - name: CHATWORK_TOKEN_KEY
          value: "xxx"
        - name: CHATWORK_ROOM_ID
          value: "x"
        - name: ECR_API
          value: "xxx"
        - name: ECR_BFF
          value: "xxx"
        - name: ECR_BFF_NGINX
          value: "xxx"
        - name: S3_BUKKEN_NAME
          value: "xxx"
        - name: K8S_NAMESPACE
          value: "dev"
        - name: AWS_ACCESS_KEY_ID
          value: "xxx"
        - name: AWS_SECRET_ACCESS_KEY
          value: "xxx"
        - name: AWS_DEFAULT_REGION
          value: "xxx"
  restartPolicy: Always
  serviceAccountName: k8s-bot
```
4. Command
```
To ChatbotAccount
deploy bff image {IMAGE_TAG} on BFF_ENV [api API_ENV]
```
or
```
To ChatbotAccount
deploy api image {IMAGE_TAG} on API_ENV
```
sample
```
[To:222222]Chatbot
deploy bff image lftv-develop on duc
```
```
[To:222222]Chatbot
deploy bff image lftv-develop on duc api staging
```
```
[To:222222]Chatbot
deploy bff image lftv-develop on duc api duc
```
```
[To:222222]Chatbot
deploy api image lftv-develop on duc
```