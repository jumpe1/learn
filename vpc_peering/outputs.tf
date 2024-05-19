output "vpc1_network" {
  value = google_compute_network.vpc1.self_link
}

output "vpc2_network" {
  value = google_compute_network.vpc2.self_link
}

output "vm1_internal_ip" {
  value = google_compute_instance.vm1.network_interface.0.network_ip
}

output "vm2_internal_ip" {
  value = google_compute_instance.vm2.network_interface.0.network_ip
}
