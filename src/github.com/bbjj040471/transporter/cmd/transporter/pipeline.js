var source = postgres({"uri": "postgres://localhost:5432/bi?sslmode=disable&user=hebaoliang&password=123456"})
var sink = postgres({"uri": "postgres://localhost:5432/admin?sslmode=disable&user=hebaoliang&password=123456"})
t.Source( source, "/public.route_base\/"  ).Transform("skip",skip({"field": "originid", "operator": ">", "match": 414})).Save( sink, "/public.route_base_t/"  )
