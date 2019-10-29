from mailjet_rest import Client
import os


def send_mail(films):
        api_key = os.environ['MJ_APIKEY_PUBLIC']
        api_secret = os.environ['MJ_APIKEY_PRIVATE']
        mailjet = Client(auth=(api_key, api_secret), version='v3.1')

        mails = [
                'slash3b@gmail.com',
                'tatiana.timoshin@gmail.com'
                ]

        data = {
        'Messages': [
                        {
                                "From": {
                                        "Email": "ilya@pirata.md",
                                        "Name": "pirata.md"
                                },
                                "To": [
                                        {
                                                "Email": "slash3b@gmail.com",
                                                "Name": "Ilya"
                                        }
                                ],
                                "Subject": "Yo! Patria has added some new films",
                                "TextPart": "Dear Subscriber, here are some new films in Patria: " + films,
                        },
                        {
                                "From": {
                                        "Email": "ilya@pirata.md",
                                        "Name": "pirata.md"
                                },
                                "To": [
                                        {
                                                "Email": "tatiana.timoshin@gmail.com",
                                                "Name": "Pus"
                                        }
                                ],
                                "Subject": "Yo! Patria has added some new films",
                                "TextPart": "Dear Subscriber, here are some new films in Patria: " + films,
                        }
                ]
        }

        result = mailjet.send.create(data=data)
        print (result.status_code)
        print (result.json())

