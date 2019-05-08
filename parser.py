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
patria_url = 'http://patria.md/?mode=normal'
response = requests.get(patria_url)
soup = BeautifulSoup(response.text, 'html.parser')

upcoming_section = soup.find('ul', 'coming-soon')
for item in upcoming_section.find_all('li'):
    title = item.find('a', 'title').contents[0]
    if ("EN" in title):

        movie_title = title.split('(')[0].rstrip()
        premiere_date = item.span.contents[0] 
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
