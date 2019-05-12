from flask import Flask, render_template, g
from datetime import datetime
import json
import sqlite3

pirata = Flask(__name__)

DATABASE = '../pirata.db'

def get_db():
    db = getattr(g, '_database', None)
    if db is None:
        db = g._database = sqlite3.connect(DATABASE)
        return db

@pirata.teardown_appcontext
def close_connection(exception):
    db = getattr(g, '_database', None)
    if db is not None:
        db.close()

@pirata.route('/')
def index():
    cn = get_db()
    cur = cn.cursor()

    cur.execute('''SELECT premiere_date,
    meta,
    title
    FROM upcoming as up 
    JOIN films as f ON f.id = up.film_id
    WHERE up.premiere_date > ?''', [datetime.today().strftime('%Y-%m-%d')])

    upcoming = []
    for item in  cur.fetchall():
        upcoming.append({
            'premiere_date': item[0],
            'meta' : json.loads(item[1]),
            'title' : item[2]
        })
    
    # define what is playing right now by checking schedule
    cur.execute('SELECT film_id, location_id, datetime FROM schedule WHERE datetime >= ? ORDER BY location_id, datetime', [datetime.today().strftime('%Y-%m-%d')])
    schedule_times = cur.fetchall()

    # find film ids
    film_ids = set()
    [film_ids.add(x[0]) for x in schedule_times]

    query = f"SELECT id, title, meta FROM films WHERE id IN ({','.join(['?']*len(film_ids))})"
    cur.execute(query, list(film_ids))

    films = cur.fetchall()
    playing = []

    for film in films:
        id = film[0] 
        times = {}
        # filter playing times for this film only
        film_times = [x[1:] for x in (filter(lambda x: x[0] == id, schedule_times))]

        # separate by patria branches 24 and 36
        for x in film_times:
            if x[0] not in times:
                times[x[0]] = []
            times[x[0]].append(x[1])

        playing.append({
            'meta': json.loads(film[2]),
            'title': film[1],
            'schedule': times
        }) 

    return render_template('index.html', upcoming=upcoming, playing=playing)

if __name__ == "__main__":
        pirata.run()

