package util

const SqlConsole = `
<html>
<head>
<script src="//ajax.googleapis.com/ajax/libs/jquery/1.9.0/jquery.min.js"></script>

<script type="text/javascript">
$(function(){
  $("#command").focus();
  $("form").submit(function(){
    $.post("sqlConsole", {command: $("#command").val()}, function(data){
      $("#result").html(data)
      $("#command").focus();
    });
	return false;
  });
});
</script>
</head>

<body bgcolor="#cdc0cd">
<form>
  Command: <textarea name="command" rows="8" cols="80" id="command"></textarea><p>
  <input type="submit" name="submit" value="submit"/>
</form>
(Query does not support datetime column like 'Last_Update')
<p>
<b>System has following tables.</b>

<pre>
CheckedOut_Book_History  Profile                  Stock_Transaction
Library_Search_Keyword   RTHK_PODCAST             trans
NHK_WORLD_Daily_News     RTHK_Radio
NHK_WORLD_News_mp3       STOCK
</pre>

<div id="result">
</div>
</html>`
