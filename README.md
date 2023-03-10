# Proxy Manager
This is a proxy manager. It users [mitmproxy](https://mitmproxy.org/) as reverse proxy. Rules should be written in Python

### Rule examples

Example or rules you can find in [mitmproxy examples](https://github.com/mitmproxy/mitmproxy/tree/main/examples)


### Manual run

1. Clone repo
2. Go to app folder `cd app`
3. Build & run server `go build . && ./main`


### Routes

- GET /proxies: Get a list of all proxies
- GET /proxies/{name}: Get details of a specific proxy
- POST /proxies: Add a new proxy
- PUT /proxies: Update an existing proxy
- DELETE /proxies/{name}: Delete a proxy

### Endpoints examples

#### POST /proxies
Create and run proxy with given settings.
```
curl --location 'http://localhost:8000/proxies' \
--form 'serviceName="Name"' \
--form 'serviceUrl="http://localhost:8000"' \
--form 'listenPort="5000"' \
--form 'proxyType="http"' \
--form 'filterFile=@"./settings.py"'
```

#### GET /proxies
Get all proxies.
```
curl --location 'http://localhost:8000/proxies/'
```

#### GET /proxies/{name}
Get proxy by its name.
```
curl --location 'http://localhost:8000/proxies/service_name_example'
```
#### DELETE /proxies/{name}
Delete proxy by its name, also it disables it.
```
curl --location --request DELETE 'http://localhost:8000/proxies/service_name_example'
```

