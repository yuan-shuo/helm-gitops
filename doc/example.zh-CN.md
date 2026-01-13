# 应用案例1 - 从0开始

因为从0开始，应用案例1的gitops链路非常长，实际使用指令并不多，大部分内容都是回显+介绍，如果实在不想看可以挑着看感兴趣的部分或者把文档丢给GPT让他为你介绍也行

## 安装 helm-gitops

这里直接从github下载release丢到环境变量目录下了，不断敲上helm太费手了

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

## 构建 helm-chart 仓库

### 创建本地 chart

就一行完事了

```bash
$ gitops create demo-chart
✅ Chart "demo-chart" created with GitOps scaffold & initial commit.
$ cd demo-chart
```

### 创建远程 Chart 仓库并推送

这里我在github上创建了一个名为“helmChart”的远程仓库，然后复制github给出的推送指令把本地chart推送上去

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

### 改动 chart 后发布

可以看到新创建的chart版本还在0.1.0

```bash
$ cat Chart.yaml | grep version: 
version: 0.1.0
```

#### 随便改一笔 chart 然后推送

然后我来到开发分支里，假设我修改了某个文件，这里就随便加条注释好了

```bash
$ gitops checkout feature/foo
creating and switching to "feature/foo"
Switched to a new branch 'feature/foo'
$ echo -e "\n# a example change content" >> values.yaml
$ tail -n 3 values.yaml
affinity: {}

# a example change content
```

然后推送到远程仓库，可以看到工具在推送前自动做了lint检查，通过了才push到远程仓库

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

#### 在 github 上 PR

这里我在github提交PR并合并到主分支了，把历史拉过来打印一下就可以看到主分支有刚才的提交信息了

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

如果在创建时直接附带 `--actions` ，其实他会在github-actions里帮你自动提交PR的，但是这个案例不会依附于工具链以外的附加功能（例如github的actions，我认为并不是所有仓库都有github那么丰富的功能）

#### 版本升级并发布

这时我们再看一下当前版本，先用本工具（会先拉一次最新的主分支然后打印最新版本），然后再用一次cat，可以看到目前版本没变化，就是 0.1.0

```bash
$ gitops version
Already on 'main'
From github.com:yuan-shuo/helmChart
 * branch            main       -> FETCH_HEAD
0.1.0
$ cat Chart.yaml | grep version:
version: 0.1.0
```

然后利用扩展，一条命令就行了

```bash
$ gitops version -m main -l patch
```

这里指定的提升等级是patch，也就是版本号第三位（末尾）+1，扩展会先回到主分支执行lint检查，自动修改 `Chart.yaml` 的版本号，生成 `index.yaml` 并打包chart，随后`add + commit + tag + push`，这些操作使用此扩展就敲下面这一行就够了，可以利用多种方式看到版本的确自动变更并且包含`index.yaml`和tgz打包文件了

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

随后回到父级目录准备构建环境仓库

```bash
$ cd ..
```

## 构建环境仓库

### 创建本地环境仓库

既然已经有一个远程的 helm-chart 的git仓库了，直接使用remote（-r）参数指定远程链接，然后指定tag为刚才推送的v0.1.1就行了，依旧一行完成

```bash
$ gitops create-env -r https://github.com/yuan-shuo/helmChart -t v0.1.1
```
来看看目录里多了些什么：`helmChart-env-non-prod`  `helmChart-env-prod`
```bash
# check new directory
$ ls
demo-chart  helmChart-env-non-prod  helmChart-env-prod
```

这俩目录其实就是自动为你创建好的环境仓库了

### 创建远程环境仓库并推送

#### 非生产环境仓库

这里我在github上创建了一个名为“helm-chart-env-repo-non-prod”的远程仓库，然后复制github给出的推送指令推送

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

#### 生产环境仓库

然后我又在github上创建了一个名为“helm-chart-env-repo-prod”的远程仓库，然后复制github给出的推送指令推送

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

### 改动环境仓库

#### 非生产环境仓库

可以看到扩展自动创建了三个环境目录，但这只是初期生成的目录，你喜欢用什么名字命名都行，自己创建也行，后续会反复证明目录名完全不影响扩展的使用

```bash
$ cd helmChart-env-non-prod
$ ls
dev  README.md  staging  test
```

##### 执行改动

假设你喜欢自定义环境目录，好，那我这里就改个名假设他是你喜欢的自定义目录，作为非生产环境仓库的改动内容

```bash
$ git checkout -b feature/env
Switched to a new branch 'feature/env'
$ mv staging custom
$ ls
custom  dev  README.md  test
```

##### 渲染 yaml

为了后续持续交付程序能够找到目录下的 yaml 同步，需要先通过 `helm + kustomize(可选)` 渲染生成yaml文件，随后放入各环境的 cd-use 目录下：

```bash
# 渲染获取yaml文件
$ gitops render-env -e custom -r https://github.com/yuan-shuo/helmChart -t v0.1.1
$ gitops render-env -e dev -r https://github.com/yuan-shuo/helmChart -t v0.1.1
$ gitops render-env -e test -r https://github.com/yuan-shuo/helmChart -t v0.1.1
# 将yaml文件放入后续cd程序用于同步的目录
$ cp custom/rendered/kustomize/demo-chart-0.1.1.yaml custom/cd-use
$ cp dev/rendered/kustomize/demo-chart-0.1.1.yaml dev/cd-use
$ cp test/rendered/kustomize/demo-chart-0.1.1.yaml test/cd-use
```

##### 提交 推送 PR tag

```bash
# push code for PR
$ git add . && git commit -m "yaml render" && git push -u origin feature/env
# after PR
$ git checkout main && git pull origin main && git tag v3.0.0 && git push origin tag v3.0.0
```

#### 生产环境仓库

##### 执行改动

这里我准备对生产环境仓库做的改动是把副本数从1调整为2

```bash
$ cd helmChart-env-prod
# 看看都有什么环境
$ ls
prod  README.md
# 进入事先准备好的 prod 环境
$ cd prod
```

然后看看values现在的副本数是几，用cat一看发现是1，随后用sed改为2

```bash
$ ls
cd-use  kustomization.yaml  patch.yaml  values.yaml
$ cat values.yaml | grep replicaCount:
replicaCount: 1
$ sed -i 's/replicaCount: 1/replicaCount: 2/' values.yaml
$ cat values.yaml | grep replicaCount:
replicaCount: 2
```

##### 渲染 yaml

渲染出 yaml 文件，然后丢到 cd-use 目录下，供后续持续交付程序同步

```bash
$ gitops render-env -e prod -r https://github.com/yuan-shuo/helmChart -t v0.1.1
$ cp prod/rendered/kustomize/demo-chart-0.1.1.yaml prod/cd-use
```

##### 提交 推送 tag

我实在懒得 `checkout + PR` 后自己审了，案例这里就直接在 main 上演示提交了

```bash
$ git add . && git commit -m "prod yaml render" && git tag v2.0.0 && git push origin main v2.0.0
```

### 环境仓库构建结果

可以看看现在都有了些什么：

```bash
# 非生产环境仓库 tag contain v3.0.0
remote = https://github.com/yuan-shuo/helm-chart-env-repo-non-prod
tag = v3.0.0
# 生产环境仓库 tag contain v2.0.0
remote = https://github.com/yuan-shuo/helm-chart-env-repo-prod
tag = v2.0.0
```

接下来就可以利用这些东西生成你的 `argocd.yaml` 了

## 构建 argocd 应用

### 生成应用 yaml 文件

也非常简单，直接用扩展两行生成双环境各自的 argo 应用 yaml 文件

```bash
# 为非生产环境生成应用 yaml 文件
$ gitops create-argo -r https://github.com/yuan-shuo/helm-chart-env-repo-non-prod -t v3.0.0 -m non-prod
# 为生产环境生成应用 yaml 文件
$ gitops create-argo -r https://github.com/yuan-shuo/helm-chart-env-repo-prod -t v2.0.0 -m prod
```

看一下当前目录都有啥，一份非生产环境应用yaml和一份生产环境应用yaml，依旧解放双手

```bash
$ ls -l | awk '{print $9}'
helm-chart-env-repo-non-prod-argocd-non-prod.yaml
helm-chart-env-repo-prod-argocd-prod.yaml
```

来看一下之前说的，把非生产环境有个环境从staging改成了custom，来看看argoyml这边能否感知：

```bash
$ cat helm-chart-env-repo-non-prod-argocd-non-prod.yaml | grep -- '- env'
      - env: custom # <-
      - env: dev
      - env: test
```

所以环境并不是锁定就那几个默认生成的，你可以随便改，里面有`kustomization.yaml`就行（因为argo靠这个文件确定这是否是一个环境目录，默认生成的`kustomization.yaml`什么都不做，纯当标记符，如果你需要使用，再去写它的内容）

### 应用 yaml 文件

`alias k = kubectl`

```bash
# 为k8s部署应用
root@dev-machine:laborant# k apply -f helm-chart-env-repo-non-prod-argo-non-prod.yaml
root@dev-machine:laborant# k apply -f helm-chart-env-repo-prod-argo-prod.yaml

# 先看看，发现prod没有同步，这是因为生成文件时指定了-m prod，默认关闭自动同步，需要手动同步
root@dev-machine:laborant# k get app -n argocd -o wide
NAME                                  SYNC STATUS   HEALTH STATUS
helm-chart-env-repo-non-prod-custom   Synced        Healthy         
helm-chart-env-repo-non-prod-dev      Synced        Healthy         
helm-chart-env-repo-non-prod-test     Synced        Healthy         
helm-chart-env-repo-prod-argo-prod    OutOfSync     Missing         

# 为生产环境手动同步 (手动同步指令在 *argo-prod.yaml 文件注释里已经提供了，直接复制过来用)
root@dev-machine:laborant# argocd app sync helm-chart-env-repo-prod-argo-prod

# 此时所有应用部署完成
root@dev-machine:laborant# k get app -n argocd -o wide
NAME                                  SYNC STATUS   HEALTH STATUS
helm-chart-env-repo-non-prod-custom   Synced        Healthy         
helm-chart-env-repo-non-prod-dev      Synced        Healthy         
helm-chart-env-repo-non-prod-test     Synced        Healthy         
helm-chart-env-repo-prod-argo-prod    Synced        Healthy         

# 可以看到 custom 应用在 non-prod 生效 (前述 non-prod env 改动)
# 同时 prod 应用副本数也变为 2 了 (前述 prod env 改动)
root@dev-machine:laborant# k get po -A | awk '{print $2 " - " $4}' | grep demo
custom-demo-chart-6bc68f495b-vzg6z - Running
dev-demo-chart-55f4cfb56b-whp9r - Running
test-demo-chart-784db6ffcc-47cd5 - Running
prod-demo-chart-66766fd86b-bmb6x - Running
prod-demo-chart-66766fd86b-jbxvh - Running
```





