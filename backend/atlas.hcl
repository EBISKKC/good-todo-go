// Environment configuration for Atlas migrations

variable "db_user" {
  type    = string
  default = getenv("POSTGRES_DB_USER")
}

variable "db_password" {
  type    = string
  default = getenv("POSTGRES_DB_PASSWORD")
}

variable "db_name" {
  type    = string
  default = getenv("POSTGRES_DB_NAME")
}

variable "db_port" {
  type    = string
  default = getenv("POSTGRES_DB_PORT")
}

variable "db_host" {
  type    = string
  default = getenv("POSTGRES_DB_HOST")
}

env "local" {
  // Main database URL (target database)
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=disable"

  // Dev database URL for schema diff comparison
  dev = "docker://postgres/17/dev?search_path=public"

  // Migration directory
  migration {
    dir = "file://internal/ent/migrate/migrations"
  }

  // Diff configuration
  diff {
    skip {
      drop_schema = true
      drop_table  = false
    }
  }

  // Format configuration
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=require"

  migration {
    dir = "file://internal/ent/migrate/migrations"
  }

  diff {
    skip {
      drop_schema = true
      drop_table  = true
    }
  }
}
