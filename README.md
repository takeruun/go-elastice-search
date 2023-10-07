# go ElasticSearch on hasura
go で　ElasticSearch を使い、hasura 経由でデータ検索するサンプルです。

## 使い方
`query searchItem` を叩く。

例：

```graphql
query {
  searchItems(where: {title: "title"}) {
    id
    server_item {
      title
    }
  }
}
```

hasura → go server → hasura の順で実行している。

### 方法
1. `remote-schema(go server)`で定義したクエリ`searchItems`を hasura で読み込ませる。

2. `remote-schema(go server)`の `type Item` と DB の `items`テーブルを紐付ける。

3. 紐付け定義したものを`remote-schema(go server)`で定義したクエリで使う。