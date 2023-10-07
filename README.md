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

![screencapture-localhost-9695-console-remote-schemas-manage-server-schema-relationships-2023-10-07-16_54_04](https://github.com/takeruun/go-elastice-search/assets/48900966/b79a8231-38f1-46d7-a3c9-de0ead1bd711)


3. 紐付け定義したものを`remote-schema(go server)`で定義したクエリで使う。

 <img width="1093" alt="スクリーンショット 2023-10-07 16 54 57" src="https://github.com/takeruun/go-elastice-search/assets/48900966/2057a1e6-7e59-4dc3-a195-638518038877">
