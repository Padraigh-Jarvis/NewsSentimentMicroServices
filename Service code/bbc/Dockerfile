#Parent image
FROM python:3 
#Maintainer 
MAINTAINER padraigh.jarvis
#Working dir
WORKDIR /BBC
#CP current dir contents to working dir
ADD . /BBC
#Install required packages
RUN pip install --trusted-host pypi.python.org -r requirements.txt
#Run app when container lanuches
CMD ["python", "bbcnews.py"]
