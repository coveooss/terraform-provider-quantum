// data "quantum_list_files" "data_files" {
//   folders   = ["/Users/lpbedard/src/GOPATH/src/github.com/coveo/terraform-provider-quantum"]
//   patterns  = ["*.go"]
//   recursive = true
// }


resource "quantum_password" "test_pw" {
    name = "quantom_password_key1"
    length = 20
    // password = "<computed>"
    // created_at = "<computed>"
    
    // Complexity
    lowercase = 2
    uppercase = 2
    numbers = 5
    specials = 2

    // Rotate every n days
    expires_in_days = 90
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
  etag = "1234"
}