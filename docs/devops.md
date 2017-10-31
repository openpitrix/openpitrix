# Set Up DevOps Environment

DevOps is recommended to use for this project. Please follow the instructions below to set up your environment. We use Jenkins with Blue Ocean plugin and deploy it on Kubernetes, also continuously deploy OpenPitrix on the Kubernetes cluster.  

- [Create Kubernetes Cluster](#create-kubernetes-cluster)
- [Deploy Jenkins](#deploy-jenkins)
- [Configure Jenkins](#devops)

## Create Kubernetes Cluster

We are using [Kubernetes on QingCloud](https://appcenter.qingcloud.com/apps/app-u0llx5j8) to create a kubernetes production environment by one click. Please follow the [instructions](https://appcenter-docs.qingcloud.com/user-guide/apps/docs/kubernetes/) to create your own cluster.

## Deploy Jenkins

* Access to Kubernetes client provided by the step above via one of the following options.
  - **OpenVPN**: Go to the left navigation tree of the QingCloud console, Networks & CND, VPC Networks; on the content of the kubernetes VPC page, choose Management Configuration, VPN Service, then you will find OpenVPN service.
  - **Port forwarding**: same as OpenVPN, but choose Port Forwarding on the kubernetes VPC page content of VPC Networks; and add a rule to forward source port to the port of ssh port of the kubernetes client, for instance, forward 10007 to 22 of the kubernetes client with the private IP being 192.168.100.7. After that, you need to open the firewall to allow the port 1007 accessible from outside. Please click the Security Group ID on the same page of the VPC, and add the downstream rule for the firewall.
  - **VNC**: If you don't want to access the client node remotely, just go to the kubernetes cluster detailed page on the QingCloud console, and click the windows icon aside of the client ID.
* Copy the [yaml file](../devops/kubernetes/jenkins-qingcloud.yaml) to the kubernetes client, and deploy
  ```
  # kubectl apply -f jenkins-qingcloud.yaml
  ```
* Access Jenkins by opening http://\<ip\>:8080 where ip depends on how you expose to outside.
  - On the kubernetes client
  ```
  # iptables -t nat -A PREROUTING -p tcp -i eth0 --dport 8080 -j DNAT --to-destination "$(kubectl get svc -n jenkins --selector=app=jenkins -o jsonpath='{.items..spec.clusterIP}')":8080
  # iptables -t nat -A POSTROUTING -p tcp --dport 8080 -j MASQUERADE
  # sysctl -w net.ipv4.conf.eth0.route_localnet=1
  ```
    Now access the kubernetes client port 8080 will be forwarded to the Jenkins service. 
  - If you use OpenVPN to access the kubernetes client, then open http://\<kubernetes client private ip\>:8080 to access Jenkins. If you use Port Forwarding to access the client, then forward the VPC port 8080 to the client port 8080 as describe above. Now open http://\<VPC EIP\>:8080 to access Jenkins 

## Configure Jenkins
  > You can refer [jenkins.io](https://jenkins.io/doc/tutorials/using-jenkins-to-build-a-java-maven-project/) about how to configure Jenkins and creae a pipeline.

* Unlock Jenkins
  - Get the Adminstrator password from the log on the kubernetes client
  ```
  # kubectl logs "$(kubectl get pods -n jenkins --selector=app=jenkins -o jsonpath='{.items..metadata.name}')" -c jenkins -n jenkins
  ```
  - Go to Jenkins console, paste the password and continue. Install suggested plugins, then create the first admin user and save & finish

## Create a pipeline
  - Fork OpenPitrix from github for your development. 
  - On the Jenkins panel, click Open Blue Ocean and start to create a new pipeline. Choose GitHub, paste your access key of GitHub, select the repository you want to create a CI/CD pipeline.
