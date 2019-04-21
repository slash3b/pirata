import requests
from datetime import datetime
from datetime import timedelta
import re
import json
import urllib
from bs4 import BeautifulSoup
from imdb import IMDb
from pprint import pprint

base = datetime.today()
dates = [base + timedelta(days=x) for x in range(0,6)]

# coming soon section
patria_url = 'http://patria.md'
response = requests.get(patria_url)
soup = BeautifulSoup(response.text, 'html.parser')

result = {}
result['upcoming'] = []
result['playing'] = {}

upcoming_section = soup.find('ul', 'coming-soon')
for item in upcoming_section.find_all('li'):
    title = item.find('a', 'title').contents[0]
    if ("EN" in title):
        movie_title = title.split('(')[0].rstrip()
        time = item.span.contents[0] 
        result['upcoming'].append({
            'title': movie_title,
            'time': time,
        })


# now playing section
cinema_ids = {"Multiplex":24, "Loteanu":36}

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
                    if movie_title not in result['playing']:
                        result['playing'][movie_title] = {}
                    if cinema_name not in result['playing'][movie_title]:
                        result['playing'][movie_title][cinema_name] = {}
                        result['playing'][movie_title][cinema_name]['schedule'] = []
                    time = link.parent.find("div", "sessions").span.contents[0]
                    result['playing'][movie_title][cinema_name]['schedule'].append(
                        datetime.strptime(movieDay + ' ' + time, '%d-%m-%Y %H:%M').strftime('%b %d, %A at %H:%M')
                    )

ia = IMDb()

imdb_keys = ['cover url', 'rating', 'canonical title', 'plot', 'synopsis', 'long imdb title', 'genres', 'runtimes']

for k, v in enumerate(result['upcoming']):
    clean_name = re.sub('(2D|3D)', '', v['title'])
    movie_data = ia.search_movie(clean_name)[0]
    ia.update(movie_data)

    for key in imdb_keys:
        if key in movie_data:
            result['upcoming'][0][key] = movie_data[key]

    # now fill in youtube data
    query_string = urllib.parse.urlencode({"search_query" : movie_data['canonical title']+' trailer'})
    html_content = urllib.request.urlopen('http://www.youtube.com/results?' + query_string)
    soup = BeautifulSoup(html_content.read().decode(), 'html.parser')
    trailer = 'http://www.youtube.com' + soup.find('div', 'yt-lockup yt-lockup-tile yt-lockup-video vve-check clearfix').find('a')['href']
    result['upcoming'][k]['trailer'] = trailer
    

for name in result['playing'].keys():
    clean_name = re.sub('(2D|3D)', '', name)
    movie_data = ia.search_movie(clean_name)[0]
    ia.update(movie_data)

    for key in imdb_keys:
        if key in movie_data:
            result['playing'][name][key] = movie_data[key]

    # now fill in youtube data
    query_string = urllib.parse.urlencode({"search_query" : movie_data['canonical title']+' trailer'})
    html_content = urllib.request.urlopen('http://www.youtube.com/results?' + query_string)
    soup = BeautifulSoup(html_content.read().decode(), 'html.parser')
    trailer = 'http://www.youtube.com' + soup.find('div', 'yt-lockup yt-lockup-tile yt-lockup-video vve-check clearfix').find('a')['href']
    result['playing'][name]['trailer'] = trailer

handle = open('result.json', 'w')
handle.write('{}')

if result:
    try:
        handle = open('result.json', 'w')
        handle.write(json.dumps(result))
    finally:
        handle.close()

