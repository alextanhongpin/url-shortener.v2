# URL Shortener, V2

# Calling Kubernetes

```
$ make up-kube
$ alias k=kubectl
```

Check services status

```bash
$ k get svc
```
Output:

```bash
NAME                    TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
kubernetes              ClusterIP   10.96.0.1        <none>        443/TCP          11d
postgres-service        NodePort    10.103.129.18    <none>        5432:31037/TCP   2m11s
url-shortener-service   NodePort    10.110.215.198   <none>        8080:32094/TCP   2m11s
```

Call the service locally:

```bash
$ curl localhost:32094/health
```

Output:

```bash
{"deployed_at":"2019-11-02T14:36:54.3052781Z","uptime":"39.1938112s","version":"195e5c1"}
```
