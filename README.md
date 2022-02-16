# NFP Server

## 起動方法
```bash
docker-compose up -d
```
Open localhost:8080


## デバック方法
1.
```bash
ngrok http 8080
```
でネットワークを外部公開する
[参考](https://qiita.com/poccariswet/items/24fac246f8760abfb51e)

2. `https://~`のurlをコピーして LINE Message APIのWebhool URLに設定する 
(https://~.ngrok.io/linebot/callback)
![image](https://user-images.githubusercontent.com/53213591/154267908-91c388aa-d1c7-4c20-8675-794f5c07ba4e.png)
![スクリーンショット 2022-02-16 21 50 23](https://user-images.githubusercontent.com/53213591/154268113-271c7de3-efca-481f-9e8e-459c2d675ef8.png)
