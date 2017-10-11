package main

const (
	sniperVersionInfo = `This is Sniper ,version 1.0 
Copyright (C) 2013 by btfak.com, et al.
This is free software; open source on github.com/btfak/sniper.`
	sniperUsageInfo = `Usage: 
   sniper [options] http[s]://hostname[:port][/path]
   sniper [options] -f urls.txt
Options:  
   -c, --concurrent     concurrent users, default is 1.
   -n, --requests       number of requests to perform.
   -r, --repetitions    number of times to run the test.
   -t, --time           testing time, 30 mean 30 seconds.
   -R, --sniperc        specify an sniperc file to get config
                        (default is $HOME/.sniperc).
   -f, --urlfile        select a specific URLS file.                     
   -p, --post           select a specific file to POST.
   -T, --content-type   set Content-Type in request
                        (default is text/plain).
   -V, --Version        print the version number.
   -h, --help           print this section.
   -C, --config         show the current config.
   -s, --plot           plot detail transactions' info
                        (true | false,default set true,
                        notice: set -t will not plot anyhow).
Copyright (C) 2013 by btfak.com, et al.
This is free software; open source on github.com/btfak/sniper.
`
	sniperNoDataInfo        = "file no data."
	sniperOpenFileErrorInfo = "open file error."
)

var plotsTemplate = `<!doctype>
<html>
<head>
  <title>Sniper Plots</title>
</head>
<body>
  	<div id="sniper" style="font-family: Courier; width: 100%%; height: 600px"></div>
  	<p>
  		<b>Display: </b>
		<input type=checkbox id=0 onClick="stockchange(this)" checked>
		<label for="trans"> transaction</label>
		<input type=checkbox id=1 onClick="stockchange(this)" unchecked>
		<label for="conn"> connection</label>
		<input type=checkbox id=2 onClick="stockchange(this)" unchecked>
		<label for="resp"> response</label>
	</p>
  <script>
	%s
  </script>
  
  <script>
  
  stockchart = new Dygraph(
    document.getElementById("sniper"),
    [%s],
    {
    	title: 'Sniper Plot',
    	labels: ['Seconds', 'Transaction(ms)','Connection(ms)','Response(ms)'],
    	ylabel: 'Latency(ms)',
      	xlabel: 'Elapse time',
      	fillGraph: true,
      	visibility: [true, false, false,false],
    	logscale: true
    }
    );
      
  function stockchange(el) {
  	stockchart.setVisibility(el.id, el.checked);
  }		
  </script>
</body>
</html>`
