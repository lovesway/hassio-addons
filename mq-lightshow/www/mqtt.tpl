{{define "content"}}
<h1>MQTT Client</h1>
{{if .MQTTConnected}}
MQTT client connected to {{.MQTTHost}}
<a href="mqtt?disconnect" class="btn btn-danger btn-sm">Disconnect</a>
{{else}}
MQTT client disconnected 
<a href="mqtt?connect" class="btn btn-primary btn-sm">Connect</a>
{{end}}

<br>
<br>
<div>
<h2>Message Log</h2>
    <ul id="log">
    </ul>
</div>
<script>
$(document).ready(function() {
    $(function(){
      return;
        setInterval(function(){
            $.getJSON( "mqtt/log/", function( data ) {
              var $log = $('#log');
              $.each( data, function( key, val ) {
                $log.prepend( "<li>" + val + "</li>" );
              });
            });

        },5000);
    });
});
</script>
{{end}}
