import requests
import datetime
import re
import json
import urllib
from bs4 import BeautifulSoup
from imdb import IMDb

base = datetime.datetime.today()
dates = [base + datetime.timedelta(days=x) for x in range(0,6)]

cinema_ids = {"Multiplex":24, "Loteanu":36}
# today = datetime.datetime.now().strftime("%d-%m-%Y")

result = {}
for cinema_name in cinema_ids:
    for date in dates:
        movieDay = date.strftime("%d-%m-%Y")
        url = "http://patria.md/beta/wp-admin/admin-ajax.php?date=" + movieDay + "&cinema=" + str(cinema_ids[cinema_name]) + "&action=flotheme_load_movies_scheduler"
        response = requests.get(url)
        soup = BeautifulSoup(response.text, 'html.parser')
        moviesHtml = soup.find('div', 'sidebar-scheduler-movies') 
        if moviesHtml:
            movies = []
            for link in soup.find_all("div", "title"):
                title = link.a.contents[0]
                if ("EN" in title):
                    movie_title = title.split('(')[0].rstrip()
                    if movie_title not in result:
                        result[movie_title] = {}
                    if cinema_name not in result[movie_title]:
                        result[movie_title][cinema_name] = {}
                        result[movie_title][cinema_name]['schedule'] = []
                    time = link.parent.find("div", "sessions").span.contents[0]
                    result[movie_title][cinema_name]['schedule'].append(
                        datetime.datetime.strptime(movieDay + ' ' + time, '%d-%m-%Y %H:%M').strftime('%b %d, %A at %H:%M')
                    )

ia = IMDb()

imdb_keys = ['cover url', 'rating', 'canonical title', 'plot', 'synopsis', 'long imdb title', 'genres', 'runtimes']
for name in result.keys():
    movie_data = ia.search_movie(name)[0]
    ia.update(movie_data)

    for key in imdb_keys:
        if key in movie_data:
            result[name][key] = movie_data[key]
    # now fill in youtube data
    query_string = urllib.parse.urlencode({"search_query" : movie_data['canonical title']})
    html_content = urllib.request.urlopen('http://www.youtube.com/results?' + query_string)
    soup = BeautifulSoup(html_content.read().decode(), 'html.parser')
    trailer = 'http://www.youtube.com' + soup.find('div', 'yt-lockup yt-lockup-tile yt-lockup-video vve-check clearfix').find('a')['href']
    result[name]['trailer'] = trailer

if result:
    try:
        handle = open('result.json', 'w')
        handle.write(json.dumps(result))
    finally:
        handle.close()
