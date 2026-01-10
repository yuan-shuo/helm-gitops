[中文版](https://github.com/yuan-shuo/helm-gitops/blob/main/README.zh-CN.md)

# Helm-GitOps

A Helm extension that provides GitOps-related auxiliary functions for Helm

## Usage

### Git-related Features

Compared to the conventional approach of manually creating Helm charts, pasting `.gitignore` files internally, and performing a series of tedious operations such as creating branches, committing, and managing version numbers, this extension provides a more comfortable simplified solution:

```bash
# 1. Create: Add GitOps skeleton on top of the original Helm chart
helm gitops create test
helm gitops create test --actions        # Also generate .github/workflows/<>.yaml


# 2. Checkout: Automatically switch to development branch (optional sync with main branch)
helm gitops checkout feature/foo
helm gitops checkout feature/foo -s       # Pull origin/main first, then create branch


# 3. Commit: add + commit + optional push & automatic PR
helm gitops commit -m "fix: foo"                    # Local commit
helm gitops commit -m "feat: xxx" --push            # Commit and push
helm gitops commit -m "ci: update" --pr --push      # Commit + push + automatic PR (includes [create-pr] commit marker)


# 4. Local check: helm lint + helm unittest
helm gitops lint


# 5. Push code: lint → push (protected branch interception)
helm gitops push                                      # Push to origin/current branch


# 6. Version management
helm gitops version # Query current version only
helm gitops version -m pr -l patch # Traditional PR mode (create release branch first → submit PR → CI auto tag)
helm gitops version -m main -l patch # Quick main branch mode (direct commit + tag + push simultaneously)
# --mode=main|pr
# --level=patch|minor|major
```

### Environment Repository Features

This tool can save time in environment repository configuration

#### Create Operation

Generate an environment repository directly based on the remote repository link of the Helm chart. The content in the files will be rendered based on information such as the remote repository link, reducing unnecessary manual writing and copying operations

```bash
# Generate environment repository directly based on remote link:
# -r/--remote specifies the Helm chart remote repository link
# -t/--tag specifies the tag of the chart repository used when creating the repository
helm gitops create-env -r https://gitee.com/yuan-shuo188/helm-test1 -t v0.1.1
```

#### Directory Tree

Only need to execute the above line to generate the following directory tree

```bash
.
├── .git
├── .gitignore
├── README.md
├── dev
│   ├── kustomization.yaml
│   ├── patch.yaml
│   └── values.yaml
├── prod
│   ├── kustomization.yaml
│   ├── patch.yaml
│   └── values.yaml
├── staging
│   ├── kustomization.yaml
│   ├── patch.yaml
│   └── values.yaml
└── test
    ├── kustomization.yaml
    ├── patch.yaml
    └── values.yaml
```

#### File Content

values.yaml is copied from the code of the corresponding tag in the remote repository. Each file has a **`directory/filename`** marker comment at the top

```yaml
# dev/values.yaml

# Default values for test-nor.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

# This will set the replicaset count more information can be found here: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/
replicaCount: 1
```

kustomization.yaml will automatically render using the remote repository link and tag, for example the following YAML. The name will be obtained using the name attribute of Chart.yaml, and it will check whether the fullnameOverride attribute of values.yaml is empty. If not empty, it will be overwritten

```yaml
# staging/kustomization.yaml

apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

# namespace: 'your_staging_namespace'

helmCharts:
- name: 'test-nor'
  repo: 'https://gitee.com/yuan-shuo188/helm-test1'
  version: 'v0.1.1'
  releaseName: 'staging'
  valuesFile: values.yaml

patchesStrategicMerge:
  - patch.yaml
```

#### View Chart Version Used by Each Environment

Just one command:

```bash
helm gitops env-version
```

The effect is as follows, so there is no need to open each environment directory one by one to find the version written somewhere in the file

```bash
$ helm gitops env-version
dev: v0.1.1
prod: v0.1.1
staging: v0.1.1
test: v0.1.1
```

### ArgoCD Features

To be developed

## Installation

### Using helm plugin install

```bash
helm plugin install https://github.com/yuan-shuo/helm-gitops
```

### Using Binary Files

- Go to: [Releases · yuan-shuo/helm-gitops](https://github.com/yuan-shuo/helm-gitops/releases) to download the binary file for your operating system

- Place the extracted `gitops` binary file in the `$HELM_PLUGIN_DIR/bin/` directory

- Grant execute permissions to the gitops binary file

```bash
chmod +x $HELM_PLUGIN_DIR/bin/gitops
```

## Requirements

* git (version >= 2.23)
* helm
* helm-unitest (optional, can be installed by running **`helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false`**)
