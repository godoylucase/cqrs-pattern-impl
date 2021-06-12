## CQRS Pattern Example

### How to use it?
To run just run
```
$ make build
$ make run
```


#### Article Tags
##### Create
```
curl -L -X POST 'localhost:8081/v1/articles' --data-raw '{
    "user_id": "some_user_id",
    "title": "humans on mars",
    "source_url": "https://humans-on-mars.com",
    "global_hash_tags": [
        "space",
        "moon",
        "rocket",
        "elonmusk"
    ],
    "paragraphs": [
        {
            "text": "foo foo",
            "hash_tags": [
                "space",
                "moon",
                "rocket"
            ]
        },
        {
            "text": "bar bar",
            "hash_tags": [
                "science",
                "moon"
            ]
        }
    ]
}'
```

##### Get Article
```
curl -L -X GET 'localhost:8081/v1/articles/<some_id>' 
```

##### Update Article
```
curl -L -X PUT 'localhost:8081/v1/articles/<some_id>' -H 'Content-Type: application/json' --data-raw '{
    "title": "earths on mars",
    "global_hash_tags": [
        "space",
        "moon",
        "rocket",
        "elonmusk",
        "bitcoin",
        "crypto"
    ],
    "paragraphs": [
        {
            "text": "bla bla",
            "hash_tags": [
                "space",
                "moon",
                "rocket"
            ]
        },
        {
            "text": "ble ble",
            "hash_tags": [
                "science",
                "moon"
            ]
        },
        {
            "text": "blo blo",
            "hash_tags": [
                "space",
                "rocket"
            ]
        }
    ]
}' 
```

#### Read Tags
##### Get articles by global hashtag
```
curl -L -X GET 'http://localhost:8082/v1/articles-by-global-hash-tags?global_hash_tags=moon'
```

##### Get user articles by source URL
```
curl -L -X GET 'http://localhost:8082/v1/user-articles-by-source-url?source_url=https://humans-on-mars.com'
```
