terraform {
  required_providers {
    dsfhub = {
      source = "imperva/dsfhub"
    }
  }
}

variable "dsfhub_host" {}
variable "dsfhub_token" {}

provider "dsfhub" {
  dsfhub_host  = var.dsfhub_host
  dsfhub_token = var.dsfhub_token
}

resource "dsfhub_classification" "this" {
  database_details {
    mongo_configuration {
      db_name           = "test_db"
      connection_string = "mongodb://user:password@1.2.3.4:27017"
    }
  }
  storage_details {
    s3_bucket_configuration {
      bucket_name       = "my-bucket-name"
      aws_region        = "us-east-1"
      access_key_id     = "my-access-key-id"
      secret_access_key = "my-secret-access-key"
    }
  }
}

