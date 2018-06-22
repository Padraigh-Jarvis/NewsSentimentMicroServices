# -*- coding: utf-8 -*-
"""
Created on Tue Apr 10 16:57:25 2018

@author: Padraigh Jarvis
"""
from newsapi import NewsApiClient
from datetime import datetime, timedelta
import pika
import json
import time

apiKey = #Enter API KEY HERE
source = " bbc-news"
query = "brexit"
outQueue = "sentimentIn"
hostName = "rabbitmq"

def myMain():

    
    
    
    
    newsApi = NewsApiClient(api_key=apiKey)
    
    toDate=datetime.now()
    fromDate=toDate-timedelta(hours=6)
    toDate = toDate.strftime("%Y-%m-%dT%H:%M:%SZ")
    fromDate = fromDate.strftime("%Y-%m-%dT%H:%M:%SZ")
    newsArticals = getNews(newsApi, fromDate, toDate)
  
    parsedArticals = parseData(newsArticals)
    sendArticals(parsedArticals)
    
                
def sendArticals(articals):
    connection = pika.BlockingConnection(pika.ConnectionParameters(host=hostName))
    channel = connection.channel()
    channel.queue_declare(queue=outQueue)
    for artical in articals:
        channel.basic_publish(exchange='',
                              routing_key=outQueue,
                              body=artical)
        print(" [x] Sent ",artical,"\n")
        time.sleep(1)
    connection.close()
def parseData(newsArticals):
    articals=[]
    for rawArtical in newsArticals["articles"]:
        data = {}
        title = rawArtical["title"]
        description = rawArtical["description"]
        articalString = title+" " + description
        articalString = articalString.replace("'","’")
        data['Source'] = 'bbc'
        data['Content'] = articalString
        
        
        utc_time = datetime.strptime(rawArtical['publishedAt'],"%Y-%m-%dT%H:%M:%SZ")
        milliseconds = (utc_time - datetime(1970, 1, 1)) // timedelta(milliseconds=1)
       
        data['Time']=milliseconds
        json_data = json.dumps(data)
        articals.append(json_data)
      
    return articals

def getNews(newsApi, fromDate, toDate):
    print("FROM:",fromDate," TO:",toDate)
    return newsApi.get_everything(q=query,
                                  sources=source,
                                  from_param=fromDate.__str__(),
                                  to=toDate.__str__(),
                                  language='en')

if __name__ == '__main__':
    myMain()
