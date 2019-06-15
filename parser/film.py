from imdb import IMDb
from bs4 import BeautifulSoup
from datetime import datetime
import re
import urllib
import json

def register(cn, title: str):
    cursor = cn.cursor()
    meta = {}
    imdb = IMDb()
    movie = imdb.search_movie(_clean_title(title))[0]

    imdb.update(movie)
    infoset = ['cover url', 'rating', 'title', 'plot', 'long imdb title', 'genres', 'runtimes']

    # update take specific object instances
    # find out how to do it properly
    for info in infoset:
        meta[info] = movie.get(info)

    # now fill in youtube data
    meta['trailer'] = _get_yt_trailer(meta['title'])

    now = datetime.now().isoformat()
    data = (title, json.dumps(meta), now)

    cursor.execute('INSERT INTO films(title, meta, register_date) VALUES(?, ?, datetime(?))', data)
    cn.commit()

    return meta

def _get_yt_trailer(title: str) -> str:
    
    query_string = urllib.parse.urlencode({"search_query" : f'{title} trailer english'})
    html_content = urllib.request.urlopen(f'http://www.youtube.com/results?{query_string}')
    soup = BeautifulSoup(html_content.read().decode(), 'html.parser')
    trailer = 'http://www.youtube.com' + soup.find('div', 'yt-lockup yt-lockup-tile yt-lockup-video vve-check clearfix').find('a')['href']

    return trailer

def _clean_title(title: str) -> str:
    #get rid of 2D and 3D
    clean_dimension = re.sub('(2D|3D)', '', title)
    return clean_dimension.split('(')[0].rstrip()

