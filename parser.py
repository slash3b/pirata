import requests
import datetime
import re
import json
from bs4 import BeautifulSoup

base = datetime.datetime.today()
dates = [base + datetime.timedelta(days=x) for x in range(0,6)]

multiplex_id = 24
# today = datetime.datetime.now().strftime("%d-%m-%Y")

result = []

for date in dates:
    movieDay = date.strftime("%d-%m-%Y")
    url = "http://patria.md/beta/wp-admin/admin-ajax.php?date=" + movieDay + "&cinema=" + str(multiplex_id) + "&action=flotheme_load_movies_scheduler"
    response = requests.get(url)
    soup = BeautifulSoup(response.text, 'html.parser')
    moviesHtml = soup.find('div', 'sidebar-scheduler-movies') 
    if moviesHtml:
        movies = []
        for link in soup.find_all("div", "title"):
            title = link.a.contents[0]
            if ("EN" in title):
                item = {
                    "title" : title,
                    "time" : link.parent.find("div", "sessions").span.contents[0],
                }
                movies.append(item)
        if movies:
            result.append({
                'date' : movieDay,
                'movies' : movies,
            })

if result:
    try:
        handle = open('result.json', 'w')
        handle.write(json.dumps(result))
    finally:
        handle.close()