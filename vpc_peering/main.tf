resource "google_compute_network" "vpc1" {
  name                    = "vpc1"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet1" {
  name          = "subnet1"
  network       = google_compute_network.vpc1.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = var.region
}

resource "google_compute_network" "vpc2" {
  name                    = "vpc2"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet2" {
  name          = "subnet2"
  network       = google_compute_network.vpc2.self_link
  ip_cidr_range = "10.1.0.0/16"
  region        = var.region
}

resource "google_compute_network_peering" "peer1" {
  name         = "peer1"
  network      = google_compute_network.vpc1.self_link
  peer_network = google_compute_network.vpc2.self_link
}

resource "google_compute_network_peering" "peer2" {
  name         = "peer2"
  network      = google_compute_network.vpc2.self_link
  peer_network = google_compute_network.vpc1.self_link
}

resource "google_compute_firewall" "allow_vpc1_to_vpc2" {
  name    = "allow-vpc1-to-vpc2"
  network = google_compute_network.vpc1.name

  allow {
    protocol = "icmp"
  }

  source_ranges = [google_compute_subnetwork.subnet2.ip_cidr_range]
}

resource "google_compute_firewall" "allow_vpc2_to_vpc1" {
  name    = "allow-vpc2-to-vpc1"
  network = google_compute_network.vpc2.name

  allow {
    protocol = "icmp"
  }

  source_ranges = [google_compute_subnetwork.subnet1.ip_cidr_range]
}

resource "google_compute_firewall" "allow-ssh-vpc1" {
  name    = "allow-ssh-vpc1"
  network = google_compute_network.vpc1.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "allow-ssh-vpc2" {
  name    = "allow-ssh-vpc2"
  network = google_compute_network.vpc2.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_instance" "vm1" {
  name         = "vm1"
  machine_type = "f1-micro"
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
    }
  }

  network_interface {
    network    = google_compute_network.vpc1.self_link
    subnetwork = google_compute_subnetwork.subnet1.self_link
    access_config {
    }
  }
}

resource "google_compute_instance" "vm2" {
  name         = "vm2"
  machine_type = "f1-micro"
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
    }
  }

  network_interface {
    network    = google_compute_network.vpc2.self_link
    subnetwork = google_compute_subnetwork.subnet2.self_link
    access_config {
    }
  }
}
