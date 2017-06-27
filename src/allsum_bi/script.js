var source = postgres({
  "uri": "postgres://rm-uf6q1kk0byn74g70zo.pg.rds.aliyuncs.com:3432/db_logistic?sslmode=disable&user=user_logistic&password=AllSum123"
  "debug": true
})

var dest = postgres({
  "uri": "postgres://rm-uf6q1kk0byn74g70zo.pg.rds.aliyuncs.com:3432/etl_test?sslmode=disable&user=user_logistic&password=AllSum123"
  "debug": true
})

t.Source("source", source, "route.route_base").Save("dest", dest, "public.route_base")
