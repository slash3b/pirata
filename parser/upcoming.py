from datetime import datetime
from parser import film

def is_registered(cn, movie_title):
    cursor = cn.cursor()

    cursor.execute('SELECT * FROM films WHERE title=?', (movie_title,))

    return cursor.fetchone() != None

def register(cn, movie_title, premiere_date):
    cursor = cn.cursor()
    
    if not is_registered(cn, movie_title):
        film.register(cn, movie_title)

    cursor.execute('SELECT * FROM films WHERE title=?', (movie_title,))
    record = cursor.fetchone()

    # record id
    film_id = record[0]

    cursor.execute('SELECT * FROM upcoming WHERE film_id=?', (film_id,))
    if cursor.fetchone() is None:
        data = (film_id, premiere_date)
        cursor.execute('INSERT INTO upcoming(film_id, premiere_date) VALUES(?, date(?))', data)
        cn.commit()
