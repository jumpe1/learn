```bash
copy terraform.tfvars.example terraform.tfvars
gcloud compute ssh vm2 --project=xxx --zone=us-central1-a
sudo apt-get update
sudo apt-get install tcpdump -y
sudo tcpdump -i any icmp

gcloud compute ssh vm1 --project=xxx --zone=us-central1-a
ping <vm2_internal_ip>
```