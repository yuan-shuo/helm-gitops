[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/helm-gitops)](https://artifacthub.io/packages/search?repo=helm-gitops)

[带有介绍的完整gitops流程使用示例](https://github.com/yuan-shuo/helm-gitops/blob/main/doc/example.zh-CN.md)

# Helm-GitOps

一个helm的扩展，为helm提供gitops相关辅助功能

## 使用

### 如果你什么都不想看

就这三条，`remote / tag` 改成自己的然后复制过去直接用

```bash
# 1.创建一个 git 初始化过的 Chart
helm gitops create my-chart

# 2.生成一个环境仓库 (基于 Chart 远程仓库 + 仓库tag)
helm gitops create-env -r https://gitee.com/yuan-shuo188/helm-test1 -t v0.1.1

# 3.生成一个argo.yaml (基于环境仓库 + 仓库tag)
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-non-prod1 -t v0.5.0 -m non-prod
```

生成内容都是基于信息生成所以省掉很多麻烦（指以人类之躯在多个远程仓库不停跳转复制肉眼检查等），自己敲 git 管理就行，如果下面的东西你都不想看，上面三条也能帮你解决大部分麻烦了。

`（标题中被 * 标记的功能也许能够帮助你，如果你只想花一点时间来看看的话）`

### chart-git 图表开发功能

相对于常规自行创建helm chart，随后在内部粘贴.gitignore等文件，自行修改创建分支、提交、版本号等一系列繁琐操作，此扩展提供了较为舒适的简化方案：

#### 创建图表

```bash
# 创建：在原 helm chart 基础上添加 GitOps 骨架
helm gitops create test
helm gitops create test --actions        # 同时生成 .github/workflows/<>.yaml
```

#### 创建分支

```bash
# 切换分支：自动进入开发分支（可选同步主分支）
helm gitops checkout feature/foo
helm gitops checkout feature/foo -s       # 先 pull origin/main 再创建分支
```

#### 提交推送（主分支保护）

```bash
# 提交代码：add + commit + 可选 push & 自动 PR
helm gitops commit -m "fix: foo"                    # 本地提交
helm gitops commit -m "feat: xxx" --push            # 提交并推送
helm gitops commit -m "ci: update" --pr --push      # 提交 + 推送 + 自动提 PR（含 [create-pr] commit 标记）
```

#### 图表检查

```bash
# 本地检查：helm lint + helm unittest
helm gitops lint
```

#### 图表推送

```bash
# 推送代码：lint → push（保护分支拦截）
helm gitops push                                      # 推送到 origin/当前分支
```

#### 版本管理 *

一行指令就能够自动完成版本更新，自动清理旧版本打包tgz，新版本tgz打包，index.yaml生成，推送

```bash
# 版本管理
helm gitops version # 仅查询当前版本
helm gitops version -m pr -l patch # 传统 PR 模式（先开 release 分支 → 提 PR → CI 自动 tag）
helm gitops version -m main -l patch # 快捷主分支模式（直接 commit + tag + 同时推送）
helm gitops version -m main -l no -s pre # v0.0.1 -> v0.0.1-pre
# --mode/-m = main|pr
# --level/-l = patch|minor|major|no(无版本数字变化, 可以搭配-s在同一vx.x.x不断更新后缀)
# --suffix/-s = <your_tag_suffix>
```

### env-repo 环境仓库功能

利用此工具可以在环境仓库配置中节省时间

#### 创建操作

直接基于helm chart的远程仓库链接生成2个环境仓库，文件中的内容会根据远程仓库链接等信息渲染得到，减少不必要的手工写入和复制等操作

```bash
# 基于 Chart 的 remote 链接直接生成环境仓库: 
# -r/--remote 指定 helm chart 远程仓库链接
# -t/--tag 指定创建仓库时使用的 chart 仓库的 tag
helm gitops create-env -r https://gitee.com/yuan-shuo188/helm-test1 -t v0.1.1
```

#### 目录树

仅需执行上面一行即可生成如下目录树，可以看到为非生产环境和生产环境各创建了一个仓库，每个仓库目录下同时包含 `.git(已经初始化过) + .gitignore`

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

#### 文件内容

values.yaml 是从远程仓库对应tag的代码中复制过来的，这样就不需要再返回chart仓库观察默认values.yaml长什么样子，同时每个文件顶部均有**`目录/文件名`**的标记注释

```yaml
# dev/values.yaml

# Default values for test-nor.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

# This will set the replicaset count more information can be found here: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/
replicaCount: 1
```

kustomization.yaml 并不需要特别修改，根据远程仓库提前协助写入的只有两行注释，帮助观察文件所在位置以及使用渲染功能（line1, lin3注释处），可以注意到这里尽管使用helm作为渲染源但并没有指定repo和tag，这会在后续参数化渲染部分解释

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

patch.yaml默认为空

```yaml
# prod/patch.yaml
```

#### 参数化渲染 *

如果你只需要helm的渲染结果，那不动kustomization.yaml（保持为空）就可以了，如果你没有kust软件也没关系，会自动跳过kust渲染，你可以在 **`<your_env>/rendered/helm/helm-chart.yaml`** 找到你的helm渲染结果

##### 指令

直接看指令：

```bash
# -e/--env 想要渲染的环境
# -r/--remote 远程仓库链接
# -t/--tag 远程仓库 tag
# -l/--use-local-cache 使用本地已有的文件渲染
# -n/--render-file-name 自定义命名渲染结果
helm gitops render-env -e prod -r https://gitee.com/yuan-shuo188/helm-test1.git -t v0.1.3
```

执行指令后便可得到一份基于helm chart，经过kustomize渲染的yaml结果文件

##### 说明

###### --remote / -r

remote使用的是git仓库链接，此部分直接省去了发布操作，**<u>你的helm chart仓库只需要有index.yaml和对应的tgz包就足够了</u>**（而且如果你使用的是本软件的版本管理功能，这两个东西会自动构建不需要你操心），程序会基于remote前往仓库raw寻找index.yaml，随后利用其内部的urls属性获得tgz文件名，随后与remote/tag拼接形成chart打包文件的下载链接，获取到对应环境下的charts目录，随后渲染出helm的yaml文件，然后利用kustomize进一步渲染（如果你没有这个软件会自动跳过，仅渲染到helm为止）

###### --use-local-cache / -l

如果使用kustomize的helm功能，让其指向一个helm chart没问题，可以，但是如果需要反复调试生成的话，不停的网络请求是没有必要的，这时你可能会把helm渲染结果文件下载到本地，然后让kustomize指向它，这时利用软件的 `-l` 参数，你就完全没必要来回折腾了，如果缓存目录（例如prod/rendered/helm/helm-chart.yaml）包含helm文件它自己会直接渲染而不提出网络请求

###### --render-file-name / -n

覆盖渲染结果文件的命名（存在默认命名但也许你会有其他命名的想法）

##### 一个小案例

先在chart仓库提交一个新版本

```bash
helm gitops version -l major -m main
```

然后用tree指令能看到index和tgz都准备好了，可以看到tag自动升级到2.0.0（major级更新）：

```bash
|-- index.yaml
|-- test-nor-2.0.0.tgz
```

然后可以使用指令自动构建环境仓库了，这里选择生产仓库用于演示

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

这时随便用什么编辑器打开 `kustomization.yaml`，能看到第三行的注释已经为你准备好渲染指令了：

```bash
helm gitops render-env -e prod -r https://gitee.com/yuan-shuo188/helm-test1 -t v2.0.0
```

执行后来用tree指令看一看目录发生什么变化：

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

可以看到已经渲染好了，多出来一个叫redered的目录，其中：

* `helm-chart.yaml` 为helm渲染结果

* `test-nor-2.0.0.yaml` 为kustomize基于 `helm-chart.yaml` 的渲染结果

这样就结束了，总共才三条指令就结束了（可以把你需要部署的yaml放入 `cd-use` 目录供持续交付程序使用）。因为渲染是参数化的，想要修改tag或者repo可以去调整那条帮你写好的注释，然后复制注释指令一键渲染就行。

关于为什么不用kustomize的helm-chart参数，直接由repo+version控制渲染源的解释：

* kust需要的是helm仓库而不是git仓库，并不是所有托管平台都有github那样丰富的功能例如pages

* kust和helm并不是特别兼容（尽管kust提供了helm功能，但是helm目前在v4+版本的改动让kust在使用旧指令调用时会报很多参数错误）。参数化渲染功能让helm渲染和kust渲染分离解耦，二者之间没有依赖关系：

`helm-chart.git -> (helm render) -> helm-chart.yaml -> (kust render) -> final.yaml`

### argocd-yaml 生成功能

使用此工具可以节省编写argocd.yaml的时间

#### 创建操作

通过指定不同的参数，可以为非开发环境和开发环境各自生成对应的yaml

```bash
# 基于环境仓库的 remote 链接直接生成argo-yaml: 
# -r/--remote 指定环境仓库远程链接 (确保仓库可达)
# -t/--tag 指定生成 argo-yaml 时使用的远程仓库的 tag
# -m/--mode 指定生成模式: non-prod|prod, 会生成适用于不同环境下的 yaml 文件
# -d/--dry-run 不生成文件, 仅打印 argo-yaml 内容
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1 -t v1.0.0 -m prod --dry-run
```

之所以使用两项环境仓库+两份argocd-yaml是为了确保两环境独立审计，仓库及tag隔离，避免提交历史杂糅。同时argocd-yaml并不指向helm chart仓库而是仅指向环境仓库，因为此前环境仓库中已经指定过helm chart作为渲染源了，所以argo不需要同时指向两个仓库造成不必要的组合混乱，最终呈现可以示意为：`argocd -> env-repo -> helm-chart`

而且由前述内容可知，argocd或者其他持续交付程序仅需指向env/cd-use目录，该目录内部存放kust+helm / helm渲染结果yaml，是一份可以直接被kubectl apply的文件，这代表argo或者其他持续交付软件不需要安装额外的插件，同时yaml文件会在 helm repo / env repo 流程中多次被审查，并以所见即所得的最终完整结果yaml存放在env/cd-use目录下

#### 目录树

执行如下指令：

```bash
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-non-prod1  -t v0.5.0 -m non-prod
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1  -t v1.0.0 -m prod
```

可以得到两份yaml文件

```
|-- helm-env-non-prod1-argocd-non-prod.yaml
|-- helm-env-prod1-argocd-prod.yaml
```

#### 文件内容

会根据环境仓库信息生成，生成的yaml提前准备好了大部分原本需要手动填写的信息，使得使用者可以聚焦问题而非查找并填写信息

##### 非生产 argo-yaml

对于非生产环境仓库可能包含多个环境值目录，例如dev，test等，当指定 `--mode/-m non-prod` 时，程序能够自行寻找远程仓库所有包含 `kustomize.yaml` 的一级目录自动添加到argo-yaml中

例如在env仓库中，额外添加一个叫做 `anthor` 的环境仓库：

```
|-- helm-test1-env-non-prod
|   |-- anthor-env
|   |   |-- kustomization.yaml # 搜寻核心
|   |   |-- patch.yaml
|   |   `-- values.yaml
```

随后打 tag=v0.5.0 后 push 代码和 tag，使用

```bash
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-non-prod1  -t v0.5.0 -m non-prod
```

生成 yaml 结果如下所示，可以看到非本程序默认生成，用户手动创建的环境目录 `anthor-env` 也被写入进argo-yaml了

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

##### 生产用 argo-yaml

仅指定prod环境作为控制源，执行如下指令后生成 yaml 文件

```bash
helm gitops create-argo -r https://github.com/yuan-shuo/helm-env1 -t v1.0.0 -m prod
```

生产模式（`--mode prod`）生成的argo-yaml中automated默认为关闭，并且在第三行生成了一条复制即用的手动同步指令注释

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

## 安装

### 使用 helm plugin install

```bash
helm plugin install https://github.com/yuan-shuo/helm-gitops/releases/download/v0.5.2/helm-gitops_0.5.2_linux_amd64.tar.gz
```

## 环境需求

* git（version>=2.23）
* helm
* helm-unitest (可选，可以通过执行 **`helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false`** 来安装)
* kustomize（可选，未安装时，渲染会跳过kust仅用helm）
