### pirata ![alt text](https://pirata.md/static/favicon.ico)

name pirata is an anagramm for patria which is romanian word and can be translated as motherland

this simple project shows films schedule of patria.md cinema in Chishinau

why? because the original site is ugly, slow and not user friendly


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
