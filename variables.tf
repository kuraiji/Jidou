variable "instance_type" {
  type = string                     # The type of the variable, in this case a string
  default = "t2.micro"                 # Default value for the variable
  description = "The type of EC2 instance" # Description of what this variable represents
}

variable "instance_name" {
  type = string
  default = "Jidou"
  description = "The name of the EC2 instance"
}

variable "instance_ami" {
  type = string
  default = "ami-00ddc330f6182b5cb"
  description = "The ami of the EC2 instance"
}