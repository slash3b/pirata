import requests
import re
import sys
from datetime import datetime
from datetime import timedelta
from bs4 import BeautifulSoup
from imdb import IMDb
from sqlite3 import connect
from parser import upcoming, playing, mail

cn = connect('pirata.db')

# initialize client with english lang preferency
client = requests.Session()
client.get('https://patria.md/lang/en')

# coming soon section
url = f'https://patria.md/films'
response = client.get(url)
soup = BeautifulSoup(response.text, 'html.parser')
movies = soup.find('div', attrs={'id':'afisha2_'}).children

for movie in movies:
    movie_title = movie.find('div', class_='name')
    if movie_title == None:
        continue
    movie_title = movie_title.contents[0]
    if ('(EN)' in movie_title):
        # grab movie title in english
        movie_title = movie_title.split('(')[0].rstrip()

        # premiere_date = datetime.strptime(premiere_date, '%d.%m.%Y').date() 
        date = movie.find('div', class_='premier')
        premiere_date = date.contents[0].split(' ')[1]
        upcoming.register(cn, movie_title, premiere_date)

# playing now section
base = datetime.today()
dates = [base + timedelta(days=x) for x in range(0,6)]

# delete future schedules. Safest way to make sure we have correct schedule
cn.cursor().execute('DELETE FROM schedule WHERE datetime >=?', (datetime.today().strftime('%Y-%m-%d'),))
cn.commit()

cinema_map = {
       'Multiplex' : 24, 
       'Loteanu' : 36
}

for date in dates:
    # parse path for one day
    # e.g. https://patria.md/films?c_date=2019-12-16
    movieDay = date.strftime('%Y-%m-%d')

    url = f'https://patria.md/films?c_date={movieDay}'
    response = client.get(url)
    soup = BeautifulSoup(response.text, 'html.parser')

    movies = soup.find('div', attrs={'id':'afisha1_'}).children

    for movie in movies:
        # exclude the fucking banner
        movie_title = movie.find('div', class_='name')
        if movie_title == None:
            continue
        movie_title = movie_title.contents[0]
        if ('(EN)' in movie_title):
            # MOVIE TITLE
            movie_title = movie_title.split('(')[0].rstrip()

            for cinema in movie.find_all('div', class_='h4'):
                # CINEMA ID
                cinema_id = cinema_map[cinema.contents[0].rstrip(':')]
                # times for this particular cinema 

                # find here all schedules
                schedule = []
                for times in cinema.find_next().find_all('a'):
                    # e.g. 16:40
                    time = times.find('span').contents[0]
                    date = datetime.strptime(movieDay + ' ' + time, '%Y-%m-%d %H:%M')
                    schedule.append(date)

                playing.register(cn, movie_title, cinema_id, schedule)

result = cn.cursor().execute('select title from films where date(register_date) = ?', (datetime.today().strftime('%Y-%m-%d'),))

films = ''
for item in result.fetchall():
    films += item[0] + ', '

if(len(films)):
    mail.send_mail(films)

# terminate connection
cn.close()

print("All is good!")
