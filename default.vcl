vcl 4.0;

backend default {
	.host = "127.0.0.1";
}

backend nginx {
    .host = "backend";
    .port = "80";
}

sub vcl_recv {
    if (req.url ~ "^/status") {
		return (synth(200));
    } else {
        set req.backend_hint = nginx;
    }
}

sub vcl_synth {
	synthetic({"
<html>
<head>
<title>Varnish + Go == <3</title>
</head>
<body>
<h1>Varnish + Go == <span style="color:pink"><3</span></h1>
<pre>
<p>VERSION="4.1"</p>
<p>HTTP_COOKIE="} + req.http.cookie + {"</p>
<p>HTTP_HOST="} + req.http.host + {"</p>
<p>HTTP_REFERER="} + req.http.referer + {"</p>
<p>HTTP_USER_AGENT="} + req.http.user-agent + {"</p>
<p>PATH="} + regsub(req.url, "\?.*", "") + {"</p>
<p>QUERY_STRING="} + regsub(req.url, "[^\?]*\??", "") + {"</p>
</pre>
</body>
</html>
"});
	return (deliver);
}
