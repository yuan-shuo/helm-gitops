[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/helm-gitops)](https://artifacthub.io/packages/search?repo=helm-gitops)

[Complete Gitops process usage example with introduction](https://github.com/yuan-shuo/helm-gitops/blob/main/doc/example.en.md)	[中文版文档](https://github.com/yuan-shuo/helm-gitops/blob/main/README.zh-CN.md)	

# Helm-GitOps

A Helm extension that provides GitOps-related auxiliary functions for Helm

## Usage

### If you don't want to read anything

Just these three, change `remote / tag` to your own, then copy and use them directly

```bash
# 1. Create a git-initialized Chart
helm gitops create my-chart

# 2. Generate an environment repository (based on Chart remote repository + repository tag)
helm gitops create-env -r https://gitee.com/yuan-shuo188/helm-test1 -t v0.1.1

# 3. Generate an argo.yaml (based on environment repository + repository tag)
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-non-prod1 -t v0.5.0 -m non-prod
```

The generated content is based on information generation, so it saves a lot of trouble (referring to constantly jumping between multiple remote repositories as a human to copy and visually check, etc.). You can just manage git yourself. If you don't want to read the following content, the above three commands can also help you solve most of the trouble.

`(Features marked with * in the title might help you, if you only want to spend a little time reading)`

### chart-git Chart Development Features

Compared to the conventional approach of manually creating Helm charts, pasting `.gitignore` files internally, and performing a series of tedious operations such as creating branches, committing, and managing version numbers, this extension provides a more comfortable simplified solution:

#### Create Chart

```bash
# Create: Add GitOps skeleton on top of the original Helm chart
helm gitops create test
helm gitops create test --actions        # Also generate .github/workflows/<>.yaml
```

#### Create Branch

```bash
# Checkout: Automatically switch to development branch (optional sync with main branch)
helm gitops checkout feature/foo
helm gitops checkout feature/foo -s       # Pull origin/main first, then create branch
```

#### Commit Push (Main Branch Protection)

```bash
# Commit: add + commit + optional push & automatic PR
helm gitops commit -m "fix: foo"                    # Local commit
helm gitops commit -m "feat: xxx" --push            # Commit and push
helm gitops commit -m "ci: update" --pr --push      # Commit + push + automatic PR (includes [create-pr] commit marker)
```

#### Chart Check

```bash
# Local check: helm lint + helm unittest
helm gitops lint
```

#### Chart Push

```bash
# Push code: lint → push (protected branch interception)
helm gitops push                                      # Push to origin/current branch
```

#### Version Management *

One command can automatically complete version updates, automatically clean up old version tgz packages, package new version tgz, generate index.yaml, and push

```bash
# Version management
helm gitops version # Query current version only
helm gitops version -m pr -l patch # Traditional PR mode (create release branch first → submit PR → CI auto tag)
helm gitops version -m main -l patch # Quick main branch mode (direct commit + tag + push simultaneously)
helm gitops version -m main -l no -s pre # v0.0.1 -> v0.0.1-pre
# --mode/-m = main|pr
# --level/-l = patch|minor|major|no (no version number change, can be combined with -s to continuously update suffix on the same v0.x.x)
# --suffix/-s = <your_tag_suffix>
```

### env-repo Environment Repository Features

This tool can save time in environment repository configuration

#### Create Operation

Generate 2 environment repositories directly based on the remote repository link of the Helm chart. The content in the files will be rendered based on information such as the remote repository link, reducing unnecessary manual writing and copying operations

```bash
# Generate environment repository directly based on Chart remote link:
# -r/--remote specifies the Helm chart remote repository link
# -t/--tag specifies the tag of the chart repository used when creating the repository
helm gitops create-env -r https://gitee.com/yuan-shuo188/helm-test1 -t v0.1.1
```

#### Directory Tree

Only need to execute the above line to generate the following directory tree. You can see that a repository is created for both non-production and production environments. Each repository directory contains `.git (already initialized) + .gitignore`

```
$ tree
.
|-- helm-test1-env-non-prod
|   |-- README.md
|   |-- dev
|   |   |-- cd-use
|   |   |-- kustomization.yaml
|   |   |-- patch.yaml
|   |   `-- values.yaml
|   |-- staging
|   |   |-- cd-use
|   |   |-- kustomization.yaml
|   |   |-- patch.yaml
|   |   `-- values.yaml
|   `-- test
|       |-- cd-use
|       |-- kustomization.yaml
|       |-- patch.yaml
|       `-- values.yaml
`-- helm-test1-env-prod
    |-- README.md
    `-- prod
        |-- cd-use
        |-- kustomization.yaml
        |-- patch.yaml
        `-- values.yaml
```

#### File Content

values.yaml is copied from the code of the corresponding tag in the remote repository, so there is no need to go back to the chart repository to see what the default values.yaml looks like. Each file has a **`directory/filename`** marker comment at the top

```yaml
# dev/values.yaml

# Default values for test-nor.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

# This will set the replicaset count more information can be found here: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/
replicaCount: 1
```

kustomization.yaml does not need special modification. Based on the remote repository, only two lines of comments are written in advance to help observe the file location and use the rendering function (line1, line3 comment locations). It can be noted that although helm is used as the rendering source here, repo and tag are not specified. This will be explained in the parameterized rendering section

```yaml
# prod/kustomization.yaml

# helm gitops render-env -e prod -r https://gitee.com/yuan-shuo188/helm-test1.git -t v0.1.3

apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - rendered/helm/helm-chart.yaml

# patchesStrategicMerge:
#   - patch.yaml
```

patch.yaml is empty by default

```yaml
# prod/patch.yaml
```

#### Parameterized Rendering *

If you only need the raw Helm-rendered output, just leave the `kustomization.yaml` untouched (or keep it empty). If Kustomize isn't installed, the plugin silently skips the Kustomize step and you'll find the plain Helm result at：

`<your_env>/rendered/helm/helm-chart.yaml`

##### Command

Just look at the command:

```bash
# -e/--env environment you want to render
# -r/--remote remote repository link
# -t/--tag remote repository tag
# -l/--use-local-cache use existing local files for rendering
# -n/--render-file-name custom naming for rendering result
helm gitops render-env -e prod -r https://gitee.com/yuan-shuo188/helm-test1.git -t v0.1.3
```

After executing the command, you can get a yaml result file rendered based on helm chart and processed by kustomize

##### Explanation

###### --remote / -r

remote uses a git repository link. This part directly saves the publishing operation. **<u>Your helm chart repository only needs to have index.yaml and the corresponding tgz package</u>** (and if you use the version management function of this software, these two things will be automatically built without you worrying). The program will go to the repository raw based on remote to find index.yaml, then use its urls attribute to get the tgz file name, then concatenate it with remote/tag to form the download link of the chart package file, get the charts directory in the corresponding environment, then render the helm yaml file, and then use kustomize for further rendering (if you don't have this software, it will automatically skip and only render to helm)

###### --use-local-cache / -l

If using kustomize's helm function, letting it point to a helm chart is fine. However, if you need to repeatedly debug and generate, continuous network requests are unnecessary. At this time, you might download the helm rendering result file locally and let kustomize point to it. At this time, using the software's `-l` parameter, you don't need to go back and forth. If the cache directory (e.g., prod/rendered/helm/helm-chart.yaml) contains helm files, it will directly render without making network requests

###### --render-file-name / -n

Override the naming of the rendering result file (there is a default naming but maybe you have other naming ideas)

##### A Small Case

First submit a new version in the chart repository

```bash
helm gitops version -l major -m main
```

Then you can use the tree command to see that index and tgz are ready. You can see that the tag automatically upgrades to 2.0.0 (major-level update):

```bash
|-- index.yaml
|-- test-nor-2.0.0.tgz
```

Then you can use the command to automatically build the environment repository. Here we choose the production repository for demonstration

```bash
helm gitops create-env -r https://gitee.com/yuan-shuo188/helm-test1 -t v2.0.0
cd helm-test1-env-prod
|-- README.md
`-- prod
    |-- cd-use
    |-- kustomization.yaml
    |-- patch.yaml
    `-- values.yaml
```

At this time, open `kustomization.yaml` with any editor, and you can see that the comment on the third line has prepared the rendering command for you:

```bash
helm gitops render-env -e prod -r https://gitee.com/yuan-shuo188/helm-test1 -t v2.0.0
```

After execution, use the tree command to see what changes have occurred in the directory:

```
|-- README.md
`-- prod
    |-- cd-use
    |-- kustomization.yaml
    |-- patch.yaml
    |-- rendered
    |   |-- helm
    |   |   `-- helm-chart.yaml
    |   `-- kustomize
    |       `-- test-nor-2.0.0.yaml
    `-- values.yaml
```

You can see that the rendering is complete, and an extra directory called rendered has been added, in which:

* `helm-chart.yaml` is the helm rendering result

* `test-nor-2.0.0.yaml` is the kustomize rendering result based on `helm-chart.yaml`

This is the end. It only took three commands in total (you can put the yaml you need to deploy into the `cd-use` directory for the continuous delivery program to use). Because the rendering is parameterized, if you want to modify the tag or repo, you can adjust the comment prepared for you, then copy the comment command and render with one click.

Explanation of why not to use kustomize's helm-chart parameter, but directly control the rendering source by repo+version:

* kustomize needs a helm repository, not a git repository. Not all hosting platforms have rich features like GitHub pages

* kustomize and helm are not particularly compatible (although kustomize provides helm functionality, helm's changes in v4+ versions cause kustomize to report many parameter errors when using old commands). The parameterized rendering function separates and decouples helm rendering and kustomize rendering. There is no dependency between the two:

`helm-chart.git -> (helm render) -> helm-chart.yaml -> (kustomize render) -> final.yaml`

### argocd-yaml Generation Features

Use this tool to save time writing argocd.yaml

#### Create Operation

By specifying different parameters, you can generate corresponding YAML for both non-production and production environments

```bash
# Generate argo-yaml directly based on the remote link of the environment repository:
# -r/--remote specifies the environment repository remote link (ensure repository is reachable)
# -t/--tag specifies the tag of the remote repository used when generating argo-yaml
# -m/--mode specifies the generation mode: non-prod|prod, will generate yaml files suitable for different environments
# -d/--dry-run does not generate files, only prints the argo-yaml content
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1 -t v1.0.0 -m prod --dry-run
```

The reason for using two environment repositories + two argocd-yaml files is to ensure independent auditing for both environments, with repository and tag isolation, avoiding mixed commit histories. At the same time, argocd-yaml does not point to the helm chart repository but only to the environment repository, because the environment repository has already specified the helm chart as the rendering source, so argo does not need to point to two repositories at the same time causing unnecessary combination confusion. The final presentation can be illustrated as: `argocd -> env-repo -> helm-chart`

Moreover, as mentioned above, argocd or other continuous delivery programs only need to point to the env/cd-use directory, which stores the kust+helm / helm rendering result yaml files. This is a file that can be directly applied by kubectl apply. This means that argo or other continuous delivery software does not need to install additional plugins. At the same time, the yaml files will be reviewed multiple times during the helm repo / env repo process, and stored in the env/cd-use directory as the final complete result yaml in a WYSIWYG manner

#### Directory Tree

Execute the following commands:

```bash
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-non-prod1  -t v0.5.0 -m non-prod
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1  -t v1.0.0 -m prod
```

You can get two YAML files

```
|-- helm-env-non-prod1-argocd-non-prod.yaml
|-- helm-env-prod1-argocd-prod.yaml
```

#### File Content

It will be generated based on the environment repository information. The generated YAML prepares most of the information that originally needed to be filled in manually, allowing users to focus on the problem rather than searching for and filling in information

##### Non-production argo-yaml

For non-production environment repositories, it may contain multiple environment value directories, such as dev, test, etc. When specifying `--mode/-m non-prod`, the program can automatically find all first-level directories containing `kustomize.yaml` in the remote repository and automatically add them to argo-yaml

For example, in the env repository, add an additional environment directory called `anthor-env`:

```
|-- helm-test1-env-non-prod
|   |-- anthor-env
|   |   |-- kustomization.yaml # search core
|   |   |-- patch.yaml
|   |   `-- values.yaml
```

After tagging v0.5.0 and pushing the code and tag, use

```bash
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-non-prod1  -t v0.5.0 -m non-prod
```

The generated YAML result is as follows. You can see that the environment directory `anthor-env` manually created by the user (not generated by this program) has also been written into argo-yaml

```yaml
# helm-env-non-prod1-argo-non-prod.yaml

apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: 'helm-env-non-prod1-argo-non-prod'
  namespace: argocd
spec:
  generators:
  - list:
      elements:
      - env: anthor-env
      - env: dev
      - env: staging
      - env: test
  template:
    metadata:
      name: 'helm-env-non-prod1-{{env}}'
    spec:
      project: default # [may need action] Adjust to the project required for the production environment
      source:
        repoURL: 'https://gitee.com/yuan-shuo188/helm-env-non-prod1'
        targetRevision: 'v0.5.0'
        path: '{{env}}/cd-use'
      destination:
        server: https://kubernetes.default.svc
        namespace: 'helm-env-non-prod1-{{env}}' # [may need action] Adjust to the namespace required for the production environment
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - CreateNamespace=true
```

##### Production argo-yaml

Only specify the prod environment as the control source. After executing the following command, a YAML file is generated

```bash
helm gitops create-argo -r https://github.com/yuan-shuo/helm-env1 -t v1.0.0 -m prod
```

In the production mode (`--mode prod`) generated argo-yaml, automated is disabled by default, and a copy-and-paste manual sync command comment is generated on the third line

```yaml
# helm-env1-argo-prod.yaml
# auto sync in prod environment is closed by default, you can use below command to sync by hand:
# argocd app sync helm-env1-argo-prod

apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: 'helm-env1-argo-prod' # [may need action] If you want to use a different name, please adjust it here
  namespace: argocd
  # annotations:
  #   # canary analysis template (optional, read by Argo Rollouts AnalysisTemplate)
  #   canary.argo.io/analysis-template: 
  #   canary.argo.io/step-weight: 
  #   canary.argo.io/step-duration: 
spec:
  project: default # [may need action] Adjust to the project required for the production environment
  source:
    repoURL: 'https://github.com/yuan-shuo/helm-env1'
    targetRevision: 'v1.0.0'
    path: prod/cd-use
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: prod # [may need action] Adjust to the namespace required for the production environment
  syncPolicy:
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
    syncOptions:
    - CreateNamespace=true
```

## Installation

### Install via helm plugin

1. Go to [Releases · yuan-shuo/helm-gitops](https://github.com/yuan-shuo/helm-gitops/releases) and pick the download link for your platform and architecture.
2. Run `helm plugin install <url>` to install the version you need.

### Install via binary

1. Go to [Releases · yuan-shuo/helm-gitops](https://github.com/yuan-shuo/helm-gitops/releases) and download the binary for your platform and architecture.
2. Move the binary to any directory in your `$PATH`.

## Requirements

* git (version >= 2.23)
* helm
* helm-unitest (optional, can be installed by running **`helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false`**)
* kustomize (optional; if not installed, rendering will skip kustomize and use Helm only)
