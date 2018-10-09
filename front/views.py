# Create your views here.
from django.http import HttpResponse
from django.template import loader

def index(request):

    template = loader.get_template('front/index.html')
    print(template)
    return
    context = {} # parse json file here and send it to the template 

    return HttpResponse(template.render(context, request))