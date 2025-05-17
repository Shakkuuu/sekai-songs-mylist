# sekai-songs-mylist

- Chartの実装と検証 ok
- 存在しないartsitIDとかsingerIDとか選ばれたときのinvalid argエラー返し ok
- badrequest ok
- 登録したらリロードされたい
- マスター管理画面整備、選択式とか ok
- enumは選択式にしたが、singerやartistもlistを取得してそれを選択肢にしたい ok
- singerのpositionによるsortがうまくいかない ok

- export PATH="$(pwd)/../../view/node_modules/.bin:$PATH"
- buf generate
- sqlc generate
- 既存のデータ抽出
  - pg_dump -U db_user_name -h localhost -p 5432 -d dbname -Fp --data-only > dbdb.sql
