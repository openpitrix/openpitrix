# Architecture

Basic idea is to decouple application repository and runtime environment. The runtime environment an application can run is by matching the labels of runtime environment and the selector of the repository where the application is from.

### Design key points:

* Application repos are labelled for GUI to show in category list, and have label selector to choose which runtime to run when user to deploy any application that belongs to the repo. 
* Runtime env is labelled. A runtime can have multiple labels.
* Repo indexer will scan configured repo list periodically and cache the metadata of the repos.
* Runtime interface will provide generic interface for application management such as create cluster etc. The specific runtime will implement the interface as a plugin.

![Arichitecture](../images/arch.png)

### Database

The project is microservice architecure oriented. The database design is different than monolithic approach. Please check the [design details](db-design.md)

