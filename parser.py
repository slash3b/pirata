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
client.get('https://cineplex.md/lang/en')

url = f'https://cineplex.md/films'
response = client.get(url)
soup = BeautifulSoup(response.text, 'html.parser')

all_movies = soup.find_all('div', attrs={'class':'movies_fimls_item'})
movie_titles_to_email = ''

for movie in all_movies:
    movie_lang = movie.find('span', class_='overlay__lang')
    if movie_lang == None:
        continue
    movie_lang = movie_lang.string

    movie_title = movie.find('h3', class_='overlay__title')
    if movie_title == None:
        continue
    movie_title = movie_title.string

    if ('(EN)' in movie_lang.string):
        upcoming_movie_response = client.get(movie["data-href"])
        upcoming_html_soup = BeautifulSoup(upcoming_movie_response.text, 'html.parser')
        movie_attributes = upcoming_html_soup.find_all('li', attrs={'class':'film_movies_info_item'})

        for attribute in movie_attributes:
            if 'Premiere date' in attribute.find('h5').string:
                premiere_date = attribute.find('p').string

                premiere_date = datetime.strptime(premiere_date, '%d.%m.%Y').date() 


                if not upcoming.is_registered(cn, movie_title):

                    upcoming.register(cn, movie_title, premiere_date)
                    movie_titles_to_email += movie_title + ', '
                break

if(len(movie_titles_to_email)):
    mail.send_mail(movie_titles_to_email)

# terminate connection
cn.close()

print("All is good!")
