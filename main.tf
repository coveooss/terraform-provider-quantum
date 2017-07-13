resource "quantum_password" "test_pw" {
    name = "quantom_password_key1"
    length = 10
  
    // Rotate every n days
    expires_in_days = 2
}

resource "aws_ssm_parameter" "quantom_password_key1" {
  name  = "quantom_password_key1"
  type  = "SecureString"
  value = "${quantum_password.test_pw.password}"
}

resource "aws_s3_bucket_object" "index" {
  bucket = "coveo-ndev-test"
  key = "quantom_password_key1"
  content = "${aws_ssm_parameter.quantom_password_key1.value}"
}