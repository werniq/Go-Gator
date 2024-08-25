#!/usr/bin/env bash

docker image rm qniw984/admission-controller-webhook:1.1.0

kubectl delete all -n webhook-demo --all
kubectl delete ValidatingWebhookConfiguration demo-webhook
kubectl delete namespace webhook-demo