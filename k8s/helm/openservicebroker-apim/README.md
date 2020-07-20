# Helm Chart for the Deployment of WSO2 API Manager Open Service Broker

## Contents

* [Prerequisites](#prerequisites)

* [Quick Start Guide](#quick-start-guide)

## Prerequisites

* Install [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git), [Helm](https://helm.sh/docs/intro/install/)
  and [Kubernetes client](https://kubernetes.io/docs/tasks/tools/install-kubectl/) in order to run the steps provided in the
  following quick start guide.<br><br>

* An already setup [Kubernetes cluster](https://kubernetes.io/docs/setup).<br><br>

* Install [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/).<br><br>

* Add the WSO2 Helm chart repository.

    ```
     helm repo add wso2 https://helm.wso2.com && helm repo update
    ```
  
## Quick Start Guide

### Install Chart From [WSO2 Helm Chart Repository](https://hub.helm.sh/charts/wso2)

##### 1. Deploy Helm chart for the deployment of WSO2 APIM Open Service Broker.

```
helm install --name <RELEASE_NAME> wso2/openservicebroker-apim --version 0.1.0 --namespace <NAMESPACE> \
--set configs.APIM_BROKER_APIM_TOKENENDPOINT='<APIM_TOKEN_ENDPOINT>' \
--set configs.APIM_BROKER_APIM_DYNAMICCLIENTENDPOINT='<APIM_DYNAMIC_CLIENT_ENDPOINT>' \
--set configs.APIM_BROKER_APIM_PUBLISHERENDPOINT='<APIM_PUBLISHER_ENDPOINT>' \
--set configs.APIM_BROKER_APIM_STOREENDPOINT='<APIM_STORE_ENDPOINT>' \
--set initContainers.initAPIM.host='<APIM_HOSTNAME>' \
--set initContainers.initAPIM.port='<APIM_PORT>' 
```

* `NAMESPACE` should be the Kubernetes Namespace in which the resources are deployed.

>Example:
To integrate with [WSO2 API manager pattern 1](https://github.com/wso2/kubernetes-apim/blob/3.1.x/advanced/am-pattern-1/README.md),
 assuming the WSO2 API Manager Open Service Broker and the WSO2 API manager pattern 1 are in the same namespace,
 
 ```
 helm install --name <RELEASE_NAME> wso2/openservicebroker-apim --namespace <NAMESPACE> \
 --set configs.APIM_BROKER_APIM_TOKENENDPOINT='https://wso2am-pattern-1-am-service:8243' \
 --set configs.APIM_BROKER_APIM_DYNAMICCLIENTENDPOINT='https://wso2am-pattern-1-am-service:9443' \
 --set configs.APIM_BROKER_APIM_PUBLISHERENDPOINT='https://wso2am-pattern-1-am-service:9443' \
 --set configs.APIM_BROKER_APIM_STOREENDPOINT='https://wso2am-pattern-1-am-service:9443' \
 --set initContainers.initAPIM.host='wso2am-pattern-1-am-service' \
 --set initContainers.initAPIM.port='9443' 
 ```
##### 2. Access WSO2 APIM Open Service Broker.
 
Default deployment will expose `openservicebroker-apim.example.com` host.
 
To access the WSO2 APIM Open Service Broker in the environment,
 
 a. Obtain the external IP (`EXTERNAL-IP`) of the Ingress resources by listing down the Kubernetes Ingresses.
 
 ```
 kubectl get ing -n <NAMESPACE>
 ```
 Output:
 
 ```
 NAME                       HOSTS                                      ADDRESS        PORTS    AGE
 <RELEASE_NAME>             openservicebroker-apim.example.com        <EXTERNAL-IP>  80, 443   3m
 ```
 
 b. Add the above host as an entry in `/etc/hosts` file as follows:
    
 ```
 <EXTERNAL-IP> openservicebroker-apim.example.com
 ``` 

### Install Chart From Source

>In the context of this document, <br>
>* `KUBERNETES_HOME` will refer to a local copy of the [`wso2/openservicebroker-apim`](https://github.com/wso2/openservicebroker-apim/)
Git repository. <br>
>* `HELM_HOME` will refer to `<KUBERNETES_HOME>/k8s/helm`. <br>

##### 1. Clone the Kubernetes Resources for WSO2 APIM Open Service Broker Git repository.

```
git clone https://github.com/wso2/openservicebroker-apim.git
```

##### 2. Deploy Helm chart for the deployment of WSO2 APIM Open Service Broker.

```
helm install --dep-up --name <RELEASE_NAME> <HELM_HOME>/openservicebroker-apim --namespace <NAMESPACE> \
--set configs.APIM_BROKER_APIM_TOKENENDPOINT='<APIM_TOKEN_ENDPOINT>' \
--set configs.APIM_BROKER_APIM_DYNAMICCLIENTENDPOINT='<APIM_DYNAMIC_CLIENT_ENDPOINT>' \
--set configs.APIM_BROKER_APIM_PUBLISHERENDPOINT='<APIM_PUBLISHER_ENDPOINT>' \
--set configs.APIM_BROKER_APIM_STOREENDPOINT='<APIM_STORE_ENDPOINT>' \
--set initContainers.initAPIM.host='<APIM_HOSTNAME>' \
--set initContainers.initAPIM.port='<APIM_PORT>' 
```
**Note:**
`NAMESPACE` should be the Kubernetes Namespace in which the resources are deployed

>Example:
To integrate with [WSO2 API manager pattern 1](https://github.com/wso2/kubernetes-apim/blob/3.1.x/advanced/am-pattern-1/README.md),
 assuming the WSO2 API Manager Open Service Broker and the WSO2 API manager pattern 1 are in the same namespace,
 
 ```
 helm install --dep-up --name <RELEASE_NAME> <HELM_HOME>/openservicebroker-apim --namespace <NAMESPACE> \
 --set configs.APIM_BROKER_APIM_TOKENENDPOINT='https://wso2am-pattern-1-am-service:8243' \
 --set configs.APIM_BROKER_APIM_DYNAMICCLIENTENDPOINT='https://wso2am-pattern-1-am-service:9443' \
 --set configs.APIM_BROKER_APIM_PUBLISHERENDPOINT='https://wso2am-pattern-1-am-service:9443' \
 --set configs.APIM_BROKER_APIM_STOREENDPOINT='https://wso2am-pattern-1-am-service:9443' \
 --set initContainers.initAPIM.host='wso2am-pattern-1-am-service' \
 --set initContainers.initAPIM.port='9443' 
 ```
  
##### 3. Access WSO2 APIM Open Service Broker.

 
Default deployment will expose `openservicebroker-apim.example.com` host.
 
To access the WSO2 APIM Open Service Broker in the environment,
 
 a. Obtain the external IP (`EXTERNAL-IP`) of the Ingress resources by listing down the Kubernetes Ingresses.
 
 ```
 kubectl get ing -n <NAMESPACE>
 ```
 Output:
 
 ```
 NAME                       HOSTS                                       ADDRESS        PORTS   AGE
 <RELEASE_NAME>            openservicebroker-apim.example.com         <EXTERNAL-IP>  80, 443   3m
 ```
 
 b. Add the above host as an entry in `/etc/hosts` file as follows:
    
 ```
 <EXTERNAL-IP> openservicebroker-apim.example.com
 ``` 
## Configuration

The following tables lists the configurable parameters of the chart and their default values.

###### Chart Dependencies

| Parameter                                                                   | Description                                                                               | Default Value               |
|-----------------------------------------------------------------------------|-------------------------------------------------------------------------------------------|-----------------------------|
| `dependencies.mysql.enabled`                                                | Enable MySQL chart as a dependency                                                        | true                        |


###### APIM Open Service Broker Deployment Configurations

| Parameter                                | Description                                                                                               | Default Value                                            |
|------------------------------------------|-----------------------------------------------------------------------------------------------------------|----------------------------------------------------------|
| `replicaCount`                           | Number of replicas for APIM Open Service Broker                                                           | 1                                                        |
| `image.repository`                       | Image name for APIM Open Service Broker                                                                   | wso2/openservicebroker-apim                              |
| `image.repository.tag`                   | Image tag for APIM Open Service Broker                                                                    | 3.1.0.1                                                     |
| `image.repository.pullPolicy`            | Refer to [doc](https://kubernetes.io/docs/concepts/containers/images/#updating-images)                    | IfNotPresent                                             |
| `imagePullSecrets`                       | Refer to [doc](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)     | []                                                       |
| `nameOverride`                           | Override name of app                                                                                      | ""                                                       |
| `fullnameOverride`                       | Override the full qualified app name                                                                      | ""                                                       |
| `service.type`                           | Service type                                                                                              | ClusterIP                                                |
| `service.port`                           | Service port                                                                                              | 8444                                                     |
| `ingress.enabled`                        | Enable Ingress                                                                                            | true                                                     |                             
| `ingress.enabled.annotations`            | Ingress annotations                                                                                       | kubernetes.io/ingress.class: nginx                       |
| `ingress.hosts`                          | Define Hosts and paths for the Ingress                                                                    | Host: openservicebroker-apim.example.com, Path: ["/"]    |
| `ingress.tls`                            | Define TLS Hosts for the Ingress                                                                          | Host: openservicebroker-apim.example.com, secretName: "" |
| `configs`                                | Configurations for the APIM Open Service Broker                                                           | A set of environment variables                           |
| `resources.requests.memory`              | The minimum amount of memory that should be allocated for a Pod                                           | 512Mi                                                    |
| `resources.requests.cpu`                 | The minimum amount of CPU that should be allocated for a Pod                                              | 500m                                                     |
| `resources.limits.memory`                | The maximum amount of memory that should be allocated for a Pod                                           | 512Mi                                                    |
| `resources.limits.cpu`                   | The maximum amount of CPU that should be allocated for a Pod                                              | 500m                                                     |
| `nodeSelector`                           | Allow the APIM service broker pods to schedule on selected nodes                                          | {}                                                       |
| `tolerations`                            | Allow the APIM service broker pods to schedule on tainted nodes (requires Kubernetes >= 1.6)              | {}                                                       |
| `affinity`                               | Allow the APIM service broker pods to schedule using affinity rules                                       | {}                                                       |
| `initContainers.initDB.host`             | Host name of the DB for Init container                                                                    | wso2openservicebroker-apim-db-service                    |
| `initContainers.initDB.port`             | Port of the DB for Init container                                                                         | 3306                                                     |
| `initContainers.initAPIM.host`           | Port of the APIM for Init container                                                                       | ""                             |
| `initContainers.initAPIM.port`           | Port of the APIM for Init container                                                                       | 9443                                                     |

###### APIM Open Service Broker Configurations

APIM Open Service Broker Configurations are listed under `configs` as environment variables.

| Environment variable                                | Description                                            | Default Value                                            |
|-----------------------------------------------------|--------------------------------------------------------|----------------------------------------------------------|
| `configs.APIM_BROKER_LOG_LEVEL`                     | Log level for the APIM Open Service Broker             | info                                                     |
| `configs.APIM_BROKER_HTTP_SERVER_AUTH_USERNAME`     | Username for the APIM Open Service Broker              | admin                                                    |
| `configs.APIM_BROKER_HTTP_SERVER_AUTH_PASSWORD`     | Password for the APIM Open Service Broker              | admin                                                    |
| `configs.APIM_BROKER_APIM_USERNAME`                 | Username for the APIM                                  | admin                                                    |
| `configs.APIM_BROKER_APIM_PASSWORD`                 | Password for the APIM                                  | admin                                                    |
| `configs.APIM_BROKER_APIM_TOKENENDPOINT`            | Token endpoint                                         | ""                                                       |
| `configs.APIM_BROKER_APIM_DYNAMICCLIENTENDPOINT`    | Dynamic client registration endpoint                   | ""                                                       |
| `configs.APIM_BROKER_APIM_PUBLISHERENDPOINT`        | APIM Publisher endpoint                                | ""                                                       |
| `configs.APIM_BROKER_APIM_STOREENDPOINT`            | APIM Store endpoint                                    | ""                                                       |
| `configs.APIM_BROKER_DB_HOST`                       | MySQL Database host name                               | wso2openservicebroker-apim-db-service                    |
| `configs.APIM_BROKER_DB_PORT`                       | MySQL Database port                                    | 3306                                                     |
| `configs.APIM_BROKER_DB_USERNAME`                   | MySQL Database Username                                | wso2carbon                                               |
| `configs.APIM_BROKER_DB_PASSWORD`                   | MySQL Database Password                                | wso2carbon                                               |
| `configs.APIM_BROKER_DB_DATABASE`                   | MySQL Database name                                    | BROKER                                                   |
