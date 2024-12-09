# gm-backend

## 開発
ローカルでgoが必要  
`src`フォルダの`Makefile`を使用する。  
  
実行
```
make run
```

## ただ動かすだけ
ローカルでdockerが必要。  
`Dockerfile`があるディレクトリの`Makefile`を使用する。
  
実行
```
make docker-build
make docker-up
```