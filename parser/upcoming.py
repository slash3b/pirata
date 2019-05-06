from datetime import datetime
from parser import film

def register(cn, movie_title, premiere_date):
    cursor = cn.cursor()

    cursor.execute('SELECT * FROM films WHERE title=?', (movie_title,))
    if cursor.fetchone() is None:
        film.register(cn, movie_title)

    cursor.execute('SELECT * FROM films WHERE title=?', (movie_title,))
    record = cursor.fetchone()

    # record id
    film_id = record[0]
    premiere_date = datetime.strptime(premiere_date, '%B %d, %Y').date()

    cursor.execute('SELECT * FROM upcoming WHERE film_id=?', (film_id,))
    if cursor.fetchone() is None:
        data = (film_id, premiere_date)
        cursor.execute('INSERT INTO upcoming(film_id, premiere_date) VALUES(?, date(?))', data)
        cn.commit()
