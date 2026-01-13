# Application Case 1 - Starting from Zero

Since we're starting from zero, the GitOps workflow for Application Case 1 is very long. In actual use, there aren't many commands to execute - most of the content is output + introduction. If you really don't want to read it, you can pick the parts that interest you, or just throw the documentation to GPT and let it explain it to you.

## Install helm-gitops

Here I directly downloaded the release from GitHub and put it in the environment variable directory. Typing "helm" repeatedly is too tedious.

```bash
$ gitops -h
Helm GitOps utilities

Usage:
  helm-gitops [command]

Available Commands:
  checkout    switch to or create a new development branch
  commit      git add & commit
  completion  Generate the autocompletion script for the specified shell
  create      create a new Helm chart with GitOps scaffold
  create-argo create a new argo yaml from an existing remote env repo
  create-env  create a new environment repository from an existing remote Helm chart
  help        Help about any command
  lint        helm lint + unittest
  push        push current development branch to remote
  render-env  create a new environment repository from an existing remote Helm chart
  version     print or bump Chart.yaml version

Flags:
  -h, --help   help for helm-gitops

Use "helm-gitops [command] --help" for more information about a command.
```

## Build helm-chart Repository

### Create Local Chart

Done in one line.

```bash
$ gitops create demo-chart
✅ Chart "demo-chart" created with GitOps scaffold & initial commit.
$ cd demo-chart
```

### Create Remote Chart Repository and Push

Here I created a remote repository named "helmChart" on GitHub, then copied the push command provided by GitHub to push the local chart up.

```bash
$ git remote add origin git@github.com:yuan-shuo/helmChart.git
$ git push -u origin main
Enumerating objects: 16, done.
Counting objects: 100% (16/16), done.
Delta compression using up to 16 threads
Compressing objects: 100% (15/15), done.
Writing objects: 100% (16/16), 6.41 KiB | 1.07 MiB/s, done.
Total 16 (delta 0), reused 0 (delta 0), pack-reused 0 (from 0)
To github.com:yuan-shuo/helmChart.git
 * [new branch]      main -> main
branch 'main' set up to track 'origin/main'.
```

### Modify Chart and Publish

You can see that the newly created chart version is still 0.1.0

```bash
$ cat Chart.yaml | grep version: 
version: 0.1.0
```

#### Make a Random Change to Chart and Push

Then I go to the development branch. Suppose I modified a file - here I'll just add a comment.

```bash
$ gitops checkout feature/foo
creating and switching to "feature/foo"
Switched to a new branch 'feature/foo'
$ echo -e "\n# a example change content" >> values.yaml
$ tail -n 3 values.yaml
affinity: {}

# a example change content
```

Then push to the remote repository. You can see that the tool automatically performed a lint check before pushing, and only pushed to the remote repository after passing.

```bash
$ gitops commit -m "feature: a comment added" --push 
running helm lint...
==> Linting .
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
running helm unittest...

### Chart [ demo-chart ] .


Charts:      1 passed, 1 total
Test Suites: 0 passed, 0 total
Tests:       0 passed, 0 total
Snapshot:    0 passed, 0 total
Time:        2.6789ms

local tests passed
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Delta compression using up to 16 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 315 bytes | 315.00 KiB/s, done.
Total 3 (delta 2), reused 0 (delta 0), pack-reused 0 (from 0)
remote: Resolving deltas: 100% (2/2), completed with 2 local objects.
remote: 
remote: Create a pull request for 'feature/foo' on GitHub by visiting:
remote:      https://github.com/yuan-shuo/helmChart/pull/new/feature/foo
remote:
To github.com:yuan-shuo/helmChart.git
 * [new branch]      HEAD -> feature/foo
```

#### PR on GitHub

Here I submitted a PR on GitHub and merged it into the main branch. Pull the history and you can see that the main branch has the commit information from just now.

```bash
$ gitops checkout main
switching to "main"
Switched to branch 'main'
$ git pull
remote: Enumerating objects: 1, done.
remote: Counting objects: 100% (1/1), done.
remote: Total 1 (delta 0), reused 0 (delta 0), pack-reused 0 (from 0)
Unpacking objects: 100% (1/1), 904 bytes | 452.00 KiB/s, done.
From github.com:yuan-shuo/helmChart
   eccd7e4..774a4f9  main       -> origin/main
Updating eccd7e4..774a4f9
Fast-forward
 values.yaml | 2 ++
 1 file changed, 2 insertions(+)
$ git log --oneline
774a4f9 (HEAD -> main, origin/main) Merge pull request #1 from yuan-shuo/feature/foo
ee70cc6 (origin/feature/foo, feature/foo) feature: a comment added
eccd7e4 helm gitops chart init
```

If you include `--actions` when creating, it will actually automatically submit a PR for you in GitHub Actions. However, this case won't rely on additional features outside the toolchain (such as GitHub's Actions - I don't think all repositories have GitHub's rich features).

#### Version Upgrade and Publish

At this point, let's check the current version again. First use this tool (it will pull the latest main branch once and then print the latest version), then use cat again. You can see that the version hasn't changed - it's still 0.1.0.

```bash
$ gitops version
Already on 'main'
From github.com:yuan-shuo/helmChart
 * branch            main       -> FETCH_HEAD
0.1.0
$ cat Chart.yaml | grep version:
version: 0.1.0
```

Then use the extension - one command is enough.

```bash
$ gitops version -m main -l patch
```

The increment level specified here is patch, which means the third digit (the end) of the version number +1. The extension will first return to the main branch to execute the lint check, automatically modify the version number in `Chart.yaml`, generate `index.yaml` and package the chart, then `add + commit + tag + push`. These operations only require typing the line below when using this extension. You can see in various ways that the version has indeed automatically changed and includes `index.yaml` and the tgz package file.

```bash
# check
$ cat Chart.yaml | grep version:
version: 0.1.1
$ cat index.yaml | grep version:
    version: 0.1.1
$ ls | grep tgz
demo-chart-0.1.1.tgz
$ git tag
v0.1.1
$ curl -I https://github.com/yuan-shuo/helmChart/tree/v0.1.1
HTTP/1.1 200 OK
Date: Tue, 13 Jan 2026 09:25:39 GMT
```

Then return to the parent directory to prepare for building the environment repository.

```bash
$ cd ..
```

## Build Environment Repository

### Create Local Environment Repository

Since we already have a remote helm-chart git repository, just use the remote (-r) parameter to specify the remote link, then specify the tag as v0.1.1 that we just pushed. Still done in one line.

```bash
$ gitops create-env -r https://github.com/yuan-shuo/helmChart -t v0.1.1
```

Let's see what's in the directory: `helmChart-env-non-prod`  `helmChart-env-prod`
```bash
# check new directory
$ ls
demo-chart  helmChart-env-non-prod  helmChart-env-prod
```

These two directories are actually the environment repositories automatically created for you.

### Create Remote Environment Repository and Push

#### Non-Production Environment Repository

Here I created a remote repository named "helm-chart-env-repo-non-prod" on GitHub, then copied the push command provided by GitHub to push.

```bash
$ cd helmChart-env-non-prod
$ git remote add origin git@github.com:yuan-shuo/helm-chart-env-repo-non-prod.git
$ git push -u origin main
Enumerating objects: 18, done.
Counting objects: 100% (18/18), done.
Delta compression using up to 16 threads
Compressing objects: 100% (13/13), done.
Writing objects: 100% (18/18), 3.32 KiB | 1.66 MiB/s, done.
Total 18 (delta 4), reused 0 (delta 0), pack-reused 0 (from 0)
remote: Resolving deltas: 100% (4/4), done.
To github.com:yuan-shuo/helm-chart-env-repo-non-prod.git
 * [new branch]      main -> main
branch 'main' set up to track 'origin/main'.
$ cd ..
```

#### Production Environment Repository

Then I created another remote repository named "helm-chart-env-repo-prod" on GitHub, then copied the push command provided by GitHub to push.

```bash
$ cd helmChart-env-prod
$ git remote add origin git@github.com:yuan-shuo/helm-chart-env-repo-prod.git
$ git push -u origin main
Enumerating objects: 10, done.
Counting objects: 100% (10/10), done.
Delta compression using up to 16 threads
Compressing objects: 100% (7/7), done.
Writing objects: 100% (10/10), 2.79 KiB | 2.79 MiB/s, done.
Total 10 (delta 0), reused 0 (delta 0), pack-reused 0 (from 0)
To github.com:yuan-shuo/helm-chart-env-repo-prod.git
 * [new branch]      main -> main
branch 'main' set up to track 'origin/main'.
$ cd ..
```

### Modify Environment Repository

#### Non-Production Environment Repository

You can see that the extension automatically created three environment directories, but these are just initially generated directories. You can name them whatever you like, create them yourself - later we will repeatedly prove that directory names have no effect on the use of the extension.

```bash
$ cd helmChart-env-non-prod
$ ls
dev  README.md  staging  test
```

##### Execute Changes

Suppose you like to customize environment directories. OK, then I'll rename one here as if it's a custom directory you like, as the change content for the non-production environment repository.

```bash
$ git checkout -b feature/env
Switched to a new branch 'feature/env'
$ mv staging custom
$ ls
custom  dev  README.md  test
```

##### Render YAML

For the subsequent continuous delivery program to find and synchronize the YAML files in the directory, you need to first generate YAML files through `helm + kustomize (optional)`, then put them in the cd-use directory of each environment:

```bash
# Render to get YAML files
$ gitops render-env -e custom -r https://github.com/yuan-shuo/helmChart -t v0.1.1
$ gitops render-env -e dev -r https://github.com/yuan-shuo/helmChart -t v0.1.1
$ gitops render-env -e test -r https://github.com/yuan-shuo/helmChart -t v0.1.1
# Put YAML files into the directory for the subsequent CD program to synchronize
$ cp custom/rendered/kustomize/demo-chart-0.1.1.yaml custom/cd-use
$ cp dev/rendered/kustomize/demo-chart-0.1.1.yaml dev/cd-use
$ cp test/rendered/kustomize/demo-chart-0.1.1.yaml test/cd-use
```

##### Commit, Push, PR, Tag

```bash
# push code for PR
$ git add . && git commit -m "yaml render" && git push -u origin feature/env
# after PR
$ git checkout main && git pull origin main && git tag v3.0.0 && git push origin tag v3.0.0
```

#### Production Environment Repository

##### Execute Changes

Here I'm preparing to change the replica count from 1 to 2 for the production environment repository.

```bash
$ cd helmChart-env-prod
# See what environments are available
$ ls
prod  README.md
# Enter the prod environment prepared in advance
$ cd prod
```

Then check what the current replica count is in values. Use cat and find it's 1, then use sed to change it to 2.

```bash
$ ls
cd-use  kustomization.yaml  patch.yaml  values.yaml
$ cat values.yaml | grep replicaCount:
replicaCount: 1
$ sed -i 's/replicaCount: 1/replicaCount: 2/' values.yaml
$ cat values.yaml | grep replicaCount:
replicaCount: 2
```

##### Render YAML

Render the YAML file, then put it in the cd-use directory for the subsequent continuous delivery program to synchronize.

```bash
$ gitops render-env -e prod -r https://github.com/yuan-shuo/helmChart -t v0.1.1
$ cp prod/rendered/kustomize/demo-chart-0.1.1.yaml prod/cd-use
```

##### Commit, Push, Tag

I'm really too lazy to `checkout + PR` and then review it myself, so this case will directly demonstrate committing on main.

```bash
$ git add . && git commit -m "prod yaml render" && git tag v2.0.0 && git push origin main v2.0.0
```

### Environment Repository Build Results

Let's see what we have now:

```bash
# Non-production environment repository tag contains v3.0.0
remote = https://github.com/yuan-shuo/helm-chart-env-repo-non-prod
tag = v3.0.0
# Production environment repository tag contains v2.0.0
remote = https://github.com/yuan-shuo/helm-chart-env-repo-prod
tag = v2.0.0
```

Next, you can use these to generate your `argocd.yaml`.

## Build argocd Application

### Generate Application YAML File

It's also very simple - just use the extension to generate the argo application YAML files for both environments in two lines.

```bash
# Generate application YAML file for non-production environment
$ gitops create-argo -r https://github.com/yuan-shuo/helm-chart-env-repo-non-prod -t v3.0.0 -m non-prod
# Generate application YAML file for production environment
$ gitops create-argo -r https://github.com/yuan-shuo/helm-chart-env-repo-prod -t v2.0.0 -m prod
```

Let's see what's in the current directory - one non-production environment application YAML and one production environment application YAML. Still hands-free.

```bash
$ ls -l | awk '{print $9}'
helm-chart-env-repo-non-prod-argocd-non-prod.yaml
helm-chart-env-repo-prod-argocd-prod.yaml
```

Let's look at what was mentioned earlier - changing one environment from staging to custom in the non-production environment. Let's see if argoyml can detect this:

```bash
$ cat helm-chart-env-repo-non-prod-argocd-non-prod.yaml | grep -- '- env'
      - env: custom # <-
      - env: dev
      - env: test
```

So the environments aren't locked to just those few default ones. You can change them as you like - as long as there's a `kustomization.yaml` inside (because argo uses this file to determine whether this is an environment directory. The default generated `kustomization.yaml` does nothing - it's purely a marker. If you need to use it, then write its content).

### Apply YAML File

`alias k = kubectl`

```bash
# Deploy applications for k8s
root@dev-machine:laborant# k apply -f helm-chart-env-repo-non-prod-argo-non-prod.yaml
root@dev-machine:laborant# k apply -f helm-chart-env-repo-prod-argo-prod.yaml

# First check - find that prod hasn't synced. This is because when generating the file we specified -m prod, which disables auto-sync by default, so manual sync is needed
root@dev-machine:laborant# k get app -n argocd -o wide
NAME                                  SYNC STATUS   HEALTH STATUS
helm-chart-env-repo-non-prod-custom   Synced        Healthy         
helm-chart-env-repo-non-prod-dev      Synced        Healthy         
helm-chart-env-repo-non-prod-test     Synced        Healthy         
helm-chart-env-repo-prod-argo-prod    OutOfSync     Missing         

# Manually sync for production environment (the manual sync command is already provided in the comments of the *argo-prod.yaml file, just copy and use it)
root@dev-machine:laborant# argocd app sync helm-chart-env-repo-prod-argo-prod

# At this point, all application deployments are complete
root@dev-machine:laborant# k get app -n argocd -o wide
NAME                                  SYNC STATUS   HEALTH STATUS
helm-chart-env-repo-non-prod-custom   Synced        Healthy         
helm-chart-env-repo-non-prod-dev      Synced        Healthy         
helm-chart-env-repo-non-prod-test     Synced        Healthy         
helm-chart-env-repo-prod-argo-prod    Synced        Healthy         

# You can see that the custom application takes effect in non-prod (the aforementioned non-prod env change)
# At the same time, the prod application replica count has also changed to 2 (the aforementioned prod env change)
root@dev-machine:laborant# k get po -A | awk '{print $2 " - " $4}' | grep demo
custom-demo-chart-6bc68f495b-vzg6z - Running
dev-demo-chart-55f4cfb56b-whp9r - Running
test-demo-chart-784db6ffcc-47cd5 - Running
prod-demo-chart-66766fd86b-bmb6x - Running
prod-demo-chart-66766fd86b-jbxvh - Running
```



