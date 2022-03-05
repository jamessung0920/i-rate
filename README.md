## 專案啟動前需執行
0. 從`.env.example`複製一份`.env`並填入設定, CHANNEL_TOKEN 填 line 的 channel access token

1. 於 db 新增貨幣資料
    ```sql
    insert into currency (currency) values ('USD');
    insert into currency (currency) values ('JPY');
    insert into currency (currency) values ('GBP');
    insert into currency (currency) values ('EUR');
    insert into currency (currency) values ('THB');
    ```

2. 觸發爬蟲開始爬，並觸發主動通知
    ```
    $ curl <https://YOUR_DOMAIN>/currency/crawl
    $ curl <https://YOUR_DOMAIN>/currency/rate/notify
    ```

## 啟動專案

```
$ docker-compose up -d
```