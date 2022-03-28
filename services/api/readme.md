incoming request
    validation maybe and processing
        converter from http to internal DTO like struct
            an internal request that knows how to transform a request and return a model
        converter from model to external representation
    a response writer that knows ho to use external representation


This is just to run docs.
```
slash3b@Nostromo ~/P/p/p/api (master)> docker run -it --rm -p 80:80 \
-v (pwd)/swagger.yaml:/usr/share/nginx/html/swagger.yaml \
-e SPEC_URL=swagger.yaml redocly/redoc
```
Validate your spec with 
```
docker run --rm -v $PWD:/spec redocly/openapi-cli lint petstore.yaml

```

Regenerate yaml to json
```
docker run --rm -v (pwd):/spec redocly/openapi-cli bundle --output openapi --ext json openapi.yaml
```
