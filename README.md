# Cloud Infrastructure - Kubernetes

This GIT contains the files created to complete each step of the lab. To see the steps and how to perform them take a look at the report found in the Documents folder.

The files are sorted by filetype. Some of the steps in the lab used different filetypes:

- **Building Kubernetes Cluster:** For this step we mainly performed actions using commands (found in the "command_history" file). We also created 2 bash scripts to check pod usage and readiness. The downloaded "components.yaml" file can be found in the respective YAML folder and the Advanced Question .txt files in the "Advanced Questions" folder.
- **Ingress Controller:** For this part we used helm to install MetalLB creating and editing the "values.yaml" file found in the respective YAML folder. 
- **Storage:** For the storage part we again used YAML files to create the StorageClass, the PersistentVolume and the Claim. 
- **Deploy an App:** The deployment of the given application was also performed by creating and editing YAML files which can be found in the respective folder. 
- **Troubleshooting:** To make it easier to edit (and mainly not having to use vim) we dumped the downloaded namespaces, deployment and pod's settings into files which too can be found in the YAML folder. These are the final files which also run on cluster.

For all of the steps we used different commands to debug, test, deploy and edit. These can be found in the "command_history" file. For further insight into the different techniques and steps we decided to use take a look into our Report found in the Documents folder