# Getir Case

Getir technoical challenge

## build and Run Project on local

With Docker:

```sh
$ docker build -t getir .
$ docker run -d -it -p 8080:8080  -rm --name getir getir
```
## pull image from dockerhub
```sh
$ docker pull nedad/challenge:getir
````
## Run Test

-  go test ./...


## APIS (deployed on AWS ECS)

- `/fetch | http://18.234.166.11:80/fetch` `post request with bellow example`
- `/in-memory | http://18.234.166.11:80/in-memory` `get and post as example bellow `



# Mongo Search

**URL** : `/fetch`

**Method** : `POST`

**Data example**

```json
{
  "startDate": "2016-01-21",
  "endDate": "2016-03-02",
  "minCount": 2900,
  "maxCount": 3000
}
```

## Success Response

**Code** : `200`

**Content example**

```json
{
  "code": 0,
  "msg": "Success",
  "records": [
    {
      "createdAt": "2016-02-19T08:35:39.409+02:00",
      "key": "GjhjVIKb",
      "totalCount": 2774
    }
  ]
}
```

# Insert Data to memory

**URL** : `/in-memory`

**Method** : `POST`

**Data example**

```json
{
    "key": "test_key",
    "value": "test_value"
}
```

## Success Response

**Code** : `201`

**Content example**

```json
{
  "key": "test_key",
  "value": "test_value"
}
```
# Get Data from memory
*  `GET /in-memory`

# Get Data
**URL** : `/in-memory?key={test_key}`

**Method** : `GET`

## Success Response

**Code** : `200`

**Content example**

```json
{
  "key": "test_key",
  "value": "test_value"
}
```