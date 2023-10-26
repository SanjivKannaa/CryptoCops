# from flask import Flask, request, jsonify
import torch
from datetime import datetime
from transformers import RobertaTokenizer, RobertaForSequenceClassification
from confluent_kafka import Consumer, KafkaError
import os
import pymongo
# from dotenv import load_dotenv

pyclient = pymongo.MongoClient(os.getenv("DB_URL"))
db = pyclient["test"]
collection = db["URL"]

def predict(text):
    # data = request.json

    # Tokenize the input text and perform inference
    inputs = tokenizer.encode_plus(
        text["message"],
        padding='max_length',
        max_length=128,
        truncation=True,
        return_tensors='pt',
        return_attention_mask=True
    )
    input_ids = inputs['input_ids'].to(device)
    attention_mask = inputs['attention_mask'].to(device)

    with torch.no_grad():
        outputs = model(input_ids, attention_mask=attention_mask)
        logits = outputs.logits
        predicted_class = torch.argmax(logits, dim=1).item()

    # You may want to use label_encoder.inverse_transform() to convert back to original labels if needed.
    result = {
        'url': text["url"],
        'predicted_class': predicted_class
    }
    print(result)
    records = collection.find({"url":text["url"]})
    if len(records) == 0 :
        result = {
            "url":text["url"],
            "data":text["message"],
            "susScore": predicted_class,
            "lastCrawled": datetime.now(),
            "crawlCount": 1
        }
        if(predicted_class>=6):
            result["isSuspicious"] = True
        try:
            collection.insert_one(result)
        except Exception as e:
            print(f"Error in inserting data into database: {e}")
    else:
        newrec = records[0]
        newrec["data"] = text["message"]
        newrec["susScore"] = predicted_class
        newrec["crawlCount"]+=1
        newrec["lastCrawled"] = datetime.now()
        if(predicted_class>=6):
            newrec["isSuspicious"] = True
        newvalues = {"$set":newrec}
        try:
            collection.update_one({"url":text["url"]}, newvalues)
        except Exception as e:
            print(f"Error in updating data into database: {e}")
        
        
    




conf = {
    'bootstrap.servers': os.getenv("KAFKA_URL"),
    'group.id': "my-consumer-group2",
    'auto.offset.reset': 'smallest'
}

consumer = Consumer(conf)


# Load the tokenizer and model
tokenizer = RobertaTokenizer.from_pretrained('fine_tuned_model_roberta')
# tokenizer = RobertaTokenizer.from_pretrained
model = RobertaForSequenceClassification.from_pretrained('fine_tuned_model_roberta')
device = 'cuda' if torch.cuda.is_available() else 'cpu'
model.to(device)

topic = os.getenv("KAFKA_TOPIC4")
consumer.subscribe([topic])
print("waiting for message")
while True:
    msg = consumer.poll(500.0)  # Wait for at most 500 second for a message
    if msg is None:
        continue

    if msg.error():
        if msg.error().code() == KafkaError._PARTITION_EOF:
            print('Reached end of partition')
        else:
            print(f'Error while consuming message: {msg.error()}')
    else:
        print(f'Received message: {msg.value().decode("utf-8")}')
        predict(msg.value())