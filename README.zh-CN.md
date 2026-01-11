# Helm-GitOps

一个helm的扩展，为helm提供gitops相关辅助功能

## 使用

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

#### 版本管理

```bash
# 版本管理
helm gitops version # 仅查询当前版本
helm gitops version -m pr -l patch # 传统 PR 模式（先开 release 分支 → 提 PR → CI 自动 tag）
helm gitops version -m main -l patch # 快捷主分支模式（直接 commit + tag + 同时推送）
# --mode=main|pr
# --level=patch|minor|major
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
.
|-- helm-test1-env-non-prod
|   |-- README.md
|   |-- dev
|   |   |-- kustomization.yaml
|   |   |-- patch.yaml
|   |   `-- values.yaml
|   |-- staging
|   |   |-- kustomization.yaml
|   |   |-- patch.yaml
|   |   `-- values.yaml
|   `-- test
|       |-- kustomization.yaml
|       |-- patch.yaml
|       `-- values.yaml
`-- helm-test1-env-prod
    |-- README.md
    `-- prod
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

kustomization.yaml 会利用远程仓库链接及tag自动渲染例如下方yaml，其中name会利用Chart.yaml的name属性进行获取，同时判断values.yaml的fullnameOverride属性是否为空，非空则覆盖

```yaml
# staging/kustomization.yaml

apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

helmCharts:
- name: 'test-nor'
  repo: 'https://gitee.com/yuan-shuo188/helm-test1'
  version: 'v0.1.1'
  releaseName: 'staging'
  valuesFile: values.yaml

patchesStrategicMerge:
  - patch.yaml
```

#### 查看各环境使用的 Chart 版本

只需一行命令即可：

```bash
helm gitops env-version
```

效果如下，这样便不再需要逐个打开各个环境目录，寻找文件中不知何处写到的版本了

```bash
$ helm gitops env-version
dev: v0.1.1
staging: v0.1.1
test: v0.1.1
```

### argocd-yaml 生成功能

使用此工具可以节省编写argocd.yaml的时间

#### 创建操作

通过指定不同的参数，可以为非开发环境和开发环境各自生成对应的yaml

```bash
# 基于环境仓库的 remote 链接直接生成argo-yaml: 
# -r/--remote 指定环境仓库远程链接
# -t/--tag 指定生成 argo-yaml 时使用的远程仓库的 tag
# -m/--mode 指定生成模式: non-prod|prod, 会生成适用于不同环境下的 yaml 文件
# -d/--dry-run 不生成文件, 仅打印 argo-yaml 内容
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1  -t v1.0.0 -m prod --dry-run
```

之所以使用两项环境仓库+两份argocd-yaml是为了确保两环境独立审计，仓库及tag隔离，避免提交历史杂糅。同时argocd-yaml并不指向helm chart仓库而是仅指向环境仓库，因为此前环境仓库中已经指定过helm chart作为渲染源了，所以argo不需要同时指向两个仓库造成不必要的组合混乱，最终呈现可以示意为：`argocd -> env-repo -> helm-chart`

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
      project: default
      source:
        repoURL: 'https://gitee.com/yuan-shuo188/helm-env-non-prod1'
        targetRevision: 'v0.5.0'
        path: '{{env}}'
      destination:
        server: https://kubernetes.default.svc
        namespace: 'helm-env-non-prod1-{{env}}'
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
helm gitops create-argo -r https://gitee.com/yuan-shuo188/helm-env-prod1  -t v1.0.0 -m prod
```

生成的argo-yaml中automated默认为关闭

```yaml
# helm-env-prod1-argo-prod.yaml

apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: 'helm-env-prod1-argo-prod'
  namespace: argocd
  # annotations:
  #   # canary analysis template (optional, read by Argo Rollouts AnalysisTemplate)
  #   canary.argo.io/analysis-template: 
  #   canary.argo.io/step-weight: 
  #   canary.argo.io/step-duration: 
spec:
  project: default
  source:
    repoURL: 'https://gitee.com/yuan-shuo188/helm-env-prod1'
    targetRevision: 'v1.0.0'
    path: prod
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: prod
  syncPolicy:
    automated: false
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
helm plugin install https://github.com/yuan-shuo/helm-gitops
```

### 使用二进制文件

- 前往：[Releases · yuan-shuo/helm-gitops](https://github.com/yuan-shuo/helm-gitops/releases)下载对应操作系统的二进制文件

- 随后将解压得到的gitops二进制文件放在`$HELM_PLUGIN_DIR/bin/` 目录下

- 给予gitops二进制文件执行权限

```bash
chmod +x $HELM_PLUGIN_DIR/bin/gitops
```

## 环境需求

* git（version>=2.23）
* helm
* helm-unitest (可选，可以通过执行 **`helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false`** 来安装)
