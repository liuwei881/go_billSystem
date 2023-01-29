# config=utf-8

import pika
credentials = pika.PlainCredentials('admin', '123456')
connection = pika.BlockingConnection(pika.ConnectionParameters(
    '10.96.141.214', 5672, '/', credentials))
channel = connection.channel()

channel.queue_declare(queue='balance')

channel.basic_publish(exchange='', routing_key='balance', body='Hello World!')
print(" [x] Sent 'Hello World!'")
connection.close()