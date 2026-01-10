go build -o bin/gitops .
./bin/gitops create test-nor
cd test-nor
git remote add origin https://gitee.com/yuan-shuo188/helm-test1.git
git push -u origin main
../bin/gitops version # 仅查询当前版本
../bin/gitops version -m main -l patch # 快捷主分支模式（直接 commit + tag + 同时推送）