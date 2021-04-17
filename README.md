Reducto
=======
A simple but performant and distributed URL shortener.

### API:
___Shorten URL___

----
Shortens provided URL and returns json data

* **URL**

  /api/v1/shorten

* **Method:**

  `POST`

* **Data Params**
  **Required:**
  `url=[string]`
  **Optional:**
  `custom=[string]`

* **Success Response:**

    * **Code:** 201 Created <br />
      **Content:** `{"longUrl":"http://foobar.com","shortUrl":"http://localhost:8081/FPvYGy"}`

* **Error Response:**

    * **Code:** 422 Unprocessable Entity <br />
      **Content:** `{ error : "User doesn't exist" }`

  OR

    * **Code:** 422 Unprocessable Entity <br />
      **Content:** `{"error":"URLValidationError: URL Domain is banned."}`

---
___Get Long URL___

----
Retrieves long URL and redirects

* **URL**

  /:shortUrl

* **Method:**

  `GET`

* **Path Params**
  **Required:**
  `shortUrl=[string]`

* **Success Response:**

    * **Code:** 302 Moved Permanently <br />
      **Content:** `<a href="http://foobar.com">Moved Permanently</a>.`

* **Error Response:**

    * **Code:** 404 Not Found <br />
      **Content:** `{"error":"Error in getSlug for f00bar: record not found"}`

  OR

    * **Code:** 422 Unprocessable Entity <br />
      **Content:** `{"error":"URLValidationError: URL Domain is banned."}`

### Installation
#### Via Docker:
```
git clone git@github.com:AnkurGel/reducto.git
docker-compose up
$> curl --location --request POST 'http://localhost:8080/api/v1/shorten' --form 'url="https://github.com/AnkurGel/reducto"'
```

### Development
* `cp config.yml.sample config.yml`
* Edit `config.yml`
* Create relevant database
* Add banned hosts in the redis manually like:   
&nbsp;&nbsp;&nbsp;`SADD urlBannedSet bit.ly tinurl.com tiny.one t.co rotf.lol goo.gl fb.me`
* Run keygen: `REDUCTO_CONFIG_PATH=config.yml go run cmd/reducto-keygen/main.go`
* Run server: `REDUCTO_CONFIG_PATH=config.yml go run cmd/reducto-server/main.go`
