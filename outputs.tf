output "instance_public_ip" {
  value = aws_instance.main.public_ip                                        # The actual value to be outputted
  description = "The public IP address of the EC2 instance" # Description of what this output represents
}

output "instance_id" {
  value = aws_instance.main.id
  description = "The id of the EC2 instance"
}