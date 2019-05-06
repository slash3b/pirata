from parser import film
from datetime import datetime

def register(cn, movie_title, schedule):
    cursor = cn.cursor()
    cursor.execute('SELECT * FROM films WHERE title=?', (movie_title,))
    if cursor.fetchone() is None:
        film.register(cn, movie_title)

    cursor.execute('SELECT * FROM films WHERE title=?', (movie_title,))
    record = cursor.fetchone()
    film_id = record[0]
    
    for time in schedule:
        cursor.execute('INSERT INTO schedule(film_id, datetime) VALUES(?, ?)', (film_id, time))
    cn.commit()
