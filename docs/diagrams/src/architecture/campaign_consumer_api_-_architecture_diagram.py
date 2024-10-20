import os, sys
from urllib.request import urlretrieve

os.chdir(os.path.dirname(sys.argv[0]))

from diagrams import Cluster, Diagram, Edge
from diagrams.k8s.compute import Pod
from diagrams.onprem.queue import Kafka
from diagrams.onprem.database import Postgresql
from diagrams.onprem.inmemory import Redis

with Diagram("campaign consumer api"):
   blueline = Edge(color="blue", style="bold")
   blackline = Edge(color="black", style="bold")
   bidirectional_edge = Edge(color="darkOrange", style="bold", forward=True, reverse=True)

   with Cluster("service"):
      consumerAPI = Pod("campaign-consumer-api")

   with Cluster("queue"):
      consumerkafkaOwner = Kafka("campaign.campaign-owner")
      consumerKafkaSlug = Kafka("campaign.campaign-slug")
      consumerKafkaRegion = Kafka("campaign.campaign-region")
      consumerKafkaMerchant = Kafka("campaign.campaign-merchant")
      consumerKafkaCampaign = Kafka("campaign.campaign")
      consumerKafkaSpent = Kafka("campaign.campaign-spent")
   
   with Cluster("cache"):
      consumerRedis = Redis("campaign")
   
   with Cluster("db"):
      consumerDb = Postgresql("campaign-consumer-db")

   # Definindo as conexões
   consumerAPI - blueline >> consumerkafkaOwner
   consumerAPI - blueline >> consumerKafkaSlug
   consumerAPI - blueline >> consumerKafkaRegion
   consumerAPI - blueline >> consumerKafkaMerchant
   consumerAPI - blueline >> consumerKafkaCampaign
   consumerAPI - blueline >> consumerKafkaSpent
   consumerAPI - bidirectional_edge >> consumerDb
   consumerAPI - bidirectional_edge >> consumerRedis