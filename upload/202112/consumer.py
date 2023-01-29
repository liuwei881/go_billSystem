# coding=utf-8

import pika

credentials = pika.PlainCredentials('admin', '123456')
connection = pika.BlockingConnection(pika.ConnectionParameters(
    '10.96.141.214', 5672, '/', credentials))
channel = connection.channel()

channel.queue_declare(queue='balance')


def callback(ch, method, properties, body):
    print(" [x] Received %r" % body)


channel.basic_consume(callback, queue='balance', no_ack=True)

print(' [*] Waiting for messages. To exit press CTRL+C')
channel.start_consuming()