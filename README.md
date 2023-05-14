# Pragmatic Live Feed Aggregator
## The Problem
###  Traffic consumption is very high and non-friendly to mobile and web clients during using Pragmatic Play's live data Socket API.
## How the app solves the problem?
### Pragmatic Live Feed Aggregator aggregates all available live tables and delivers to web clients on demand in one big batch.

---

# RUN in Docker Compose - Locally
### For running the app locally you just need to create `.env` file in the project root directory with the following environment variables in it:
```dotenv
PRAGMATIC_FEED_WS_URL= The vendor WS URL
CASINO_ID= Casino's ID
TABLE_IDS= Comma-seperated tables IDs
CURRENCY_IDS= Comma-seperated Currency IDs
REDIS_PORT= Port to run and connect to Redis DB
SERVER_PORT= Port to run the HTTP server on
PUSHER_CHANNEL_ID= Pusher.com's channel ID
PUSHER_PERIOD_MINUTES= Period to push data into the pusher channel
PUSHER_APP_ID= Pusher.com's app ID
PUSHER_KEY= Pusher.com's app key
PUSHER_SECRET= Pusher.com's app secret
PUSHER_CLUSTER= Pusher.com's app cluster
```
### And run the following command in the project's root directory:
```shell
docker-compose up
```

---

# See the result in your favorite Web Browser
### After running the app in Docker Compose, you need to open your favorite web browser and go to the following links:

1. All endpoints - SWAGGER UI
   ``http://localhost:[PORT]/v1/swagger``

2. Get the Pragmatic Live Feed aggregated data as a one big batch: \
   ```http://localhost:[PORT]/v1/tables``` \
   Where `PORT` is `SERVER_PORT` from the `.env` file. \
   API success response:
   ```json
   {
     "data": [
       {
         "tableAndCurrencyID": "100:200",
         "pragmaticTable": {
           "totalSeatedPlayers": 0,
           "last20Results": [
             {
               "time": "",
               "result": 0,
               "color": "",
               "gameId": "",
               "powerUpList": [],
               "powerUpMultipliers": []
             }
           ],
           "tableId": "",
           "tableName": "",
           "newTable": false,
           "languageSpecificTableInfo": "",
           "tableImage": "",
           "tableLimits": {
             "ranges": [],
             "minBet": 0,
             "maxBet": 0,
             "maxPlayers": 0
           },
           "dealer": {
             "name": ""
           },
           "tableOpen": false,
           "tableType": "",
           "tableSubtype": "",
           "currency": ""
         }
       }
     ]
   }
   ```

3. Check the previous endpoint health: \
   ```http://localhost:[PORT]/v1/healthcheck``` \
   Where `PORT` is `SERVER_PORT` from the `.env` file. \
   API success response:
    ```json
    {
       "status": "available",
       "system_info": {
       "version": "1.0.0"
       }
    }
    ```
