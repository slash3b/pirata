FROM python:3-alpine

ADD . /pirata

WORKDIR /pirata

RUN apk add --update --no-cache alpine-sdk \
                     libxml2-dev \
                     libxslt-dev \
                     openssl-dev \
                     libffi-dev \
                     zlib-dev \
                     py-pip

RUN pip install -r requirements.txt
