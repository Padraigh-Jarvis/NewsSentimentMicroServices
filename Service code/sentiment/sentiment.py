# -*- coding: utf-8 -*-
"""
Created on Wed Mar  7 15:35:29 2018

@author: Padraigh Jarvis
"""
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import threading
import pika
import json

hostName ='rabbitmq'
inqueueName='sentimentIn'
outqueueName='toDatabase'

class myThread (threading.Thread):
   def __init__(self, jsonObj):
      threading.Thread.__init__(self)
      self.jsonObj = jsonObj
   def run(self):
      handleMessage(self.jsonObj)


def sendSentament(jsonObj):
    connection = pika.BlockingConnection(pika.ConnectionParameters(host=hostName))
    channel = connection.channel()

    channel.queue_declare(queue=outqueueName)
    channel.basic_publish(exchange='',
                      routing_key=outqueueName,
                      body=jsonObj)
    print("Send",jsonObj)
    connection.close()


def callback(ch, method, properties, body):
    print("\n [x] Received:",body.decode('utf-8'),"\n")
    thread = myThread(body)
    thread.start()
    


def handleMessage(jsonObj):
    jsonDecoded = json.loads(jsonObj)
    sentiment = analize(jsonDecoded['Content'])
    jsonDecoded['Sentiment']=sentiment
    jsonObj= json.dumps(jsonDecoded)
    sendSentament(jsonObj)
    
    
    
def analize(content):
    analyzer = SentimentIntensityAnalyzer()        
    vs = analyzer.polarity_scores(content)
    return str(vs["compound"])
   
    
    
def setUpQueue():
    connection = pika.BlockingConnection(pika.ConnectionParameters(hostName))
    channel = connection.channel()
    channel.queue_declare(queue=inqueueName)
    
    channel.basic_consume(callback,
                      queue=inqueueName,
                      no_ack=True)
    return channel

if __name__ == '__main__':
    
    channel = setUpQueue()
    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()

        
