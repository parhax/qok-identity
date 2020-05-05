## Getting Started

For running this Go app you should first build its image by running

`docker-compose up --build `

 **Note** : 
> some dependencies which should be downloaded for building this image are banned for IPs with IRR location , so please use some proxy or vpn before running the above command


3 HTTP endpoints have been provided in this code which you can test and examine it via Postman application or curl like following codes:

**Register**
```
curl --location --request POST 'http://127.0.0.1:8585/register' \
--header 'Content-Type: application/json' \
--data-raw '{
	"username" : "my_sample_3",
	"password" : "123123"
}'

```

***

**Login**
```
curl --location --request POST 'http://127.0.0.1:8585/login' \
--header 'Content-Type: text/plain' \
--data-raw '{
	"username" : "my_sample_2",
	"password" : "123123"
}'
```
***

**UserInfo**
```
curl --location --request GET 'http://127.0.0.1:8585/user_info' \
--header 'Authorization: {jwt which got from login response}
```
