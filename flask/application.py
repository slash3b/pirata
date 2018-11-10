from flask import Flask, render_template
import json

app = Flask(__name__)

@app.route('/')
def index():
    with open('../result.json') as file:
        schedule = json.load(file)

    return render_template('index.html', schedule=schedule)