# Create your views here.
from django.http import HttpResponse
from django.template import loader
import json

def index(request):

    template = loader.get_template('front/index.html')

    # grab json and convert it to array and that is it
    
    with open('result.json') as file:
        schedule = json.load(file)

    # print(films)
    # return HttpResponse(films)

    context = {
        'schedule' : schedule,
    } # parse json file here and send it to the template 

    return HttpResponse(template.render(context, request))