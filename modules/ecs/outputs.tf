output "ec2_url" {
  value = aws_instance.ec2.public_ip
}

output "ec2_dns" {
  value = aws_instance.ec2.public_dns
}

output "dns_name" {
  value = "${cloudflare_dns_record.dns_record.name}.kuraiji.me"
}