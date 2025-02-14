# ptp-events-v1-v2-example
This repository includes a sample PTP events consumer application, to showcase the required changes between the PTP Events APIv1 and v2.

You can find detailed documentation on the PTP Events API in the relevant [OpenShift documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/4.16/html-single/networking/index#ptp-cloud-events-consumer-dev-reference-v2).

Please checkout the `v1` and `v2` branches to check the differences between the applications.

*NOTE*: the example consumer application shown here is overly simplified for readability purposes, and will not behave correctly in a real
production environment. For example, it does not make any checks that the events producer is available, nor it checks for disconnects or
reconnects. For a proper example, with full Events API v1 and v2 support, please refer to the [cloud-event-proxy repository](https://github.com/redhat-cne/cloud-event-proxy/tree/main/examples).

## Compiling and testing

To compile and push a container image including the generated binary, just type:

```
make all
```

Make sure you modify the makefile before, to include the right container registry, namespace and image name you want to use.

You can then create the required namespace and deployment by running:

```
oc apply -f deployment/namespace.yaml
oc apply -f deployment/deployment.yaml
```
