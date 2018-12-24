from flask import Flask, render_template
import json

pirata = Flask(__name__)

@pirata.route('/')
def index():
    with open('../result.json') as file:
        schedule = json.load(file)

    return render_template('index.html', schedule=schedule)