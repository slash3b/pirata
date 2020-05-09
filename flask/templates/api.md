Humble pirata.md API just for fun!

## Upcoming movies endpoint: 
`[GET] /api/upcoming` — a list of soon to be played movies

**Output example:**
```
[
  {
    "meta": {
      "cover url": "https://m.media-amazon.com/images/M/MV5BZGRlNTY3NGYtM2YzZS00N2YyLTg0ZDYtNmY2ZDg2NDM3N2JlXkEyXkFqcGdeQXVyNTI4MzE4MDU@._V1_SY150_CR0,0,101,150_.jpg",
      "genres": [
        "Action",
        "Adventure",
        "Sci-Fi"
      ],
      "long imdb title": "Black Widow (2020)",
      "plot": [
        "A film about Natasha Romanoff in her quests between the films Civil War and Infinity War."
      ],
      "rating": null,
      "runtimes": null,
      "title": "Black Widow",
      "trailer": "http://www.youtube.com/watch?v=RxAtuMu_ph4"
    },
    "premiere_date": "2020-11-05",
    "title": "Black Widow"
  }
]
```

## Playing movies endpoint: 
`[GET] /api/playing` — return all playing movies right now in json format  
In schedule section 24 stands for Multiplex and 36 for Loteanu locations. 
**Output example:** 
```
[
  {
    "meta": {
      "cover url": "https://m.media-amazon.com/images/M/MV5BOTVjMmFiMDUtOWQ4My00YzhmLWE3MzEtODM1NDFjMWEwZTRkXkEyXkFqcGdeQXVyMTkxNjUyNQ@@._V1_SY150_CR0,0,101,150_.jpg", 
      "genres": [
        "Action", 
        "Adventure", 
        "Comedy", 
        "Fantasy"
      ], 
      "long imdb title": "Jumanji: The Next Level (2019)", 
      "plot": [
        "In Jumanji: The Next Level, the gang is back but the game has changed. As they return to rescue one of their own, the players will have to brave parts unknown from arid deserts to snowy mountains, to escape the world's most dangerous game.::Sony Pictures Entertainment", 
        "The gang is back but the game has changed. As they return to Jumanji to rescue one of their own, they discover that nothing is as they expect. The players will have to brave parts unknown and unexplored, from the arid deserts to the snowy mountains, in order to escape the world's most dangerous game.", 
        "\"When Spencer goes back into the fantastical world of Jumanji, pals Martha, Fridge and Bethany re-enter the game to bring him home. But everything about Jumanji is about to change, as they soon discover more obstacles and more danger to overcome.::krmanirethnam"
      ], 
      "rating": 7.0, 
      "runtimes": [
        "123"
      ], 
      "title": "Jumanji: The Next Level", 
      "trailer": "http://www.youtube.com/watch?v=rBxcF-r9Ibs"
    }, 
    "schedule": {
      "24": [
        "2021-12-12 19:30:00", 
        "2021-12-12 19:30:00", 
        "2021-12-12 19:30:00", 
        "2021-12-12 19:30:00", 
        "2021-12-12 19:30:00", 
        "2021-12-12 19:30:00"
      ]
    }, 
    "title": "Jumanji: The next level"
  }
]
```
