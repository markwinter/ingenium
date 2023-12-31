# Broker

Reference: https://knative.dev/docs/eventing/broker/kafka-broker/

## Installation

### Install Kafka Cluster

```BASH
kubectl apply -f kafka-crd.yaml
kubectl apply -f kafka.yaml
```

### Install Kafka Broker for Knative Eventing

```BASH
kubectl apply -f kafka-controller.yaml
kubectl apply -f kafka-broker.yaml
```

### Create a Broker

```BASH
kubectl apply -f namespace.yaml
kubectl apply -f broker-config.yaml
kubectl apply -f broker.yaml
```