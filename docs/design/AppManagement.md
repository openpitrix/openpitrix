# Summary 
Openpitrix is designed to make user easily deploy their service into kubernetes. 
it mainly contain two parts,App and Cluster,you can read about these two part to deep into openpitirx!
## App life cycle
![](https://github.com/openpitrix/openpitrix/tree/master/docs/images/AppManagement.png)
## AppManagement RPC
```go
rpc SyncRepo (SyncRepoRequest) returns (SyncRepoResponse)
// Get index.yaml in helm repository, read detail informations of app and version,then save to database

rpc CreateApp (CreateAppRequest) returns (CreateAppResponse)
// Create App wit version_package

rpc ValidatePackage (ValidatePackageRequest) returns (ValidatePackageResponse)
// Validate package, call ValidatePackage before CreateApp, extract app_name and version_name

rpc GetAppStatistics (GetAppStatisticsRequest) returns (GetAppStatisticsResponse)
// Get summary

rpc DescribeApps (DescribeAppsRequest) returns (DescribeAppsResponse)
// Query unreleased app, unreleased means active = 0

rpc DescribeActiveApps (DescribeAppsRequest) returns (DescribeAppsResponse)
// Query released app, released means active = 1

rpc ModifyApp (ModifyAppRequest) returns (ModifyAppResponse)
// Modify unreleased app, cannot modify released app, or release version to update information of app

rpc UploadAppAttachment(UploadAppAttachmentRequest) returns (UploadAppAttachmentResponse)
// Upload app attachment, such as icon or screenshot

rpc DeleteApps (DeleteAppsRequest) returns (DeleteAppsResponse)
// Delete app only if all of versions have been deleted

rpc CreateAppVersion (CreateAppVersionRequest) returns (CreateAppVersionResponse)
// create app_version

rpc DescribeAppVersions (DescribeAppVersionsRequest) returns (DescribeAppVersionsResponse)
// Query unreleased app_version

rpc DescribeActiveAppVersions (DescribeAppVersionsRequest) returns (DescribeAppVersionsResponse)
// Query released app_version


rpc DescribeAppVersionAudits (DescribeAppVersionAuditsRequest) returns (DescribeAppVersionAuditsResponse)
// DescribeAppVersionAudits 是查询app_version相关的状态切换的审计记录
rpc DescribeAppVersionReviews (DescribeAppVersionReviewsRequest) returns (DescribeAppVersionReviewsResponse)
// DescribeAppVersionReviews 是查询app_version审核相关的操作记录

rpc ModifyAppVersion (ModifyAppVersionRequest) returns (ModifyAppVersionResponse)
// Modify app_version

rpc GetAppVersionPackage (GetAppVersionPackageRequest) returns (GetAppVersionPackageResponse)
// GetAppVersionPackage,get whole package
rpc GetAppVersionPackageFiles (GetAppVersionPackageFilesRequest) returns (GetAppVersionPackageFilesResponse)
// GetAppVersionPackageFiles,get file in package

rpc SubmitAppVersion (SubmitAppVersionRequest) returns (SubmitAppVersionResponse)
// Submit app_version to audit
rpc CancelAppVersion (CancelAppVersionRequest) returns (CancelAppVersionResponse)
// Cancel app_version to audit

rpc SuspendAppVersion (SuspendAppVersionRequest) returns (SuspendAppVersionResponse)
// Suspend app_version
rpc RecoverAppVersion (RecoverAppVersionRequest) returns (RecoverAppVersionResponse)
// Recover app_version

rpc DeleteAppVersion (DeleteAppVersionRequest) returns (DeleteAppVersionResponse)
// Delete app_version,not all app_version can be deleted,it depends on its state

rpc ReviewAppVersion (ReviewAppVersionRequest) returns (ReviewAppVersionResponse)
// Start to audit app_version
rpc PassAppVersion (PassAppVersionRequest) returns (PassAppVersionResponse)
// Pass app_version
rpc RejectAppVersion (RejectAppVersionRequest) returns (RejectAppVersionResponse)
// Reject app_version

rpc ReleaseAppVersion (ReleaseAppVersionRequest) returns (ReleaseAppVersionResponse)
// Release app_version
```
## Cluster 
Cluster means app deployed in kubernetes.these rpc interface all design in Vmbased style,
Cluster module need to be designed to knative and discard Vmbased

```go

rpc AddNodeKeyPairs (AddNodeKeyPairsRequest) returns (AddNodeKeyPairsResponse);
rpc DeleteNodeKeyPairs (DeleteNodeKeyPairsRequest) returns (DeleteNodeKeyPairsResponse);
// Create key pair
rpc CreateKeyPair (CreateKeyPairRequest) returns (CreateKeyPairResponse)
// Get key pairs, support filter with these fields(key_pair_id, name, owner), default return all key pairs
rpc DescribeKeyPairs (DescribeKeyPairsRequest) returns (DescribeKeyPairsResponse)
// Batch delete key pairs
rpc DeleteKeyPairs (DeleteKeyPairsRequest) returns (DeleteKeyPairsResponse)
// Batch attach key pairs to node
rpc AttachKeyPairs (AttachKeyPairsRequest) returns (AttachKeyPairsResponse)
//Batch detach key pairs from node
rpc DetachKeyPairs (DetachKeyPairsRequest) returns (DetachKeyPairsResponse)
// Get subnets
rpc DescribeSubnets (DescribeSubnetsRequest) returns (DescribeSubnetsResponse)
// Create cluster
rpc CreateCluster (CreateClusterRequest) returns (CreateClusterResponse) 
// Create debug cluster
rpc CreateDebugCluster (CreateClusterRequest) returns (CreateClusterResponse) 
rpc ModifyCluster (ModifyClusterRequest) returns (ModifyClusterResponse);
rpc ModifyClusterNode (ModifyClusterNodeRequest) returns (ModifyClusterNodeResponse);
// Modify cluster attributes
rpc ModifyClusterAttributes (ModifyClusterAttributesRequest) returns (ModifyClusterAttributesResponse) 
// Modify node attributes in the cluster
rpc ModifyClusterNodeAttributes (ModifyClusterNodeAttributesRequest) returns (ModifyClusterNodeAttributesResponse)
rpc AddTableClusterNodes (AddTableClusterNodesRequest) returns (google.protobuf.Empty);
rpc DeleteTableClusterNodes (DeleteTableClusterNodesRequest) returns (google.protobuf.Empty);
// Batch delete clusters
rpc DeleteClusters (DeleteClustersRequest) returns (DeleteClustersResponse) 
// Upgrade cluster
rpc UpgradeCluster (UpgradeClusterRequest) returns (UpgradeClusterResponse) 
// Rollback cluster
rpc RollbackCluster (RollbackClusterRequest) returns (RollbackClusterResponse)
// Resize cluster
rpc ResizeCluster (ResizeClusterRequest) returns (ResizeClusterResponse) 
// Batch add nodes to cluster
rpc AddClusterNodes (AddClusterNodesRequest) returns (AddClusterNodesResponse) 
// Batch delete nodes from cluster
rpc DeleteClusterNodes (DeleteClusterNodesRequest) returns (DeleteClusterNodesResponse)
rpc UpdateClusterEnv (UpdateClusterEnvRequest) returns (UpdateClusterEnvResponse)
// Get clusters, can filter with these fields(cluster_id, app_id, version_id, status, runtime_id, frontgate_id, owner, cluster_type), default return all clusters
rpc DescribeClusters (DescribeClustersRequest) returns (DescribeClustersResponse) 
// Get debug clusters, can filter with these fields(cluster_id, app_id, version_id, status, runtime_id, frontgate_id, owner, cluster_type), default return all debug clusters
rpc DescribeDebugClusters (DescribeClustersRequest) returns (DescribeClustersResponse) 
// Get app clusters, can filter with these fields(cluster_id, app_id, version_id, status, runtime_id, frontgate_id, owner, cluster_type), default return all app clusters
rpc DescribeAppClusters (DescribeAppClustersRequest) returns (DescribeAppClustersResponse) runtime_id, frontgate_id, owner, cluster_type), default return all debug app clusters
rpc DescribeDebugAppClusters (DescribeAppClustersRequest) returns (DescribeAppClustersResponse) 
// Get nodes in cluster, can filter with these fields(cluster_id, node_id, status, owner)
rpc DescribeClusterNodes (DescribeClusterNodesRequest) returns (DescribeClusterNodesResponse)
// Batch stop clusters
rpc StopClusters (StopClustersRequest) returns (StopClustersResponse) 
// Batch start cluster
rpc StartClusters (StartClustersRequest) returns (StartClustersResponse) 
// Batch recover clusters
rpc RecoverClusters (RecoverClustersRequest) returns (RecoverClustersResponse)
// Batch cease clusters
rpc CeaseClusters (CeaseClustersRequest) returns (CeaseClustersResponse) 
rpc GetClusterStatistics (GetClusterStatisticsRequest) returns (GetClusterStatisticsResponse)
```