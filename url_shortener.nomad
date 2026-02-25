job "url-shortener"{
    datacenters=["dc1"]
    group "app-group"{
        network{
            port "http"{
                static=8080
            }
        }

        count=1

        task "app"{
            driver="raw_exec"
            config{
                command="/Users/timmareddydeekshitha/url_shortener/url_shortener"
            }
            env {
                DATABASE_DSN = "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"
            }
            resources{
                cpu=10
                memory=32
            }
        }
    }
}