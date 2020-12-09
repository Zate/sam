FROM python:3.8-alpine

RUN apk add --no-cache build-base zlib-dev libxml2-dev libxslt-dev libjpeg libmagic jpeg-dev \
                       zlib-dev \
                       freetype-dev \
                       lcms2-dev \
                       openjpeg-dev \
                       tiff-dev \
                       tk-dev \
                       tcl-dev && \
    /usr/local/bin/python -m pip install --upgrade pip && \
    pip install --no-cache https://download.splunk.com/misc/appinspect/splunk-appinspect-latest.tar.gz && \
    apk del build-base

WORKDIR /src

# ENTRYPOINT ["/bin/sh", "-c"]

CMD ["splunk-appinspect", "inspect", "/src/"]
