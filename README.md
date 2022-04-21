### pirata

##### How to run locally:
- install pipenv
- run `pipenv install`
- run `pipenv shell`
- cd to `flask` directory
- do `export FLASK_ENV=development`
- then `flask run`

Command to convert "API documentation" to html is

```bash
./to_md.py templates/api.md -o templates/api.html
```

##### build parser image
docker build . -t pirata:latest

##### run dev image
docker run --volume (pwd):/pirata --rm -it pirata /bin/sh
