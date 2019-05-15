import requests
import re
from datetime import datetime
from datetime import timedelta
from bs4 import BeautifulSoup
from imdb import IMDb
from sqlite3 import connect
from parser import upcoming, playing

cn = connect('pirata.db')

# coming soon section
patria_url = 'http://patria.md/movies'
response = requests.get(patria_url)
soup = BeautifulSoup(response.text, 'html.parser')

# since html markup is a piece of shit like the rest of the patria we take all the items from the second 'movies-page' div
upcoming_movies = soup.find_all('div', 'page movies-page')[1].find_all('div', 'movies-item')
for item in upcoming_movies:
    attrs = item.find('a').attrs
    if ("EN" in attrs['title']):
        movie_title = attrs['title'].split('(')[0].rstrip()
        movie_url = attrs['href']
        # in order to get premiere date we have to visit each page, that sucks
        response = requests.get(movie_url)
        soup = BeautifulSoup(response.text, 'html.parser')
        # todo: way to improve is to make request and find out if we already have that film in the database ?
        premiere_date = soup.find_all('div', class_='premiere')[-1].contents[1] 
        upcoming.register(cn, movie_title, premiere_date)

# playing now section
base = datetime.today()
dates = [base + timedelta(days=x) for x in range(0,6)]

# delete future schedules. Safest way to make sure we have correct schedule
cn.cursor().execute('DELETE FROM schedule WHERE datetime >=?', (datetime.today().strftime('%Y-%m-%d'),))
cn.commit()

for date in dates:
    movieDay = date.strftime("%d-%m-%Y")
    # cinema 24 stands for Mall cinema
    # cinema 36 stands for Loteanu cinema
    for cinema_id in ['24', '36']:
        url = "http://patria.md/beta/wp-admin/admin-ajax.php?date=" + movieDay + "&cinema=" + cinema_id + "&action=flotheme_load_movies_scheduler"
        response = requests.get(url)
        soup = BeautifulSoup(response.text, 'html.parser')
        moviesHtml = soup.find('div', 'sidebar-scheduler-movies') 
        if moviesHtml:
            for link in soup.find_all("div", "title"):
                title = link.a.contents[0]
                if ("EN" in title):
                    # skip new lines and return only html elements
                    schedule_content =  [x for x in link.parent.find("div", "sessions").children if x != '\n']
                    schedule = []
                    for item in schedule_content:
                        if item:
                            # get time string
                            time = item.contents[0].strip('/ ')
                            date = datetime.strptime(movieDay + ' ' + time, '%d-%m-%Y %H:%M')
                            schedule.append(date)
                    movie_title = title.split('(')[0].rstrip()
                    playing.register(cn, movie_title, cinema_id, schedule)

# terminate connection
cn.close()

print("All is good!")
